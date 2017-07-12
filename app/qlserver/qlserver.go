// qlserver.go
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/labstack/echo"
	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

var sr = naprr.NewStreamReader()
var nd *naprr.NAPLANData
var executor *graphql.Executor

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// whitespace stripping routine for schema efficiency
//
func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}
	resolvers["NAPEvent/TestDisruptionList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		disruptionEvents := []interface{}{}

		// log.Printf("params: %#v\n\n", params)

		if napEvent, ok := params.Source.(xml.NAPEvent); ok {
			for _, event := range napEvent.TestDisruptionList.TestDisruption {
				disruptionEvents = append(disruptionEvents, event)
			}
		}
		return disruptionEvents, nil
	}

	resolvers["NAPEvent/Adjustment"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// keep the NAPEvent in play for downstream resolvers
		// otherwise embedded structs will be returned as anonymous by reflection
		// but impossible to resolve from type interface{} with no concrete
		// type behind them.
		return params.Source, nil
	}

	resolvers["Adjustment/PNPCodeList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		pnpCodes := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napEvent, ok := params.Source.(xml.NAPEvent); ok {
			return napEvent.Adjustment.PNPCodelist.PNPCode, nil
		}
		return pnpCodes, nil
	}

	resolvers["School/events"] = func(params *graphql.ResolveParams) (interface{}, error) {
		events := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if sd, ok := params.Source.(*naprr.SchoolData); ok {
			if params.Args["onlyDisruptions"].(bool) {
				for _, event := range sd.Events {
					if event.TestDisruptionList.TestDisruption != nil {
						events = append(events, event)
					}
				}
			} else {
				for _, event := range sd.Events {
					events = append(events, event)
				}
			}
		}
		return events, nil
	}

	resolvers["NAPResponseSet/DomainScore"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// keep the ResponseSet in play for downstream resolvers
		// otherwise embedded structs will be returned as anonymous by reflection
		// but impossible to resolve from type interface{} with no concrete
		// type behind them.
		return params.Source, nil
	}

	resolvers["NAPResponseSet/TestletList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		testletList := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			return napResponse.TestletList.Testlet, nil
		}
		return testletList, nil

	}

	resolvers["NAPResponseSet_ItemResponse/SubscoreList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		subscoreList := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_ItemResponse); ok {
			return napResponse.SubscoreList.Subscore, nil
		}
		return subscoreList, nil

	}

	resolvers["NAPResponseSet_ItemResponse/Item"] = func(params *graphql.ResolveParams) (interface{}, error) {

		linkedItem := ""
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_ItemResponse); ok {
			return nd.Items[napResponse.ItemRefID], nil
		}
		return linkedItem, nil

	}

	resolvers["NAPResponseSet_Testlet/ItemResponseList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		itemList := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_Testlet); ok {
			return napResponse.ItemResponseList.ItemResponse, nil
		}
		return itemList, nil

	}

	resolvers["DomainScore/PlausibleScaledValueList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		scaledValues := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			return napResponse.DomainScore.PlausibleScaledValueList.PlausibleScaledValue, nil
		}
		return scaledValues, nil

	}

	resolvers["DomainScore/RawScore"] = func(params *graphql.ResolveParams) (interface{}, error) {
		rawScore := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			rawScore = napResponse.DomainScore.RawScore
		}
		return rawScore, nil
	}

	resolvers["DomainScore/ScaledScoreValue"] = func(params *graphql.ResolveParams) (interface{}, error) {
		scaledScoreValue := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			scaledScoreValue = napResponse.DomainScore.ScaledScoreValue
		}
		return scaledScoreValue, nil
	}

	resolvers["DomainScore/ScaledScoreLogitValue"] = func(params *graphql.ResolveParams) (interface{}, error) {
		scaledScoreLogitValue := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			scaledScoreLogitValue = napResponse.DomainScore.ScaledScoreValue
		}
		return scaledScoreLogitValue, nil
	}

	resolvers["DomainScore/ScaledScoreStandardError"] = func(params *graphql.ResolveParams) (interface{}, error) {
		scaledScoreStandardError := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			scaledScoreStandardError = napResponse.DomainScore.ScaledScoreStandardError
		}
		return scaledScoreStandardError, nil
	}

	resolvers["DomainScore/ScaledScoreLogitStandardError"] = func(params *graphql.ResolveParams) (interface{}, error) {
		scaledScoreLogitStandardError := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			scaledScoreLogitStandardError = napResponse.DomainScore.ScaledScoreLogitStandardError
		}
		return scaledScoreLogitStandardError, nil
	}

	resolvers["DomainScore/StudentDomainBand"] = func(params *graphql.ResolveParams) (interface{}, error) {
		studentDomainBand := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			studentDomainBand = napResponse.DomainScore.StudentDomainBand
		}
		return studentDomainBand, nil
	}

	resolvers["DomainScore/StudentProficiency"] = func(params *graphql.ResolveParams) (interface{}, error) {
		studentProficiency := ""
		if napResponse, ok := params.Source.(xml.NAPResponseSet); ok {
			studentProficiency = napResponse.DomainScore.StudentDomainBand
		}
		return studentProficiency, nil
	}

	resolvers["School/responses"] = func(params *graphql.ResolveParams) (interface{}, error) {
		responses := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if sd, ok := params.Source.(*naprr.SchoolData); ok {
			for _, response := range sd.Responses {
				responses = append(responses, response)
			}
		}
		return responses, nil
	}

	resolvers["School/students"] = func(params *graphql.ResolveParams) (interface{}, error) {
		students := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if sd, ok := params.Source.(*naprr.SchoolData); ok {
			for _, student := range sd.Students {
				students = append(students, student)
			}
		}
		return students, nil
	}

	resolvers["School/score_summaries"] = func(params *graphql.ResolveParams) (interface{}, error) {
		summaries := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if sd, ok := params.Source.(*naprr.SchoolData); ok {
			for _, summary := range sd.ScoreSummaries {
				summaries = append(summaries, summary)
			}
		}
		return summaries, nil
	}

	resolvers["RegistrationRecord/OtherIdList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		otherIDs := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napRegistrationRecord, ok := params.Source.(xml.RegistrationRecord); ok {
			return napRegistrationRecord.OtherIdList.OtherId, nil
		}
		return otherIDs, nil
	}

	resolvers["NaplanQuery/getSchoolData"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// get the school data
		schools := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		schoolIDs := params.Args["acaraIDs"].([]interface{})
		for _, id := range schoolIDs {
			if schoolID, ok := id.(string); ok {
				sd := sr.GetSchoolData(schoolID)
				schools = append(schools, sd)
			}
		}

		return schools, nil
	}

	return resolvers
}

func buildSchema() string {
	dat, err := ioutil.ReadFile("naplan_schema.graphql")
	if err != nil {
		log.Fatalln("Unable to load schema from file: ", err)
	}
	return string(dat)
}

func buildExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(buildSchema(), "NaplanQuery", "", buildResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	executor.ResolveType = func(value interface{}) string {
		log.Printf("resolve: %#v\n\n", value)
		switch value.(type) {
		case *xml.NAPEvent:
			return "NAPEvent"
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
		case xml.XMLAttributeStruct:
			return "XMLAttributeStruct"
		case xml.NAPTestScoreSummary:
			return "NAPTestScoreSummary"
		}
		return ""
	}

	return executor
}

//
// simple web server to support gql queries & web ui (graphiql)
//
func graphQLHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}
	// log.Printf("variables: %v\n\n", variables)
	result, err := executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

func main() {

	// get the NAPLAN Test Data once globally
	nd = sr.GetNAPLANData(naprr.META_STREAM)

	executor = buildExecutor()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/", "public")
	e.File("/sifql", "public/index.html")

	e.POST("/graphql", graphQLHandler)

	e.Logger.Fatal(e.Start(":1329"))
}
