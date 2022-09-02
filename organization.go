package graphqlgotask

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/graphql-go/graphql"
)

type (
	Organization struct {
		ID                  string
		Name                string
		OrganizationMembers []User
	}
)

const (
	argOrganizationID           = "organizationID"
	tableNameOrganization       = "Organization"
	tableNameOrganizationMember = "OrganizationMember"
	columnNameOrganizationID    = "id"
	columnNameOrganizationName  = "name"
)

func NewOrganizationField(db *goqu.Database) *graphql.Field {
	return &graphql.Field{
		Type: NewOrganizationObject(),
		Args: graphql.FieldConfigArgument{argOrganizationID: &graphql.ArgumentConfig{Type: graphql.String}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			organizationID := p.Args[argOrganizationID]
			user := new(Organization)
			found, err := db.Select(columnNameOrganizationID).From(tableNameOrganization).Where(goqu.Ex{columnNameOrganizationID: organizationID}).ScanStruct(user)
			if err != nil {
				return nil, err
			}

			if !found {
				return nil, errors.New("user not found")
			}

			return user, nil
		},
	}
}

func NewOrganizationObject() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Organization",
		Fields: graphql.Fields{
			"organizationID": &graphql.Field{
				Type: NonNullString,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*Organization); ok {
						return user.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*Organization); ok {
						return user.Name, nil
					}
					return nil, nil
				},
			},
			//"tasks": &graphql.Field{
			//	Type: graphql.NewList(NewTaskObject),
			//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//		if user, ok := p.Source.(*Organization); ok {
			//			return user.tasks, nil
			//		}
			//		return nil, nil
			//	},
			//},
		},
	})
}
