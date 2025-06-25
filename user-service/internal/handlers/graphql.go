package handlers

import (
	"net/http"

	"user-service/internal/graph"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func Handler(c *gin.Context) {
	var body struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := graphql.Do(graphql.Params{
		Schema:         graph.Schema,
		RequestString:  body.Query,
		VariableValues: body.Variables,
		Context:        c,
	})
	status := http.StatusOK
	if len(res.Errors) > 0 {
		status = http.StatusBadRequest
	}
	c.JSON(status, res)
}
