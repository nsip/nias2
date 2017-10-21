package naprrql

import (
	"log"

	"github.com/playlyfe/go-graphql"
)

func buildNAPExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildNAPSchema(), "NaplanData", "", buildNAPResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	return executor
}

func buildNAPResolvers() map[string]interface{} {

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

func buildISRPrintExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildISRPrintSchema(), "ISRPrint", "", buildISRPrintResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	return executor

}

func buildItemPrintExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildItemPrintSchema(), "ItemPrint", "", buildItemPrintResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	return executor

}
