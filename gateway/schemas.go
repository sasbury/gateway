package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	utils "github.com/honeyscience/honey-utils-go"
	"github.com/nautilus/graphql"
)

func ParseRemoteSchemas() ([]*graphql.RemoteSchema, error) {
	// build up the list of remote schemas
	remoteSchemas := []*graphql.RemoteSchema{}

	files, err := ioutil.ReadDir("../schema/remote")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("../schema/remote/%s", file.Name())
		rawSchema, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		rawSchemaString := string(rawSchema)

		parsedSchema, err := graphql.LoadSchema(rawSchemaString)
		if err != nil {
			return nil, err
		}

		urlEnvString := fmt.Sprintf("INTERNAL_API_URL_%s", strings.TrimSuffix(file.Name(), ".graphql"))
		url, err := utils.EnvString(strings.ToUpper(urlEnvString))
		if err != nil {
			return nil, err
		}

		remoteSchema := &graphql.RemoteSchema{
			Schema: parsedSchema,
			URL:    url,
		}

		remoteSchemas = append(remoteSchemas, remoteSchema)
	}

	return remoteSchemas, nil
}
