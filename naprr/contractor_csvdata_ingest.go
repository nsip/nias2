package naprr

import (
	"bufio"
	//ms "github.com/mitchellh/mapstructure"
	"github.com/nats-io/nuid"
	"strconv"
	//"github.com/nats-io/go-nats-streaming"
	"fmt"
	"github.com/nsip/nias2/lib"
	nxml "github.com/nsip/nias2/xml"
	"github.com/twinj/uuid"
	"log"
	"path/filepath"
	"sync"
)

func (di *DataIngest) RunYr3Writing() {

	uuid.Init()
	csvFiles := parsePearsonCSVFileDirectory()

	var wg sync.WaitGroup

	for _, csvFile := range csvFiles {
		wg.Add(1)
		go di.ingestPearsonResultsFile(csvFile, &wg)
	}

	wg.Wait()

	di.finaliseTransactions()

	di.sc.Close()

	log.Println("All data files read, ingest complete.")
}

func parsePearsonCSVFileDirectory() []string {
	files, _ := filepath.Glob("./in/Pearson/*.txt")
	if len(files) == 0 {
		log.Fatalln("No results data zip files found in input folder.")
	}
	return files
}

type fixedLengthField struct {
	name   string
	start  int
	length int
}

var pearsonStruct = []fixedLengthField{
	fixedLengthField{name: "booklet_id", start: 1, length: 8},
	fixedLengthField{name: "vsn", start: 9, length: 10},
	fixedLengthField{name: "fname", start: 19, length: 22},
	fixedLengthField{name: "lname", start: 41, length: 25},
	fixedLengthField{name: "dob_dm", start: 81, length: 4},
	fixedLengthField{name: "dob_year_1", start: 85, length: 2},
	fixedLengthField{name: "dob_year_2", start: 87, length: 2},
	fixedLengthField{name: "gender", start: 89, length: 1},
	fixedLengthField{name: "indigenous_status", start: 90, length: 1},
	fixedLengthField{name: "lbote", start: 91, length: 1},
	fixedLengthField{name: "school_vcaa_code", start: 92, length: 5},
	fixedLengthField{name: "home_group", start: 97, length: 12},
	fixedLengthField{name: "test_type", start: 109, length: 4},
	fixedLengthField{name: "paper_number", start: 113, length: 1},
	fixedLengthField{name: "participation_linked", start: 114, length: 1},
	fixedLengthField{name: "student_attendance", start: 115, length: 1},
	fixedLengthField{name: "form_confirmation_status", start: 116, length: 1},
	fixedLengthField{name: "contractor_booklet_serial", start: 117, length: 8},
	fixedLengthField{name: "aud", start: 125, length: 1},
	fixedLengthField{name: "txt", start: 126, length: 1},
	fixedLengthField{name: "ide", start: 127, length: 1},
	fixedLengthField{name: "cas", start: 128, length: 1},
	fixedLengthField{name: "voc", start: 129, length: 1},
	fixedLengthField{name: "coh", start: 130, length: 1},
	fixedLengthField{name: "par", start: 131, length: 1},
	fixedLengthField{name: "sen", start: 132, length: 1},
	fixedLengthField{name: "pun", start: 133, length: 1},
	fixedLengthField{name: "spe", start: 134, length: 1},
	fixedLengthField{name: "name_check_first_name", start: 400, length: 25},
	fixedLengthField{name: "name_check_last_name", start: 425, length: 25},
	fixedLengthField{name: "contractor_internal_identifier", start: 450, length: 8},
	fixedLengthField{name: "reason_for_exemption", start: 458, length: 1},
	fixedLengthField{name: "reason_for_withholding", start: 459, length: 1},
	fixedLengthField{name: "forms_status", start: 460, length: 3},
	fixedLengthField{name: "booklet_image_batch", start: 464, length: 8},
	fixedLengthField{name: "booklet_image_serial", start: 472, length: 6},
	fixedLengthField{name: "special_provision", start: 501, length: 20},
	fixedLengthField{name: "special_provision_other_text", start: 521, length: 25},
}

