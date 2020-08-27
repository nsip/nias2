package napcomp

// napcomp takes input file from student registration
// and output file from results reporting
// and does a comparison to find students not accounted for
// in both files.

import (
	gocsv "encoding/csv"
	goxml "encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ms "github.com/mitchellh/mapstructure"
	"github.com/wildducktheories/go-csv"
	"gopkg.in/fatih/set.v0"

	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/naprrql"
	"github.com/nsip/nias2/xml"
	"github.com/syndtr/goleveldb/leveldb"
)

var resultsKeys, registrationKeys *set.Set

//
// iterate & load any r/r data files and
// registration data files provided
//
func IngestData() {
	// ingest the data
	log.Println("invoking data ingest...")

	clearDBWorkingDirectory()

	registrationKeys = set.New(set.ThreadSafe).(*set.Set)
	resultsKeys = set.New(set.ThreadSafe).(*set.Set)

	log.Println("reading results data files...")
	resultsFiles := parseResultsFileDirectory()
	for _, resultsFile := range resultsFiles {
		ingestResultsFile(resultsFile)
	}

	log.Println("reading registration data files...")
	registrationFiles := parseRegistrationFileDirectory()
	for _, regFile := range registrationFiles {
		ingestRegistrationFile(regFile)
	}

}

//
// reads in student information from a results reporting dataset file
//
func ingestResultsFile(resultsFilePath string) {

	naprrql.SetDBReadWrite()
	db := naprrql.GetDB()
	ge := naprrql.GobEncoder{}

	// open the data file for streaming read
	xmlFile, err := naprrql.OpenResultsFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results file")
	}

	log.Printf("Reading data file [%s]", resultsFilePath)

	batch := new(leveldb.Batch)

	decoder := goxml.NewDecoder(xmlFile)
	totalStudents := 0
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case goxml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...handle by type
			switch inElement {
			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				sp.Flatten() // align structure to registration record, explicit other-ids
				gsp, err := ge.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}

				key := makeComparisonKey(&sp)
				// log.Println("Result Key: ", key)

				// store object in db
				// {StudentPersonal-id} = object
				batch.Put([]byte("res:"+key), gsp)
				// keep the key for comparisons
				resultsKeys.Add(key)

				totalStudents++

			}
		default:
		}
	}

	// commit database entries
	batcherr := db.Write(batch, nil)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
	}

	log.Println("Data file read complete...")
	log.Printf("No. report data students found: %d \n", totalStudents)
	log.Printf("ingestion complete for [%s]", resultsFilePath)

}

//
// reads in student information from a registraiton data file
//
func ingestRegistrationFile(regFilePath string) {

	db := naprrql.GetDB()
	ge := naprrql.GobEncoder{}

	log.Printf("Reading data file [%s]", regFilePath)

	batch := new(leveldb.Batch)

	regFile, err := os.Open(regFilePath)
	if err != nil {
		log.Fatalln("Unable to open: ", regFilePath)
	}
	defer regFile.Close()
	reader := csv.WithIoReader(regFile)
	defer reader.Close()

	totalStudents := 0
	for record := range reader.C() {
		regr := &xml.RegistrationRecord{}
		r := lib.RemoveBlanks(record.AsMap())
		decode_err := ms.Decode(r, regr)
		regr.Unflatten() // make equivalent to xml record for other-ids
		if decode_err != nil {
			continue
		}
		gsp, err := ge.Encode(regr)
		if err != nil {
			log.Println("Unable to gob-encode studentpersonal: ", err)
		}
		key := makeComparisonKey(regr)
		// log.Println("Reg Key: ", key)

		// store object in db
		// {StudentPersonal-id} = object
		batch.Put([]byte("reg:"+key), gsp)
		// keep the key for comparisons
		registrationKeys.Add(key)

		totalStudents++
	}

	// commit database entries
	batcherr := db.Write(batch, nil)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
	}

	log.Println("Data file read complete...")
	log.Printf("No. registration data students found: %d \n", totalStudents)
	log.Printf("ingestion complete for [%s]", regFilePath)

}

//
// create key to use in record comparisons from selected
// data fields
//
func makeComparisonKey(r *xml.RegistrationRecord) string {

	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s",
		r.FamilyName,
		r.GivenName,
		r.MiddleName,
		// r.PreferredName,
		r.LocalId,
		// r.StateProvinceId,
		// r.DiocesanId,
		// r.NationalId,
		r.PlatformId,
		r.ASLSchoolId,
		r.BirthDate,
		// r.SchoolLocalId,
	)

	return key
}

