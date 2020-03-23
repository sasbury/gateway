package gateway

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

// Merger is an interface for structs that are capable of taking a list of schemas and returning something that resembles
// a "merge" of those schemas.
type GatewayMerger interface {
	Merge(*ast.Schema, *ast.Schema) (*ast.Schema, error)
}

// MergerFunc is a wrapper of a function of the same signature as Merger.Merge
type GatewayMergerFunc func(*ast.Schema, *ast.Schema) (*ast.Schema, error)

// Merge invokes and returns the wrapped function
func (m GatewayMergerFunc) Merge(remoteSchema *ast.Schema, gatewaySchema *ast.Schema) (*ast.Schema, error) {
	return m(remoteSchema, gatewaySchema)
}

// mergeSchemas takes in a bunch of schemas and merges them into one. Following the strategies outlined here:
// https://github.com/nautilus/gateway/blob/master/docs/mergingStrategies.md
func mergeGateway(remoteSchema *ast.Schema, gatewaySchema *ast.Schema) (*ast.Schema, error) {

	// merging the schemas has to happen in 2 passes per definnition. The first groups definitions into different
	// categories. Interfaces need to be validated before we can add the types that implement them
	types := map[string][]*ast.Definition{}
	directives := map[string][]*ast.DirectiveDefinition{}
	interfaces := map[string][]*ast.Definition{}

	// add each type declared by the source schema to the one we are building up
	for name, definition := range gatewaySchema.Types {
		// if the definition is an interface
		if definition.Kind == ast.Interface {
			// ad it to the list
			interfaces[name] = append(interfaces[name], definition)
		} else {
			types[name] = append(types[name], definition)
		}
	}

	// add each directive to the list
	for name, definition := range gatewaySchema.Directives {
		directives[name] = append(directives[name], definition)
	}

	// merge each interface into one
	for name, definitions := range interfaces {
		for _, definition := range definitions {
			// look up if the type is already registered in the aggregate
			_, exists := remoteSchema.Types[name]

			// if we haven't seen it before
			if !exists {
				// use the declaration that we got from the new schema
				remoteSchema.Types[name] = definition

				remoteSchema.AddPossibleType(name, definition)

				// we're done with this definition
				continue
			}

			// if err := mergeInterfaces(remoteSchema, previousDefinition, definition); err != nil {
			// 	return nil, err
			// }
		}
	}

	possibleTypesSet := map[string]Set{}

	// merge each definition of each type into one
	for name, definitions := range types {
		if _, exists := possibleTypesSet[name]; !exists {
			possibleTypesSet[name] = Set{}
		}
		for _, definition := range definitions {
			// look up if the type is already registered in the aggregate
			previousDefinition, exists := remoteSchema.Types[name]

			// if we haven't seen it before
			if !exists {
				// use the declaration that we got from the new schema
				remoteSchema.Types[name] = definition

				if definition.Kind == ast.Union {
					for _, possibleType := range definition.Types {
						for _, typedef := range types[possibleType] {
							if !possibleTypesSet[name].Has(typedef.Name) {
								possibleTypesSet[name].Add(typedef.Name)
								remoteSchema.AddPossibleType(name, typedef)
							}
						}
					}
				} else {
					// register the type as an implementer of itself
					remoteSchema.AddPossibleType(name, definition)
				}

				// each interface that this type implements needs to be registered
				for _, iface := range definition.Interfaces {
					remoteSchema.AddPossibleType(iface, definition)
					remoteSchema.AddImplements(definition.Name, remoteSchema.Types[definition.Name])
				}

				// we're done with this type
				continue
			}

			// we only want one copy of the internal stuff
			if strings.HasPrefix(definition.Name, "__") {
				continue
			}

			// unify handling of errors for merging
			var err error

			switch definition.Kind {
			case ast.Object:
				err = mergeObjectTypesGateway(remoteSchema, previousDefinition, definition)
			}

			if err != nil {
				return nil, err
			}
		}
	}

	// merge each directive definition together
	for name, definitions := range directives {
		for _, definition := range definitions {
			// look up if the type is already registered in the aggregate
			_, exists := remoteSchema.Directives[name]
			// if we haven't seen it before
			if !exists {
				// use the declaration that we got from the new schema
				remoteSchema.Directives[name] = definition
			}
		}
	}

	// for now, just use the query type as the query type
	queryType, _ := remoteSchema.Types["Query"]
	mutationType, _ := remoteSchema.Types["Mutation"]
	subscriptionType, _ := remoteSchema.Types["Subscription"]

	remoteSchema.Query = queryType
	remoteSchema.Mutation = mutationType
	remoteSchema.Subscription = subscriptionType

	// we're done here
	return remoteSchema, nil
}

func mergeObjectTypesGateway(schema *ast.Schema, previousDefinition *ast.Definition, newDefinition *ast.Definition) error {
	// the fields in the aggregate
	previousFields := previousDefinition.Fields

	// we have to add the fields in the source definition with the one in the aggregate
	for _, newField := range newDefinition.Fields {
		// look up if we already know about this field
		field := previousFields.ForName(newField.Name)
		fmt.Println(newField.Name)

		// if the field is new, add it
		if field == nil {
			// and they aren't equal
			previousFields = append(previousFields, newField)
		}

	}
	previousDefinition.Fields = previousFields
	return nil
}
