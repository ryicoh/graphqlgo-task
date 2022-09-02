package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
	"github.com/ryicoh/graphqlgotask"
)

func main() {
	db := newDB()

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(*graphqlgotask.NewRootQuery(db))}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	tokenStore := new(tokenStore)
	tokenStore.kvs = map[string]string{"aaaa": "a49846af-6fad-42ab-ac0b-fb5a88e7a369"}
	http.Handle("/graphql", handler.New(&handler.Config{
		Schema:       &schema,
		Pretty:       true,
		GraphiQL:     true,
		RootObjectFn: newRootObjectFn(tokenStore),
	}))
	fmt.Println("server started")
	http.ListenAndServe("127.0.0.1:8084", nil)
}

func newDB() *goqu.Database {
	dialect := goqu.Dialect("postgres")

	pgDb, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable port=25432 password=postgrespassword")

	if err != nil {
		panic(err.Error())
	}
	db := dialect.DB(pgDb)

	return db
}

type tokenStore struct {
	kvs map[string]string
}

func (t *tokenStore) GetUserID(token string) (userID string, err error) {
	userID, ok := t.kvs[token]
	if !ok {
		return "", errors.New("user not found")
	}
	return userID, nil
}

type AuthStore interface {
	GetUserID(token string) (userID string, err error)
}

func newRootObjectFn(store AuthStore) func(ctx context.Context, r *http.Request) map[string]interface{} {
	return func(ctx context.Context, r *http.Request) map[string]interface{} {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			userID, err := store.GetUserID(token)
			if err != nil {
				return nil
			}
			return map[string]interface{}{"userID": userID}
		}
		return nil
	}
}
