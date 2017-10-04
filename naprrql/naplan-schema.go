package naprrql

import (
	"io/ioutil"
	"log"
)

func buildNAPSchema() string {
	dat, err := ioutil.ReadFile("./gql_schemas/naplan_schema.graphql")
	if err != nil {
		log.Fatalln("Unable to load schema from file: ", err)
	}
	return string(dat)
}

func buildISRPrintSchema() string {
	dat, err := ioutil.ReadFile("./gql_schemas/isrprint_schema.graphql")
	if err != nil {
		log.Fatalln("Unable to load schema from file: ", err)
	}
	return string(dat)
}
