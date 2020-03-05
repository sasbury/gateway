package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

var Schema = `
	schema {
		query: Query
	}

	interface Node {
		id: ID!
	}

	type UserUser implements Node {
		id: ID!
		recentPurchases: [ProductProduct]
	}

	type ProductProduct implements Node {
		id: ID!
		name: String!
		price: String!
	}

	type Query {
		node(id: ID!): Node
		allProducts: [ProductProduct!]!
	}
`

// the products by id
var products = map[string]*ProductProduct{
	"p1": {
		id:    "p1",
		name:  "Tide",
		price: "12.2",
	},
	"p2": {
		id:    "p2",
		name:  "Cheer",
		price: "11.1",
	},
}

type UserUser struct {
	id              graphql.ID
	recentPurchases []*ProductProduct
}

func (u *UserUser) ID() graphql.ID {
	return u.id
}

func (u *UserUser) RecentPurchases() *[]*ProductProduct {
	return &u.recentPurchases
}

// type resolvers

type ProductProduct struct {
	id    graphql.ID
	name  string
	price string
}

func (u *ProductProduct) ID() graphql.ID {
	return u.id
}

func (u *ProductProduct) Name() string {
	return u.name
}

func (u *ProductProduct) Price() string {
	return u.price
}

type Node interface {
	ID() graphql.ID
}

type NodeResolver struct {
	node Node
}

func (n *NodeResolver) ID() graphql.ID {
	return n.node.ID()
}

func (n *NodeResolver) ToProductProduct() (*ProductProduct, bool) {
	Product, ok := n.node.(*ProductProduct)
	return Product, ok
}

func (n *NodeResolver) ToUserUser() (*UserUser, bool) {
	user, ok := n.node.(*UserUser)
	return user, ok
}

// query resolvers

type queryB struct{}

func (q *queryB) Node(args struct{ ID string }) *NodeResolver {
	switch args.ID[0] {
	case 'u':
		u := q.getUser(args.ID)
		if u != nil {
			log.Printf("returning user node for %s with id %q", args.ID, u.ID())
			return &NodeResolver{u}
		}
		return nil
	case 'p':
		p := q.getProduct(args.ID)
		if p != nil {
			return &NodeResolver{p}
		}
		return nil
	}
	return nil
}

func (q *queryB) getProduct(id string) *ProductProduct {
	product := products[id]

	if product != nil {
		log.Printf("Returning product %q\n", product)
		return product
	} else {
		log.Printf("No product matched %q\n", id)
		return nil
	}
}

func (q *queryB) getUser(id string) *UserUser {
	productSlice := []*ProductProduct{}

	if id == "u1" {
		for _, product := range products {
			productSlice = append(productSlice, product)
		}
	}

	log.Printf("Returning %d products for user %s\n", len(productSlice), id)

	return &UserUser{
		id:              graphql.ID(id),
		recentPurchases: productSlice,
	}
}

func (q *queryB) AllProducts() []*ProductProduct {
	productSlice := []*ProductProduct{}

	for _, product := range products {
		productSlice = append(productSlice, product)
	}

	log.Printf("Returning %d products\n", len(productSlice))

	return productSlice
}

func main() {
	// attach the schema to the resolver object
	schema := graphql.MustParseSchema(Schema, &queryB{})

	// make sure we add the Product info to the execution context
	relayH := &relay.Handler{Schema: schema}
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
		relayH.ServeHTTP(rw, req)
	})

	log.Println("Products fragment running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