func pearsonLineScan(s string) map[string]string {
	ret := make(map[string]string)
	for _, f := range pearsonStruct {
		ret[f.name] = s[f.start-1 : f.start-1+f.length]
	}
	return lib.RemoveBlanks(ret)
}

func pearson2sifSex(s string) string {
	switch s {
	case "M":
		return "1"
	case "F":
		return "2"
	default:
		return "9"
	}
}

func pearson2sifIndigenousStatus(s string) string {
	switch s {
	case "T":
		return "2"
	case "K":
		return "1"
	case "N":
		return "4"
	case "B":
		return "3"
	case "U":
		return "9"
	default:
		return "9"
	}
}

func pearson2sifParticipationCode(s string) string {
	switch s {
	case "P":
		return "P"
	case "C":
		// catch up
		return "P"
	case "A":
		return "A"
	case "W":
		return "W"
	case "E":
		return "E"
	case "H":
		// withheld
		return "S"
	case "L":
		return "X"
	default:
		return "P"
	}
}

func pearson2sifParticipationText(s string) string {
	switch s {
	case "P":
		return "Present"
	case "C":
		return "Present (Catch-Up)"
	case "A":
		return "Absent"
	case "W":
		return "Withdrawn"
	case "E":
		return "Exempt"
	case "H":
		// withheld
		return "Sanctioned Abandonment (Withheld)"
	case "L":
		return "No Longer Enrolled"
	default:
		return "Present"
	}
}

func pearson2sifExemptionReason(s string) string {
	switch s {
	case "1":
		return "Disability"
	case "2":
		return "Language"
	default:
		return ""
	}
}

func pearson2sifTestDisruption(s string) string {
	switch s {
	case "1":
		return "Illness"
	case "2":
		return "Parental removal"
	default:
		return ""
	}
}

func pearson2sifAdjustment(s string) []string {
	ret := make([]string, 0)
	for i := 0; i < len(s); i = i + 2 {
		s1 := s[i : i+2]
		switch s1 {
		case "ET":
			ret = append(ret, "ETA") // conventional choice against ETB or ETC
		case "RB":
			ret = append(ret, "RBK")
		case "SS":
			ret = append(ret, "SUP") // conventional choice of NAPLAN Support for Separate Supervision
		case "OS":
			ret = append(ret, "OSS")
		case "SC":
			ret = append(ret, "SCR")
		case "TY":
			ret = append(ret, "AST") // conventional choice of "assistive technology" for "typed response/attachment"
		case "CT":
			ret = append(ret, "AST")
		case "SP":
			ret = append(ret, "SUP")
		case "CO":
			ret = append(ret, "COL")
		case "SR":
			ret = append(ret, "AIA")
		case "OT":
			ret = append(ret, "AST") // conventional choice of Assistive Technology for Other  ---- make it pearson-other
		}
	}
	return ret
}

func wrapMessage(regr interface{}, i int, txid string, route string) *lib.NiasMessage {
	msg := &lib.NiasMessage{}
	msg.Body = regr
	msg.SeqNo = strconv.Itoa(i)
	msg.TxID = txid
	msg.MsgID = nuid.Next()
	// msg.Target = VALIDATION_PREFIX
	msg.Route = []string{route}
	return msg
}

