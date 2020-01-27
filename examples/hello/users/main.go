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

	type User implements Node {
		id: ID!
		name: String!
	}

	type Query {
		node(id: ID!): Node
		allUsers: [User!]!
	}
`

// the users by id
var users = map[string]*User{
	"u1": {
		id:        "u1",
		name: "Alec",
	},
	"u2": {
		id:        "u2",
		name: "Stephen",
	},
}

// type resolvers

type User struct {
	id        graphql.ID
	name string
}

func (u *User) ID() graphql.ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
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

func (n *NodeResolver) ToUser() (*User, bool) {
	user, ok := n.node.(*User)
	return user, ok
}

// query resolvers

type queryA struct{}

func (q *queryA) Node(args struct{ ID string }) *NodeResolver {
	user := users[args.ID]

	if user != nil {
		log.Printf("Returning user %q\n", user)
		return &NodeResolver{user}
	} else {
		log.Printf("No user matched %q\n", args.ID)
		return nil
	}
}

func (q *queryA) AllUsers() []*User {
	// build up a list of all the users
	userSlice := []*User{}

	for _, user := range users {
		userSlice = append(userSlice, user)
	}

	log.Printf("Returning %d users %q\n", len(userSlice), userSlice)

	return userSlice
}

func main() {
	// attach the schema to the resolver object
	schema := graphql.MustParseSchema(Schema, &queryA{})

	// make sure we add the user info to the execution context
	relayH := &relay.Handler{Schema: schema}
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request){
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
		relayH.ServeHTTP(rw, req)
	})

	log.Println("Users fragment running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
