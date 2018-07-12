// query-report-writer.go

package naprrql

import (
	"context"
	"encoding/json"
	"log"
	"os"

	nx "github.com/nsip/nias2/xml"
	"github.com/pkg/errors"
	//"github.com/tidwall/sjson"
)

//
// create and run an xml printing pipeline.
//
// Pipeline streams requests at a school by school
// feeding results to the ouput xml file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//

var codeframe [][]byte = nil

func runXMLPipeline(schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	varsc, errc, err := xmlParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	jsonc, errc, err := xmlQueryExecutor(ctx, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/xml"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	xmlFileName := "sif.xml"
	outFileName := outFileDir + "/" + xmlFileName
	errc, err = xmlFileSink(ctx, outFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("XML file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

func runXMLPipelinePerSchool(school string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	varsc, errc, err := xmlParametersSource(ctx, school)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	jsonc, errc, err := xmlQueryExecutor(ctx, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/xml/" + school + "/"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	xmlFileName := "sif.xml"
	outFileName := outFileDir + "/" + xmlFileName
	errc, err = xmlFileSink(ctx, outFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("XML file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn
//
func xmlParametersSource(ctx context.Context, schools ...string) (<-chan systemQueryParams, <-chan error, error) {

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
// query executor transform stage takes query params in, excutes gql query
// and writes results to output chaneel
//
func xmlQueryExecutor(ctx context.Context, in <-chan systemQueryParams) (<-chan []byte, <-chan error, error) {
	out := make(chan []byte)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		// there are two queries: XML independent of schools (the codeframe and test objects), and XML dependent on schools
		json, err := xmlCodeframeQuery(map[string]interface{}{})
		if err != nil {
			// Handle an error that occurs during the goroutine.
			errc <- err
			return
		}
		for _, result := range json {
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- result:
			case <-ctx.Done():
				return
			}

		}
		for params := range in {
			vars := map[string]interface{}{"acaraIDs": []string{params.schoolAcaraID}}
			json, err := xmlQuery(vars)
			if err != nil {
				errc <- err
				return
			}
			for _, result := range json {
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}

			}
		}
	}()
	return out, errc, nil
}

func xmlCodeframeQuery(params map[string]interface{}) ([][]byte, error) {
	if codeframe != nil {
		return codeframe, nil
	}
	ret := make([][]byte, 0)

	ids := getIdentifiers("NAPCodeFrame:")
	objs, err := getObjects(ids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, obj := range objs {
		o, _ := obj.(nx.NAPCodeFrame)
		b, err := json.Marshal(TypedObject{NAPCodeFrame: &o})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	ids = getIdentifiers("NAPTest:")
	objs, err = getObjects(ids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, obj := range objs {
		o, _ := obj.(nx.NAPTest)
		b, err := json.Marshal(TypedObject{NAPTest: &o})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	ids = getIdentifiers("NAPTestlet:")
	objs, err = getObjects(ids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, obj := range objs {
		o, _ := obj.(nx.NAPTestlet)
		b, err := json.Marshal(TypedObject{NAPTestlet: &o})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	ids = getIdentifiers("NAPTestItem:")
	objs, err = getObjects(ids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, obj := range objs {
		o, _ := obj.(nx.NAPTestItem)
		b, err := json.Marshal(TypedObject{NAPTestItem: &o})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}
	codeframe = ret
	return ret, err

}

func xmlQuery(params map[string]interface{}) ([][]byte, error) {
	// get the acara ids from the request params
	acaraids := make([]string, 0)
	for _, a_id := range params["acaraIDs"].([]string) {
		acaraids = append(acaraids, a_id)
	}

	// get the sif refid for each of the acarids supplied
	refids := make([]string, 0)
	for _, acaraid := range acaraids {
		refid := getIdentifiers(acaraid + ":")
		if len(refid) > 0 {
			refids = append(refids, refid...)
		}

	}

	ret := make([][]byte, 0)

	siObjects, err := getObjects(refids)
	for _, sio := range siObjects {
		si, _ := sio.(nx.SchoolInfo)
		b, err := json.Marshal(TypedObject{SchoolInfo: &si})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)

	}

	// now construct the composite keys
	school_summary_keys := make([]string, 0)
	for _, refid := range refids {
		school_summary_keys = append(school_summary_keys, refid+":NAPTestScoreSummary:")
	}

	summ_refids := make([]string, 0)
	for _, summary_key := range school_summary_keys {
		ids := getIdentifiers(summary_key)
		for _, id := range ids {
			summ_refids = append(summ_refids, id)
		}
	}
	summaries, err := getObjects(summ_refids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, summary := range summaries {
		summ, _ := summary.(nx.NAPTestScoreSummary)
		b, err := json.Marshal(TypedObject{NAPTestScoreSummary: &summ})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	studentids := make([]string, 0)
	for _, acaraid := range acaraids {
		key := "student_by_acaraid:" + acaraid
		studentRefIds := getIdentifiers(key)
		studentids = append(studentids, studentRefIds...)
	}
	students, err := getObjects(studentids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, student := range students {
		s, _ := student.(nx.RegistrationRecord)
		b, err := json.Marshal(TypedObject{StudentPersonal: &s})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	responseids := make([]string, 0)
	for _, studentid := range studentids {
		key := "responseset_by_student:" + studentid
		responseRefId := getIdentifiers(key)
		responseids = append(responseids, responseRefId...)
	}
	responses, err := getObjects(responseids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, response := range responses {
		r, _ := response.(nx.NAPResponseSet)
		b, err := json.Marshal(TypedObject{NAPStudentResponseSet: &r})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}

	eventids := make([]string, 0)
	for _, studentid := range studentids {
		key := studentid + ":NAPEventStudentLink:"
		eventRefId := getIdentifiers(key)
		eventids = append(eventids, eventRefId...)
	}
	events, err := getObjects(eventids)
	if err != nil {
		return [][]byte{}, err
	}
	for _, event := range events {
		e, _ := event.(nx.NAPEvent)
		b, err := json.Marshal(TypedObject{NAPEventStudentLink: &e})
		if err != nil {
			return [][]byte{}, err
		}
		ret = append(ret, b)
	}
	return ret, err

}
