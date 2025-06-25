package graph

import (
	"user-service/internal/database"
	"user-service/internal/model"
	"user-service/internal/utils"

	"github.com/graphql-go/graphql"
)

var userTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "UserType",
	Values: graphql.EnumValueConfigMap{
		"MANAGER": &graphql.EnumValueConfig{Value: "MANAGER"},
		"MEMBER":  &graphql.EnumValueConfig{Value: "MEMBER"},
	},
})

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.ID},
		"username": &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
		"role":     &graphql.Field{Type: userTypeEnum},
	},
})

var userMutResponse = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserMutationResponse",
	Fields: graphql.Fields{
		"code":    &graphql.Field{Type: graphql.String},
		"success": &graphql.Field{Type: graphql.Boolean},
		"message": &graphql.Field{Type: graphql.String},
		"user":    &graphql.Field{Type: userType},
		"errors":  &graphql.Field{Type: graphql.NewList(graphql.String)},
	},
})

var authMutResponse = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthMutationResponse",
	Fields: graphql.Fields{
		"code":         &graphql.Field{Type: graphql.String},
		"success":      &graphql.Field{Type: graphql.Boolean},
		"message":      &graphql.Field{Type: graphql.String},
		"accessToken":  &graphql.Field{Type: graphql.String},
		"refreshToken": &graphql.Field{Type: graphql.String},
		"user":         &graphql.Field{Type: userType},
		"errors":       &graphql.Field{Type: graphql.NewList(graphql.String)},
	},
})

var createUserInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateUserInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"username": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"email":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"role":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(userTypeEnum)},
	},
})

var userInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"email":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Args: graphql.FieldConfigArgument{
				"role": &graphql.ArgumentConfig{Type: graphql.NewNonNull(userTypeEnum)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				role := p.Args["role"].(string)
				var users []model.User
				if err := database.DB.Where("role = ?", role).Find(&users).Error; err != nil {
					return nil, err
				}
				return users, nil
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				var u model.User
				if err := database.DB.First(&u, "id = ?", id).Error; err != nil {
					return nil, err
				}
				return u, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userMutResponse,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createUserInput)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				inp := p.Args["input"].(map[string]interface{})
				pwHash, err := utils.HashPassword(inp["password"].(string))
				if err != nil {
					return map[string]interface{}{"code": "500", "success": false, "message": "hash error"}, nil
				}
				u := model.User{
					Username: inp["username"].(string),
					Email:    inp["email"].(string),
					Password: pwHash,
					Role:     inp["role"].(string),
				}
				if err := database.DB.Create(&u).Error; err != nil {
					return map[string]interface{}{
						"code":    "400",
						"success": false,
						"errors":  []string{err.Error()},
					}, nil
				}
				return map[string]interface{}{
					"code":    "200",
					"success": true,
					"message": "User created",
					"user":    u,
				}, nil
			},
		},
		"login": &graphql.Field{
			Type: authMutResponse,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(userInput)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				inp := p.Args["input"].(map[string]interface{})
				var u model.User
				if err := database.DB.Where("email = ?", inp["email"].(string)).First(&u).Error; err != nil {
					return map[string]interface{}{"code": "400", "success": false, "message": "user not found"}, nil
				}
				if err := utils.CheckPassword(u.Password, inp["password"].(string)); err != nil {
					return map[string]interface{}{"code": "400", "success": false, "message": "invalid credentials"}, nil
				}
				access, _ := utils.GenerateAccessToken(u.ID.String())
				refresh, _ := utils.GenerateRefreshToken(u.ID.String())
				return map[string]interface{}{
					"code":         "200",
					"success":      true,
					"message":      "logged in",
					"accessToken":  access,
					"refreshToken": refresh,
					"user":         u,
				}, nil
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
