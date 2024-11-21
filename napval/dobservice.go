// dobservice.go
package napval

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
)

const OOB string = "out_of_band"

// date of birth validator service

// implementation of the id service
type DOBService struct {
	layout                                                                 string
	yr3start, yr3end, yr5start, yr5end, yr7start, yr7end, yr9start, yr9end time.Time
}

// create a new id service instance
func NewDOBService(tstyr string) (*DOBService, error) {
	dob := DOBService{}

	// set up the calendar windows for testing against
	layout := "2006-01-02" // reference date format for date fields
	dob.layout = layout
	// create the testing 'now' date window baseline year
	tnow, _ := time.Parse(layout, tstyr+"-01-01")
	dob.yr3start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-9)+"-01-01")
	dob.yr3end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-8)+"-07-31")

	dob.yr5start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-11)+"-01-01")
	dob.yr5end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-10)+"-07-31")

	dob.yr7start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-13)+"-01-01")
	dob.yr7end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-12)+"-07-31")

	dob.yr9start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-15)+"-01-01")
	dob.yr9end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-14)+"-07-31")

	log.Println("\n==========================================\nNOTE:\n")
	log.Println("Baseline year for creating test/year level ranges is: ", tstyr)
	// log.Println("To change this, pass a year as comand line parameter -tstyr")
	// log.Println("e.g. \n\n> dobvalidator(.exe) -tstyr 2018")
	log.Println("\n==========================================\n")
	log.Println("Date Validation Ranges:\n")
	log.Println("Year 3 range: ", dob.yr3start.Format(layout), " - ", dob.yr3end.Format(layout))
	log.Println("Year 5 range: ", dob.yr5start.Format(layout), " - ", dob.yr5end.Format(layout))
	log.Println("Year 7 range: ", dob.yr7start.Format(layout), " - ", dob.yr7end.Format(layout))
	log.Println("Year 9 range: ", dob.yr9start.Format(layout), " - ", dob.yr9end.Format(layout))
	log.Println("\n==========================================\n")

	return &dob, nil
}

// implement the nias Service interface
func (dob *DOBService) HandleMessage(req *lib.NiasMessage) ([]lib.NiasMessage, error) {

	responses := make([]lib.NiasMessage, 0)

	rr, ok := req.Body.(xml.RegistrationRecord)
	if !ok {
		log.Println("IDService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}
	if rr.BirthDate == "" {
		return responses, nil
	}

	dateStr := rr.BirthDate
	date, err := time.Parse("02-01-2006", rr.BirthDate)
	if err == nil {
		dateStr = date.Format("2006-01-02")
		log.Println("Date converted to " + dateStr)
	}

	t, err := time.Parse(dob.layout, dateStr)
	// log.Println("Provided birth date is: ", t)
	if err != nil {
		// log.Println("unable to parse date: ", err)
		desc := "Date provided does not parse correctly for yyyy-mm-dd"
		ve := ValidationError{
			Description:  desc,
			Field:        "BirthDate",
			OriginalLine: req.SeqNo,
			Vtype:        "date",
			Severity:     "error",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	} else {

		yrlvl := rr.YearLevel
		tstlvl := rr.TestLevel
		desc := ""
		field := "BirthDate"
		ok := true
		severity := "error"
		matched, _ := regexp.MatchString("^([KFP0-9]|1[012]|UG|11MINUS|12PLUS|CC|K[34]|PS|UG(JunSec|Pri|Sec|SnrSec))$", yrlvl)
		tooyoung, _ := regexp.MatchString("^([KFP]|CC|K[34]|PS))$", yrlvl)
		switch {
		case !matched:
			// will be rejected in schema
			ok = true
		case tooyoung:
			// log.Println("student is primary")
			desc = "Year level supplied is " + yrlvl + ", does not match expected test level " + rr.TestLevel
			field = field + "/YearLevel"
			ok = false
		case strings.Contains(yrlvl, "UG") && tstlvl != dob.calculateYearLevel(t):
			// log.Println("student is ungraded")
			desc = "Year level supplied is UG, will result in SRM warning flag for test level " + rr.TestLevel
			field = field + "/TestLevel"
			ok = false
			severity = "warning"
		case yrlvl == "0":
			// log.Println("student is in year 0!!")
			desc = "Year level supplied is 0, does not match expected test level " + rr.TestLevel
			field = field + "/TestLevel"
			ok = false
		case dob.calculateYearLevel(t) == OOB:
			// log.Println("Age derived year level not in any NAPLAN window")
			desc = "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges"
			field = field + "/YearLevel"
			ok = false
			severity = "warning"
		case yrlvl != tstlvl && !strings.Contains(yrlvl, "UG"):
			desc = "Year Level " + yrlvl + " does not match Test level " + tstlvl
			field = field + "/TestLevel"
			ok = false
		default:
			field = "BirthDate"
			if yrlvl != "" && yrlvl != dob.calculateYearLevel(t) {
				// log.Println("Student is in wrong yr level: ", yrlvl)
				desc = "Student Year Level (yr " + yrlvl + ") does not match year level derived from BirthDate (yr " + dob.calculateYearLevel(t) + ")"
				field = field + "/" + "YearLevel"
				ok = false
				severity = "warning"
			}
			if tstlvl != "" && tstlvl != dob.calculateYearLevel(t) {
				// log.Println("Student is in wrong test level: ", tstlvl)
				desc = "Student Test Level (yr " + tstlvl + ") does not match year level derived from BirthDate (yr " + dob.calculateYearLevel(t) + ")"
				field = field + "/" + "TestLevel"
				ok = false
				severity = "warning"
			}
		}

		if !ok {
			ve := ValidationError{
				Description:  desc,
				Field:        field,
				OriginalLine: req.SeqNo,
				Vtype:        "date",
				Severity:     severity,
			}
			r := lib.NiasMessage{}
			r.TxID = req.TxID
			r.SeqNo = req.SeqNo
			// r.Target = VALIDATION_PREFIX
			r.Body = ve
			responses = append(responses, r)
		}

	}

	return responses, nil

}

// check whether a date is within a range
func inDateRange(start, end, check time.Time) bool {
	// assume date range is inclusive
	if start.Equal(check) || end.Equal(check) {
		return true
	}
	return check.After(start) && check.Before(end)
}

// derive the year level from the given date of birth
func (dob *DOBService) calculateYearLevel(t time.Time) string {

	yr := OOB
	if inDateRange(dob.yr3start, dob.yr3end, t) {
		// log.Println(t, "Matches Yr3: is between", yr3start, "and", yr3end, ".")
		return "3"
	}
	if inDateRange(dob.yr5start, dob.yr5end, t) {
		// log.Println(t, "Matches Yr3: is between", yr5start, "and", yr5end, ".")
		return "5"
	}
	if inDateRange(dob.yr7start, dob.yr7end, t) {
		// log.Println(t, "Matches Yr3: is between", yr7start, "and", yr7end, ".")
		return "7"
	}
	if inDateRange(dob.yr9start, dob.yr9end, t) {
		// log.Println(t, "Matches Yr3: is between", yr9start, "and", yr9end, ".")
		return "9"
	}
	// log.Println("Detected student year from dob is: ", yr)

	return yr

}
