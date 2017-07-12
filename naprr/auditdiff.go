package naprr

import (
	"log"

	"github.com/nsip/nias2/lib"
	nxml "github.com/nsip/nias2/xml"
	"github.com/twinj/uuid"
	//"sync"
	"github.com/RoaringBitmap/roaring"
	ms "github.com/mitchellh/mapstructure"
	"github.com/wildducktheories/go-csv"
	"os"
	"path/filepath"
	"strconv"
)

func (di *DataIngest) RunRegRecords() {
	uuid.Init()
	di.PSIBitmap = roaring.NewBitmap()
	regFiles := parseRegRecordsDirectory()
	for _, file := range regFiles {
		di.ingestRegCSVFile(file)
	}
	if len(regFiles) == 0 {
		// send out signal that there is nothing in the Yr3W queues
		eot := lib.TxStatusUpdate{TxComplete: true}
		geot, err := di.ge.Encode(eot)
		if err != nil {
			log.Println("Unable to gob-encode tx complete message: ", err)
		}
		di.sc.Publish(REGISTRATION_STUDENT_RECORDS, geot)
	}

	di.finaliseTransactions()

	log.Printf("NAPLAN Registration records ingestion complete\n")
}

func parseRegRecordsDirectory() []string {
	files, _ := filepath.Glob("./in/Reg/*.csv")
	if len(files) == 0 {
		log.Println("No writing data files found in Registration record input folder.")
	}
	return files
}

// encode PSI as a 32 bit integer. PSI is 9 digits plus system letter + checksum; it just fits, uint32 having 10 digits
func Psi2uint32(psi string) uint32 {
	if len(psi) < 2 {
		return 0
	}
	num := psi[1 : len(psi)-1]
	ret, err := strconv.Atoi(num)
	// last bit is whether prefix of psi is 'R' or 'D'
	if psi[0] == 'R' {
		ret += 1
	}
	if err != nil {
		ret = 0
	}
	return uint32(ret)
}

func (di *DataIngest) ingestRegCSVFile(resultsFilePath string) {
	log.Printf("Opening results data file for registration records: %s...", resultsFilePath)
	src, err := os.Open(resultsFilePath)
	if err != nil {
		log.Println("Unable to open: ", resultsFilePath)
	}
	defer src.Close()
	reader := csv.WithIoReader(src)
	defer reader.Close()

	totalStudents := 0
	i := 0
	//txid := nuid.Next()
	for record := range reader.C() {
		i = i + 1
		regr := &nxml.RegistrationRecord{}
		r := lib.RemoveBlanks(record.AsMap())
		decode_err := ms.Decode(r, regr)
		regr.Unflatten()
		if decode_err != nil {
			continue
		}
		gsp, err := di.ge.Encode(regr)
		if err != nil {
			log.Println("Unable to gob-encode studentpersonal: ", err)
		}
		di.sc.Publish(REGISTRATION_STUDENT_RECORDS, gsp)
		totalStudents++
	}
	log.Println("Finished reading data file...")
	// post end of stream message to responses queue
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := di.ge.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}
	di.sc.Publish(REGISTRATION_STUDENT_RECORDS, geot)

}
