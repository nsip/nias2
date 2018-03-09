// isrprint-resolvers.go

// resolver functions to support creation of isr printing report
package naprrql

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

// convenience type to hold isr print item data
type ISRPrintItem struct {
	StudentLocalID   string
	StudentPSI       string
	StudentFirstName string
	StudentLastName  string
	SchoolID         string
	SchoolName       string
	R_Score          float64
	R_Mean           float64
	R_Comment        string
	W_Score          float64
	W_Mean           float64
	W_Comment        string
	S_Score          float64
	S_Mean           float64
	S_Comment        string
	G_Score          float64
	G_Mean           float64
	G_Comment        string
	N_Score          float64
	N_Mean           float64
	N_Comment        string
}

// type to hold isr print item data, plus containers for Student and Attendances
type ISRPrintItemExpanded struct {
	Student         xml.RegistrationRecord
	SchoolID        string
	SchoolName      string
	R_Score         float64
	R_Mean          float64
	R_Comment       string
	R_Participation string
	W_Score         float64
	W_Mean          float64
	W_Comment       string
	W_Participation string
	S_Score         float64
	S_Mean          float64
	S_Comment       string
	S_Participation string
	G_Score         float64
	G_Mean          float64
	G_Comment       string
	G_Participation string
	N_Score         float64
	N_Mean          float64
	N_Comment       string
	N_Participation string
}

//
// if no isr item exists for this student create one and fill with
// known information
//
func (isrpi *ISRPrintItem) initialiseISRItem(schoolInfo xml.SchoolInfo, event xml.NAPEvent) {

	isrpi.SchoolID = schoolInfo.LocalId
	isrpi.SchoolName = schoolInfo.SchoolName

	student := getStudent(event.SPRefID)
	student.Flatten()
	isrpi.StudentFirstName = student.GivenName
	isrpi.StudentLastName = student.FamilyName
	isrpi.StudentLocalID = student.LocalId
	isrpi.StudentPSI = student.PlatformId

}

//
// if no isr item exists for this student create one and fill with
// known information
//
func (isrpi *ISRPrintItemExpanded) initialiseISRItemExpanded(schoolInfo xml.SchoolInfo, event xml.NAPEvent) {

	isrpi.SchoolID = schoolInfo.LocalId
	isrpi.SchoolName = schoolInfo.SchoolName

	student := getStudent(event.SPRefID)
	student.Flatten()
	isrpi.Student = student
}

//
// allocates the domain score from the response
// and the school mean from the score summary
// to the appropriate values in the ISRPrintItem
//
func (isrpi *ISRPrintItem) allocateDomainScoreAndMean(event *xml.NAPEvent,
	testLookup map[string]xml.NAPTest,
	summaryLookup map[string]xml.NAPTestScoreSummary) {

	test, ok := testLookup[event.TestID]
	if !ok {
		test = xml.NAPTest{}
	}

	summary, ok := summaryLookup[event.TestID]
	if !ok {
		summary = xml.NAPTestScoreSummary{}
	}

	resp := getResponseDomainScore(test.TestID, event.SPRefID)

	domain := strings.ToLower(test.TestContent.TestDomain)
	switch {
	case strings.Contains(domain, "gramm"):
		isrpi.G_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.G_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
	case strings.Contains(domain, "num"):
		isrpi.N_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.N_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
	case strings.Contains(domain, "read"):
		isrpi.R_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.R_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
	case strings.Contains(domain, "writ"):
		isrpi.W_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.W_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
	case strings.Contains(domain, "spell"):
		isrpi.S_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.S_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
	default:
		log.Println("Unknown test domain supplied (allocateDomainScore): ", domain)
	}

}

func (isrpi *ISRPrintItemExpanded) allocateDomainScoreAndMeanAndParticipation(event *xml.NAPEvent,
	testLookup map[string]xml.NAPTest,
	summaryLookup map[string]xml.NAPTestScoreSummary) {

	test, ok := testLookup[event.TestID]
	if !ok {
		test = xml.NAPTest{}
	}

	summary, ok := summaryLookup[event.TestID]
	if !ok {
		summary = xml.NAPTestScoreSummary{}
	}

	resp := getResponseDomainScore(test.TestID, event.SPRefID)

	domain := strings.ToLower(test.TestContent.TestDomain)
	switch {
	case strings.Contains(domain, "gramm"):
		isrpi.G_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.G_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		isrpi.G_Participation = event.ParticipationCode
	case strings.Contains(domain, "num"):
		isrpi.N_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.N_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		isrpi.N_Participation = event.ParticipationCode
	case strings.Contains(domain, "read"):
		isrpi.R_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.R_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		isrpi.R_Participation = event.ParticipationCode
	case strings.Contains(domain, "writ"):
		isrpi.W_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.W_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		isrpi.W_Participation = event.ParticipationCode
	case strings.Contains(domain, "spell"):
		isrpi.S_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
		isrpi.S_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		isrpi.S_Participation = event.ParticipationCode
	default:
		log.Println("Unknown test domain supplied (allocateDomainScore): ", domain)
	}

}

