// system-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"gopkg.in/fatih/set.v0"
)

//
// printing pipeline for system-level reports
//

//
// helper type to capture query parameters
//
type systemQueryParams struct {
	schoolAcaraID string
}

//
// create and run a system-level report printing pipeline.
//
// Pipeline streams requests a given report, school by school
// feeding results to the ouput csv file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//
func runSystemReportPipeline(queryFileName string, query string, schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	// check if system query needs to iterate schools, if not
	// pass single dummy school to fire the query
	if !strings.Contains(query, "$acaraIDs") {
		schools = []string{"no-op"}
	}
	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/system_reports"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := deriveCSVFileName(queryFileName)
	outFileName := outFileDir + "/" + csvFileName
	mapFileName := deriveMapFileName(queryFileName)
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("System report file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

// Run a single query for each school, then run a series of reports off that query
func runQASystemSingleQueryReportPipeline(query string, queryFileNames []string, schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	// check if system query needs to iterate schools, if not
	// pass single dummy school to fire the query
	if !strings.Contains(query, "$acaraIDs") {
		schools = []string{"no-op"}
	}
	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	jsonc1 := make([]chan gjson.Result, len(queryFileNames))
	for i, _ := range jsonc1 {
		jsonc1[i] = make(chan gjson.Result, 0)
	}
	go func() {
		for i := range jsonc {
			for _, c := range jsonc1 {
				c <- i
			}
		}
		for _, c := range jsonc1 {
			close(c)
		}
	}()
	for i, queryFileName := range queryFileNames {
		csvFileName := deriveCSVFileName(queryFileName)
		outFileName := outFileDir + "/" + csvFileName
		mapFileName := deriveMapFileName(queryFileName)
		errc1, err := csvFileSink(ctx, outFileName, mapFileName, jsonc1[i])
		if err != nil {
			return err
		}
		errcList = append(errcList, errc1)

		log.Println("QA file writing... " + outFileName)
	}
	return WaitForPipeline(errcList...)
}

// Run the query once with the argument schools
func runQASystemReportAllSchoolsPipeline(queryFileName string, query string, schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	// check if system query needs to iterate schools, if not
	// pass single dummy school to fire the query
	if !strings.Contains(query, "$acaraIDs") {
		schools = []string{"no-op"}
	}
	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutorNoBreakdown(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := deriveCSVFileName(queryFileName)
	outFileName := outFileDir + "/" + csvFileName
	mapFileName := deriveMapFileName(queryFileName)
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("QA file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

func runQAErdsReportPipeline(querymap map[string]string, schools []string) error {

	erds_query := `query NAPDomainScores($acaraIDs: [String]) {
  domain_scores_event_report_by_school(acaraIDs: $acaraIDs) {
    Student {
      FamilyName
      GivenName
      BirthDate
      OtherIdList {
        OtherId {
          Type
          Value
        }
      }
    }
    Test {
      TestContent {
        TestLevel
        TestDomain
      }
    }
    Event {
       EventID
       ParticipationCode
       TestDisruptionList {
         TestDisruption {
           Event
         }
       }
    }
    Response {
       ResponseID
       PathTakenForDomain
       ParallelTest
       DomainScore {
         RawScore
       }
    }
    SchoolDetails {
      ACARAId
      SchoolName
    }
  }
}`

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, erds_query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	jsonc1 := make([]chan gjson.Result, len(querymap))
	for i, _ := range jsonc1 {
		jsonc1[i] = make(chan gjson.Result, 0)
	}
	go func() {
		for i := range jsonc {
			for _, c := range jsonc1 {
				c <- i
			}
		}
		for _, c := range jsonc1 {
			close(c)
		}
	}()

	i := 0
	for queryFileName, _ := range querymap {
		var jsonc2 <-chan gjson.Result
		// for now, I'm merely hardcoding the query names
		// These are transforms on the CSV
		if strings.HasSuffix(queryFileName, "systemTestAttempts.gql") {
			jsonc2, errc, _ = systemTestAttempts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemTestTypeImpacts.gql") {
			jsonc2, errc, _ = systemTestTypeImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemTestIncidents.gql") {
			jsonc2, errc, _ = systemTestIncidents(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemParticipationCodeImpacts.gql") {
			jsonc2, errc, _ = systemParticipationCodeImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemTestCompleteness.gql") {
			jsonc2, errc, _ = systemTestCompleteness(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemObjectFrequency.gql") {
			jsonc2, errc, _ = systemObjectFrequency(jsonc1[i])
			errcList = append(errcList, errc)
		} else {
			jsonc2 = nil
		}
		if jsonc2 != nil {
			csvFileName := deriveCSVFileName(queryFileName)
			outFileName := outFileDir + "/" + csvFileName
			mapFileName := deriveMapFileName(queryFileName)
			errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc2)
			if err != nil {
				return err
			}
			errcList = append(errcList, errc)

			log.Println("ERDS QA file writing... " + outFileName)
		}
		i++
	}
	return WaitForPipeline(errcList...)
}

func runQAItemRespReportPipeline(querymap map[string]string, schools []string) error {

	itemresp_query := `query NAPItemResults($acaraIDs: [String]) {
  item_results_report_by_school(acaraIDs: $acaraIDs) {
    Test {
      TestContent {
        LocalId
        TestName
        TestLevel
        TestDomain
        TestYear
        StagesCount
        TestType
      }
    }
    TestItem {
      ItemID
      TestItemContent {
        NAPTestItemLocalId
        ItemName
        ItemType
        Subdomain
        WritingGenre
        ItemSubstitutedForList {
          SubstituteItem {
            SubstituteItemRefId
            LocalId
            PNPCodeList {
              PNPCode
            }
          }
        }
      }
    }
    Student {
      BirthDate
      Sex
      IndigenousStatus
      LBOTE
      YearLevel
      ASLSchoolId
      OtherIdList {
        OtherId {
          Type
          Value
        }
      }
    }
    ParticipationCode
    Response {
      PathTakenForDomain
      ParallelTest
      PSI
      TestletList {
        Testlet {
          NapTestletLocalId
          TestletScore
          ItemResponseList {
            ItemResponse {
              LocalID
              Response
              ResponseCorrectness
              Score
              LapsedTimeItem
              SequenceNumber
              ItemWeight
              SubscoreList {
                Subscore {
                  SubscoreType
                  SubscoreValue
                }
              }
            }
          }
        }
      }
    }
  }
}`
	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, itemresp_query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	jsonc1 := make([]chan gjson.Result, len(querymap))
	for i, _ := range jsonc1 {
		jsonc1[i] = make(chan gjson.Result, 0)
	}
	go func() {
		for i := range jsonc {
			for _, c := range jsonc1 {
				c <- i
			}
		}
		for _, c := range jsonc1 {
			close(c)
		}
	}()

	i := 0
	for queryFileName, _ := range querymap {
		var jsonc2 <-chan gjson.Result
		// for now, I'm merely hardcoding the query names
		// These are transforms on the CSV
		if strings.HasSuffix(queryFileName, "systemTestTypeItemImpacts.gql") {
			jsonc2, errc, _ = systemTestTypeItemImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemParticipationCodeItemImpacts.gql") {
			jsonc2, errc, _ = systemParticipationCodeItemImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if strings.HasSuffix(queryFileName, "systemItemCounts.gql") {
			jsonc2, errc, _ = systemItemCounts(jsonc1[i])
			errcList = append(errcList, errc)
		} else {
			jsonc2 = nil
		}
		if jsonc2 != nil {
			csvFileName := deriveCSVFileName(queryFileName)
			outFileName := outFileDir + "/" + csvFileName
			mapFileName := deriveMapFileName(queryFileName)
			errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc2)
			if err != nil {
				return err
			}
			errcList = append(errcList, errc)

			log.Println("ITEM RESP QA report file writing... " + outFileName)
		}
		i++
	}
	return WaitForPipeline(errcList...)
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn, for the given year level.
//
func systemParametersSource(ctx context.Context, schools ...string) (<-chan systemQueryParams, <-chan error, error) {

	{ //check input variables, handle errors before goroutine starts
		if len(schools) == 0 {
			return nil, nil, errors.Errorf("no schools provided")
		}
	}

	out := make(chan systemQueryParams)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for schoolIndex, school := range schools {
			if school == "" {
				// Handle an error that occurs during the goroutine.
				errc <- errors.Errorf("school %v is empty string", schoolIndex+1)
				return
			}
			vars := systemQueryParams{schoolAcaraID: school}
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- vars:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

//
// query executor transform stage takes query params in, executes gql query
// and writes results to output channel
//
func systemQueryExecutor(ctx context.Context, query, url string, in <-chan systemQueryParams) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for params := range in {
			vars := map[string]interface{}{"acaraIDs": []string{params.schoolAcaraID}}
			json, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			for _, result := range json.Array() {
				// Send the data to the output channel but return early
				// if the context has been cancelled.
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
			// handle edge cases where data is not returned as an array
			if len(json.Array()) == 0 {
				select {
				case out <- json:
				case <-ctx.Done():
					return
				}
			}

		}
	}()
	return out, errc, nil
}

// As above, but do not break up acaraID arguments
func systemQueryExecutorNoBreakdown(ctx context.Context, query, url string, in <-chan systemQueryParams) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		acaraIDs := make([]string, 0)
		for params := range in {
			acaraIDs = append(acaraIDs, params.schoolAcaraID)
		}
		vars := map[string]interface{}{"acaraIDs": acaraIDs}
		json, err := gql.DoQuery(url, query, vars)
		if err != nil {
			// Handle an error that occurs during the goroutine.
			errc <- err
			return
		}
		for _, result := range json.Array() {
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- result:
			case <-ctx.Done():
				return
			}
		}
		// handle edge cases where data is not returned as an array
		if len(json.Array()) == 0 {
			select {
			case out <- json:
			case <-ctx.Done():
				return
			}
		}

	}()
	return out, errc, nil
}

func systemTestAttempts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			if record.Get("Event.ParticipationCode").String() == "S" {
				out <- record
			}
		}
	}()
	return out, errc, nil
}

func systemTestTypeImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			if record.Get("Test.TestContent.TestDomain").String() == "Writing" &&
				(len(record.Get("Response.PathTakenForDomain").String()) > 0 ||
					len(record.Get("Response.ParallelTest").String()) > 0) {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Writing test with adaptive structure"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if record.Get("Test.TestContent.TestDomain").String() != "Writing" &&
				(record.Get("Event.ParticipationCode").String() == "P" ||
					record.Get("Event.ParticipationCode").String() == "S") &&
				(len(record.Get("Response.PathTakenForDomain").String()) == 0 ||
					len(record.Get("Response.ParallelTest").String()) == 0) {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Non-Writing test with non-adaptive structure"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			}
		}
	}()
	return out, errc, nil
}

func systemTestIncidents(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			disruptions := record.Get("Event.TestDisruptionList.TestDisruption").Array()
			if len(disruptions) > 0 {
				m := record.Value().(map[string]interface{})
				for _, e := range disruptions {
					m["TestDisruption"] = e.Get("Event").String()
					j, _ := json.Marshal(m)
					out <- gjson.Parse(string(j))
				}
			}
		}
	}()
	return out, errc, nil
}

func systemParticipationCodeImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			participationCode := record.Get("Event.ParticipationCode").String()
			rawscore := record.Get("Response.DomainScore.RawScore").String()
			rawscore_num, _ := strconv.Atoi(rawscore)
			if (len(record.Get("Response.PathTakenForDomain").String()) > 0 ||
				len(record.Get("Response.ParallelTest").String()) > 0) &&
				(participationCode != "P" && participationCode != "S") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Adaptive pathway without student undertaking test"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if len(rawscore) > 0 &&
				(participationCode != "P" && participationCode != "R") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Scored test with status other than P or R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if len(rawscore) > 0 &&
				participationCode == "R" &&
				rawscore_num != 0 {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Scored test with status other than P or R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if len(rawscore) == 0 &&
				(participationCode == "P" || participationCode == "R") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Unscored test with status of P or R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			}

		}
	}()
	return out, errc, nil
}

func systemTestCompleteness(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		counts := make(map[string]map[string]map[string]map[string]*set.Set)
		for record := range in {
			acaraid := record.Get("SchoolDetails.ACARAId").String()
			testdomain := record.Get("Test.TestContent.TestDomain").String()
			testlevel := record.Get("Test.TestContent.TestLevel").String()
			participationcode := record.Get("Event.ParticipationCode").String()
			psi := record.Get("Student.OtherIdList.OtherId.#[Type==NAPPlatformStudentId].Value").String()
			responseid := record.Get("Response.ResponseID").String()

			if _, ok := counts[acaraid]; !ok {
				counts[acaraid] = make(map[string]map[string]map[string]*set.Set)
			}
			if _, ok := counts[acaraid][testdomain]; !ok {
				counts[acaraid][testdomain] = make(map[string]map[string]*set.Set)
			}
			if _, ok := counts[acaraid][testdomain][testlevel]; !ok {
				counts[acaraid][testdomain][testlevel] = make(map[string]*set.Set)
				counts[acaraid][testdomain][testlevel]["P_attempts"] = set.New()
				counts[acaraid][testdomain][testlevel]["S_attempts"] = set.New()
				counts[acaraid][testdomain][testlevel]["responses"] = set.New()
			}
			if participationcode == "P" {
				counts[acaraid][testdomain][testlevel]["P_attempts"].Add(psi)
			} else if participationcode == "S" {
				counts[acaraid][testdomain][testlevel]["S_attempts"].Add(psi)
			}
			if (participationcode == "P" || participationcode == "S") && len(responseid) > 0 {
				counts[acaraid][testdomain][testlevel]["responses"].Add(psi)
			}
		}

		result := make(map[string]string)
		for k, v1 := range counts {
			for k1, v2 := range v1 {
				for k2, v3 := range v2 {
					result["ACARAID"] = k
					result["TestDomain"] = k1
					result["TestLevel"] = k2
					result["P_Attempts_Count"] = strconv.Itoa(v3["P_attempts"].Size())
					result["S_Attempts_Count"] = strconv.Itoa(v3["S_attempts"].Size())
					result["Responses_Count"] = strconv.Itoa(v3["responses"].Size())
					attempts := set.Union(v3["P_attempts"], v3["S_attempts"])
					result["Attempts_With_No_Response"] = set.Difference(attempts, v3["responses"]).String()
					result["Responses_With_No_Attempt"] = set.Difference(v3["responses"], attempts).String()
					j, _ := json.Marshal(result)
					out <- gjson.Parse(string(j))
				}
			}
		}
	}()
	return out, errc, nil
}

