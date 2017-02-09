// reports.go
package main

import (
	"github.com/nsip/nias2/naprr"
	"log"
	"sync"
)

var sr = naprr.NewStreamReader()
var rg = naprr.NewReportGenerator()

func main() {

	var wg sync.WaitGroup

	schools := sr.GetSchoolDetails()
	nd := sr.GetNAPLANData()

	for _, subslice := range schools {
		for _, school := range subslice {
			wg.Add(1)
			go createSchoolReports(nd, school.ACARAId, &wg)
		}
	}

	wg.Add(1)
	go createTestReports(nd, &wg)

	// block until all reports generated
	wg.Wait()
	log.Println("All reports generated")

}

// generate school-level data reports
func createSchoolReports(nd *naprr.NAPLANData, acaraid string, wg *sync.WaitGroup) {
	sd := sr.GetSchoolData(acaraid)
	rg.GenerateParticipationData(nd, sd)
	log.Println("Participation data created for: ", acaraid)
	rg.GenerateSchoolScoreSummaryData(nd, sd)
	log.Println("Score summary data created for: ", acaraid)
	rg.GenerateDomainScoreData(nd, sd)
	log.Println("Domain scores data created for: ", acaraid)

	wg.Done()
}

// generate test-level reports
func createTestReports(nd *naprr.NAPLANData, wg *sync.WaitGroup) {
	rg.GenerateCodeframeData(nd)
	log.Println("Codeframe data created.")
	wg.Done()
}
