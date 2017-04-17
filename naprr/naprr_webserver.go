// naprr_webserver.go

package naprr

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ResultsReportingServer struct {
	sr *StreamReader
}

func NewResultsReportingServer() *ResultsReportingServer {
	return &ResultsReportingServer{sr: NewStreamReader()}
}

func (rrs *ResultsReportingServer) Run() {

	// create the web service framework
	e := echo.New()

	//
	// returns the list of known schools
	//
	e.GET("/naprr/schooldetails", func(c echo.Context) error {

		sds := rrs.sr.GetSchoolDetails()

		return c.JSON(http.StatusAccepted, sds)

	})

	//
	// return the school score summary report for the given school
	//
	e.GET("/naprr/scoresummary/:acaraid", func(c echo.Context) error {

		acaraID := c.Param("acaraid")
		ssumm := rrs.sr.GetScoreSummaryData(acaraID)

		return c.JSON(http.StatusAccepted, ssumm)

	})

	//
	// return the domain scores for the given school
	//
	e.GET("/naprr/domainscores/:acaraid", func(c echo.Context) error {

		acaraID := c.Param("acaraid")
		dscores := rrs.sr.GetDomainScoreData(acaraID)

		return c.JSON(http.StatusAccepted, dscores)

	})

	//
	// return the participation data for the given school
	//
	e.GET("/naprr/participation/:acaraid", func(c echo.Context) error {

		acaraID := c.Param("acaraid")
		ptcpn := rrs.sr.GetParticipationData(acaraID)

		return c.JSON(http.StatusAccepted, ptcpn)

	})

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

		localFileName := "./out/" + acaraID + "/" + fileName
		// log.Println("csv local file: ", localFileName)

		return c.File(localFileName)

	})

	//
	// download the codeframe report
	//
	e.GET("/naprr/downloadreport/codeframe", func(c echo.Context) error {

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+"codeframe.csv")
		c.Response().Header().Set("Content-Type", "text/csv")

		localFileName := "./out/codeframe.csv"
		// log.Println("csv local file: ", localFileName)

		return c.File(localFileName)

	})

	//
	// get the schoolinfo object for the given acaraid
	//
	e.GET("/naprr/schoolinfo/:acaraid", func(c echo.Context) error {

		acaraID := c.Param("acaraid")
		sd := rrs.sr.GetSchoolData(acaraID)

		return c.JSON(http.StatusAccepted, sd.SchoolInfos[acaraID])

	})

	//
	// get the test codeframe data
	//
	e.GET("/naprr/codeframe", func(c echo.Context) error {

		cfds := rrs.sr.GetCodeFrameData()

		return c.JSON(http.StatusAccepted, cfds)
	})

	// static resources
	e.Static("/", "public")
	e.Static("/reports", "out") // access to pre-generated reports

	// homepage
	e.File("/", "public/index.html")
	e.File("/naprr", "public/index.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("Starting web-ui services...")
	port := "1328" //DefaultValidationConfig.WebServerPort
	log.Println("Service is listening on localhost:" + port)

	e.Logger.Fatal(e.Start(":" + port))

}
