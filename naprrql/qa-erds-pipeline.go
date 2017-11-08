// qa-erds-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/tidwall/gjson"
	"gopkg.in/fatih/set.v0"
)

//
// printing pipeline for QA reports depending on the ERDS query
// (event-response-data-set)
//

func runQAErdsReportPipeline(schools []string) error {

	reports_path := "./reporting_templates/qa/"

	reports := []string{"systemTestAttempts.gql",
		"systemParticipationCodeImpacts.gql",
		"systemTestTypeImpacts.gql",
		"systemObjectFrequency.gql",
		"systemTestCompleteness.gql",
		"systemTestIncidents.gql",
	}

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
	out_error_rpt_FileDir := "./out/qa/error_reports"
	err = os.MkdirAll(out_error_rpt_FileDir, os.ModePerm)
	if err != nil {
		return err
	}
	jsonc1, errc, err := splitter(ctx, jsonc, len(reports))
	errcList = append(errcList, errc)

	for i, queryFileName := range reports {
		var jsonc2 <-chan gjson.Result
		// These are transforms on the CSV
		if queryFileName == "systemTestAttempts.gql" {
			jsonc2, errc, _ = qaTestAttempts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemTestTypeImpacts.gql" {
			jsonc2, errc, _ = qaTestTypeImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemTestIncidents.gql" {
			jsonc2, errc, _ = qaTestIncidents(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemParticipationCodeImpacts.gql" {
			jsonc2, errc, _ = qaParticipationCodeImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemTestCompleteness.gql" {
			jsonc2, errc, _ = qaTestCompleteness(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemObjectFrequency.gql" {
			jsonc2, errc, _ = qaObjectFrequency(jsonc1[i])
			errcList = append(errcList, errc)
		} else {
			jsonc2 = nil
		}
		if jsonc2 != nil {
			csvFileName := deriveCSVFileName(queryFileName)
			outFileName := out_error_rpt_FileDir + "/" + csvFileName
			mapFileName := deriveMapFileName(reports_path + queryFileName)
			errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc2)
			if err != nil {
				return err
			}
			errcList = append(errcList, errc)

			log.Println("ERDS QA file writing... " + outFileName)
		}
	}
	return WaitForPipeline(errcList...)
}

func qaTestAttempts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaTestTypeImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaTestIncidents(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaParticipationCodeImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaTestCompleteness(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaObjectFrequency(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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
