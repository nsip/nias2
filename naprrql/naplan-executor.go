package naprrql

import (
	"log"

	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func buildExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildSchema(), "NaplanData", "", buildResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	executor.ResolveType = func(value interface{}) string {
		log.Printf("resolve: %#v\n\n", value)
		switch value.(type) {
		case xml.NAPEvent:
			return "NAPEvent"
		case xml.TestDisruptionList:
			return "TestDisruptionList"
		case xml.TestDisruption:
			return "TestDisruption"
		case xml.Adjustment:
			return "Adjustment"
		case xml.PNPCodelist:
			return "PNPCodeList"
		case xml.NAPResponseSet:
			return "NAPResponseSet"
		case xml.NAPResponseSet_Testlet:
			return "NAPResponseSet_Testlet"
		case xml.NAPResponseSet_ItemResponse:
			return "NAPResponseSet_ItemResponse"
		case xml.NAPTestItem:
			return "NAPTestItem"
		case xml.TestItemContent:
			return "TestItemContent"
		case xml.NAPResponseSet_Subscore:
			return "NAPResponseSet_Subscore"
		case xml.RegistrationRecord:
			return "RegistrationRecord"
		case xml.OtherIdList:
			return "OtherIdList"
		case xml.XMLAttributeStruct:
			return "XMLAttributeStruct"
		case xml.NAPTestScoreSummary:
			return "NAPTestScoreSummary"
		}
		return ""
	}

	return executor
}
