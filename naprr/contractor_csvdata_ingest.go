package naprr

import (
	"bufio"
	ms "github.com/mitchellh/mapstructure"
	"github.com/nats-io/nuid"
	"strconv"
	//"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	nxml "github.com/nsip/nias2/xml"
	"log"
	"path/filepath"
	"sync"
)

/*
type DataIngest struct {
	sc stan.Conn
	ge GobEncoder
	sr *StreamReader
}

func NewDataIngest() *DataIngest {
	di := DataIngest{sc: CreateSTANConnection(), ge: GobEncoder{}, sr: NewStreamReader()}
	return &di
}
*/

func (di *DataIngest) RunCSV() {

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
	fixedLengthField{name: "second_name", start: 66, length: 15},
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
	log.Println("%v\n", ret)
	return lib.RemoveBlanks(ret)
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

	//defer file.Close()
	i := 0
	txid := nuid.Next()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := pearsonLineScan(scanner.Text())
		i = i + 1

		regr := nxml.RegistrationRecord{}
		decode_err := ms.Decode(r, &regr)
		if decode_err != nil {
			log.Fatalln("unable to open results data file: ", decode_err)
		}

		msg := &lib.NiasMessage{}
		msg.Body = regr
		msg.SeqNo = strconv.Itoa(i)
		msg.TxID = txid
		msg.MsgID = nuid.Next()
		// msg.Target = VALIDATION_PREFIX
		msg.Route = []string{"pearson"}

		//publish(msg)

	}
	// post end of stream message to responses queue
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := di.ge.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}
	di.sc.Publish("responses", geot)
	di.sc.Publish("meta", geot)

	di.assignResponsesToSchools(ss_link)

	log.Println("response assignment complete")

	log.Printf("ingestion complete for %s", resultsFilePath)

	wg.Done()

}
