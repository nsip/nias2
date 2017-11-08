// qa-item-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strconv"

	"encoding/json"
	"github.com/tidwall/gjson"
)

//
// printing pipeline for QA reports depending on the item report query
//

func runQAItemRespReportPipeline(schools []string) error {

	reports := []string{"systemTestTypeItemImpacts.gql",
		"systemParticipationCodeItemImpacts.gql",
		"systemItemCounts.gql",
	}
	reports_path := "./reporting_templates/qa/"

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
	out_error_rpt_FileDir := "./out/qa/error_reports"
	err = os.MkdirAll(out_error_rpt_FileDir, os.ModePerm)
	if err != nil {
		return err
	}

	jsonc1, errc, err := splitter(ctx, jsonc, len(reports))
	errcList = append(errcList, errc)

	for i, queryFileName := range reports {
		var jsonc2 <-chan gjson.Result
		// for now, I'm merely hardcoding the query names
		// These are transforms on the CSV
		if queryFileName == "systemTestTypeItemImpacts.gql" {
			jsonc2, errc, _ = qaTestTypeItemImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemParticipationCodeItemImpacts.gql" {
			jsonc2, errc, _ = qaParticipationCodeItemImpacts(jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemItemCounts.gql" {
			jsonc2, errc, _ = qaItemCounts(jsonc1[i])
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

			log.Println("ITEM RESP QA report file writing... " + outFileName)
		}
	}
	return WaitForPipeline(errcList...)
}

func qaTestTypeItemImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaParticipationCodeItemImpacts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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

func qaItemCounts(in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
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
			participationcode := record.Get("ParticipationCode").String()
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
