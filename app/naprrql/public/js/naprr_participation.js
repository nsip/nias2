// naprr_participation.js

// instantiate interaction listeners
$(document).ready(function() {

    // initiate modals
    $('.modal').modal();

    // create the report when button is clicked
    $("#btn-participation").on("click", function(event) {
        hideReport();
        createParticipationReport();
        showReport();
    });

});

// 
// gql query to support search
// 
function participationQuery() {
    return `
query ParticipationData($acaraIDs: [String]) {
  participation_report_by_school(acaraIDs: $acaraIDs) {
    Student {
      FamilyName
      GivenName
      ClassGroup
      Homegroup
      OfflineDelivery
      EducationSupport
      HomeSchooledStudent
    }
    School {
      SchoolName
      SchoolSector
      SchoolType
      IndependentSchool
      ACARAId
      LocalId
    }
    EventInfos {
      Event {
        EventID
        SPRefID
        PSI
        SchoolRefID
        SchoolID
        TestID
        NAPTestLocalID
        SchoolSector
        System
        SchoolGeolocation
        ReportingSchoolName
        JurisdictionID
        ParticipationCode
        ParticipationText
        Device
        Date
        StartTime
        LapsedTimeTest
        ExemptionReason
        PersonalDetailsChanged
        PossibleDuplicate
        DOBRange
        TestDisruptionList {
          TestDisruption {
            Event
          }
        }
        Adjustment {
          BookletType
          PNPCodelist {
            PNPCode
          }
        }
      }
      Test {
        TestContent {
          TestLevel
          TestDomain
        }
      }
    }
    Summary {
      Domain
      ParticipationCode
    }
  }
}
`

}




// 
// Fetch data & display for participation
// 
function createParticipationReport() {


    // order by level & alpha for test name
    sortParticipationData(participationData);

    // create ui elements
    // 
    createParticipationTitle();
    // filters
    createParticipationFilters();
    initFilterHandlers();

    // table
    createParticipationTableHeader();
    createParticipationTableBody(participationData);
    initParticipationTableHandler();

    // create report download link
    createParticipationDownloadLink();
    initParticipationDownloadLinkHandler();

}


// 
// sort data - not required for now
// 
function sortParticipationData(data) {

    // console.log(data);

    data.sort(function(a, b) {

        var compA = (a.EventInfos[0].Test.TestContent.TestLevel || '').toUpperCase() +
          (a.Student.FamilyName || '').toUpperCase() + 
          (a.Student.GivenName || '').toUpperCase() ;
        var compB = (b.EventInfos[0].Test.TestContent.TestLevel || '').toUpperCase() +
          (b.Student.FamilyName || '').toUpperCase() + 
          (b.Student.GivenName || '').toUpperCase() ;

        return (compA < compB) ? -1 : (compA > compB) ? 1 : 0;
    });

}

// 
// set the report title
// 
function createParticipationTitle() {
    var title = $("#report-title");
    title.empty();
    title.text("Student Participation");
}

// 
// create option selectors to filter data table
// 
function createParticipationFilters() {
    $("#report-filters").empty();
    // $("#report-filters").remove();
    $("#report-filters").append(createYearLevelFilter());
}

// 
// create main report table header
// 
function createParticipationTableHeader() {

    $("#report-table-hdr").empty();
    // $("#report-table-hdr").remove();

    var hdr = $("<tr></tr>");
    var hdr_row = $("<th>Level</th>" +
        "<th>Name</th>" +
        "<th>G and P</th>" +
        "<th>Numeracy</th>" +
        "<th>Reading</th>" +
        "<th>Spelling</th>" +
        "<th>Writing</th>");

    hdr.append(hdr_row);
    $("#report-table-hdr").append(hdr);

}


// 
// create the main tabular data diplay
// 
function createParticipationTableBody(data) {

    $("#report-table-body").empty();
    // $("#report-table-body").remove();

    $.each(data, function(index, pds) {

        var event = pds.EventInfos[0];
        var test = event.Test;
        // var summary = pds.Summary;

        // create map from summary info
        summary = {};
        jQuery.each(pds.Summary, function(index, summaryItem) {
            summary[summaryItem.Domain] = summaryItem.ParticipationCode
        });



        var $row = $("<tr/>");
        $row.append("<td>" + test.TestContent.TestLevel + "</td>" +
            // "<td>" + event.Event.PSI + "</td>" +
            "<td>" + pds.Student.GivenName + " " + pds.Student.FamilyName +  "</td>" +
            "<td class='domain'>" + hideNull(summary['Grammar and Punctuation']) + "</td>" +
            "<td class='domain'>" + hideNull(summary['Numeracy']) + "</td>" +
            "<td class='domain'>" + hideNull(summary['Reading']) + "</td>" +
            "<td class='domain'>" + hideNull(summary['Spelling']) + "</td>" +
            "<td class='domain'>" + hideNull(summary['Writing']) + "</td>"
        );
        $row.data("pdsdata", pds);
        $row.attr("yr-level", test.TestContent.TestLevel);
        $("#report-table-body").append($row);
        $row = null;
    });

    // colour any non-P values
    $("#report-table-body td.domain").each(function() {
        if ($(this).text() !== 'P') {
            $(this).css('background-color', '#e3f2fd');
        }
    });

}

// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initParticipationTableHandler() {
    // remove any existing handlers
    $('#report-content').off("click");
    // respond to table row selections
    $('#report-content').on("click", '#report-table-body tr', function(event) {
        var $row = event.currentTarget;
        var pds = jQuery.data($row, "pdsdata");
        createExtendedDataParticipation(pds);
        showExtendedData();
    });

}

// 
// when user selects line item in table
// show the extended data in a modal
// footer popup
// 
function createExtendedDataParticipation(pdsdata) {
    $("#ed-content").empty();

    $("#ed-modal").css("max-height", "85%");

    // 
    // title
    // 
    $("#ed-title").text('Student Participation');

    // 
    // student/school summary info 
    // 
    var event = pdsdata.EventInfos[0];
    var test = event.Test;

    var topRow = $("<div class='row'><div class='col s12'></div></div>");
    var student = pdsdata.Student;
    var school = pdsdata.School;

    topRow.append("<p>" +
        "Student: " + student.GivenName + " " + student.FamilyName +
        ", PSI: " + event.Event.PSI +
        ", Homegroup: " + student.Homegroup +
        ", Class: " + student.ClassGroup +
        ", Offline: " + student.OfflineDelivery +
        ", Ed. Support: " + student.EducationSupport +
        ", Home-Schooled: " + student.HomeSchooledStudent +
        "</p>");
    topRow.append("<p>" +
        "School: " + school.SchoolName +
        ", Sector: " + school.SchoolSector +
        ", Type: " + school.SchoolType +
        ", Independent: " + school.IndependentSchool +
        ", ACARA ID: " + school.ACARAId +
        ", Local ID: " + school.LocalId +
        "</p>");


    $("#ed-content").append(topRow);

    // 
    // participation summary
    // 
    var ptcpTable = $("<table class='bordered'></table>");
    var tblHdr = $("<thead></thead>");
    tblHdr.append("<tr>" +
        "<th>Domain</th>" +
        "<th>Test Date</th>" +
        "<th>Test Time</th>" +
        "<th>Duration</th>" +
        "<th>P. Code</th>" +
        "<th>Exemption</th>" +
        "<th>Disruptions</th>" +
        "<th>PNP Code</th>" +
        "<th>Duplicate</th>" +
        "<th>Details Changed</th>" +
        "</tr>");
    ptcpTable.append(tblHdr);

    var tblBody = $("<tbody></tbody>");

    var domains = ["Grammar and Punctuation", "Numeracy", "Reading", "Spelling", "Writing"];

    for (i = 0; i < domains.length; i++) {
        var eventInfo = getEventInfo(domains[i], pdsdata);
        var $row = $("<tr/>");
        $row.append("<td>" + domains[i] + "</td>" +
            "<td>" + hideNull(eventInfo.Event.Date) + "</td>" +
            "<td>" + hideNull(eventInfo.Event.StartTime) + "</td>" +
            "<td>" + hideNull(eventInfo.Event.LapsedTimeTest) + "</td>" +
            "<td>" + hideNull(eventInfo.Event.ParticipationCode) +
            " (" + hideNull(eventInfo.Event.ParticipationText) + ")</td>" +
            "<td>" + hideNull(eventInfo.Event.ExemptionReason) + "</td>" +
            "<td>" + unpackList(eventInfo.Event.TestDisruptionList) + "</td>" +
            "<td>" + unpackList(eventInfo.Event.Adjustment.PNPCodelist) + "</td>" +
            "<td>" + unpackBool(eventInfo.Event.PossibleDuplicate) + "</td>" +
            "<td>" + unpackBool(eventInfo.Event.PersonalDetailsChanged) + "</td>"
        );

        tblBody.append($row);
        $row = null;
    }

    ptcpTable.append(tblBody);
    $("#ed-content").append(ptcpTable);

}

// 
// helper function to retrieve event info for 
// given domain from within nested structure
// 
function getEventInfo(domainName, data) {

    ei = {};
    ei.Event = {};
    ei.Event.Adjustment = {};

    jQuery.each(data.EventInfos, function(index, eventInfo) {
        if (eventInfo.Test.TestContent.TestDomain == domainName) {
            ei = eventInfo;
        }
    });

    return ei;


}

// 
// add a download link to the csv report for this school
// 
function createParticipationDownloadLink() {

    $('#download-report').empty();
    $('#download-report').append("<a id='csv-download'>Download report as CSV file</a>");

}

// 
// handle the download
// 
function initParticipationDownloadLinkHandler() {

    var reportURL = "/naprr/downloadreport/" + currentASLId + "/schoolParticipation.csv";

    $('#csv-download').off("click");

    $('#csv-download').on("click", function(event) {
        window.location.href = reportURL;

    });


}
