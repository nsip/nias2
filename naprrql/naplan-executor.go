package naprrql

import (
	"log"

	"github.com/playlyfe/go-graphql"
)

func buildExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildSchema(), "NaplanData", "", buildResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	return executor
}

func buildResolvers() map[string]interface{} {

	allResolvers := make(map[string]interface{})

	naplanResolvers := buildNaplanResolvers()
	reportResolvers := buildReportResolvers()

	for k, v := range naplanResolvers {
		allResolvers[k] = v
	}

	for k, v := range reportResolvers {
		allResolvers[k] = v
	}

	return allResolvers

}