func systemObjectFrequency(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		counts := make(map[string]map[string][]string)
		for record := range in {
			psi := record.Get("Student.OtherIdList.OtherId.#[Type==NAPPlatformStudentId].Value").String()
			participationcode := record.Get("Event.ParticipationCode").String()
			eventcode := record.Get("Test.TestContent.TestLevel").String() + ":" + record.Get("Test.TestContent.TestDomain").String()
			responseid := record.Get("Response.ResponseID").String()

			if _, ok := counts[psi]; !ok {
				counts[psi] = make(map[string][]string)
				counts[psi]["events"] = make([]string, 0)
				counts[psi]["events_with_response"] = make([]string, 0)
				counts[psi]["responses"] = make([]string, 0)
			}
			counts[psi]["events"] = append(counts[psi]["events"], eventcode)
			if participationcode == "P" || participationcode == "R" || participationcode == "S" {
				counts[psi]["events_with_response"] = append(counts[psi]["events_with_response"], eventcode)
			}
			if len(responseid) > 0 {
				counts[psi]["responses"] = append(counts[psi]["responses"], eventcode)
			}
		}

		result := make(map[string]string)
		for k, v := range counts {
			result["PSI"] = k
			result["Events_Count"] = strconv.Itoa(len(v["events"]))
			result["PRS_Events_Count"] = strconv.Itoa(len(v["events_with_response"]))
			result["Responses_Count"] = strconv.Itoa(len(v["responses"]))
			result["Events"] = strings.Join(v["events"], ";")
			result["PRS_Events"] = strings.Join(v["events_with_response"], ";")
			result["Responses"] = strings.Join(v["responses"], ";")
			j, _ := json.Marshal(result)
			out <- gjson.Parse(string(j))
		}
	}()
	return out, errc, nil
}

func systemTestTypeItemImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			if record.Get("Test.TestContent.TestDomain").String() == "Writing" &&
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Score").String()) > 0 &&
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()) == 0 {
				m := record.Value().(map[string]interface{})
				m["Error"] = "No subscores for Writing test"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if record.Get("Test.TestContent.TestDomain").String() != "Writing" &&
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()) > 0 {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Subscores for Non-Writing test"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			}
		}
	}()
	return out, errc, nil
}

func systemParticipationCodeItemImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			participationcode := record.Get("ParticipationCode").String()
			lapsedtimeitem := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.LapsedTimeItem").String()
			testletscore := record.Get("Response.TestletList.Testlet.0.TestletScore").String()
			testletscore_num, _ := strconv.Atoi(testletscore)
			itemscore := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Score").String()
			itemscore_num, _ := strconv.Atoi(itemscore)

			if (len(lapsedtimeitem) > 0 ||
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Response").String()) > 0) &&
				participationcode != "P" && participationcode != "S" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Response captured without student writing test"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if (len(testletscore) > 0 ||
				len(itemscore) > 0 ||
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()) > 0) &&
				participationcode != "P" && participationcode != "R" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Scored test with status other than P or R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if ((len(testletscore) > 0 && testletscore_num != 0) ||
				(len(itemscore) > 0 && itemscore_num != 0) ||
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()) > 0) &&
				participationcode == "R" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Non-zero Scored test with status of R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if (len(testletscore) == 0 || len(itemscore) == 0) &&
				(participationcode == "R" || participationcode == "P") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Unscored test with status of P or R"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			} else if len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()) == 0 &&
				record.Get("Test.TestContent.TestDomain").String() == "Writing" &&
				participationcode == "P" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Unscored writing test with status of P"
				j, _ := json.Marshal(m)
				out <- gjson.Parse(string(j))
			}
		}
	}()
	return out, errc, nil
}