//
// create .csv reports
//
func WriteReports() {

	clearReportsDirectory()
	db := naprrql.GetDB()
	ge := naprrql.GobEncoder{}

	log.Println("generating difference reports...")

	log.Println()
	regbutnotres := set.Difference(registrationKeys, resultsKeys)
	log.Printf("Students registred  but not in results: %d", regbutnotres.Size())

	resbutnotreg := set.Difference(resultsKeys, registrationKeys)
	log.Printf("Students in results but not in registration: %d", resbutnotreg.Size())

	symdiff := set.SymmetricDifference(resultsKeys, registrationKeys)
	log.Printf("Total students not in both files: %d", symdiff.Size())
	log.Println()

	// registered not in results...
	f, err := os.Create("./out/RegisteredButNotInResults.csv")
	if err != nil {
		log.Fatalln("Cannot open file to publish report: ", err)
	}
	defer f.Close()
	w := gocsv.NewWriter(f)
	// header
	hdr := xml.RegistrationRecord{}
	w.Write(hdr.GetHeaders())
	for _, key := range set.StringSlice(regbutnotres) {
		k := "reg:" + key
		val, err := db.Get([]byte(k), nil)
		if err != nil {
			log.Println("Cannot retrieve object for id: ", k)
		}
		var rrObj interface{}
		err = ge.Decode(val, &rrObj)
		if err != nil {
			log.Println("Cannot gob-decode object.")
		}
		sp, ok := rrObj.(xml.RegistrationRecord)
		if !ok {
			log.Println("Cannot assert object as Student Personal.")
		}
		w.Write(sp.GetSlice())
	}
	w.Flush()

	// results but not registered...
	f, err = os.Create("./out/ResultsButNotInRegister.csv")
	if err != nil {
		log.Fatalln("Cannot open file to publish report: ", err)
	}
	defer f.Close()
	w = gocsv.NewWriter(f)
	// header
	w.Write(hdr.GetHeaders())
	for _, key := range set.StringSlice(resbutnotreg) {
		k := "res:" + key
		val, err := db.Get([]byte(k), nil)
		if err != nil {
			log.Println("Cannot retrieve object for id: ", k)
		}
		var rrObj interface{}
		err = ge.Decode(val, &rrObj)
		if err != nil {
			log.Println("Cannot gob-decode object.")
		}
		sp, ok := rrObj.(xml.RegistrationRecord)
		if !ok {
			log.Println("Cannot assert object as Student Personal.")
		}
		w.Write(sp.GetSlice())
	}
	w.Flush()

	log.Println("reports generated to /out folder.")
}

//
// look for results data files
//
func parseResultsFileDirectory() []string {

	files := make([]string, 0)

	zipFiles, err := filepath.Glob("./in/results/*.zip")
	xmlFiles, err := filepath.Glob("./in/results/*.xml")

	files = append(files, zipFiles...)
	files = append(files, xmlFiles...)
	if len(files) == 0 {
		log.Fatalln("No results data *.zip *.xml.zip or *.xml files found in input folder /in/results.", err)
	}

	return files

}

//
// look for registration data files
//
func parseRegistrationFileDirectory() []string {

	files := make([]string, 0)

	csvFiles, err := filepath.Glob("./in/registration/*.csv")

	files = append(files, csvFiles...)
	if len(files) == 0 {
		log.Fatalln("No registration data *.csv files found in input folder /in/registration.", err)
	}

	return files

}

//
// ensure clean shutdown of data store
//
func CloseDB() {
	log.Println("Closing datastore...")
	naprrql.GetDB().Close()
	log.Println("Datastore closed.")
}

//
// remove working files of datastore
//
func clearDBWorkingDirectory() {

	// remove existing logs and recreate the directory
	err := os.RemoveAll("kvs")
	if err != nil {
		log.Println("Error trying to reset datastore working directory: ", err)
	}
	createDBWorkingDirectory()
}

//
// remove reports working directory
//
func clearReportsDirectory() {
	// remove existing logs and recreate the directory
	err := os.RemoveAll("out")
	if err != nil {
		log.Println("Error trying to reset reports directory: ", err)
	}
	createReportsDirectory()

}

//
// create folder for .csv reports
//
func createReportsDirectory() {
	err := os.Mkdir("out", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create reports directory: ", err)
	}

}

//
// create folder for datastore
//
func createDBWorkingDirectory() {
	err := os.Mkdir("kvs", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create datastore working directory: ", err)
	}

}
