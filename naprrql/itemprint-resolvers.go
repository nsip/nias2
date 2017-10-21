// isrprint-resolvers.go

// resolver functions to support creation of isr printing report
package naprrql

import (
	"errors"
	"log"

	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

// create the resolver functions for item result printing
//
func buildItemPrintResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	//
	// resolver for item printing results, line per student.
	//

	resolvers["ItemPrint/item_results_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {
		reqErr := checkItemReportParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara id from the request params
		acaraid := params.Args["acaraID"].(string)
		log.Println("Generating item results for ACARA ID: ", acaraid)

		// get students for the schools
		studentids := make([]string, 0)
		key := "student_by_acaraid:" + acaraid
		studentRefIds := getIdentifiers(key)
		studentids = append(studentids, studentRefIds...)

		// get responses for student
		responseids := make([]string, 0)
		for _, studentid := range studentids {
			key := "responseset_by_student:" + studentid
			responseRefId := getIdentifiers(key)
			responseids = append(responseids, responseRefId...)
		}

		// get responses
		responses, err := getObjects(responseids)
		if err != nil {
			return []interface{}{}, err
		}

		// convenience map to avoid revisiting db for tests
		testLookup := make(map[string]xml.NAPTest) // key string = test refid
		// get tests for yearLevel
		for _, yrLvl := range []string{"3", "5", "7", "9"} {
			tests, err := getTestsForYearLevel(yrLvl)
			if err != nil {
				return nil, err
			}
			// convenience map to avoid revisiting db for tests
			testLookup := make(map[string]xml.NAPTest) // key string = test refid
			for _, test := range tests {
				t := test
				testLookup[t.TestID] = t
			}
		}

		// construct RDS by including referenced test
		results := make([]ItemResponseDataSet, 0)
		for _, response := range responses {
			resp, _ := response.(xml.NAPResponseSet)

			students, err := getObjects([]string{resp.StudentID})
			student, ok := students[0].(xml.RegistrationRecord)
			if err != nil || !ok {
				return []interface{}{}, err
			}
			student.Flatten()
			test := testLookup[resp.TestID]

			for _, testlet := range resp.TestletList.Testlet {
				for _, item_response := range testlet.ItemResponseList.ItemResponse {

					resp1 := resp // pruned copy of response
					resp1.TestletList.Testlet = make([]xml.NAPResponseSet_Testlet, 1)
					resp1.TestletList.Testlet[0] = testlet
					resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse = make([]xml.NAPResponseSet_ItemResponse, 1)
					resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse[0] = item_response

					items, err := getObjects([]string{item_response.ItemRefID})
					item, ok := items[0].(xml.NAPTestItem)
					if err != nil || !ok {
						return []interface{}{}, err
					}

					irds := ItemResponseDataSet{TestItem: item, Response: resp1, Student: student, Test: test}
					results = append(results, irds)
				}
			}
		}
		return results, nil
	}

	return resolvers
}

//
// simple check for query params that should have been passed as part
// of the api call
//
func checkItemReportParams(params *graphql.ResolveParams) error {

	//check presence of schoolAcarid
	if len(params.Args) < 1 {
		log.Println("no query params")
		return errors.New("Required variables: acaraID not supplied, query aborting")
	}

	_, present := params.Args["acaraID"]
	if !present {
		log.Println("no acaraid parameter provided")
		return errors.New("Required variable: 'acaraID' not provided.")
	}

	acaraid, ok := params.Args["acaraID"].(string)
	if !ok || acaraid == "" {
		log.Println("Cannot interpret variable acaraID as meaningful string")
		return errors.New("Required variable: 'acaraID' not provided.")
	}

	return nil

}