func systemItemCounts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		counts := make(map[string]map[string]map[string]map[string]map[string]int)
		for record := range in {
			testname := record.Get("Test.TestContent.TestName").String()
			testdomain := record.Get("Test.TestContent.TestDomain").String()
			testlevel := record.Get("Test.TestContent.TestLevel").String()
			itemlocalid := record.Get("TestItem.TestItemContent.NAPTestItemLocalId").String()
			participationcode := record.Get("Event.ParticipationCode").String()
			if record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.ResponseCorrectness").String() == "Not In Path" {
				continue
			}
			if participationcode != "P" && participationcode != "S" {
				continue
			}

			if _, ok := counts[testname]; !ok {
				counts[testname] = make(map[string]map[string]map[string]map[string]int)
			}
			if _, ok := counts[testname][testdomain]; !ok {
				counts[testname][testdomain] = make(map[string]map[string]map[string]int)
			}
			if _, ok := counts[testname][testdomain][testlevel]; !ok {
				counts[testname][testdomain][testlevel] = make(map[string]map[string]int)
			}
			if _, ok := counts[testname][testdomain][testlevel][itemlocalid]; !ok {
				counts[testname][testdomain][testlevel][itemlocalid] = make(map[string]int)
				counts[testname][testdomain][testlevel][itemlocalid]["substitute"] = 0
				counts[testname][testdomain][testlevel][itemlocalid]["count"] = 0
			}
			counts[testname][testdomain][testlevel][itemlocalid]["count"] = counts[testname][testdomain][testlevel][itemlocalid]["count"] + 1
			if len(record.Get("TestItem.TestItemContent.ItemSubstitutedForList.SubstituteItem").Array()) > 0 {
				counts[testname][testdomain][testlevel][itemlocalid]["substitute"] = 1
			}
		}

		result := make(map[string]string)
		for k, v := range counts {
			for k1, v1 := range v {
				for k2, v2 := range v1 {
					for k3, v3 := range v2 {
						result["TestName"] = k
						result["TestLevel"] = k1
						result["TestDomain"] = k2
						result["TestItemLocalId"] = k3
						if v3["substitute"] == 1 {
							result["Substitute"] = "true"
						} else {
							result["Substitute"] = "false"
						}
						result["Count"] = strconv.Itoa(v3["count"])
						j, _ := json.Marshal(result)
						out <- gjson.Parse(string(j))
					}
				}
			}
		}
	}()
	return out, errc, nil
}
