package graphqlgotask

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/graphql-go/graphql"
)

type (
	Task struct {
		ID        string
		Name      string
		Body      string
		ProjectID string `db:"projectID"`
	}
)

const (
	argTaskID = "taskID"
)

func NewTaskField(db *goqu.Database) *graphql.Field {
	return &graphql.Field{
		Type: NewTaskObject(),
		Args: graphql.FieldConfigArgument{argTaskID: &graphql.ArgumentConfig{Type: graphql.String}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			taskID := p.Args[argTaskID]
			task := new(Task)
			found, err := db.Select("Task.id", "Task.name").From("Task").
				Where(goqu.Ex{"Task.id": taskID}).ScanStruct(task)
			if err != nil {
				return nil, err
			}

			if !found {
				return nil, errors.New("task not found")
			}

			return task, nil
		},
	}
}

func NewTaskObject() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"taskID": &graphql.Field{
				Type: NonNullString,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(*Task); ok {
						return task.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(*Task); ok {
						return task.Name, nil
					}
					return nil, nil
				},
			},
			"body": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(*Task); ok {
						return task.Body, nil
					}
					return nil, nil
				},
			},
			"ProjectID": &graphql.Field{
				Type: NonNullString,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if task, ok := p.Source.(*Task); ok {
						return task.ProjectID, nil
					}
					return nil, nil
				},
			},
		},
	})
}