func (di *DataIngest) ingestPearsonResultsFile(resultsFilePath string, wg *sync.WaitGroup) {

	// create a connection to the streaming server
	log.Println("Connecting to STAN server...")

	// map to hold student-school links temporarily
	// so student responses can be assigned to correct schools
	ss_link := make(map[string]string)

	// simple list of schools
	// schools := make([]SchoolDetails, 0)

	// open the data file for streaming read
	log.Printf("Opening results data file %s...", resultsFilePath)
	file, err := openDataFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results data file: ", err)
	}

	log.Println("Reading data file...")
	testRefId := uuid.NewV4().String()
	naptest := nxml.NAPTest{TestID: testRefId}
	//naptest.TestContent = nxml.TestContent{LocalId: r["booklet_id"], TestLevel: "3", TestDomain: "Writing"}
	naptest.TestContent = nxml.TestContent{TestLevel: "3", TestDomain: "Writing"}
	gt, err := di.ge.Encode(naptest)
	if err != nil {
		log.Println("Unable to gob-encode nap test: ", err)
	}
	di.sc.Publish("meta_yr3w", gt)

	// we will set up 1 testlet and 10 items, which are assessed in 1 rubric each
	testletRefId := uuid.NewV4().String()
	naptestlet := nxml.NAPTestlet{TestletID: testletRefId,
		NAPTestRefId: testRefId,
	}
	naptestlet.TestItemList.TestItem = make([]nxml.NAPTestlet_TestItem, 10)
	for i := range naptestlet.TestItemList.TestItem {
		naptestlet.TestItemList.TestItem[i] = nxml.NAPTestlet_TestItem{TestItemRefId: uuid.NewV4().String(),
			SequenceNumber: string(i)}
	}
	gtl, err := di.ge.Encode(naptestlet)
	if err != nil {
		log.Println("Unable to gob-encode nap testlet: ", err)
	}
	di.sc.Publish("meta_yr3w", gtl)

	rubrics := [10]string{"audience", "text structure", "ideas", "character and setting", "vocabulary",
		"cohesion", "paragraphing", "sentence structure", "punctuation", "spelling"}
	rubricsabbr := [10]string{"AUD", "TXT", "IDE", "CAS", "VOC", "COH", "PAR", "SEN", "PUN", "SPE"}
	maxscores := [10]string{"6", "4", "5", "4", "5", "4", "2", "6", "5", "6"}

	for i := range naptestlet.TestItemList.TestItem {
		naptestitem := nxml.NAPTestItem{ItemID: naptestlet.TestItemList.TestItem[i].TestItemRefId}
		naptestitem.TestItemContent = nxml.TestItemContent{ItemName: rubrics[i]}
		gti, err := di.ge.Encode(naptestlet)
		if err != nil {
			log.Println("Unable to gob-encode nap test item: ", err)
		}
		di.sc.Publish("meta_yr3w", gti)
	}

	//defer file.Close()
	totalStudents := 0
	i := 0
	//txid := nuid.Next()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := pearsonLineScan(scanner.Text())
		if r == nil {
			break
		}
		i = i + 1
		studentRefId := uuid.NewV4().String()
		regr := nxml.RegistrationRecord{BirthDate: fmt.Sprintf("%s%s-%s-%s", r["dob_year_1"], r["dob_year_2"], r["dob_dm"][0:2], r["dob_dm"][2:]),
			RefId:            studentRefId,
			TAAId:            r["booklet_id"],
			StateProvinceId:  r["vsn"],
			LocalId:          r["vsn"],
			FamilyName:       r["lname"],
			MainSchoolFlag:   "01",
			GivenName:        r["fname"],
			Sex:              pearson2sifSex(r["gender"]),
			IndigenousStatus: pearson2sifIndigenousStatus(r["gender"]),
			LBOTE:            r["lbote"],
			ClassGroup:       r["home_group"],
			SchoolLocalId:    r["school_vcaa_code"],
			ASLSchoolId:      r["asl_school_code"], // may be missing
		}
		// stopgap
		if r["asl_school_code"] != "" {
			r["asl_school_code"] = r["school_vcaa_code"]
		}
		regr.Unflatten()
		gsp, err := di.ge.Encode(regr)
		if err != nil {
			log.Println("Unable to gob-encode studentpersonal: ", err)
		}
		// store linkage locally
		ss_link[regr.RefId] = regr.ASLSchoolId
		di.sc.Publish("studentAndResults", gsp)
		totalStudents++
		event := nxml.NAPEvent{EventID: uuid.NewV4().String(),
			SPRefID:           studentRefId,
			SchoolID:          regr.ASLSchoolId,
			TestID:            testRefId,
			ParticipationCode: pearson2sifParticipationCode(r["student_attendance"]),
			ParticipationText: pearson2sifParticipationText(r["student_attendance"]),
			ExemptionReason:   pearson2sifExemptionReason(r["reason_for_exemption"]),
		}
		event.Adjustment.BookletType = r["special_provision_other_text"]
		if r["special_provision"] != "" {
			event.Adjustment.PNPCodelist = struct {
				PNPCode []string `xml:"PNPCode"`
			}{PNPCode: pearson2sifAdjustment(r["special_provision"])}
		}
		if r["reason_for_withholding"] != "" {
			event.TestDisruptionList.TestDisruption = make([]struct {
				Event string `xml:"Event"`
			}, 1)
			event.TestDisruptionList.TestDisruption[0].Event = pearson2sifTestDisruption(r["reason_for_withholding"])
		}
		ge, err := di.ge.Encode(event)
		if err != nil {
			log.Println("Unable to gob-encode nap event link: ", err)
		}
		di.sc.Publish("studentAndResults", ge)

		response := nxml.NAPResponseSet{ResponseID: uuid.NewV4().String(),
			StudentID: studentRefId,
			TestID:    testRefId,
		}
		response.TestletList.Testlet = make([]nxml.NAPResponseSet_Testlet, 1)
		response.TestletList.Testlet[0] = nxml.NAPResponseSet_Testlet{NapTestletRefId: testletRefId}
		response.TestletList.Testlet[0].ItemResponseList.ItemResponse = make([]nxml.NAPResponseSet_ItemResponse, 10)

		for i := range naptestlet.TestItemList.TestItem {
			response.TestletList.Testlet[0].ItemResponseList.ItemResponse[i] = nxml.NAPResponseSet_ItemResponse{ItemRefID: naptestlet.TestItemList.TestItem[i].TestItemRefId}

			response.TestletList.Testlet[0].ItemResponseList.ItemResponse[i].SubscoreList.Subscore = make([]nxml.NAPResponseSet_Subscore, 1)
		}
		for i := 0; i < 10; i++ {
			response = rubricPopulate(i, rubricsabbr[i], rubrics[i], maxscores[i], r, response)
		}
		gr, err := di.ge.Encode(response)
		if err != nil {
			log.Println("Unable to gob-encode student response set: ", err)
		}
		di.sc.Publish("studentAndResults", gr)

	}
	log.Println("Finished reading data file...")
	// post end of stream message to responses queue
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := di.ge.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}
	di.sc.Publish("studentAndResults", geot)
	di.sc.Publish("meta_yr3w", geot)

	//di.assignResponsesToSchools(ss_link)

	//log.Println("response assignment complete")

	log.Printf("ingestion complete for %s", resultsFilePath)

	wg.Done()

}

func rubricPopulate(seqPos int, abbr string, rubric string, maxscore string, r map[string]string, response nxml.NAPResponseSet) nxml.NAPResponseSet {
	if r[abbr] == "" {
		response.TestletList.Testlet[0].ItemResponseList.ItemResponse[seqPos].ResponseCorrectness = "NotAttempted"
	} else {
		if r[abbr] == maxscore {
			response.TestletList.Testlet[0].ItemResponseList.ItemResponse[seqPos].ResponseCorrectness = "Correct"
		} else {
			response.TestletList.Testlet[0].ItemResponseList.ItemResponse[seqPos].ResponseCorrectness = "Incorrect"
		}
		response.TestletList.Testlet[0].ItemResponseList.ItemResponse[0].SubscoreList.Subscore[seqPos].SubscoreType = rubric
		response.TestletList.Testlet[0].ItemResponseList.ItemResponse[0].SubscoreList.Subscore[seqPos].SubscoreValue = r[abbr]
	}
	return response
}