//
// create the resolver functions for isr printing
//
func buildISRPrintResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	//
	// resolver for isr printing results, line per student.
	//
	resolvers["ISRPrint/reportItems"] = func(params *graphql.ResolveParams) (interface{}, error) {

		// ISRPrintItems := make([]ISRPrintItem, 0)
		isrPrintItems := make([]ISRPrintItem, 0)

		// validate input params
		reqErr := checkISRReportParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara id from the request params
		acaraid := params.Args["schoolAcaraID"].(string)
		// log.Println("ACARA ID: ", acaraid)

		// get the test year level from the request params
		yrLvl := params.Args["testYrLevel"].(string)
		// log.Println("Yr Level: ", yrLvl)

		// get the school info for the acarid supplied
		schoolInfo, err := getSchoolInfo(acaraid)
		if err != nil {
			return isrPrintItems, err
		}
		// log.Println("School: ", schoolInfo.SchoolName)

		// get tests for yearLevel
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
		studentISRItems := make(map[string]ISRPrintItem) // index is string = student refid
		for _, event := range events {
			_, present := studentISRItems[event.SPRefID]
			if !present {
				s := ISRPrintItem{}
				s.initialiseISRItem(schoolInfo, event)
				studentISRItems[event.SPRefID] = s
			}
			s := studentISRItems[event.SPRefID]
			s.allocateDomainScoreAndMean(&event, testLookup, summaryLookup)
			studentISRItems[event.SPRefID] = s
		}

		// once collated return the isr print items
		for _, isrpi := range studentISRItems {
			isrPrintItems = append(isrPrintItems, isrpi)
		}

		return isrPrintItems, nil

	}

	resolvers["ISRPrint/reportItemsExpanded"] = func(params *graphql.ResolveParams) (interface{}, error) {

		isrPrintItems := make([]ISRPrintItemExpanded, 0)

		// validate input params
		reqErr := checkISRReportParams(params)
		if reqErr != nil {
			return nil, reqErr
		}

		// get the acara id from the request params
		acaraid := params.Args["schoolAcaraID"].(string)

		// get the test year level from the request params
		yrLvl := params.Args["testYrLevel"].(string)

		// get the school info for the acarid supplied
		schoolInfo, err := getSchoolInfo(acaraid)
		if err != nil {
			return isrPrintItems, err
		}
		// log.Println("School: ", schoolInfo.SchoolName)

		// get tests for yearLevel
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
		studentISRItems := make(map[string]ISRPrintItemExpanded) // index is string = student refid
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

		// once collated return the isr print items
		for _, isrpi := range studentISRItems {
			isrPrintItems = append(isrPrintItems, isrpi)
		}

		return isrPrintItems, nil

	}

	return resolvers
}

//
// returns set of score summary objects for the test/school combination
//
func getScoreSummaries(tests []xml.NAPTest, schoolRefID string) ([]xml.NAPTestScoreSummary, error) {

	summaries := make([]xml.NAPTestScoreSummary, 0)

	for _, test := range tests {
		keyPrefix := test.TestID + ":NAPTestScoreSummary:" + schoolRefID + ":"
		summaryObjects, err := getObjects(getIdentifiers(keyPrefix))
		if err != nil {
			return summaries, err
		}

		for _, summaryObj := range summaryObjects {
			summary, ok := summaryObj.(xml.NAPTestScoreSummary)
			if !ok {
				log.Printf("Unable to type assert summary object: %#v", summaryObj)
				continue //ignore and move on to next

			}
			summaries = append(summaries, summary)
		}

	}

	return summaries, nil
}

