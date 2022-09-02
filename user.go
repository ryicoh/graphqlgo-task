package graphqlgotask

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/graphql-go/graphql"
)

type (
	User struct {
		ID                  string
		Name                string
		Tasks               []*Task
		OrganizationMembers []*Organization
	}
)

const (
	argUserID = "userID"
)

func NewUserField(db *goqu.Database) *graphql.Field {
	return &graphql.Field{
		Type: NewUserObject(db),
		Args: graphql.FieldConfigArgument{argUserID: &graphql.ArgumentConfig{Type: graphql.String}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Args[argUserID]
			user := new(User)
			found, err := db.Select("User.id", "User.name").From("User").
				Where(goqu.Ex{"User.id": userID}).ScanStruct(user)
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

func NewUserObject(db *goqu.Database) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"userID": &graphql.Field{
				Type: NonNullString,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*User); ok {
						return user.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*User); ok {
						return user.Name, nil
					}
					return nil, nil
				},
			},
			"organizationMembers": &graphql.Field{
				Type: graphql.NewList(NewOrganizationObject()),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*User)
					if !ok {
						return nil, nil
					}

					organizations := make([]*Organization, 0)
					if err := db.Select("Organization.id", "Organization.name").
						From(tableNameOrganization).
						Join(goqu.T(tableNameOrganizationMember), goqu.On(goqu.Ex{"OrganizationMember.organizationID": goqu.I("Organization.id")})).
						Where(goqu.Ex{"OrganizationMember.userID": user.ID}).
						ScanStructs(&organizations); err != nil {
						return nil, err
					}

					return organizations, nil
				},
			},
			"tasks": &graphql.Field{
				Type: graphql.NewList(NewTaskObject()),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*User)
					if !ok {
						return nil, nil
					}

					tasks := make([]*Task, 0)
					if err := db.Select("Task.id", "Task.name", "Task.body", "Task.projectID").From("Task").
						Where(goqu.L(`"Task"."projectID" in 
(SELECT "Project"."id" FROM "Project" WHERE "Project"."organizationID" in 
(SELECT "OrganizationMember"."organizationID" FROM "OrganizationMember" WHERE "OrganizationMember"."userID" = ?))`, user.ID)).
						ScanStructs(&tasks); err != nil {
						return nil, err
					}
					return tasks, nil
				},
			},
		},
	})
}
