// report-resolvers.go

//
// resolver logic for composite reporting objects
//
package naprrql

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/nats-io/nuid"
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

var lem = jargon.NewLemmatizer(contractions.Dictionary, 3)
var tokenRe = regexp.MustCompile("[a-zA-Z0-9]")
var hyphens = regexp.MustCompile("-+")

func countwords(html string) int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0
	}
	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})
	// Jargon lemmatiser tokenises text, and resolves contractions
	// We will resolve hyphenated compounds ourselves
	tokens := jargon.Tokenize(strings.NewReader(doc.Text()))
	lemmas := lem.Lemmatize(tokens)
	wc := 0
	for {
		lemma := lemmas.Next()
		if lemma == nil {
			break
		}
		wordpart := hyphens.Split(lemma.String(), -1)
		for _, w := range wordpart {
			if len(tokenRe.FindString(w)) > 0 {
				wc = wc + 1
			}
		}
	}
	return wc
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

// The results set coming out of the NAP includes empty testlets with no items in them. These testlets need to be removed from the result set,
// which is assuming three or four contentful testlets for the purposes of reporting (particularly domain_scores_summary_event_report_by_school,
// which genereates the nswPrintAll report, which iterates through testlets 0...3 for each domain.
func TrimEmptyTestlets(r xml.NAPResponseSet) xml.NAPResponseSet {
	delete_testlets := make(map[int]bool)
	for i, testlet := range r.TestletList.Testlet {
		if len(testlet.ItemResponseList.ItemResponse) == 0 {
			delete_testlets[i] = true
		}
	}
	if len(delete_testlets) > 0 {
		newtestlets := make([]xml.NAPResponseSet_Testlet, 0)
		for i, testlet := range r.TestletList.Testlet {
			if _, ok := delete_testlets[i]; !ok {
				newtestlets = append(newtestlets, testlet)
			}
		}
		r.TestletList.Testlet = newtestlets
	}
	return r
}

func domain_scores_event_report_by_school(params *graphql.ResolveParams, yrLvl string, domain string) (interface{}, error) {

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
			if len(yrLvl) > 0 && test.TestContent.TestLevel != yrLvl {
				continue
			}
			if len(domain) > 0 && test.TestContent.TestDomain != domain {
				continue
			}

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

func domain_scores_summary_event_report_by_school(params *graphql.ResolveParams, yrLvl string, domain string) (interface{}, error) {

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

	// get school infos and school summaries
	schools := make(map[string]xml.SchoolInfo)                             // key string = acara id
	schoolsummaries := make(map[string]map[string]xml.NAPTestScoreSummary) // key string = acara id + test refid

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
			schools[acaraid] = si
			summaryrefid := getIdentifiers(si.RefId + ":NAPTestScoreSummary:")
			ssObjects, err := getObjects(summaryrefid)
			if err != nil {
				return []interface{}{}, err
			}
			for _, sso := range ssObjects {
				ss := sso.(xml.NAPTestScoreSummary)
				if _, ok := schoolsummaries[acaraid]; !ok {
					schoolsummaries[acaraid] = make(map[string]xml.NAPTestScoreSummary)
				}
				schoolsummaries[acaraid][ss.NAPTestRefId] = ss
			}
		}
	}

	studentObjs, err := getObjects(studentids)
	if err != nil {
		return []interface{}{}, err
	}

	// iterate students and assemble Event/Response Data Set
	results := make([]EventResponseSummaryAllDomainsDataSet, 0)
	for _, studentObj := range studentObjs {
		student, _ := studentObj.(xml.RegistrationRecord)
		if len(yrLvl) > 0 && student.TestLevel != yrLvl {
			continue
		}
		studentEventIds := getIdentifiers(student.RefId + ":NAPEventStudentLink:")
		if len(studentEventIds) < 1 {
			// log.Println("no events found for student: ", student.RefId)
			continue
		}
		eventObjs, err := getObjects(studentEventIds)
		if err != nil {
			return []interface{}{}, err
		}
		perdomain_slice := make([]EventResponseSummaryPerDomain, 0)
		schoolid := ""
		for _, eventObj := range eventObjs {
			event := eventObj.(xml.NAPEvent)
			// report assumes that the school is the same across all domains; it may not be
			if _, ok := schools[event.SchoolID]; ok {
				schoolid = event.SchoolID
			}
			pnpcodelistmap := pnpcodelistmap(event)
			testObj, err := getObjects([]string{event.TestID})
			if err != nil {
				return []interface{}{}, err
			}
			test := testObj[0].(xml.NAPTest)
			if len(domain) > 0 && test.TestContent.TestDomain != domain {
				continue
			}
			responseIds := getIdentifiers(test.TestID + ":NAPStudentResponseSet:" + student.RefId)
			response := xml.NAPResponseSet{}
			if len(responseIds) > 0 {
				responseObjs, err := getObjects(responseIds)
				if err != nil {
					return []interface{}{}, err
				}
				// response = TrimEmptyTestlets(responseObjs[0].(xml.NAPResponseSet))
				response = responseObjs[0].(xml.NAPResponseSet)
			}
			perdomain := EventResponseSummaryPerDomain{Domain: test.TestContent.TestDomain,
				Event:          event,
				Test:           test,
				PNPCodeListMap: pnpcodelistmap,
				Summary:        schoolsummaries[event.SchoolID][test.TestID],
				Response:       response,
			}
			perdomain_slice = append(perdomain_slice, perdomain)
		}
		erds := EventResponseSummaryAllDomainsDataSet{Student: student,
			School:                        schools[schoolid],
			EventResponseSummaryPerDomain: perdomain_slice,
		}
		results = append(results, erds)

	}
	return results, nil
}

