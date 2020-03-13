package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	utils "github.com/honeyscience/honey-utils-go"
	"github.com/nautilus/graphql"
)

var GATEWAY_SCHEMA_FILE_NAME = "gateway.graphql"

func ParseRemoteSchemas() ([]*graphql.RemoteSchema, *graphql.RemoteSchema, error) {
	// build up the list of remote schemas
	remoteSchemas := []*graphql.RemoteSchema{}
	var legacySchema *graphql.RemoteSchema

	files, err := ioutil.ReadDir("../schema/remote")
	if err != nil {
		return nil, legacySchema, err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("../schema/remote/%s", file.Name())
		rawSchema, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, legacySchema, err
		}
		rawSchemaString := string(rawSchema)

		parsedSchema, err := graphql.LoadSchema(rawSchemaString)
		if err != nil {
			return nil, legacySchema, err
		}

		urlEnvString := fmt.Sprintf("INTERNAL_API_URL_%s", strings.TrimSuffix(file.Name(), ".graphql"))
		url, err := utils.EnvString(strings.ToUpper(urlEnvString))
		if err != nil {
			return nil, legacySchema, err
		}

		schema := &graphql.RemoteSchema{
			Schema: parsedSchema,
			URL:    url,
		}

		if file.Name() == GATEWAY_SCHEMA_FILE_NAME {
			legacySchema = schema
		} else {
			remoteSchemas = append(remoteSchemas, schema)
		}
	}

	return remoteSchemas, legacySchema, nil
}
