package main

// Alternative naprr server offering graphql data feeds and KV store backend
// to try and elimnate issues with tcp, memory and disk access
// on Win32 that hamper all other deployments so far.

import (
"bufio"
"sort"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/nsip/nias2/naprrql"
	"github.com/nsip/nias2/version"

"github.com/kylelemons/godebug/diff"
)

var ingest = flag.Bool("ingest", false, "Loads data from results file. Exisitng data is overwritten.")
var report = flag.Bool("report", false, "Creates .csv reports. Existing reports are overwritten")

// var isrprint = flag.Bool("isrprint", false, "Creates .csv files for use in isr printing")
var itemprint = flag.Bool("itemprint", false, "Creates .csv files reporting item results for each student against items")

var writingextract = flag.Bool("writingextract", false, "Creates .csv file extract of all writing items, for input into marking systems")
var psiexceptions = flag.String("psiexceptions", "-", "File containing list of PSIs to ignore in generating writing extract")
var psiwhitelist = flag.String("psiwhitelist", "-", "File containing list of PSIs to use in generating writing extract")
var qa = flag.Bool("qa", false, "Creates .csv files for QA checking of NAPLAN results")
var vers = flag.Bool("version", false, "Reports version of NIAS distribution")
var xml = flag.Bool("xml", false, "Reexports redacted xml of RRD dataset")

var csvdiff = flag.Bool("csvdiff", false, "Compare two directories full of CSV files recursively")

//
// deprecated as at jan 2019 - all users want wordcount, removing flag avoids
// conditional logic in nap-writing-print utility
//
// var wordcount = flag.Bool("wordcount", false, "Include wordcount in writingextract report")

var pnpadd = flag.String("pnpadd", "", "Adds PNP codes from CSV file into RRD dataset")

func main() {

	// runtime.GOMAXPROCS(16) // optional performance improvement for larger systems.

	flag.Parse()

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		closeDB()
		os.Exit(1)
	}()

	if *vers {
		fmt.Printf("NIAS: Version %s\n", version.TagName)
		os.Exit(1)
	}

	// ingest results data, rebuild reports, and exit to save memory
	if *ingest {
		ingestData()
		if !naprrql.DataUnfit() {
			// startWebServer(true)
			// writeReports()
		} else {
			log.Println(`
				--
				Data is unfit for reporting,
				please correct before running further naprr services
				--
				`)
		}
		// shut down
		compactDB()
		closeDB()
		os.Exit(1)
	}

	// for all other operations check if the database has data
	// before proceeding
	if naprrql.DatabaseIsEmpty() {
		log.Println("\n\nNo data in datastore...\n\nPlease run naprrql -ingest to load results data\n\n")
		closeDB()
		os.Exit(1)
	}

        // create the csv reports
        if *csvdiff {
		dir := flag.Args()
		if len(dir) != 2 {
			log.Fatalln("Must specify two directory names to csvdiff: naprrql -csvdiff dir1 dir2\n")
		}
		csvdiff_comparison_recursive(dir[0], dir[1])
                os.Exit(1)
        }

	// create the csv reports
	if *report {
		// launch web-server
		startWebServer(true)
		writeReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	// create the writing item report
	if *writingextract {
		// launch web-server
		if *psiexceptions != "-" && *psiwhitelist != "-" {
			log.Fatalln("\nCannot specify both psiexceptions and psiwhitelist, they are mutually exclusive!")
		}
		startWebServer(true)
		filename := ""
		blacklist := true
		if *psiexceptions != "-" {
			filename = *psiexceptions
		} else if *psiwhitelist != "-" {
			filename = *psiwhitelist
			blacklist = false
		}
		writeWritingExtractReports(filename, blacklist)
		// shut down
		closeDB()
		os.Exit(1)
	}

	// create the writing item report
	if *xml {
		// launch web-server
		startWebServer(true)
		writeXMLReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	// add PNP codes into ingested RRD
	if len(*pnpadd) > 0 {
		// launch web-server
		startWebServer(true)
		addEventCSV(*pnpadd)
		writeXMLReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	/*
		// create the isr printing reports
		if *isrprint {
		 launch web-server
			startWebServer(true)
			writeISRPrintingReports()
		// shut down
			closeDB()
			os.Exit(1)
		}
	*/

	// create the item reports
	if *itemprint {
		// launch web-server
		startWebServer(true)
		writeItemPrintingReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	// create the QA reports
	if *qa {
		// launch web-server
		startWebServer(true)
		writeQAPrintingReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	// otherwise just start the webserver
	startWebServer(false)

	// wait for shutdown
	for {
		runtime.Gosched()
	}

}

//
// iterate & load any r/r data files provided
//
func ingestData() {
	// ingest the data
	log.Println("invoking data ingest...")
	clearDBWorkingDirectory()
	resultsFiles := parseResultsFileDirectory()
	for _, resultsFile := range resultsFiles {
		naprrql.IngestResultsFile(resultsFile)
	}
}

//
// launch the webserver
//
func startWebServer(silent bool) {
	go naprrql.RunQLServer()
	if !silent {
		fmt.Printf("\n\nBrowse to following locations:\n\n")
		fmt.Printf("\n\tFor qa reporting user interface:\n\n\t\thttp://localhost:1329/ui\n")
		fmt.Printf("\n\tFor data query explorer:\n\n\t\thttp://localhost:1329/sifql\n")
		fmt.Printf("\n\tFor data model viewer:\n\n\t\thttp://localhost:1329/datamodel\n\n ")
	}

}

func csvdiff_comparison_recursive(dir1, dir2 string){
if _, err := os.Stat(dir1); os.IsNotExist(err) {
	                        log.Fatalln("Directory " + dir1 + " does not exist!\n")
}
if _, err := os.Stat(dir2); os.IsNotExist(err) {
	                        log.Fatalln("Directory " + dir2 + " does not exist!\n")
}
if _, err := os.Stat("compare"); os.IsNotExist(err) {
err = os.Mkdir("compare", os.ModePerm)
        if !os.IsExist(err) && err != nil {
                log.Fatalln("Error trying to create comparison directory: ", err)
        }
}
f, err:=os.Create("compare/comparison" + time.Now().Format("20060102150405") + ".txt")
defer f.Close()



err = filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
path2:=dir2 + strings.TrimPrefix(path, dir1)
if _, err := os.Stat(path2); os.IsNotExist(err) {
log.Printf("File %s does not exist, to compare with file %s\n", path2, path)
return nil
}
csvdiff_comparison(path, path2, f)		
}
		return nil
})
if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir1, err)
		return
	}

}

