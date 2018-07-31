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
		"avgRating": &graphql.Field{
			Type: graphql.Float,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				toilet, isOk := params.Source.(Toilet)
				if isOk {
					var data []*ToiletRating
					err := db.Model(&toilet).Related(&data).Error
					if err != nil {
						return nil, err
					}
					var total, count = 0, 0
					for _, rating := range data {
						total += rating.Rating
						count += 1
					}
					return total / count, nil
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

var mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createToilet": &graphql.Field{
			Type: toilet,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"location": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				location, _ := p.Args["location"].(string)

				newToilet := &Toilet{
					Name: name,
					Location: location,
				}
				err := db.Save(&newToilet).Error
				if err != nil {
					return nil, err
				}

				return newToilet, nil
			},
		},
		"rateToilet": &graphql.Field{
			Type: toiletRating,
			Args: graphql.FieldConfigArgument{
				"uid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"rating": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"comment": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var toilet Toilet
				err := db.First(&toilet, p.Args["uid"]).Error
				if err != nil {
					return nil, err
				}

				rating, _ := p.Args["rating"].(int)
				comment, _ := p.Args["comment"].(string)

				newRating := &ToiletRating{
					Toilet: toilet,
					Rating: rating,
					Comment: comment,
				}
				err = db.Save(&newRating).Error
				if err != nil {
					return nil, err
				}

				return newRating, nil
			},
		},
	},
})

func runServer() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
		Mutation: mutation,
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
