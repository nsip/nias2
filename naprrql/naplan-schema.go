package naprrql

import (
	"io/ioutil"
	"log"
)

func buildSchema() string {
	dat, err := ioutil.ReadFile("naplan_schema.graphql")
	if err != nil {
		log.Fatalln("Unable to load schema from file: ", err)
	}
	return string(dat)
}