func csvdiff_comparison(file1, file2 string, rept *os.File) {
	csv1 := readlines(file1)
	csv2 := readlines(file2)
	sort.Strings(csv1)
	sort.Strings(csv2)
	difference := strings.Split(diff.Diff(strings.TrimSpace(strings.Join(csv1, "\n")), strings.TrimSpace(strings.Join(csv2, "\n"))), "\n")
	if(len(difference) > 1) {
		rept.WriteString("Comparison, " + file1 + ", " + file2 + "\n\n")
		for _, line := range difference {
			if !strings.HasPrefix(line, " ") {
				rept.WriteString(line + "\n")
			}
		}
				rept.WriteString("\n")
	}
}

func readlines(filename string) []string {
file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	ret := make([]string, 0)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ret = append(ret, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	return ret
}


//
// create .csv reports
//
func writeReports() {
	clearReportsDirectory()
	log.Println("generating reports...")
	naprrql.GenerateReports()
	log.Println("reports generated...")
}

//
// create isr printing reports
//
func writeISRPrintingReports() {
	clearISRPrintingDirectory()
	log.Println("generating isr printing reports...")
	naprrql.GenerateISRPrintReports()
	log.Println("isr printing reports generated...")
}

//
// create item printing reports
//
func writeItemPrintingReports() {
	// clearReportsDirectory() // - check this
	log.Println("generating item printing reports...")
	naprrql.GenerateItemPrintReports()
	log.Println("item printing reports generated...")
}

//
// create writing extract printing reports
//
func writeWritingExtractReports(psi_exceptions string, blacklist bool) {
	log.Println("generating Writing item extract reports...")
	naprrql.GenerateWritingExtractReports(psi_exceptions, blacklist)
	log.Println("Writing item extract reports generated...")
}

//
// create XML printing reports
//
func writeXMLReports() {
	log.Println("generating XML reports...")
	naprrql.GenerateXMLReports()
	log.Println("Writing XML generated...")
}

//
// create XML printing reports
//
func addEventCSV(filename string) {
	log.Println("Adding PNP values to dataset...")
	naprrql.AddEventCSV(filename)
	log.Println("Added PNP values to dataset...")
}

// create QA reports
//
func writeQAPrintingReports() {
	clearQADirectory()
	log.Println("generating QA reports...")
	naprrql.GenerateQAReports()
	log.Println("QA reports generated...")
}

//
// look for results data files
//
func parseResultsFileDirectory() []string {

	files := make([]string, 0)

	zipFiles, _ := filepath.Glob("./in/*.zip")
	xmlFiles, _ := filepath.Glob("./in/*.xml")

	files = append(files, zipFiles...)
	files = append(files, xmlFiles...)
	if len(files) == 0 {
		log.Fatalln("No results data *.zip *.xml.zip or *.xml files found in input folder /in.")
	}

	return files

}

//
// ensure clean shutdown of data store
//
func closeDB() {
	log.Println("Closing datastore...")
	naprrql.GetDB().Close()
	log.Println("Datastore closed.")
}

//
// run compaction to minimise database size
//
func compactDB() {
	log.Println("Compacting datastore...")
	naprrql.CompactDatastore()
	log.Println("Datastore compaction completed.")
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
	err := os.RemoveAll("out/school_reports")
	err = os.RemoveAll("out/system_reports")
	if err != nil {
		log.Println("Error trying to reset reports directory: ", err)
	}
	createReportsDirectory()
}

//
// remove the isr printing report directory
//
func clearISRPrintingDirectory() {
	// remove existing logs and recreate the directory
	err := os.RemoveAll("out/isr_printing")
	if err != nil {
		log.Println("Error trying to reset isr printing reports directory: ", err)
	}
	createISRPrintingDirectory()
}

// remove the isr printing report directory
//
func clearQADirectory() {
	// remove existing logs and recreate the directory
	err := os.RemoveAll("out/qa")
	if err != nil {
		log.Println("Error trying to reset QA reports directory: ", err)
	}
	createQADirectory()
}

//
// create folder for .csv reports
//
func createReportsDirectory() {
	err := os.Mkdir("out", os.ModePerm)
	err = os.Mkdir("out/school_reports", os.ModePerm)
	err = os.Mkdir("out/system_reports", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create reports directory: ", err)
	}

}

//
// create the folder for isr printing reports
//
func createISRPrintingDirectory() {
	err := os.Mkdir("out", os.ModePerm)
	err = os.Mkdir("out/isr_printing", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create isr printing reports directory: ", err)
	}

}

//
// create the folder for QA reports
//
func createQADirectory() {
	err := os.Mkdir("out", os.ModePerm)
	err = os.Mkdir("out/qa", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create QA reports directory: ", err)
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
