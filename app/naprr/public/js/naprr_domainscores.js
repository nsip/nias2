// naprr_domainscores.js


// instantiate interaction listeners
$(document).ready(function()
{

    // initiate modals
    $('.modal').modal();

    // create the report when button is clicked
    $("#btn-domainscores").on("click", function(event)
    {
        hideReport();
        createDomainScoresReport();
        showReport();
    });

});

// 
// Fetch data & display for domain scores
// 
function createDomainScoresReport()
{


    // order by level & alpha for test name
    sortDomainScoresData(domainscoresData);

    // create ui elements
    // 
    createDomainScoresTitle();
    // filters
    createDomainScoresFilters();
    initFilterHandlers();

    // table
    createDomainScoresTableHeader();
    createDomainScoresTableBody(domainscoresData);
    initDomainScoresTableHandler();

    // create report download link
    createDomainScoresDownloadLink();
    initDomainScoresDownloadLinkHandler();

}

// 
// sort data - by year level & test domain for now
// 
function sortDomainScoresData(data)
{

    // console.log(data);

    data.sort(function(a, b)
    {
        var compA = a.Test.TestContent.TestLevel.toUpperCase() +
            a.Test.TestContent.TestDomain.toUpperCase() +
            a.Response.DomainScore.StudentDomainBand.toUpperCase();

        var compB = b.Test.TestContent.TestLevel.toUpperCase() +
            b.Test.TestContent.TestDomain.toUpperCase() +
            b.Response.DomainScore.StudentDomainBand.toUpperCase();

        return (compA < compB) ? -1 : (compA > compB) ? 1 : 0;
    });

}

// 
// set the report title
// 
function createDomainScoresTitle()
{
    var title = $("#report-title");
    title.empty();
    title.text("Student Domain Scores");
}

// 
// create option selectors to filter data table
// 
function createDomainScoresFilters()
{
    $("#report-filters").empty();
    $("#report-filters").append(createYearLevelFilter());
    $("#report-filters").append(createDomainFilter());

}

// 
// create main report table header
// 
function createDomainScoresTableHeader()
{

    $("#report-table-hdr").empty();

    var hdr = $("<tr></tr>");
    var hdr_row = $("<th>Level</th>" +
        "<th>Domain</th>" +
        "<th>PSI</th>" +
        "<th>Raw Score</th>" +
        "<th>Scaled Score</th>" +
        "<th>Scaled Score Std. Error</th>" +
        "<th>Domain Band</th>" +
        "<th>Proficiency</th>");

    hdr.append(hdr_row);
    $("#report-table-hdr").append(hdr);

}


// 
// create the main tabular data diplay
// 
function createDomainScoresTableBody(data)
{

    $("#report-table-body").empty();
    $("#report-table-body tr").remove();

    $.each(data, function(index, rds)
    {
        var $row = $("<tr/>");
        $row.append("<td>" + rds.Test.TestContent.TestLevel + "</td>" +
            "<td>" + rds.Test.TestContent.TestDomain + "</td>" +
            "<td>" + rds.Response.PSI + "</td>" +
            "<td>" + rds.Response.DomainScore.RawScore + "</td>" +
            "<td>" + rds.Response.DomainScore.ScaledScoreValue + "</td>" +
            "<td>" + rds.Response.DomainScore.ScaledScoreStandardError + "</td>" +
            "<td>" + rds.Response.DomainScore.StudentDomainBand + "</td>" +
            "<td>" + rds.Response.DomainScore.StudentProficiency + "</td>"
        );
        $row.data("rdsdata", rds);
        $row.attr("yr-level", rds.Test.TestContent.TestLevel);
        $row.attr("domain", rds.Test.TestContent.TestDomain);
        $("#report-table-body").append($row);
        $row = null;
    });
}

// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initDomainScoresTableHandler()
{
    // remove any existing handlers
    $('#report-content').off("click");
    // respond to table row selections
    $('#report-content').on("click", '#report-table-body tr', function(event)
    {
        var $row = event.currentTarget;
        var rds = jQuery.data($row, "rdsdata");
        createExtendedDataDomainScores(rds);
        showExtendedData();
    });

}

// 
// when user selects line item in table
// show the extended data in a modal
// footer popup
// 
function createExtendedDataDomainScores(rdsdata)
{
    $("#ed-content").empty();

    $("#ed-modal").css("max-height", "80%");

    $("#ed-title").text('Student Domain Score: ' +
        rdsdata.Test.TestContent.TestDomain + ' Yr ' +
        rdsdata.Test.TestContent.TestLevel);

    var topRow = $("<div class='row'><div class='col s12'></div></div>");
    topRow.append("Test Year:" + rdsdata.Test.TestContent.TestYear +
        ", Test Type: " + rdsdata.Test.TestContent.TestType +
        ", Test Stages: " + rdsdata.Test.TestContent.StagesCount +
        ", Local Id:" + rdsdata.Test.TestContent.LocalId);

    topRow.append(getStudentInfoSummaryLine(rdsdata.Response.PSI));


    $("#ed-content").append(topRow);

    var threecol = $("<div class='row'></div>");
    var col1 = $("<div id='ed-col1' class='col s6'></div>");
    col1.append("<h5>Graph</h5>");

    // 
    // score details
    // 
    var col2 = $("<div id='ed-col2' class='col s3'></div>");
    col2.append("<h5>Domain Score</h5>");
    col2.append("<p>Raw Score: " +
        rdsdata.Response.DomainScore.RawScore + "</p>");
    col2.append("<p>Scaled Score: " +
        rdsdata.Response.DomainScore.ScaledScoreValue + "</p>");
    col2.append("<p>Scaled Score Logit: " +
        rdsdata.Response.DomainScore.ScaledScoreLogitValue + "</p>");
    col2.append("<p>Scaled Score Standard Error: " +
        rdsdata.Response.DomainScore.ScaledScoreStandardError);
    col2.append("<p>Scaled Score Standard Error Logit: " +
        rdsdata.Response.DomainScore.ScaledScoreLogitStandardError + "</p>");

    // 
    // extended score info
    // 
    var col3 = $("<div id='ed-col2' class='col s3'></div>");
    col3.append("<h5>Score Info</h5>");
    col3.append("<p>Calibration Sample: " +
        rdsdata.Response.CalibrationSampleFlag + "</p>");
    col3.append("<p>Equating Sample: " +
        rdsdata.Response.EquatingSampleFlag + "</p>");
    col3.append("<p>Path Taken: " +
        rdsdata.Response.PathTakenForDomain + "</p>");
    col3.append("<p>Parallel Test Path: " +
        rdsdata.Response.ParallelTest + "</p>");
    col3.append("<p>Student Domain Band: " +
        rdsdata.Response.DomainScore.StudentDomainBand + "</p>");
    col3.append("<p>Proficiency: " +
        rdsdata.Response.DomainScore.StudentProficiency + "</p>");


    threecol.append(col1);
    threecol.append(col2);
    threecol.append(col3);

    $("#ed-content").append(threecol);

    col1.append("<div id='graphContainer'></div>");
    createDomainScoreGraph(rdsdata);


    // 
    // add score bands display
    // 

    $("#ed-content").append(createTestBandsDisplay(rdsdata));

}

// 
// add a download link to the csv report for this school
// 
function createDomainScoresDownloadLink()
{

    $('#download-report').empty();
    $('#download-report').append("<a id='csv-download'>Download report as CSV file</a>");

}

// 
// handle the download
// 
function initDomainScoresDownloadLinkHandler()
{

    var reportURL = "/naprr/downloadreport/" + currentASLId + "/domain_scores.csv";

    $('#csv-download').off("click");

    $('#csv-download').on("click", function(event)
    {
        window.location.href = reportURL;

    });


}
