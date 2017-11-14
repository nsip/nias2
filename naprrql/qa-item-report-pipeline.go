// qa-item-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strconv"

	"encoding/json"
	"github.com/tidwall/gjson"
	"gopkg.in/fatih/set.v0"
)

//
// printing pipeline for QA reports depending on the item report query
//

func itemresults_query() string {
	return `query NAPItemResults($acaraIDs: [String]) {
  item_results_report_by_school(acaraIDs: $acaraIDs) {
    Test {
      TestContent {
        LocalId
        TestName
        TestLevel
        TestDomain
        TestYear
        TestType
      }
    }
    Testlet {
      TestletContent {
        LocalId
        Node
        LocationInStage
        TestletName
        LocalId
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
          }
        }
      }
    }
    Student {
      BirthDate
      Sex
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
}

func codeframeQuery() string {
	return `query CodeFrame {
  codeframe_report {
    Test {
      TestContent {
	LocalId
        TestName
        TestType
        TestLevel
        TestDomain
      }
    }
    Testlet {
      TestletContent {
        LocalId
        Node
        LocationInStage
        TestletName
        LocalId
      }
    }
    Item {
      ItemID
      TestItemContent {
        NAPTestItemLocalId
        Subdomain
        MarkingType
        ItemSubstitutedForList {
          SubstituteItem {
            SubstituteItemRefId
	    LocalId
          }
        }
      }
    }
  }
}`
}

func runQAItemRespReportPipeline(schools []string) error {

	reports := []string{
		"systemTestTypeItemImpacts.gql",
		"systemParticipationCodeItemImpacts.gql",
		"systemItemCounts.gql",
		"itemExpectedResponses.gql",
		"itemPrinting.gql",
	}
	reports_path := "./reporting_templates/qa/"

	itemresp_query := itemresults_query()
	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	varsc1, errc, err := systemParametersSource(ctx, schools...)
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

	// get codeframe
	codeframec, errc, err := systemQueryExecutor(ctx, codeframeQuery(), DEF_GQL_URL, varsc1)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	//cf, sub := qaReadCodeframe(ctx, codeframec)

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
	codeframec1, errc, err := splitter(ctx, codeframec, 2)
	errcList = append(errcList, errc)
	for i, queryFileName := range reports {
		// for now, I'm merely hardcoding the query names
		// These are transforms on the CSV
		output_prefix := out_error_rpt_FileDir
		var jsonc2 <-chan gjson.Result
		if queryFileName == "systemTestTypeItemImpacts.gql" {
			jsonc2, errc, _ = qaTestTypeItemImpacts(ctx, jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemParticipationCodeItemImpacts.gql" {
			jsonc2, errc, _ = qaParticipationCodeItemImpacts(ctx, jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "systemItemCounts.gql" {
			jsonc2, errc, _ = qaItemCounts(ctx, codeframec1[0], jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "itemExpectedResponses.gql" {
			// block on reading codeframe
			jsonc2, errc, _ = qaItemExpectedResponses(ctx, codeframec1[1], jsonc1[i])
			errcList = append(errcList, errc)
		} else if queryFileName == "itemPrinting.gql" {
			jsonc2, errc, _ = qaItemResponses(ctx, jsonc1[i])
			output_prefix = outFileDir
			errcList = append(errcList, errc)
		} else {
			jsonc2 = nil
		}
		/*if jsonc2 != nil*/ {
			csvFileName := deriveCSVFileName(queryFileName)
			outFileName := output_prefix + "/" + csvFileName
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

func qaTestTypeItemImpacts(ctx context.Context, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		var j []byte
		for record := range in {
			j = nil
			subscores := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()
			if record.Get("Test.TestContent.TestDomain").String() == "Writing" &&
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Score").String()) > 0 &&
				len(subscores) == 0 {
				m := record.Value().(map[string]interface{})
				m["Error"] = "No subscores for Writing test"
				j, _ = json.Marshal(m)
			} else if record.Get("Test.TestContent.TestDomain").String() != "Writing" &&
				len(subscores) > 0 {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Subscores for Non-Writing test"
				j, _ = json.Marshal(m)
			}
			if j != nil {
				select {
				case out <- gjson.ParseBytes(j):
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, errc, nil
}

func qaParticipationCodeItemImpacts(ctx context.Context, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		var j []byte
		for record := range in {
			participationcode := record.Get("ParticipationCode").String()
			lapsedtimeitem := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.LapsedTimeItem").String()
			testletscore := record.Get("Response.TestletList.Testlet.0.TestletScore").String()
			testletscore_num, _ := strconv.Atoi(testletscore)
			itemscore := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Score").String()
			itemscore_num, _ := strconv.Atoi(itemscore)
			subscores := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.SubscoreList.Subscore").Array()
			j = nil
			if (len(lapsedtimeitem) > 0 ||
				len(record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.Response").String()) > 0) &&
				participationcode != "P" && participationcode != "S" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Response captured without student writing test"
				j, _ = json.Marshal(m)
			} else if (len(testletscore) > 0 ||
				len(itemscore) > 0 ||
				len(subscores) > 0) &&
				participationcode != "P" && participationcode != "R" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Scored test with status other than P or R"
				j, _ = json.Marshal(m)
			} else if ((len(testletscore) > 0 && testletscore_num != 0) ||
				(len(itemscore) > 0 && itemscore_num != 0) || len(subscores) > 0) &&
				participationcode == "R" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Non-zero Scored test with status of R"
				j, _ = json.Marshal(m)
			} else if (len(testletscore) == 0) &&
				(participationcode == "R" || participationcode == "P") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Missing testlet score with status of P or R"
				j, _ = json.Marshal(m)
			} else if (len(itemscore) == 0) &&
				(participationcode == "R" || participationcode == "P") {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Unscored test with status of P or R"
				j, _ = json.Marshal(m)
			} else if len(subscores) == 0 &&
				record.Get("Test.TestContent.TestDomain").String() == "Writing" &&
				participationcode == "P" {
				m := record.Value().(map[string]interface{})
				m["Error"] = "Unscored writing test with status of P"
				j, _ = json.Marshal(m)
			}
			if j != nil {
				select {
				case out <- gjson.ParseBytes(j):
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, errc, nil
}

func qaItemCountsRowInit(counts map[string]map[string]map[string]map[string]map[string]int, testname string, testdomain string,
	testlevel string, itemlocalid string, hasSubstituteItems bool) map[string]map[string]map[string]map[string]map[string]int {
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
	if hasSubstituteItems {
		counts[testname][testdomain][testlevel][itemlocalid]["substitute"] = 1
	}
	return counts
}

func qaItemCounts(ctx context.Context, codeframe <-chan gjson.Result, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		counts := make(map[string]map[string]map[string]map[string]map[string]int)
		for record := range codeframe {
			testname := record.Get("Test.TestContent.TestName").String()
			testdomain := record.Get("Test.TestContent.TestDomain").String()
			testlevel := record.Get("Test.TestContent.TestLevel").String()
			itemlocalid := record.Get("Item.TestItemContent.NAPTestItemLocalId").String()
			counts = qaItemCountsRowInit(counts, testname, testdomain, testlevel, itemlocalid,
				len(record.Get("Item.TestItemContent.ItemSubstitutedForList.SubstituteItem").Array()) > 0)
		}
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
			counts = qaItemCountsRowInit(counts, testname, testdomain, testlevel, itemlocalid,
				len(record.Get("TestItem.TestItemContent.ItemSubstitutedForList.SubstituteItem").Array()) > 0)
			counts[testname][testdomain][testlevel][itemlocalid]["count"] = counts[testname][testdomain][testlevel][itemlocalid]["count"] + 1
		}

		result := make(map[string]string)
		for k, v := range counts {
			for k1, v1 := range v {
				for k2, v2 := range v1 {
					for k3, v3 := range v2 {
						result["TestName"] = k
						result["TestLevel"] = k2
						result["TestDomain"] = k1
						result["TestItemLocalId"] = k3
						if v3["substitute"] == 1 {
							result["Substitute"] = "true"
						} else {
							result["Substitute"] = "false"
						}
						result["Count"] = strconv.Itoa(v3["count"])
						j, _ := json.Marshal(result)
						select {
						case out <- gjson.ParseBytes(j):
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}()
	return out, errc, nil
}

func qaItemResponses(ctx context.Context, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			select {
			case out <- record:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

type QaItemExpectedResponseType struct {
	PSI                   string
	TestName              string
	ParticipationCode     string
	TestletName           map[string]string
	ExpectedItemsCount    map[string]int
	FoundItems            map[string]map[string]int
	ExpectedItemsNotFound map[string]string
	FoundItemsNotExpected map[string]string
}

func QaItemExpectedResponseTypeNew() QaItemExpectedResponseType {
	s := QaItemExpectedResponseType{}
	s.PSI = ""
	s.TestName = ""
	s.ParticipationCode = ""
	s.TestletName = make(map[string]string)
	s.ExpectedItemsCount = make(map[string]int)
	s.FoundItems = make(map[string]map[string]int)
	s.FoundItems["Correct"] = make(map[string]int)
	s.FoundItems["Incorrect"] = make(map[string]int)
	s.FoundItems["NotAttempted"] = make(map[string]int)
	s.FoundItems["NotInPath"] = make(map[string]int)
	s.ExpectedItemsNotFound = make(map[string]string)
	s.FoundItemsNotExpected = make(map[string]string)
	return s
}

func qaReadCodeframe(ctx context.Context, codeframe <-chan gjson.Result) (map[string]map[string]*set.Set, map[string]*set.Set) {
	done := make(chan bool, 1)
	// read codeframe
	cf := make(map[string]map[string]*set.Set)
	// track substitute items
	sub := make(map[string]*set.Set)
	go func() {
		for record := range codeframe {
			testid := record.Get("Test.TestContent.LocalId").String()
			testletid := record.Get("Testlet.TestletContent.LocalId").String()
			substitutes := record.Get("Item.TestItemContent.ItemSubstitutedForList.SubstituteItem").Array()
			itemid := record.Get("Item.TestItemContent.NAPTestItemLocalId").String()
			//log.Printf("%s\t%+v\n", itemid, substitutes)
			if _, ok := cf[testid]; !ok {
				cf[testid] = make(map[string]*set.Set)
			}
			if _, ok := cf[testid][testletid]; !ok {
				cf[testid][testletid] = set.New()
			}
			if len(substitutes) == 0 {
				cf[testid][testletid].Add(itemid)
			} else {
				sub[itemid] = set.New()
				for _, s := range substitutes {
					//substitutes[0].Get("LocalId").String()
					sub[itemid].Add(s.Get("LocalId").String())
				}
			}
		}
		done <- true
	}()
	<-done
	//log.Printf("%+v\n", sub)
	return cf, sub
}

//func qaItemExpectedResponses(ctx context.Context, cf map[string]map[string]*set.Set, sub map[string]string, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
func qaItemExpectedResponses(ctx context.Context, codeframec <-chan gjson.Result, in <-chan gjson.Result) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		cf, sub := qaReadCodeframe(ctx, codeframec)
		defer close(out)
		defer close(errc)
		// we are assuming the records come in for item results report in sorted order
		curr_testletid := ""
		curr_testid := ""
		curr_psi := ""
		curr_locationinstage := ""
		locationinstage := 0
		locationinstage_str := ""
		testitems := set.New()
		result := QaItemExpectedResponseTypeNew()
		for record := range in {
			psi := record.Get("Response.PSI").String()
			testletid := record.Get("Testlet.TestletContent.LocalId").String()
			testletname := record.Get("Testlet.TestletContent.TestletName").String()
			testid := record.Get("Test.TestContent.LocalId").String()
			testname := record.Get("Test.TestContent.TestName").String()
			itemid := record.Get("TestItem.TestItemContent.NAPTestItemLocalId").String()
			correctness := record.Get("Response.TestletList.Testlet.0.ItemResponseList.ItemResponse.0.ResponseCorrectness").String()
			participationcode := record.Get("ParticipationCode").String()
			// we're assuming LocationInStage is only 1, 2, 3
			// locationinstage := record.Get("Testlet.TestletContent.LocationInStage").String()
			// We're just incrementing testlets
			//if len(locationinstage) == 0 {
			//	locationinstage = "1"
			//}

			if !testitems.IsEmpty() && (testletid != curr_testletid || testid != curr_testid || psi != curr_psi) {
				result = checkExpectedItems(result, cf, sub, curr_testid, curr_testletid, curr_locationinstage, testitems)
			}
			if psi != curr_psi || testid != curr_testid {
				if psi != curr_psi {
					//log.Println(psi)
				}
				if !testitems.IsEmpty() {
					j, _ := json.Marshal(result)
					select {
					case out <- gjson.ParseBytes(j):
					case <-ctx.Done():
						return
					}
				}
				result = QaItemExpectedResponseTypeNew()
				result.PSI = psi
				result.TestName = testname
				result.ParticipationCode = participationcode
			}
			if testletid != curr_testletid || testid != curr_testid || psi != curr_psi {
				if testid != curr_testid || psi != curr_psi {
					locationinstage = 1
				} else {
					locationinstage++
				}
				locationinstage_str = strconv.Itoa(locationinstage)
				result.TestletName[locationinstage_str] = testletname
				if _, ok := cf[testid][testletid]; ok {
					result.ExpectedItemsCount[locationinstage_str] = cf[testid][testletid].Size()
				} else {
					result.ExpectedItemsCount[locationinstage_str] = 0
					log.Printf("NOT IN CODEFRAME: %s\t%s\n", testid, testletid)
				}

				testitems.Clear()
			}
			switch correctness {
			case "Correct":
				result.FoundItems["Correct"][locationinstage_str]++
			case "Incorrect":
				result.FoundItems["Incorrect"][locationinstage_str]++
			case "NotAttempted":
				result.FoundItems["NotAttempted"][locationinstage_str]++
			case "NotInPath":
				result.FoundItems["NotInPath"][locationinstage_str]++
			}
			testitems.Add(itemid)
			curr_testletid = testletid
			curr_testid = testid
			curr_psi = psi
			curr_locationinstage = locationinstage_str
			//log.Printf("%s\t%s\t%s\t%s\t%d\n", psi, testname, testletname, itemid, locationinstage)
			// log.Printf("%+v\n", result)
		}
		result = checkExpectedItems(result, cf, sub, curr_testid, curr_testletid, curr_locationinstage, testitems)
		j, _ := json.Marshal(result)
		select {
		case out <- gjson.ParseBytes(j):
		case <-ctx.Done():
			return
		}
	}()
	return out, errc, nil
}

func checkExpectedItems(result QaItemExpectedResponseType, cf map[string]map[string]*set.Set, sub map[string]*set.Set, curr_testid string, curr_testletid string, curr_locationinstage string, testitems *set.Set) QaItemExpectedResponseType {
	if expected_testitems, ok := cf[curr_testid][curr_testletid]; ok {

		foundNotExp := set.New()
		testitems_expansion := set.New()
		//log.Printf("Expected: %+v\n", expected_testitems)
		//log.Printf("Found   : %+v\n", testitems)
		for _, t := range set.StringSlice(testitems) {
			if s, ok := sub[t]; ok {
				testitems_expansion.Merge(s)
			} else {
				testitems_expansion.Add(t)
			}
		}
		//log.Printf("Expanded : %+v\n", testitems_expansion)
		result.ExpectedItemsNotFound[curr_locationinstage] = set.Difference(expected_testitems, testitems_expansion).String()
		for _, t := range set.StringSlice(testitems) {
			if s, ok := sub[t]; ok {
				found := false
				for _, s1 := range s.List() {
					if expected_testitems.Has(s1) {
						found = true
					}
				}
				if !found {
					foundNotExp.Add(t)
				}
			} else {
				if !expected_testitems.Has(t) {
					foundNotExp.Add(t)
				}
			}
		}
		//log.Printf("FoundNEx : %+v\n", foundNotExp)
		result.ExpectedItemsNotFound[curr_locationinstage] = foundNotExp.String()
	} else {
		result.ExpectedItemsNotFound[curr_locationinstage] = ""
		result.FoundItemsNotExpected[curr_locationinstage] = testitems.String()
	}
	return result
}