// Code item responses in EventResponseSummaryAllDomainsDataSet slice, according to ACARA coding conventions:
// if Multiple Choice, and Item Response is known, 1 2 3 4 for multiple choice response A B C D, 9 for no response
// else: 1 for correct response, 0 for incorrect response, 9 for no response
// Coding is injected into Item Response field ResponseCorrectness, overwriting value already there
func acaraCodeItemResponses(erds []EventResponseSummaryAllDomainsDataSet) interface{} {
	// map of Test Item ID to whether they are multiple choice items or not
	itemmap := make(map[string]bool)

	testitem_ids := getIdentifiers("NAPTestItem:")
	testitems, err := getObjects(testitem_ids)
	if err != nil {
		for _, io := range testitems {
			item, _ := io.(xml.NAPTestItem)
			itemmap[item.ItemID] = item.TestItemContent.ItemType == "MC"
		}
	}
	for i, _ := range erds {
		for j, _ := range erds[i].EventResponseSummaryPerDomain {
			for k, _ := range erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet {
				for l, _ := range erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse {
					r := erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l]
					if r.ResponseCorrectness == "NotAttempted" || r.ResponseCorrectness == "NotInPath" {
						erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "9"
					} else if itemmap[r.ItemRefID] && len(r.Response) > 0 {
						switch r.Response {
						case "A":
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "1"
						case "B":
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "2"
						case "C":
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "3"
						case "D":
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "4"
						default:
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "7"
						}
					} else {
						if r.ResponseCorrectness == "Correct" {
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "1"
						} else {
							erds[i].EventResponseSummaryPerDomain[j].Response.TestletList.Testlet[k].ItemResponseList.ItemResponse[l].ResponseCorrectness = "0"
						}
					}
				}
			}
		}
	}
	return erds
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
		if err != nil {
			return []interface{}{}, err
		}
		summary_datasets := make([]ScoreSummaryDataSet, 0)
		for _, summary := range summaries {
			summ, _ := summary.(xml.NAPTestScoreSummary)
			testid := []string{summ.NAPTestRefId}
			obj, err := getObjects(testid)
			if err != nil {
				return []interface{}{}, err
			}
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

	// Because there are no Writing Yr 3 enrolments in the platform at all, we are cheating by providing the Numeracy enrolments instead
	resolvers["NaplanData/domain_scores_event_report_by_school_writing_yr3"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_event_report_by_school(params, "3", "Numeracy")
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_writing_yr5"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_event_report_by_school(params, "5", "Writing")
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_writing_yr7"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_event_report_by_school(params, "7", "Writing")
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_writing_yr9"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_event_report_by_school(params, "9", "Writing")
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_acara_yr3"] = func(params *graphql.ResolveParams) (interface{}, error) {
		e, err := domain_scores_summary_event_report_by_school(params, "3", "")
		return acaraCodeItemResponses(e.([]EventResponseSummaryAllDomainsDataSet)), err
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_acara_yr5"] = func(params *graphql.ResolveParams) (interface{}, error) {
		e, err := domain_scores_summary_event_report_by_school(params, "5", "")
		return acaraCodeItemResponses(e.([]EventResponseSummaryAllDomainsDataSet)), err
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_acara_yr7"] = func(params *graphql.ResolveParams) (interface{}, error) {
		e, err := domain_scores_summary_event_report_by_school(params, "7", "")
		return acaraCodeItemResponses(e.([]EventResponseSummaryAllDomainsDataSet)), err
	}
	resolvers["NaplanData/domain_scores_event_report_by_school_acara_yr9"] = func(params *graphql.ResolveParams) (interface{}, error) {
		e, err := domain_scores_summary_event_report_by_school(params, "9", "")
		return acaraCodeItemResponses(e.([]EventResponseSummaryAllDomainsDataSet)), err
	}
	resolvers["NaplanData/domain_scores_event_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_event_report_by_school(params, "", "")
	}
	resolvers["NaplanData/domain_scores_summary_event_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return domain_scores_summary_event_report_by_school(params, "", "")
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
			log.Printf("No match for Codeframe %#v\n", codeframeIds)
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
				log.Printf("No match for Test %s\n", codeframe.NAPTestRefId)
				return []interface{}{}, err
			}
			test, _ := testObj[0].(xml.NAPTest)
			for _, cf_testlet := range codeframe.TestletList.Testlet {
				tlObj, err := getObjects([]string{cf_testlet.NAPTestletRefId})
				if err != nil {
					log.Printf("No match for Testlet %s\n", cf_testlet.NAPTestletRefId)
					return []interface{}{}, err
				}
				tl, _ := tlObj[0].(xml.NAPTestlet)
				for _, cf_item := range cf_testlet.TestItemList.TestItem {
					tiObj, err := getObjects([]string{cf_item.TestItemRefId})
					if err != nil {
						log.Printf("No match for Test Item %s\n", cf_item.TestItemRefId)
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

		// get school names and student records
		schoolnames := make(map[string]string)
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

			for idx, testlet := range resp.TestletList.Testlet {
				// if there are no item responses included under a testlet, generate an entry with a dummy item response,
				// so that the discrepancy can be picked up downstream by QA reports
				if len(testlet.ItemResponseList.ItemResponse) == 0 {
					log.Printf("Empty response testlet\n")
					resp.TestletList.Testlet[idx].ItemResponseList.ItemResponse = append(testlet.ItemResponseList.ItemResponse,
						xml.NAPResponseSet_ItemResponse{LocalID: "empty_response_testlet"})
				}
			}
			for _, testlet := range resp.TestletList.Testlet {
				for _, item_response := range testlet.ItemResponseList.ItemResponse {

					resp1 := resp // pruned copy of response
					resp1.TestletList.Testlet = make([]xml.NAPResponseSet_Testlet, 1)
					resp1.TestletList.Testlet[0] = testlet
					resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse = make([]xml.NAPResponseSet_ItemResponse, 1)
					resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse[0] = item_response

					var item xml.NAPTestItem
					var ok bool
					if item_response.LocalID == "empty_response_testlet" {
						item = xml.NAPTestItem{}
						item.TestItemContent = xml.TestItemContent{NAPTestItemLocalId: "empty_response_testlet"}
					} else {
						items, err := getObjects([]string{item_response.ItemRefID})
						item, ok = items[0].(xml.NAPTestItem)
						if err != nil || !ok {
							return []interface{}{}, err
						}
					}

					testlets, err := getObjects([]string{testlet.NapTestletRefId})
					tl, ok := testlets[0].(xml.NAPTestlet)
					if err != nil || !ok {
						return []interface{}{}, err
					}

					irds := ItemResponseDataSet{TestItem: item, Response: resp1,
						Student: student, Test: test, Testlet: tl,
						SchoolDetails:     SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
						ParticipationCode: event.ParticipationCode}
					results = append(results, irds)
				}
			}
		}
		return results, nil
	}

	resolvers["NaplanData/item_writing_results_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

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

		// get school names and student records
		schoolnames := make(map[string]string)
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

			if test.TestContent.TestDomain != "Writing" {
				continue
			}

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
						SchoolDetails:     SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
						ParticipationCode: event.ParticipationCode}
					results = append(results, irds)
				}
			}
		}
		return results, nil
	}

	// same as the above, but adds AnonymisedId; iterates all events, not just responses;
	// and leaves out open events (P participation status; no response container)
	resolvers["NaplanData/writing_item_for_marking_report_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {
		wordcount := 0
		reqErr := checkRequiredParams(params)
		if reqErr != nil {
			log.Println("writing_item_for_marking_report_by_school #1", reqErr)
			return nil, reqErr
		}
		// get the acara ids from the request params
		acaraids := make([]string, 0)
		for _, a_id := range params.Args["acaraIDs"].([]interface{}) {
			acaraid, _ := a_id.(string)
			acaraids = append(acaraids, acaraid)
		}

		// get school names and student records
		schoolnames := make(map[string]string)
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)

			schoolrefid := getIdentifiers(acaraid + ":")
			siObjects, err := getObjects(schoolrefid)
			if err != nil {
				log.Println("writing_item_for_marking_report_by_school #2", err)
				return []interface{}{}, err
			}
			for _, sio := range siObjects {
				si, _ := sio.(xml.SchoolInfo)
				schoolnames[acaraid] = si.SchoolName
			}
		}
		// convenience map to avoid revisiting db for tests
		testLookup := make(map[string]xml.NAPTest) // key string = test refid
		// get tests for yearLevel
		tests, err := getObjects(getIdentifiers("NAPTest:"))
		if err != nil {
			log.Println("writing_item_for_marking_report_by_school #3", err)
			return nil, err
		}
		for _, test := range tests {
			t := test.(xml.NAPTest)
			testLookup[t.TestID] = t
		}

		// construct RDS by including referenced test
		results := make([]ItemResponseDataSetWordCount, 0)
		//log.Printf("School %+v: %d Writing tests for %d students\n", acaraids, len(testLookup), len(studentids))
		for testid, test := range testLookup {
			if test.TestContent.TestDomain != "Writing" {
				continue
			}
			for _, studentid := range studentids {
				var ok bool

				students, err := getObjects([]string{studentid})
				student, ok := students[0].(xml.RegistrationRecord)
				if err != nil || !ok {
					log.Println("writing_item_for_marking_report_by_school #4", err)
					return []interface{}{}, err
				}
				student.Flatten()
				student.OtherIdList.OtherId = append(student.OtherIdList.OtherId, xml.XMLAttributeStruct{Type: "AnonymisedId", Value: nuid.New().Next()})

				eventsRefId := getIdentifiers("event_by_student_test:" + studentid + ":" + testid + ":")
				events, err := getObjects(eventsRefId)
				var event xml.NAPEvent
				if len(events) == 0 {
					continue
				} else {
					event, ok = events[0].(xml.NAPEvent)
					if err != nil || !ok {
						log.Println("writing_item_for_marking_report_by_school #5", err)
						return []interface{}{}, err
					}
				}
				/*
					for i, e := range events {
						log.Printf("%d, %s : %s : %s\n", i, e.(xml.NAPEvent).PSI, e.(xml.NAPEvent).ParticipationCode, e.(xml.NAPEvent).NAPTestLocalID)
					}
				*/

				responseRefId := getIdentifiers(testid + ":NAPStudentResponseSet:" + studentid)
				responses, err := getObjects(responseRefId)
				if err != nil {
					/*
						log.Println("writing_item_for_marking_report_by_school #6", err)
						return []interface{}{}, err
					*/
					// deem this an open event, which has not yet had a response registered; ignore
					continue
				}
				var resp xml.NAPResponseSet
				if len(responses) == 0 {
					ok = false
				} else {
					/* ok iff there is an actual response recorded -- even if it is just &nbsp; */
					resp, ok = responses[0].(xml.NAPResponseSet)
					if ok {
						ok = (len(resp.TestletList.Testlet) != 0)
						if ok {
							ok = (len(resp.TestletList.Testlet[0].ItemResponseList.ItemResponse) != 0)
						}
					}
					if ok {
						ok = (len(resp.TestletList.Testlet[0].ItemResponseList.ItemResponse[0].Response) != 0)
					}
				}
				if ok {
					for _, testlet := range resp.TestletList.Testlet {
						for _, item_response := range testlet.ItemResponseList.ItemResponse {
							wordcount = countwords(item_response.Response)

							resp1 := resp // pruned copy of response
							resp1.TestletList.Testlet = make([]xml.NAPResponseSet_Testlet, 1)
							resp1.TestletList.Testlet[0] = testlet
							resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse = make([]xml.NAPResponseSet_ItemResponse, 1)
							resp1.TestletList.Testlet[0].ItemResponseList.ItemResponse[0] = item_response

							items, err := getObjects([]string{item_response.ItemRefID})
							item, ok := items[0].(xml.NAPTestItem)
							if err != nil || !ok {
								log.Println("writing_item_for_marking_report_by_school #7", err)
								return []interface{}{}, err
							}

							testlets, err := getObjects([]string{testlet.NapTestletRefId})
							tl := testlets[0].(xml.NAPTestlet)
							if err != nil || !ok {
								log.Println("writing_item_for_marking_report_by_school #8", err)
								return []interface{}{}, err
							}
							irds := ItemResponseDataSetWordCount{TestItem: item, Response: resp1,
								Student: student, Test: test, Testlet: tl,
								SchoolDetails:     SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
								ParticipationCode: event.ParticipationCode,
								WordCount:         strconv.Itoa(wordcount)}
							results = append(results, irds)
						}
					}
				} else {
					if event.ParticipationCode == "P" && len(event.StartTime) == 0 {
						// no response recorded; consider this to be an open event, and ignore it
					} else {
						irds := ItemResponseDataSetWordCount{TestItem: xml.NAPTestItem{}, Response: xml.NAPResponseSet{},
							Student: student, Test: test, Testlet: xml.NAPTestlet{},
							SchoolDetails:     SchoolDetails{ACARAId: event.SchoolID, SchoolName: schoolnames[event.SchoolID]},
							ParticipationCode: event.ParticipationCode,
							WordCount:         "0"}
						results = append(results, irds)
					}
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
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
		for _, event := range events {
			t, _ := event.(xml.NAPEvent)
			results = checkGuid(results, t.EventID, "event", t.SPRefID, "student", ids)
			results = checkGuid(results, t.EventID, "event", t.SchoolRefID, "school", ids)
			results = checkGuid(results, t.EventID, "event", t.TestID, "test", ids)
		}
		events = nil

		responses, err := getObjects(response_ids)
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
		for _, testitem := range testitems {
			t, _ := testitem.(xml.NAPTestItem)
			for _, s := range t.TestItemContent.ItemSubstitutedForList.SubstituteItem {
				results = checkGuid(results, t.ItemID, "testitem", s.SubstituteItemRefId, "testitem", ids)
			}
		}
		testitems = nil

		testlets, err := getObjects(testlet_ids)
		if err != nil {
			return []interface{}{}, err
		}
		for _, testlet := range testlets {
			t, _ := testlet.(xml.NAPTestlet)
			results = checkGuid(results, t.TestletID, "testlet", t.NAPTestRefId, "test", ids)
			for _, ti := range t.TestItemList.TestItem {
				results = checkGuid(results, t.TestletID, "testlet", ti.TestItemRefId, "testitem", ids)
			}
		}
		testlets = nil

		summarys, err := getObjects(scoresummary_ids)
		if err != nil {
			return []interface{}{}, err
		}
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
		if err != nil {
			return []interface{}{}, err
		}
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
	// resolver for isr printing results, line per student.
	//
	resolvers["NaplanData/isrReportItems"] = func(params *graphql.ResolveParams) (interface{}, error) {

		isrPrintItems := make([]ISRPrintItem, 0)

		// validate input params
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

		/*
		  // get the test year level from the request params
		  yrLvl := params.Args["testYrLevel"].(string)
		  // log.Println("Yr Level: ", yrLvl)
		*/

		studentISRItems := make(map[string]ISRPrintItem) // index is string = student refid
		for _, acaraid := range acaraids {
			// get the school info for the acarid supplied
			schoolInfo, err := getSchoolInfo(acaraid)
			if err != nil {
				return isrPrintItems, err
			}
			// log.Println("School: ", schoolInfo.SchoolName)

			// get tests for yearLevel
			for _, yrLvl := range []string{"3", "5", "7", "9"} {

				tests, err := getTestsForYearLevel(yrLvl)
				if err != nil {
					return isrPrintItems, err
				}
				// convenience map to avoid revisiting db for tests
				testLookup := make(map[string]xml.NAPTest) // key string = test refid
				for _, test := range tests {
					t := test
					testLookup[t.TestID] = t
				}

				// get events for each test at this school
				events := make([]xml.NAPEvent, 0)
				events, err = getTestEvents(tests, schoolInfo.RefId)
				if err != nil {
					return isrPrintItems, err
				}

				// get score summaries for tests at this school
				summaries := make([]xml.NAPTestScoreSummary, 0)
				summaries, err = getScoreSummaries(tests, schoolInfo.RefId)
				if err != nil {
					return isrPrintItems, nil
				}
				// convenience map to avoid revisiting db for summaries
				summaryLookup := make(map[string]xml.NAPTestScoreSummary)
				for _, summary := range summaries {
					s := summary
					summaryLookup[s.NAPTestRefId] = s
				}

				// iterate events creating collated items for students
				for _, event := range events {
					_, present := studentISRItems[event.SPRefID]
					if !present {
						s := ISRPrintItem{}
						s.initialiseISRItem(schoolInfo, event, yrLvl)
						studentISRItems[event.SPRefID] = s
					}
					s := studentISRItems[event.SPRefID]
					s.allocateDomainScoreAndMean(&event, testLookup, summaryLookup)
					studentISRItems[event.SPRefID] = s
				}
			}
		}

		// once collated return the isr print items
		for _, isrpi := range studentISRItems {
			isrPrintItems = append(isrPrintItems, isrpi)
		}

		return isrPrintItems, nil

	}
	resolvers["NaplanData/isrReportItemsExpanded"] = func(params *graphql.ResolveParams) (interface{}, error) {

		isrPrintItems := make([]ISRPrintItemExpanded, 0)

		// validate input params
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
		studentISRItems := make(map[string]ISRPrintItemExpanded) // index is string = student refid

		/*
		  // get the test year level from the request params
		  yrLvl := params.Args["testYrLevel"].(string)
		*/
		for _, acaraid := range acaraids {
			// get the school info for the acarid supplied
			schoolInfo, err := getSchoolInfo(acaraid)
			if err != nil {
				return isrPrintItems, err
			}
			// log.Println("School: ", schoolInfo.SchoolName)

			// get tests for yearLevel
			for _, yrLvl := range []string{"3", "5", "7", "9"} {
				tests, err := getTestsForYearLevel(yrLvl)
				if err != nil {
					return isrPrintItems, err
				}
				// convenience map to avoid revisiting db for tests
				testLookup := make(map[string]xml.NAPTest) // key string = test refid
				for _, test := range tests {
					t := test
					testLookup[t.TestID] = t
				}

				// get events for each test at this school
				events := make([]xml.NAPEvent, 0)
				events, err = getTestEvents(tests, schoolInfo.RefId)
				if err != nil {
					return isrPrintItems, err
				}

				// get score summaries for tests at this school
				summaries := make([]xml.NAPTestScoreSummary, 0)
				summaries, err = getScoreSummaries(tests, schoolInfo.RefId)
				if err != nil {
					return isrPrintItems, nil
				}
				// convenience map to avoid revisiting db for summaries
				summaryLookup := make(map[string]xml.NAPTestScoreSummary)
				for _, summary := range summaries {
					s := summary
					summaryLookup[s.NAPTestRefId] = s
				}

				// iterate events creating collated items for students
				for _, event := range events {
					_, present := studentISRItems[event.SPRefID]
					if !present {
						s := ISRPrintItemExpanded{}
						s.initialiseISRItemExpanded(schoolInfo, event)
						studentISRItems[event.SPRefID] = s
					}
					s := studentISRItems[event.SPRefID]
					s.allocateDomainScoreAndMeanAndParticipation(&event, testLookup, summaryLookup)
					studentISRItems[event.SPRefID] = s
				}
			}
		}

		// once collated return the isr print items
		for _, isrpi := range studentISRItems {
			isrPrintItems = append(isrPrintItems, isrpi)
		}

		return isrPrintItems, nil

	}

	resolvers["NaplanData/test_year_level_discrepancy_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

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

		siObjects, err := getObjects(studentids)
		studentPersonals := make([]xml.RegistrationRecord, 0)
		for _, sio := range siObjects {
			si, _ := sio.(xml.RegistrationRecord)
			if si.TestLevel != si.YearLevel {
				studentPersonals = append(studentPersonals, si)
			}
		}
		return studentPersonals, err
	}

	resolvers["NaplanData/student_event_acara_id_discrepancy_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

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

		results := make([]EventResponseDataSet, 0)

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

		siObjects, err := getObjects(studentids)
		for _, sio := range siObjects {
			student, _ := sio.(xml.RegistrationRecord)
			studentEventIds := getIdentifiers(student.RefId + ":NAPEventStudentLink:")
			if len(studentEventIds) < 1 {
				continue
			}
			eventObjs, err := getObjects(studentEventIds)
			if err != nil {
				return []interface{}{}, err
			}
			for _, eventObj := range eventObjs {
				event := eventObj.(xml.NAPEvent)
				if student.ASLSchoolId != event.SchoolID {
					results = append(results, EventResponseDataSet{Event: event, Student: student})
				}
			}
		}
		return results, err
	}

	resolvers["NaplanData/extraneous_characters_student_report"] = func(params *graphql.ResolveParams) (interface{}, error) {
		log.Println("Launching")
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

		results := make([]xml.RegistrationRecord, 0)

		// get students for the schools
		studentids := make([]string, 0)
		for _, acaraid := range acaraids {
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)
		}

		re := regexp.MustCompile("[^a-zA-Z' -]")

		siObjects, err := getObjects(studentids)
		for _, sio := range siObjects {
			student, _ := sio.(xml.RegistrationRecord)
			if len(re.FindString(student.FamilyName)) > 0 || len(re.FindString(student.GivenName)) > 0 ||
				len(re.FindString(student.MiddleName)) > 0 || len(re.FindString(student.PreferredName)) > 0 {
				results = append(results, student)
			}
		}
		return results, err

	}

	// This report hard codes the number of expected testlets for each domain and year level
	resolvers["NaplanData/missing_testlets_by_school"] = func(params *graphql.ResolveParams) (interface{}, error) {

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
				if event.ParticipationCode != "P" {
					// ignore "F", those testlets are not adaptive, so they won't have the same testlet count
					continue
				}
				test := testLookup[event.TestID]
				responseIds := getIdentifiers(test.TestID + ":NAPStudentResponseSet:" + student.RefId)
				response := xml.NAPResponseSet{}
				if len(responseIds) > 0 {
					responseObjs, err := getObjects(responseIds)
					if err != nil {
						return []interface{}{}, err
					}

					testlevel := test.TestContent.TestLevel
					var expectednodes int
					switch test.TestContent.TestDomain {
					case "Writing":
						expectednodes = 1
					case "Numeracy":
						if testlevel == "7" || testlevel == "9" {
							expectednodes = 4
						} else {
							expectednodes = 3
						}
					case "Reading":
						expectednodes = 3
					case "Grammar and Punctuation":
						expectednodes = 1
					case "Spelling":
						expectednodes = 3
					default:
						expectednodes = -1
					}
					for _, responseObj := range responseObjs {

						response = responseObj.(xml.NAPResponseSet)
						testletnodes := strings.Split(response.PathTakenForDomain, ":")
						if len(testletnodes) == expectednodes {
							continue
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
			}
		}
		return results, nil

	}

	resolvers["NaplanData/pnp_events_report"] = func(params *graphql.ResolveParams) (interface{}, error) {

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
		// iterate students and assemble Event/Response Data Set
		results := make([]xml.NAPEvent, 0)
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
				if len(event.Adjustment.PNPCodelist.PNPCode) > 0 {
					results = append(results, event)
				}
			}
		}
		return results, nil

	}

	resolvers["NaplanData/homeschooled_student_tests_report"] = func(params *graphql.ResolveParams) (interface{}, error) {

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

		results := make([]EventResponseSummaryDataSet, 0)

		for _, acaraid := range acaraids {
			schoolrefid := getIdentifiers(acaraid + ":")
			siObjects, err := getObjects(schoolrefid)
			if err != nil {
				return []interface{}{}, err
			}
			var school xml.SchoolInfo
			school = xml.SchoolInfo{}
			for _, sio := range siObjects {
				school, _ = sio.(xml.SchoolInfo)
			}

			// get students for the schools
			studentids := make([]string, 0)
			key := "student_by_acaraid:" + acaraid
			studentRefIds := getIdentifiers(key)
			studentids = append(studentids, studentRefIds...)

			siObjects, err = getObjects(studentids)
			for _, sio := range siObjects {
				student, _ := sio.(xml.RegistrationRecord)
				if student.HomeSchooledStudent != "Y" {
					continue
				}
				studentEventIds := getIdentifiers(student.RefId + ":NAPEventStudentLink:")
				if len(studentEventIds) < 1 {
					continue
				}
				eventObjs, err := getObjects(studentEventIds)
				if err != nil {
					return []interface{}{}, err
				}
				for _, eventObj := range eventObjs {
					event := eventObj.(xml.NAPEvent)
					test := testLookup[event.TestID]
					results = append(results, EventResponseSummaryDataSet{Event: event,
						Test:           test,
						Student:        student,
						Response:       xml.NAPResponseSet{},
						Summary:        xml.NAPTestScoreSummary{},
						School:         school,
						PNPCodeListMap: pnpcodelistmap(event)})
				}
			}
		}
		return results, nil
	}

	return resolvers
}
