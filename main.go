package main

import (
	"fmt"
	"log"
	"net/http"

	gu "github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/handler"
)

type (
	User struct {
		userID gu.UUID
		name   string
		tasks  []Task
	}

	Task struct {
		taskID gu.UUID
		name   string
		userID gu.UUID
	}
)

func coerceUUID(value interface{}) interface{} {
	v, ok := value.(gu.UUID)
	if !ok {
		return gu.Nil
	}

	return v
}

var (
	uuid = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "uuid",
		Description: "uuid",
		Serialize:   coerceUUID,
		ParseValue:  coerceUUID,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				u, err := gu.Parse(valueAST.Value)
				if err != nil {
					return nil
				}
				return u
			}
			return nil
		},
	})

	user1 = User{gu.MustParse("4624DBD4-F795-4FD9-9C02-3FB0539B6808"), "user1taro", []Task{}}
	user2 = User{gu.MustParse("F3CF9793-35A7-43EF-896A-977387EF563B"), "user2taro", []Task{}}
	task1 = Task{gu.MustParse("AC1EBA9D-8E63-4E8C-9B89-365F718D92E5"), "task1todo", user1.userID}
	task2 = Task{gu.MustParse("E48A1F07-2DE5-49CB-A31B-F69FD0466C80"), "task2todo", user2.userID}

	taskObject     *graphql.Object
	userObject     *graphql.Object
	userWhereInput *graphql.InputObject
)

func init() {
	user1.tasks = []Task{task1, task2}
	user2.tasks = []Task{task2}

	taskObject = graphql.NewObject(graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"taskID": &graphql.Field{
				Type: graphql.NewNonNull(uuid),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(Task); ok {
						return task.taskID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(Task); ok {
						return task.name, nil
					}
					return nil, nil
				},
			},
			"userID": &graphql.Field{
				Type: graphql.NewNonNull(uuid),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(Task); ok {
						return task.userID, nil
					}
					return nil, nil
				},
			},
		},
	})

	userObject = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"userID": &graphql.Field{
				Type: graphql.NewNonNull(uuid),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.userID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.name, nil
					}
					return nil, nil
				},
			},
			"tasks": &graphql.Field{
				Type: graphql.NewList(taskObject),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.tasks, nil
					}
					return nil, nil
				},
			},
		},
	})

	userInputWhereObjectEq := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UserInputWhereEq",
		Fields: graphql.InputObjectConfigFieldMap{
			"_eq": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(uuid),
			},
		},
	})
	userWhereInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UserInputWhere",
		Fields: graphql.InputObjectConfigFieldMap{
			"userID": &graphql.InputObjectFieldConfig{
				Type: userInputWhereObjectEq,
			},
		},
	})
}

func main() {
	fields := graphql.Fields{
		"Tasks": &graphql.Field{
			Type: graphql.NewList(taskObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return []Task{task1, task2}, nil
			},
		},
		"Task": &graphql.Field{
			Type: taskObject,
			Args: graphql.FieldConfigArgument{
				"taskID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(uuid)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch p.Args["taskID"] {
				case task1.userID:
					return task1, nil
				case task2.taskID:
					return task2, nil
				}
				return nil, nil
			},
		},
		"Users": &graphql.Field{
			Type: userObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return []User{user1, user2}, nil
			},
		},
		"User": &graphql.Field{
			Type: userObject,
			Args: graphql.FieldConfigArgument{"where": &graphql.ArgumentConfig{Type: userWhereInput}},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Args["where"].(map[string]any)["userID"].(map[string]any)["_eq"].(gu.UUID)
				switch userID {
				case user1.userID:
					return user1, nil
				case user2.userID:
					return user2, nil
				}
				return nil, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("server started")
	http.ListenAndServe("127.0.0.1:8084", nil)
}
