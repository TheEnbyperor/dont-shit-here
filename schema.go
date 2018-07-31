package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
	"log"
)

var toiletRating = graphql.NewObject(graphql.ObjectConfig{
	Name: "ToiletRating",
	Fields: graphql.Fields{
		"uid": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				toilet, isOk := params.Source.(ToiletRating)
				if isOk {
					return toilet.ID, nil
				}
				return nil, nil
			},
		},
		"rating": &graphql.Field{
			Type: graphql.Int,
		},
		"comment": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var toilet = graphql.NewObject(graphql.ObjectConfig{
	Name: "Toilet",
	Fields: graphql.Fields{
		"uid": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				toilet, isOk := params.Source.(Toilet)
				if isOk {
					return toilet.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"ratings": &graphql.Field{
			Type: graphql.NewList(toiletRating),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				toilet, isOk := params.Source.(Toilet)
				if isOk {
					var data []*ToiletRating
					err := db.Model(&toilet).Related(&data).Error
					if err != nil {
						return nil, err
					}
					return data, nil
				}
				return nil, nil
			},
		},
	},
})

var query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"toilets": &graphql.Field{
			Type: graphql.NewList(toilet),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var data []*Toilet
				err := db.Select(&data).Error
				if err != nil {
					return nil, err
				}
				return data, nil
			},
		},
		"toilet": &graphql.Field{
			Type: toilet,
			Args: graphql.FieldConfigArgument{
				"uid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var toilet Toilet
				err := db.First(&toilet, p.Args["uid"]).Error
				if err != nil {
					return nil, err
				}
				return toilet, nil
			},
		},
	},
})

func runServer() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})

	if err != nil {
		log.Fatalln(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	log.Println("Starting serve on :8080")
	http.ListenAndServe(":8080", nil)
}
