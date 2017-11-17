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

func checkGuid(issues []GuidCheckDataSet, objectguid string, objecttype string, guid string, expected_type string, idmap map[string]string) []GuidCheckDataSet {
	found, ok := idmap[guid]
	if ok {
		if found != expected_type {
			issues = append(issues, GuidCheckDataSet{ObjectName: objectguid,
				ObjectType:    objecttype,
				Guid:          guid,
				ShouldPointTo: expected_type,
				PointsTo:      found})
		}
	} else {
		issues = append(issues, GuidCheckDataSet{ObjectName: objectguid,
			ObjectType:    objecttype,
			Guid:          guid,
			ShouldPointTo: expected_type,
			PointsTo:      "nil"})
	}
	return issues
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
					for _, responseObj := range responseObjs {
						response = responseObj.(xml.NAPResponseSet)
						erds := EventResponseDataSet{Student: student,
							Event:         event,
							Test:          test,
							Response:      response,
							SchoolDetails: SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
						}
						results = append(results, erds)
					}
				} else {

					erds := EventResponseDataSet{Student: student,
						Event:         event,
						Test:          test,
						Response:      response,
						SchoolDetails: SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
					}
					results = append(results, erds)
				}
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
						Test:           test,
						Testlet:        tl,
						Item:           ti,
						SequenceNumber: "unknown",
					}
					cfd.SequenceNumber = cfd.GetItemLocationInTestlet(ti.ItemID)

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

		// get the codeframe

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

					testlets, err := getObjects([]string{testlet.NapTestletRefId})
					tl := testlets[0].(xml.NAPTestlet)
					if err != nil || !ok {
						return []interface{}{}, err
					}

					irds := ItemResponseDataSet{TestItem: item, Response: resp1,
						Student: student, Test: test, Testlet: tl,
						ParticipationCode: event.ParticipationCode}
					results = append(results, irds)
				}
			}
		}
		return results, nil
	}

	resolvers["NaplanData/orphan_school_summary_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
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

		// all school summary IDs ingested
		all_school_summary_keys := getIdentifiers("NAPTestScoreSummary:")

		// all school summary IDs linked to one of the acaraIDs
		linked_school_summary_keys := make([]string, 0)
		for _, refid := range refids {
			linked_school_summary_keys = append(linked_school_summary_keys, refid+":NAPTestScoreSummary:")
		}

		summ_refids := make([]string, 0)
		for _, summary_key := range linked_school_summary_keys {
			ids := getIdentifiers(summary_key)
			for _, id := range ids {
				summ_refids = append(summ_refids, id)
			}
		}
		seen := map[string]bool{}
		for _, x := range summ_refids {
			seen[x] = true
		}
		orphans := make([]string, 0)
		for _, x := range all_school_summary_keys {
			if !seen[x] {
				orphans = append(orphans, x)
			}
		}
		log.Printf("Found: %d orphan score summaries\n", len(orphans))

		summaries, err := getObjects(orphans)
		summary_datasets := make([]xml.NAPTestScoreSummary, 0)
		for _, summary := range summaries {
			summ, _ := summary.(xml.NAPTestScoreSummary)
			summary_datasets = append(summary_datasets, summ)
		}

		return summary_datasets, err

	}

	resolvers["NaplanData/orphan_event_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
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

		// all event IDs ingested
		all_keys := getIdentifiers("NAPEventStudentLink:")

		// all event IDs linked to one of the acaraIDs
		linked_keys := make([]string, 0)
		for _, refid := range refids {
			linked_keys = append(linked_keys, refid+":NAPEventStudentLink:")
		}

		summ_refids := make([]string, 0)
		for _, summary_key := range linked_keys {
			ids := getIdentifiers(summary_key)
			for _, id := range ids {
				summ_refids = append(summ_refids, id)
			}
		}
		seen := map[string]bool{}
		for _, x := range summ_refids {
			seen[x] = true
		}
		orphans := make([]string, 0)
		for _, x := range all_keys {
			if !seen[x] {
				orphans = append(orphans, x)
			}
		}
		log.Printf("Found: %d orphan events\n", len(orphans))

		events, err := getObjects(orphans)
		summary_datasets := make([]xml.NAPEvent, 0)
		for _, summary := range events {
			summ, _ := summary.(xml.NAPEvent)
			summary_datasets = append(summary_datasets, summ)
		}

		return summary_datasets, err

	}

	resolvers["NaplanData/orphan_student_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
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

		// all student IDs ingested
		all_keys := getIdentifiers("StudentPersonal:")

		linked_keys := make([]string, 0)
		for _, acaraid := range acaraids {
			linked_keys = append(linked_keys, "student_by_acaraid:"+acaraid+":")
		}

		summ_refids := make([]string, 0)
		for _, summary_key := range linked_keys {
			ids := getIdentifiers(summary_key)
			for _, id := range ids {
				summ_refids = append(summ_refids, id)
			}
		}
		seen := map[string]bool{}
		for _, x := range summ_refids {
			seen[x] = true
		}
		orphans := make([]string, 0)
		for _, x := range all_keys {
			if !seen[x] {
				orphans = append(orphans, x)
			}
		}
		log.Printf("Found: %d orphan students\n", len(orphans))

		students, err := getObjects(orphans)
		summary_datasets := make([]xml.RegistrationRecord, 0)
		for _, summary := range students {
			summ, _ := summary.(xml.RegistrationRecord)
			summary_datasets = append(summary_datasets, summ)
		}

		return summary_datasets, err

	}

	resolvers["NaplanData/guid_check_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// all student IDs ingested
		test_ids := getIdentifiers("NAPTest:")
		testlet_ids := getIdentifiers("NAPTestlet:")
		testitem_ids := getIdentifiers("NAPTestItem:")
		scoresummary_ids := getIdentifiers("NAPTestScoreSummary:")
		event_ids := getIdentifiers("NAPEventStudentLink:")
		response_ids := getIdentifiers("NAPStudentResponseSet:")
		codeframe_ids := getIdentifiers("NAPCodeFrame:")
		school_ids := getIdentifiers("SchoolInfo:")
		student_ids := getIdentifiers("StudentPersonal:")

		ids := make(map[string]string)
		for _, x := range test_ids {
			ids[x] = "test"
		}
		for _, x := range testlet_ids {
			ids[x] = "testlet"
		}
		for _, x := range testitem_ids {
			ids[x] = "testitem"
		}
		for _, x := range scoresummary_ids {
			ids[x] = "scoresummary"
		}
		for _, x := range event_ids {
			ids[x] = "event"
		}
		for _, x := range response_ids {
			ids[x] = "response"
		}
		for _, x := range codeframe_ids {
			ids[x] = "codeframe"
		}
		for _, x := range school_ids {
			ids[x] = "school"
		}
		for _, x := range student_ids {
			ids[x] = "student"
		}
		results := make([]GuidCheckDataSet, 0)

		codeframes, err := getObjects(codeframe_ids)
		for _, codeframe := range codeframes {
			t, _ := codeframe.(xml.NAPCodeFrame)
			results = checkGuid(results, t.RefId, "codeframe", t.NAPTestRefId, "test", ids)
			for _, tl := range t.TestletList.Testlet {
				results = checkGuid(results, t.RefId, "codeframe", tl.NAPTestletRefId, "testlet", ids)
				for _, ti := range tl.TestItemList.TestItem {
					results = checkGuid(results, t.RefId, "codeframe", ti.TestItemRefId, "testitem", ids)
				}
			}
		}
		codeframes = nil

		events, err := getObjects(event_ids)
		for _, event := range events {
			t, _ := event.(xml.NAPEvent)
			results = checkGuid(results, t.EventID, "event", t.SPRefID, "student", ids)
			results = checkGuid(results, t.EventID, "event", t.SchoolRefID, "school", ids)
			results = checkGuid(results, t.EventID, "event", t.TestID, "test", ids)
		}
		events = nil

		responses, err := getObjects(response_ids)
		for _, response := range responses {
			t, _ := response.(xml.NAPResponseSet)
			results = checkGuid(results, t.ResponseID, "response", t.StudentID, "student", ids)
			results = checkGuid(results, t.ResponseID, "response", t.TestID, "test", ids)
			for _, tl := range t.TestletList.Testlet {
				results = checkGuid(results, t.ResponseID, "response", tl.NapTestletRefId, "testlet", ids)
				for _, ti := range tl.ItemResponseList.ItemResponse {
					results = checkGuid(results, t.ResponseID, "response", ti.ItemRefID, "testitem", ids)
				}
			}
		}
		responses = nil

		testitems, err := getObjects(testitem_ids)
		for _, testitem := range testitems {
			t, _ := testitem.(xml.NAPTestItem)
			for _, s := range t.TestItemContent.ItemSubstitutedForList.SubstituteItem {
				results = checkGuid(results, t.ItemID, "testitem", s.SubstituteItemRefId, "testitem", ids)
			}
		}
		testitems = nil

		testlets, err := getObjects(testlet_ids)
		for _, testlet := range testlets {
			t, _ := testlet.(xml.NAPTestlet)
			results = checkGuid(results, t.TestletID, "testlet", t.NAPTestRefId, "test", ids)
			for _, ti := range t.TestItemList.TestItem {
				results = checkGuid(results, t.TestletID, "testlet", ti.TestItemRefId, "testitem", ids)
			}
		}
		testlets = nil

		summarys, err := getObjects(scoresummary_ids)
		for _, summary := range summarys {
			t, _ := summary.(xml.NAPTestScoreSummary)
			results = checkGuid(results, t.SummaryID, "scoresummary", t.SchoolInfoRefId, "school", ids)
			results = checkGuid(results, t.SummaryID, "test", t.NAPTestRefId, "test", ids)
		}
		summarys = nil

		log.Printf("%d GUID mismatches\n", len(results))
		return results, err

	}

	resolvers["NaplanData/codeframe_check_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// all student IDs ingested
		test_ids := getIdentifiers("NAPTest:")
		testlet_ids := getIdentifiers("NAPTestlet:")
		testitem_ids := getIdentifiers("NAPTestItem:")
		codeframe_ids := getIdentifiers("NAPCodeFrame:")

		results := make([]CodeframeCheckDataSet, 0)
		seen := make(map[string]bool)

		codeframes, err := getObjects(codeframe_ids)
		for _, codeframe := range codeframes {
			t, _ := codeframe.(xml.NAPCodeFrame)
			seen[t.NAPTestRefId] = true
			for _, tl := range t.TestletList.Testlet {
				seen[tl.NAPTestletRefId] = true
				for _, ti := range tl.TestItemList.TestItem {
					seen[ti.TestItemRefId] = true
				}
			}
		}
		codeframes = nil
		for _, t := range test_ids {
			if _, ok := seen[t]; !ok {
				tests, err := getObjects([]string{t})
				if err != nil {
					results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: "nil", ObjectType: "test", Test: xml.NAPTest{}, Testlet: xml.NAPTestlet{}, TestItem: xml.NAPTestItem{}})
				} else {
					t1 := tests[0].(xml.NAPTest)
					results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: t1.TestContent.LocalId, ObjectType: "test", Test: t1, Testlet: xml.NAPTestlet{}, TestItem: xml.NAPTestItem{}})
				}
			}
		}
		for _, t := range testlet_ids {
			if _, ok := seen[t]; !ok {
				testlets, err := getObjects([]string{t})
				if err != nil {
					results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: "nil", ObjectType: "testlet", Test: xml.NAPTest{}, Testlet: xml.NAPTestlet{}, TestItem: xml.NAPTestItem{}})
				} else {
					t1 := testlets[0].(xml.NAPTestlet)
					tests, err := getObjects([]string{t1.NAPTestRefId})
					if err != nil {
						test := tests[0].(xml.NAPTest)
						results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: t1.TestletContent.LocalId, ObjectType: "testlet", Test: test, Testlet: t1, TestItem: xml.NAPTestItem{}})
					} else {
						results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: t1.TestletContent.LocalId, ObjectType: "testlet", Test: xml.NAPTest{}, Testlet: t1, TestItem: xml.NAPTestItem{}})
					}
				}
			}
		}
		for _, t := range testitem_ids {
			if _, ok := seen[t]; !ok {
				testitems, err := getObjects([]string{t})
				if err != nil {
					results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: "nil", ObjectType: "testitem", Test: xml.NAPTest{}, Testlet: xml.NAPTestlet{}, TestItem: xml.NAPTestItem{}})
				} else {
					t1 := testitems[0].(xml.NAPTestItem)
					results = append(results, CodeframeCheckDataSet{ObjectID: t, LocalID: t1.TestItemContent.NAPTestItemLocalId, ObjectType: "testitem", Test: xml.NAPTest{}, Testlet: xml.NAPTestlet{}, TestItem: t1})
					// note that TestItems don't link back to tests, and are in fact not unique to them. If they are not in a codeframe, we cannot trace them
				}
			}
		}

		return results, err

	}

	return resolvers
}
