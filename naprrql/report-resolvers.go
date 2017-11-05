// report-resolvers.go

//
// resolver logic for composite reporting objects
//
package naprrql

import (
	"errors"
	"log"

	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

//
// simple check for query params that should have been passed as part
// of the api call
//
func checkRequiredParams(params *graphql.ResolveParams) error {

	//check presence of acarids
	if len(params.Args) < 1 {
		log.Println("no query params")
		return errors.New("School selection variable: acaraIDs not supplied, query aborting")
	}
	val, present := params.Args["acaraIDs"]
	if !present || len(val.([]interface{})) < 1 {
		log.Println("not enough query params")
		return errors.New("School selection variable: acaraIDs not supplied, query aborting")
	}
	for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
		acaraid, _ := a_id.(string)
		if acaraid == "" {
			return errors.New("School acaraid cannot be blank/empty, query aborting")
		}

	}

	return nil

}

func buildReportResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	//
	// shorthand lookup objects for basic school info
	//
	resolvers["NaplanData/school_details"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("SchoolDetails:"))
	}

	//
	// resolver for score summary report object
	//
	resolvers["NaplanData/score_summary_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get the sif refid for each of the acarids supplied
		refids := make([]string, 0)
		for _, acaraid := range acaraids {
			refid := getIdentifiers(acaraid + ":")
			if len(refid) > 0 {
				refids = append(refids, refid...)
			}

		}

		// now construct the composite keys
		school_summary_keys := make([]string, 0)
		for _, refid := range refids {
			school_summary_keys = append(school_summary_keys, refid+":NAPTestScoreSummary:")
		}

		summ_refids := make([]string, 0)
		for _, summary_key := range school_summary_keys {
			ids := getIdentifiers(summary_key)
			for _, id := range ids {
				summ_refids = append(summ_refids, id)
			}
		}

		summaries, err := getObjects(summ_refids)
		summary_datasets := make([]ScoreSummaryDataSet, 0)
		for _, summary := range summaries {
			summ, _ := summary.(xml.NAPTestScoreSummary)
			testid := []string{summ.NAPTestRefId}
			obj, _ := getObjects(testid)
			test, _ := obj[0].(xml.NAPTest)
			sds := ScoreSummaryDataSet{Summ: summ, Test: test}
			summary_datasets = append(summary_datasets, sds)
		}

		return summary_datasets, err

	}

	resolvers["NaplanData/school_infos_by_acaraid"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get the sif refid for each of the acarids supplied
		refids := make([]string, 0)
		for _, acaraid := range acaraids {
			refid := getIdentifiers(acaraid + ":")
			if len(refid) > 0 {
				refids = append(refids, refid...)
			}

		}

		// get the school infos from the datastore
		siObjects, err := getObjects(refids)
		schoolInfos := make([]xml.SchoolInfo, 0)
		for _, sio := range siObjects {
			si, _ := sio.(xml.SchoolInfo)
			schoolInfos = append(schoolInfos, si)
		}

		return schoolInfos, err

	}

	resolvers["NaplanData/students_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

		return getObjects(studentids)

	}

	resolvers["NaplanData/domain_scores_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

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

		// construct RDS by including referenced test
		results := make([]ResponseDataSet, 0)
		for _, response := range responses {
			resp, _ := response.(xml.NAPResponseSet)
			// domain score entries will be null if response not completed e.g. abandoned
			if resp.DomainScore.RawScore == "" {
				continue
			}
			tests, err := getObjects([]string{resp.TestID})
			test, ok := tests[0].(xml.NAPTest)
			if err != nil || !ok {
				return []interface{}{}, err
			}
			students, err := getObjects([]string{resp.StudentID})
			student, ok := students[0].(xml.RegistrationRecord)
			if err != nil || !ok {
				return []interface{}{}, err
			}
			rds := ResponseDataSet{Test: test, Response: resp, Student: student}
			results = append(results, rds)
		}

		return results, nil
	}

	resolvers["NaplanData/domain_scores_event_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get school names
		schoolnames := make(map[string]string)
		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)

			schoolrefid := getIdentifiers(acaraid + ":")
			siObjects, err := getObjects(schoolrefid)
			if err != nil {
				return []interface{}{}, err
			}
			for _, sio := range siObjects {
				si, _ := sio.(xml.SchoolInfo)
				schoolnames[acaraid] = si.SchoolName
			}
		}

		studentObjs, err := getObjects(studentids)
		if err != nil {
			return []interface{}{}, err
		}

		// iterate students and assemble Event/Response Data Set
		results := make([]EventResponseDataSet, 0)
		for _, studentObj := range studentObjs {
			student, _ := studentObj.(xml.RegistrationRecord)
			studentEventIds := getIdentifiers(student.RefId + ":NAPEventStudentLink:")
			if len(studentEventIds) < 1 {
				// log.Println("no events found for student: ", student.RefId)
				continue
			}
			eventObjs, err := getObjects(studentEventIds)
			if err != nil {
				return []interface{}{}, err
			}
			for _, eventObj := range eventObjs {
				event := eventObj.(xml.NAPEvent)
				testObj, err := getObjects([]string{event.TestID})
				if err != nil {
					return []interface{}{}, err
				}
				test := testObj[0].(xml.NAPTest)

				responseIds := getIdentifiers(test.TestID + ":NAPStudentResponseSet:" + student.RefId)
				response := xml.NAPResponseSet{}
				if len(responseIds) > 0 {
					responseObjs, err := getObjects(responseIds)
					if err != nil {
						return []interface{}{}, err
					}
					response = responseObjs[0].(xml.NAPResponseSet)
				}

				erds := EventResponseDataSet{Student: student,
					Event:         event,
					Test:          test,
					Response:      response,
					SchoolDetails: SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
				}
				results = append(results, erds)
			}
		}

		return results, nil

	}

	resolvers["NaplanData/participation_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

		studentObjs, err := getObjects(studentids)
		if err != nil {
			return []interface{}{}, err
		}

		// iterate students and assemble ParticipationDataSets
		results := make([]ParticipationDataSet, 0)
		for _, studentObj := range studentObjs {
			student, _ := studentObj.(xml.RegistrationRecord)
			studentEventIds := getIdentifiers(student.RefId + ":NAPEventStudentLink:")
			if len(studentEventIds) < 1 {
				// log.Println("no events found for student: ", student.RefId)
				continue
			}
			eventObjs, err := getObjects(studentEventIds)
			if err != nil {
				return []interface{}{}, err
			}
			eventInfos := make([]EventInfo, 0)
			for _, eventObj := range eventObjs {
				event := eventObj.(xml.NAPEvent)
				testObj, err := getObjects([]string{event.TestID})
				if err != nil {
					return []interface{}{}, err
				}
				test := testObj[0].(xml.NAPTest)
				eventInfo := EventInfo{Test: test, Event: event}
				eventInfos = append(eventInfos, eventInfo)
			}
			schoolKey := eventInfos[0].Event.SchoolRefID
			schoolObj, err := getObjects([]string{schoolKey})
			if err != nil {
				return []interface{}{}, err
			}
			school, _ := schoolObj[0].(xml.SchoolInfo)
			// construct summary
			summaries := make([]ParticipationSummary, 0)
			for _, event := range eventInfos {
				summary := ParticipationSummary{
					Domain:            event.Test.TestContent.TestDomain,
					ParticipationCode: event.Event.ParticipationCode}
				summaries = append(summaries, summary)
			}
			pds := ParticipationDataSet{Student: student,
				School:     school,
				EventInfos: eventInfos,
				Summary:    summaries,
			}
			results = append(results, pds)
		}

		return results, nil

	}

	resolvers["NaplanData/codeframe_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// get the codeframe objects
		codeframes := make([]xml.NAPCodeFrame, 0)
		codeframeIds := getIdentifiers("NAPCodeFrame:")
		codeFrameObjs, err := getObjects(codeframeIds)
		if err != nil {
			return []interface{}{}, err
		}
		for _, codeframeObj := range codeFrameObjs {
			codeFrame, _ := codeframeObj.(xml.NAPCodeFrame)
			codeframes = append(codeframes, codeFrame)
		}

		cfds := make([]CodeFrameDataSet, 0)
		for _, codeframe := range codeframes {
			testObj, err := getObjects([]string{codeframe.NAPTestRefId})
			if err != nil {
				return []interface{}{}, err
			}
			test, _ := testObj[0].(xml.NAPTest)
			for _, cf_testlet := range codeframe.TestletList.Testlet {
				tlObj, _ := getObjects([]string{cf_testlet.NAPTestletRefId})
				if err != nil {
					return []interface{}{}, err
				}
				tl, _ := tlObj[0].(xml.NAPTestlet)
				for _, cf_item := range cf_testlet.TestItemList.TestItem {
					tiObj, err := getObjects([]string{cf_item.TestItemRefId})
					if err != nil {
						return []interface{}{}, err
					}
					ti, _ := tiObj[0].(xml.NAPTestItem)
					// log.Printf("\t\t%s", ti.TestItemContent.ItemName)
					cfd := CodeFrameDataSet{
						Test:    test,
						Testlet: tl,
						Item:    ti,
					}
					cfds = append(cfds, cfd)
				}
			}
		}

		return cfds, nil

	}

	resolvers["NaplanData/item_results_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			return nil, reqErr
		}
		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

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

			eventsRefId := getIdentifiers("event_by_student_test:" + resp.StudentID + ":" + resp.TestID + ":")
			events, err := getObjects(eventsRefId)
			event, ok := events[0].(xml.NAPEvent)
			if err != nil || !ok {
				return []interface{}{}, err
			}

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

					irds := ItemResponseDataSet{TestItem: item, Response: resp1,
						Student: student, Test: test, ParticipationCode: event.ParticipationCode}
					results = append(results, irds)
				}
			}
		}
		return results, nil
	}

	return resolvers
}
