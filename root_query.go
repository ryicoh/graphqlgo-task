package graphqlgotask

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/graphql-go/graphql"
)

func NewRootQuery(db *goqu.Database) *graphql.ObjectConfig {
	return &graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			//	"Task": &graphql.Field{
			//		Type: taskObject,
			//		Args: graphql.FieldConfigArgument{
			//			"taskID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(uuid)},
			//		},
			//		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//			switch p.Args["taskID"] {
			//			case task1.userID:
			//				return task1, nil
			//			case task2.taskID:
			//				return task2, nil
			//			}
			//			return nil, nil
			//		},
			//	},
			"User": NewUserField(db),
		},
	}
}
