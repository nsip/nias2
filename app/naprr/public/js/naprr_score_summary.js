// naprr_score_summary.js
// 
// display & interactions for school score summary report
// 


// instantiate interaction listeners
$(document).ready(function()
{

    // initiate modals
    $('.modal').modal();

    // create the report when button is clicked
    $("#btn-scoresummary").on("click", function(event)
    {
        hideReport();
        createScoreSummaryReport();
        showReport();
    });

});


// handler routines for the score summary reports & extended data

// 
// Fetch data & display for score summary
// 
function createScoreSummaryReport()
{


    // order by level & alpha for test name
    sortScoreSummaryData(scoresummaryData);

    // create ui elements
    // 
    createScoreSummaryTitle();
    // filters
    createScoreSummaryFilters();
    initFilterHandlers();

    // table
    createScoreSummaryTableHeader();
    createScoreSummaryTableBody(scoresummaryData);
    initScoreSummaryTableHandler();

    // create report download link
    createScoreSummaryDownloadLink();
    initScoreSummaryDownloadLinkHandler();

}

// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initScoreSummaryTableHandler()
{

    // remove any existing handlers
    $('#report-content').off("click");
    // respond to table row selections
    $('#report-content').on("click", '#report-table-body tr', function(event)
    {
        var $row = event.currentTarget;
        var ssd = jQuery.data($row, "ssdata");
        createExtendedDataScoreSummary(ssd);
        showExtendedData();
    });

}


// 
// create the main tabular data diplay
// 
function createScoreSummaryTableBody(data)
{

    $("#report-table-body").empty();

    $.each(data, function(index, ss)
    {
        var $row = $("<tr/>");
        $row.append("<td>" + ss.Test.TestContent.TestLevel + "</td>" +
            "<td>" + ss.Test.TestContent.TestDomain + "</td>" +
            "<td>" + ss.Summ.DomainSchoolAverage + "</td>" +
            "<td>" + ss.Summ.DomainJurisdictionAverage + "</td>" +
            "<td>" + ss.Summ.DomainNationalAverage + "</td>" +
            "<td>" + ss.Summ.DomainTopNational60Percent + "</td>" +
            "<td>" + ss.Summ.DomainBottomNational60Percent + "</td>"
        );
        $row.data("ssdata", ss);
        $row.attr("yr-level", ss.Test.TestContent.TestLevel);
        $("#report-table-body").append($row);
        $row = null;
    });


}


// 
// set the report title
// 
function createScoreSummaryTitle()
{
    var title = $("#report-title");
    title.empty();
    title.text("School Score Summary");
}

// 
// create main report table header
// 
function createScoreSummaryTableHeader()
{
    $("#report-table-hdr").empty();

    var hdr = $("<tr></tr>");
    var hdr_row = $("<th>Level</th>" +
        "<th>Domain</th>" +
        "<th>School Average</th>" +
        "<th>Jurisd. Average</th>" +
        "<th>National Average</th>" +
        "<th>Top National 60%</th>" +
        "<th>Bottom National 60%</th>");

    hdr.append(hdr_row);
    $("#report-table-hdr").append(hdr);

}


// 
// when user selects line item in table
// show the extended data in a modal
// footer popup
// 
function createExtendedDataScoreSummary(ssdata)
{
    $("#ed-content").empty();

    $("#ed-modal").css("max-height", "80%");

    $("#ed-title").text('School Score Summary: ' +
        ssdata.Test.TestContent.TestDomain + ' Yr ' +
        ssdata.Test.TestContent.TestLevel);

    var topRow = $("<div class='row'><div class='col s12'></div></div>");
    topRow.append("Test Year:" + ssdata.Test.TestContent.TestYear +
        ", Test Type: " + ssdata.Test.TestContent.TestType +
        ", Test Stages: " + ssdata.Test.TestContent.StagesCount +
        ", Local Id:" + ssdata.Test.TestContent.LocalId +
        ", School: " + $('input.autocomplete').val()
    );

    topRow.append(getSchoolInfoSummaryLine(ssdata.Summ.SchoolACARAId));

    $("#ed-content").append(topRow);

    var twocol = $("<div class='row'></div>");
    var col1 = $("<div id='ed-col1' class='col s9'></div>");
    col1.append("<h5>Graph</h5>");


    var col2 = $("<div id='ed-col2' class='col s3'></div>");
    col2.append("<h5>Score</h5>");
    col2.append("<p>School Average: " +
        ssdata.Summ.DomainSchoolAverage + "</p>");
    col2.append("<p>Jurisdiction Average: " +
        ssdata.Summ.DomainJurisdictionAverage + "</p>");
    col2.append("<p>National Average: " +
        ssdata.Summ.DomainNationalAverage + "</p>");
    col2.append("<p>National Top 60%: " +
        ssdata.Summ.DomainTopNational60Percent + "</p>");
    col2.append("<p>National Bottom 60%: " +
        ssdata.Summ.DomainBottomNational60Percent + "</p>");

    twocol.append(col1);
    twocol.append(col2);

    $("#ed-content").append(twocol);

    col1.append("<div id='graphContainer'></div>");
    createScoreSummaryGraph(ssdata);

    // 
    // add score bands display
    // 
    $("#ed-content").append(createTestBandsDisplay(ssdata));
}

// 
// sort data - by year level & test name for now
// 
function sortScoreSummaryData(data)
{

    data.sort(function(a, b)
    {
        var compA = a.Test.TestContent.TestLevel.toUpperCase() +
            a.Test.TestContent.TestDomain.toUpperCase();

        var compB = b.Test.TestContent.TestLevel.toUpperCase() +
            b.Test.TestContent.TestDomain.toUpperCase();

        return (compA < compB) ? -1 : (compA > compB) ? 1 : 0;
    });

}

// 
// create option selectors to filter data table
// 
function createScoreSummaryFilters()
{

    $("#report-filters").empty();
    $("#report-filters").append(createYearLevelFilter());

}

// 
// add a download link to the csv report for this school
// 
function createScoreSummaryDownloadLink()
{

    $('#download-report').empty();
    $('#download-report').append("<a id='csv-download'>Download report as CSV file</a>");

}

// 
// handle the download
// 
function initScoreSummaryDownloadLinkHandler()
{

    var reportURL = "/naprr/downloadreport/" + currentASLId + "/score_summary.csv";

    $('#csv-download').off("click");

    $('#csv-download').on("click", function(event)
    {
        window.location.href = reportURL;

    });


}
