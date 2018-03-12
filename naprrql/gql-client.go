// gql-client.go
package naprrql

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/tidwall/gjson"
)

var DEF_INTERACTIVE_GQL_URL = "http://localhost:1329/graphql"
var DEF_GQL_URL = "http://localhost:1329/report"
var DEF_ISR_URL = "http://localhost:1329/isrprint"
var DEF_ITEM_URL = "http://localhost:1329/itemprint"

//
// this module provides a golang client that
// accesses gql servers for query support; used for e.g. reporting.
//
// Gql queries are sent, returns an array of data from the query
// or an error.
//

//
// Create a client to run gql queries & return json results/errors
//
type GqlClient struct{}

//
// Returns a new gql client
//
func NewGqlClient() *GqlClient {
	return &GqlClient{}
}

//
// make a gql query and return the json payload
//
// gqlURL - the endpoint url for the query service i.e. http://loclahost:1329/graphql
//
// query - the query itself
//
// variables - map of named parameters passed to the query i.e. acaraID = "12345"
//
// returns a gjson.Result - the data array contianed in the response
//
// returns error - either errors from the call attempt or errors within the response
//
func (gqlc *GqlClient) DoQuery(gqlURL string, query string, variables map[string]interface{}) (gjson.Result, error) {

	emptyResponse := gjson.Result{}

	// assemble the request & client to invoke it
	gqlReq := GQLRequest{Query: query, Variables: variables}
	jsonStr, err := json.Marshal(gqlReq)
	req, err := http.NewRequest("POST", gqlURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// connect to webserver, exp. backoff to ensure
	// server has time to come up before querying
	var resp *http.Response
	operation := func() error {
		resp, err = client.Do(req)
		return err // or an error
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		// if we got here connection was made but there's a
		// different error, so pass up the stack
		return emptyResponse, err
	}
	defer resp.Body.Close()

	// process the returned json, content taken from the
	// element in the returned payload identified in the query
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyResponse, err
	}
	resultRoot := getDataPath(query)
	response := gjson.GetBytes(body, resultRoot)

	return response, nil

}

//
// parses a gql query to find the root node at which the
// data array will be returned
//
func getDataPath(query string) string {

	defRoot := "data"

	tokens := strings.SplitAfterN(query, "{", 3)

	// not an array query, flat members e.g. object counts
	if len(tokens) < 3 {
		return defRoot
	}

	// remove variables clause from first line of query
	// to find root element
	line1_tokens := strings.SplitN(tokens[1], "{", 2)
	query_tokens := strings.SplitN(line1_tokens[0], "(", 2)

	queryRoot := strings.TrimSpace(query_tokens[0])

	return defRoot + "." + queryRoot

}
