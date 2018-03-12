// sifql-server.go
//
// simple web server to support gql queries & web ui (graphiql)
//
package naprrql

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/playlyfe/go-graphql"
)

var nap_executor *graphql.Executor
var report_executor *graphql.Executor

// var isr_executor *graphql.Executor
// var item_executor *graphql.Executor

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// the core graphql handler routine
//
func graphQLHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := nap_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

//
// the graphql handler routine for reports: kept separate since it will be static,
// whereas the main graphql handler can have its contents experimented with
//
func reportHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := report_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

//
// specialist handler for creating isr printing files
//
/*
func isrPrintHandlerDISABLE(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := isr_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}
*/

//
// specialist handler for creating item result printing files
//
/*
func itemPrintHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := item_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}
*/

//
// launches the server
//
func RunQLServer() {

	nap_executor = buildNAPExecutor()
	report_executor = buildNAPExecutor()
	// isr_executor = buildISRPrintExecutor()
	// item_executor = buildItemPrintExecutor()

	e := echo.New()

	e.Use(middleware.Gzip())

	// routes to html/css/javascript resources
	e.Static("/", "public")
	e.File("/sifql", "public/ql_index.html")
	e.File("/ui", "public/ui_index.html")
	e.File("/datamodel", "public/vy_index.html")

	// the main graphql handler
	e.POST("/graphql", graphQLHandler)

	// the graphql handler for reports
	e.POST("/report", reportHandler)

	// special handler for isr printing
	// e.POST("/isrprint", isrPrintHandler)

	// special handler for item printing
	// e.POST("/itemprint", itemPrintHandler)

	//
	// download the requested pre-generated csv file.
	//
	e.GET("/naprr/downloadreport/:acaraid/:filename", func(c echo.Context) error {

		acaraID := c.Param("acaraid")
		fileName := c.Param("filename")
		uniqFileName := acaraID + "_" + fileName

		// log.Println("csvfile return name: ", uniqFileName)

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+uniqFileName)
		c.Response().Header().Set("Content-Type", "text/csv")

		localFileName := "./out/school_reports/" + acaraID + "/" + fileName
		// log.Println("csv local file: ", localFileName)

		return c.File(localFileName)

	})

	//
	// download the codeframe report
	//
	e.GET("/naprr/downloadreport/codeframe", func(c echo.Context) error {

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+"codeframe.csv")
		c.Response().Header().Set("Content-Type", "text/csv")

		localFileName := "./out/system_reports/systemCodeframe.csv"
		// log.Println("csv local file: ", localFileName)

		return c.File(localFileName)

	})

	e.Logger.Fatal(e.Start(":1329"))
}
