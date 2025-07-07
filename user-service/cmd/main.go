package main

import (
	"fmt"
	"user-service/internal/http/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"pkg/config"
	"pkg/database"
	"pkg/logger"
	"user-service/graphql"
	"user-service/graphql/generated"
	"user-service/internal/user"
)

func main() {
	config.Init()
	logger.Init()
	defer logger.Log.Sync()
	database.Connect()
	// use flyway instead
	//database.DB.AutoMigrate(&user.User{})

	repo := user.NewRepository()
	service := user.NewService(repo)
	resolver := &graphql.Resolver{UserService: service}

	env := config.Cfg.ENV
	host := config.Cfg.Host
	port := config.Cfg.UserServicePort

	addr := fmt.Sprintf("%s:%s", host, port)
	fullURL := fmt.Sprintf("http://%s/", addr)

	logger.Log.Info("application starting",
		zap.String("env", env),
		zap.String("url", fullURL),
	)

	r := gin.Default()
	r.Use(middleware.GinCookieMiddleware())
	r.Use(graphql.AuthMiddleware())
	/* r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}) */

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Options{})
	srv.Use(extension.Introspection{})

	r.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL", "/graphql")(c.Writer, c.Request)
	})

	//r.Run(host + ":" + port)
	if err := r.Run(addr); err != nil {
		logger.Log.Fatal("failed to run server",
			zap.Error(err),
			zap.String("addr", addr),
		)
	}
}
