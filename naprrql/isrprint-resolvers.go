// isrprint-resolvers.go

// resolver functions to support creation of isr printing report
package naprrql

import (
	"errors"
	"log"
	//"strconv"
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
	YearLevel        string
	SchoolID         string
	SchoolName       string
	R_Score          string
	R_Mean           string
	R_Comment        string
	W_Score          string
	W_Mean           string
	W_Comment        string
	S_Score          string
	S_Mean           string
	S_Comment        string
	G_Score          string
	G_Mean           string
	G_Comment        string
	N_Score          string
	N_Mean           string
	N_Comment        string
}

// type to hold isr print item data, plus containers for Student and Attendances
type ISRPrintItemExpanded struct {
	Student         xml.RegistrationRecord
	SchoolID        string
	SchoolName      string
	R_Score         string
	R_Mean          string
	R_Comment       string
	R_Participation string
	W_Score         string
	W_Mean          string
	W_Comment       string
	W_Participation string
	S_Score         string
	S_Mean          string
	S_Comment       string
	S_Participation string
	G_Score         string
	G_Mean          string
	G_Comment       string
	G_Participation string
	N_Score         string
	N_Mean          string
	N_Comment       string
	N_Participation string
	G_Pathway       string
	R_Pathway       string
	W_Pathway       string
	S_Pathway       string
	N_Pathway       string
	G_Stddev        string
	R_Stddev        string
	W_Stddev        string
	S_Stddev        string
	N_Stddev        string
	G_DomainBand    string
	R_DomainBand    string
	W_DomainBand    string
	S_DomainBand    string
	N_DomainBand    string
}

//
// if no isr item exists for this student create one and fill with
// known information
//
func (isrpi *ISRPrintItem) initialiseISRItem(schoolInfo xml.SchoolInfo, event xml.NAPEvent, yrLvl string) {

	isrpi.SchoolID = schoolInfo.LocalId
	isrpi.SchoolName = schoolInfo.SchoolName

	student := getStudent(event.SPRefID)
	student.Flatten()
	isrpi.StudentFirstName = student.GivenName
	isrpi.StudentLastName = student.FamilyName
	isrpi.StudentLocalID = student.LocalId
	isrpi.StudentPSI = student.PlatformId
	isrpi.YearLevel = yrLvl
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

	resp, _ := getResponseDomainScore(test.TestID, event.SPRefID)

	domain := strings.ToLower(test.TestContent.TestDomain)
	switch {
	case strings.Contains(domain, "gramm"):
		/*
			isrpi.G_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.G_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.G_Score = resp.ScaledScoreValue
		isrpi.G_Mean = summary.DomainSchoolAverage
	case strings.Contains(domain, "num"):
		/*
			isrpi.N_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.N_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.N_Score = resp.ScaledScoreValue
		isrpi.N_Mean = summary.DomainSchoolAverage
	case strings.Contains(domain, "read"):
		/*
			isrpi.R_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.R_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.R_Score = resp.ScaledScoreValue
		isrpi.R_Mean = summary.DomainSchoolAverage
	case strings.Contains(domain, "writ"):
		/*
			isrpi.W_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.W_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.W_Score = resp.ScaledScoreValue
		isrpi.W_Mean = summary.DomainSchoolAverage
	case strings.Contains(domain, "spell"):
		/*
			isrpi.S_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.S_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.S_Score = resp.ScaledScoreValue
		isrpi.S_Mean = summary.DomainSchoolAverage
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

	resp, pathway := getResponseDomainScore(test.TestID, event.SPRefID)

	domain := strings.ToLower(test.TestContent.TestDomain)
	switch {
	case strings.Contains(domain, "gramm"):
		/*
			isrpi.G_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.G_Stddev, _ = strconv.ParseFloat(resp.ScaledScoreStandardError, 32)
			isrpi.G_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.G_Score = resp.ScaledScoreValue
		isrpi.G_Stddev = resp.ScaledScoreStandardError
		isrpi.G_Mean = summary.DomainSchoolAverage
		isrpi.G_Participation = event.ParticipationCode
		isrpi.G_Pathway = pathway
		isrpi.G_DomainBand = resp.StudentDomainBand
	case strings.Contains(domain, "num"):
		/*
			isrpi.N_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.N_Stddev, _ = strconv.ParseFloat(resp.ScaledScoreStandardError, 32)
			isrpi.N_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.N_Score = resp.ScaledScoreValue
		isrpi.N_Stddev = resp.ScaledScoreStandardError
		isrpi.N_Mean = summary.DomainSchoolAverage
		isrpi.N_Participation = event.ParticipationCode
		isrpi.N_Pathway = pathway
		isrpi.N_DomainBand = resp.StudentDomainBand
	case strings.Contains(domain, "read"):
		/*
			isrpi.R_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.R_Stddev, _ = strconv.ParseFloat(resp.ScaledScoreStandardError, 32)
			isrpi.R_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.R_Score = resp.ScaledScoreValue
		isrpi.R_Stddev = resp.ScaledScoreStandardError
		isrpi.R_Mean = summary.DomainSchoolAverage
		isrpi.R_Participation = event.ParticipationCode
		isrpi.R_Pathway = pathway
		isrpi.R_DomainBand = resp.StudentDomainBand
	case strings.Contains(domain, "writ"):
		/*
			isrpi.W_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.W_Stddev, _ = strconv.ParseFloat(resp.ScaledScoreStandardError, 32)
			isrpi.W_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.W_Score = resp.ScaledScoreValue
		isrpi.W_Stddev = resp.ScaledScoreStandardError
		isrpi.W_Mean = summary.DomainSchoolAverage
		isrpi.W_Participation = event.ParticipationCode
		isrpi.W_Pathway = pathway
		isrpi.W_DomainBand = resp.StudentDomainBand
	case strings.Contains(domain, "spell"):
		/*
			isrpi.S_Score, _ = strconv.ParseFloat(resp.ScaledScoreValue, 32)
			isrpi.S_Stddev, _ = strconv.ParseFloat(resp.ScaledScoreStandardError, 32)
			isrpi.S_Mean, _ = strconv.ParseFloat(summary.DomainSchoolAverage, 32)
		*/
		isrpi.S_Score = resp.ScaledScoreValue
		isrpi.S_Stddev = resp.ScaledScoreStandardError
		isrpi.S_Mean = summary.DomainSchoolAverage
		isrpi.S_Participation = event.ParticipationCode
		isrpi.S_Pathway = pathway
		isrpi.S_DomainBand = resp.StudentDomainBand
	default:
		log.Println("Unknown test domain supplied (allocateDomainScore): ", domain)
	}

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
// returns a domain score object and Pathway for this student for this test
//
func getResponseDomainScore(testID string, studentRefID string) (xml.DomainScore, string) {

	ds := xml.DomainScore{}

	dsObjs, err := getObjects(getIdentifiers(testID + ":NAPStudentResponseSet:" + studentRefID))
	if err != nil || len(dsObjs) == 0 {
		// log.Println("Unable to find student domain score for id: ", testID, studentRefID)
		return ds, ""
	}

	response, ok := dsObjs[0].(xml.NAPResponseSet)
	if !ok {
		log.Printf("Unable to assert response record for object: %#v", dsObjs[0])
	}

	return response.DomainScore, response.PathTakenForDomain

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