//
// retrieve all test events for this test at this school
//
func getTestEvents(tests []xml.NAPTest, schoolRefID string) ([]xml.NAPEvent, error) {

	events := make([]xml.NAPEvent, 0)

	for _, test := range tests {
		keyPrefix := test.TestID + ":NAPEventStudentLink:" + schoolRefID + ":"
		eventObjects, err := getObjects(getIdentifiers(keyPrefix))
		if err != nil {
			return events, err
		}

		for _, eventObj := range eventObjects {
			event, ok := eventObj.(xml.NAPEvent)
			if !ok {
				log.Printf("Unable to type assert event object: %#v", eventObj)
				continue //ignore and move on to next

			}
			events = append(events, event)
		}
	}

	return events, nil

}

//
// returns a student personal for the given refid
//
func getStudent(studentRefID string) xml.RegistrationRecord {

	student := xml.RegistrationRecord{}

	studentObjs, err := getObjects([]string{studentRefID})
	if err != nil || len(studentObjs) == 0 {
		log.Println("Unable to find student record for id: ", studentRefID)
		return student
	}

	student, ok := studentObjs[0].(xml.RegistrationRecord)
	if !ok {
		log.Printf("Unable to assert student record for object: %#v", studentObjs[0])
	}

	return student
}

//
// returns a domain score object for this student for this test
//
func getResponseDomainScore(testID string, studentRefID string) xml.DomainScore {

	ds := xml.DomainScore{}

	dsObjs, err := getObjects(getIdentifiers(testID + ":NAPStudentResponseSet:" + studentRefID))
	if err != nil || len(dsObjs) == 0 {
		// log.Println("Unable to find student domain score for id: ", testID, studentRefID)
		return ds
	}

	response, ok := dsObjs[0].(xml.NAPResponseSet)
	if !ok {
		log.Printf("Unable to assert response record for object: %#v", dsObjs[0])
	}

	return response.DomainScore

}

//
// convenience method to lookup a full school info from an acaraid
//
func getSchoolInfo(acaraid string) (xml.SchoolInfo, error) {

	si := xml.SchoolInfo{}

	schoolInfoObjects, err := getObjects(getIdentifiers(acaraid + ":"))
	if err != nil {
		return si, err
	}

	if len(schoolInfoObjects) > 1 {
		log.Println("Unexpected search result: More than one school info found for acaraid: ", acaraid)
	}

	for _, schoolInfoObj := range schoolInfoObjects {
		xsi, ok := schoolInfoObj.(xml.SchoolInfo)
		if !ok {
			log.Printf("Unable to type assert schoolinfo object: %#v", schoolInfoObj)
			continue //ignore and move on to next
		}
		si = xsi
	}

	return si, nil

}

//
// for the given yr lvl get the relevant tests from the overall naplan set
//
func getTestsForYearLevel(yrLvl string) ([]xml.NAPTest, error) {

	tests := make([]xml.NAPTest, 0)

	testObjects, err := getObjects(getIdentifiers("NAPTest:"))
	if err != nil {
		return tests, err
	}

	for _, testObj := range testObjects {
		test, ok := testObj.(xml.NAPTest)
		if !ok {
			log.Printf("Unable to type assert test object: %#v", testObj)
			continue //ignore and move on to next
		}
		if test.TestContent.TestLevel == yrLvl {
			tests = append(tests, test)
		}
	}

	/*
		if len(tests) == 0 {
			return tests, errors.New("No tests found for supplied year level")
		}
	*/
	return tests, nil

}

//
// simple check for query params that should have been passed as part
// of the api call
//
func checkISRReportParams(params *graphql.ResolveParams) error {
	//check presence of schoolAcarid
	if len(params.Args) < 2 {
		log.Println("no query params")
		return errors.New("Required variables: schoolAcaraID and testYrLevel not supplied, query aborting")
	}

	_, present := params.Args["schoolAcaraID"]
	if !present {
		log.Println("no acaraid parameter provided")
		return errors.New("Required variable: 'schoolAcaraID' not provided.")
	}

	acaraid, ok := params.Args["schoolAcaraID"].(string)
	if !ok || acaraid == "" {
		log.Println("Cannot interpret variable schoolAcaraID as meaningful string")
		return errors.New("Required variable: 'schoolACARAID' not provided.")
	}

	_, present = params.Args["testYrLevel"]
	if !present {
		log.Println("no testYrLevel parameter provided")
		return errors.New("Required variable: 'testYrLevel' not provided.")
	}

	yrlvl, ok := params.Args["testYrLevel"].(string)
	if !ok || yrlvl == "" {
		log.Println("Cannot interpret variable testYrLevel as meaningful string")
		return errors.New("Required variable: 'testYrLevel' not provided.")
	}

	return nil

}
