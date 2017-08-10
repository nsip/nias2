// naprr_codeframe.js

// instantiate interaction listeners
$(document).ready(function()
{

    // initialise modals plugin
    $('.modal').modal();

    // initialise collapsible plugin
    $('.collapsible').collapsible();

    // create the report when button is clicked
    $("#btn-codeframe").on("click", function(event)
    {
        hideReport();
        // alert("Sorry, this report still under construction.");
        createCodeframeReport();
        showReport();
    });

});


// 
// gql query for data
// 
function codeframeQuery() {
    return `
query NAPData {
  codeframe_report {
    Test {
      TestID
      TestContent {
        TestLevel
        TestDomain
        TestYear
        StagesCount
        TestType
        LocalId
      }
    }
    Testlet {
      TestletContent {
        TestletName
        LocationInStage
        Node
        TestletMaximumScore
      }
    }
    Item {
      TestItemContent {
        NAPTestItemLocalId
        ItemName
        ItemType
        Subdomain
        WritingGenre
        ItemDescriptor
        ReleasedStatus
        MarkingType
        MultipleChoiceOptionCount
        CorrectAnswer
        MaximumScore
        ItemDifficulty
        ItemDifficultyLogit5
        ItemDifficultyLogit62
        ItemDifficultyLogit5SE
        ItemDifficultyLogit62SE
        ItemProficiencyBand
        ItemProficiencyLevel
        ExemplarURL
        ItemSubstitutedForList {
          SubstituteItem {
            SubstituteItemRefId
            LocalId
            PNPCodeList {
              PNPCode
            }
          }
        }
        ContentDescriptionList {
          ContentDescription
        }
        StimulusList {
          Stimulus {
            LocalId
            TextGenre
            TextType
            WordCount
            TextDescriptor
            Content
          }
        }
        NAPWritingRubricList {
          NAPWritingRubric {
            RubricType
            Descriptor
            ScoreList {
              Score {
                MaxScoreValue
                ScoreDescriptionList {
                  ScoreDescription {
                    ScoreValue
                    Descriptor
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`;
}

// 
// Fetch data & display for codeframe
// 
function createCodeframeReport()
{


    // order by level & alpha for test name
    sortCodeframeData(codeframeData);

    // create ui elements
    // 
    createCodeframeTitle();

    // filters
    createCodeframeFilters();
    initFilterHandlers();

    // table
    createCodeframeTableHeader();
    createCodeframeTableBody(codeframeData);
    initCodeframeTableHandler();

    // create report download link
    createCodeframeDownloadLink();
    initCodeframeDownloadLinkHandler();

}

// 
// sort data - by year level & test, testlet node & stage 
// 
// 
function sortCodeframeData(data)
{

    // console.log(data);

    data.sort(function(a, b)
    {

        var compA = a.Test.TestContent.TestLevel.toUpperCase() +
            a.Test.TestContent.TestDomain.toUpperCase() +
            a.Testlet.TestletContent.Node.toUpperCase() +
            a.Testlet.TestletContent.LocationInStage.toUpperCase() +
            a.Item.TestItemContent.Subdomain.toUpperCase();

        var compB = b.Test.TestContent.TestLevel.toUpperCase() +
            b.Test.TestContent.TestDomain.toUpperCase() +
            b.Testlet.TestletContent.Node.toUpperCase() +
            b.Testlet.TestletContent.LocationInStage.toUpperCase() +
            b.Item.TestItemContent.Subdomain.toUpperCase();


        return (compA < compB) ? -1 : (compA > compB) ? 1 : 0;
    });

}

// 
// set the report title
// 
function createCodeframeTitle()
{
    var title = $("#report-title");
    title.empty();
    title.text("NAPLAN Codeframe");
}

// 
// create option selectors to filter data table
// 
function createCodeframeFilters()
{
    $("#report-filters").empty();
    $("#report-filters").append(createYearLevelFilter());
    $("#report-filters").append(createDomainFilter());

}

// 
// create main report table header
// 
function createCodeframeTableHeader()
{

    $("#report-table-hdr").empty();

    var hdr = $("<tr></tr>");
    var hdr_row = $("<th>Level</th>" +
        "<th>Domain</th>" +
        "<th>Subdomain</th>" +
        "<th>Node</th>" +
        "<th>Testlet Name</th>" +
        "<th>Item Name</th>");

    hdr.append(hdr_row);
    $("#report-table-hdr").append(hdr);

}

// 
// create the main tabular data diplay
// 
function createCodeframeTableBody(data)
{

    $("#report-table-body").empty();
    $("#report-table-body tr").remove();

    $.each(data, function(index, cfds)
    {
        var $row = $("<tr/>");
        $row.append("<td>" + cfds.Test.TestContent.TestLevel + "</td>" +
            "<td>" + cfds.Test.TestContent.TestDomain + "</td>" +
            "<td>" + cfds.Item.TestItemContent.Subdomain + "</td>" +
            "<td>" + cfds.Testlet.TestletContent.Node + "</td>" +
            "<td>" + cfds.Testlet.TestletContent.TestletName + "</td>" +
            "<td>" + cfds.Item.TestItemContent.ItemName + "</td>"
        );
        $row.data("cfdsdata", cfds);
        $row.attr("yr-level", cfds.Test.TestContent.TestLevel);
        $row.attr("domain", cfds.Test.TestContent.TestDomain);
        $("#report-table-body").append($row);
        $row = null;
    });
}


// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initCodeframeTableHandler()
{
    // remove any existing handlers
    $('#report-content').off("click");
    // respond to table row selections
    $('#report-content').on("click", '#report-table-body tr', function(event)
    {
        var $row = event.currentTarget;
        var cfds = jQuery.data($row, "cfdsdata");
        createExtendedDataCodeframe(cfds);
        showExtendedData();
    });

}

// 
// when user selects line item in table
// show the extended data in a modal
// footer popup
// 
function createExtendedDataCodeframe(cfdsdata)
{
    $("#ed-content").empty();

    $("#ed-modal").css("max-height", "90%");

    $("#ed-title").text('Codeframe: ' +
        cfdsdata.Test.TestContent.TestDomain + ' Yr ' +
        cfdsdata.Test.TestContent.TestLevel);

    var topRow = $("<div class='row'><div class='col s12'></div></div>");
    topRow.append("Test Year:" + cfdsdata.Test.TestContent.TestYear +
        ", Test Type: " + cfdsdata.Test.TestContent.TestType +
        ", Test Stages: " + cfdsdata.Test.TestContent.StagesCount +
        ", Local Id:" + cfdsdata.Test.TestContent.LocalId);

    topRow.append("<br/>" + getTestletSummaryLine(cfdsdata.Testlet));

    topRow.append("<br/>Item: " + cfdsdata.Item.TestItemContent.ItemName +
        ", Subdomain: " + cfdsdata.Item.TestItemContent.Subdomain +
        ", Released: " + unpackBool(cfdsdata.Item.TestItemContent.ReleasedStatus) +
        ", Prof. Band|Level: " + cfdsdata.Item.TestItemContent.ItemProficiencyBand + "|" +
        cfdsdata.Item.TestItemContent.ItemProficiencyLevel +
        ", Type: " + cfdsdata.Item.TestItemContent.ItemType +
        ", MarkingType: " + cfdsdata.Item.TestItemContent.MarkingType
    );


    $("#ed-content").append(topRow);

    // 
    // scoring details
    // 
    var scoringRow = $("<div class='row'></div>");

    var scoringRowTitle = $("<div class='col s2'><h5>Scoring:</h5></div>");
    scoringRow.append(scoringRowTitle);

    var scoringRowContent = $("<div class='col s10'></div>");
    var scoringRowTable = $("<table></table>");

    scoringRowTable.append("<thead><tr>" +
        "<th>Max Score</th>" +
        "<th>Correct Answer</th>" +
        "<th>Difficulty</th>" +
        "<th>0.5 Prob. Logit</th>" +
        "<th>0.5 Prob. Logit Std. Err.</th>" +
        "<th>0.62 Prob. Logit</th>" +
        "<th>0.62 Prob. Logit Std. Err.</th>" +
        "</tr></thead>");

    scoringRowTable.append("<tbody><tr>" +
        "<td>" + cfdsdata.Item.TestItemContent.MaximumScore + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.CorrectAnswer + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.ItemDifficulty + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.ItemDifficultyLogit5 + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.ItemDifficultyLogit5SE + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.ItemDifficultyLogit62 + "</td>" +
        "<td>" + cfdsdata.Item.TestItemContent.ItemDifficultyLogit62SE + "</td>" +
        "</tr></tbody>");

    scoringRowContent.append(scoringRowTable);
    scoringRow.append(scoringRowContent);

    // console.log(scoringRow);


    // 
    // content info
    // 
    var contentinfoRow = $("<div class='row'></div>");

    var contentinfoRowTitle = $("<div class='col s2'><h5>Content:</h5></div>");
    contentinfoRow.append(contentinfoRowTitle);

    var contentinfoRowContent = $("<div class='col s10'></div>");
    var contentinfoRowTable = $("<table></table>");

    // exemplar & content descriptors
    contentinfoRowContent.append("<p>Exemplar: " +
        "<a class='exemplarLink' href='" + cfdsdata.Item.TestItemContent.ExemplarURL + "'" +
        "onclick='javascript:openExemplarLink(this.href);return false;'>" +
        cfdsdata.Item.TestItemContent.ExemplarURL + "</a>" +
        ", Content Descriptors: " +
        unpackList(cfdsdata.Item.TestItemContent.ContentDescriptionList) +
        "</p>");


    // 
    // stimulus list
    // 
    var stimuli = cfdsdata.Item.TestItemContent.StimulusList.Stimulus;
    // console.log(stimuli);

    if (stimuli != null)
    {
        contentinfoRowTable.append("<thead><tr>" +
            "<th>Text Genre</th>" +
            "<th>Text Type</th>" +
            "<th>Word Count</th>" +
            "<th>Text Descriptor</th>" +
            "<th>Content</th>" +
            "</tr></thead>");

        var tblBody = $("<tbody></tbody>");

        $.each(stimuli, function(index, stimulus)
        {
            tblBody.append("<tr>" +
                "<td>" + stimulus.TextGenre + "</td>" +
                "<td>" + stimulus.TextType + "</td>" +
                "<td>" + stimulus.WordCount + "</td>" +
                "<td>" + stimulus.TextDescriptor + "</td>" +
                "<td>" + "<a class='exemplarLink' href='" + stimulus.Content + "'" +
                "onclick='javascript:openExemplarLink(this.href);return false;'>" +
                stimulus.Content + "</a>" + "</td>" +
                "</tr>");
        });

        contentinfoRowTable.append(tblBody);
        contentinfoRowContent.append(contentinfoRowTable);

    }

    contentinfoRow.append(contentinfoRowContent);

    // 
    // Writing Rubrics
    // 
    var rubricRow = $("<div class='row'></div>");
    var rubrics = cfdsdata.Item.TestItemContent.NAPWritingRubricList.NAPWritingRubric;

    if (rubrics != null)
    {
        

        var rubricRowTitle = $("<div class='col s2'><h5>Writing Rubric:</h5></div>");
        rubricRow.append(rubricRowTitle);

        var rubricRowContent = $("<div class='col s10'></div>");
        var rubricRowTable = $("<table></table>");

        


        var rubricTypes = {};

        $.each(rubrics, function(index, rubric)
        {

            // rubricRowContent.append("<p>=============================</br>" +
            //     rubric.RubricType + "-->" + rubric.Descriptor + "-->" +
            //     "</p>");

            var rubricTypeArray;
            rubricTypeArray = rubricTypes[rubric.RubricType];
            if (rubricTypeArray == null)
            {
                rubricTypeArray = [];
                rubricTypes[rubric.RubricType] = rubricTypeArray;
            }
            rubricTypeArray.push(rubric);



            // add to array by type


            // var scores = rubric.ScoreList.Score;
            // console.log(scores);
            // $.each(scores, function(index, score)
            // {
            //     rubricRowContent.append("<p>" +
            //         "Max-Score ---->" + score.MaxScoreValue +
            //         "</p>");

            //     var scoreDescriptions = score.ScoreDescriptionList.ScoreDescription;
            //     console.log(scoreDescriptions);
            //     $.each(scoreDescriptions, function(index, scoreDescription)
            //     {
            //         rubricRowContent.append("<p>" +
            //             scoreDescription.Descriptor + "-->" + scoreDescription.ScoreValue +
            //             "</p>");

            //     });

            // });

        });

        console.log(rubricTypes);



        var tblHdr = $("<thead></thead>");
        var tblHdrRow = $("<tr></tr>");

        var tblBody = $("<tbody></tbody>");
        var tblBodyRow = $("<tr></tr>");

        $.each(rubricTypes, function(key, rubricTypeArray)
        {
            var hdrText = $("<th>" + key + "</th>");

            var typeInfoColumn = $("<td></td>");

            $.each(rubricTypeArray, function(index, rubric)
            {
                var scores = rubric.ScoreList.Score;

                // var scoreList = $("<ul class='collapsible' data-collapsible='accordion'></ul>");
                var scoreList = $("<ul></ul>");
                typeInfoColumn.append(scoreList);
                var slList = $("<li></li>");
                scoreList.append(slList);
                // var slContent = slList.append("<div class='collapsible-header'>" + rubric.Descriptor + "</div>");

                // $.each(scores, function(index, score)
                // {


                var slContent = slList.append("<div>Max: " + scores[0].MaxScoreValue + "</div>");


                // var scoreDescriptions = score.ScoreDescriptionList.ScoreDescription;
                // console.log(scoreDescriptions);

                // $.each(scoreDescriptions, function(index, scoreDescription)
                // {

                // slContent.append("<div class='collapsible-body'>" +
                // scoreDescription.Descriptor + ": " + scoreDescription.ScoreValue +
                // "</div>");

                // });

                // });


                tblHdrRow.append(hdrText);
                tblHdr.append(tblHdrRow);

                tblBodyRow.append(typeInfoColumn);
                // tblBodyRow.append(scoreList);
                tblBody.append(tblBodyRow);
                return false;

            });

        });


        rubricRowTable.append(tblHdr);
        rubricRowTable.append(tblBody);

        console.log(rubricRowTable);

        rubricRowContent.append(rubricRowTable);

        rubricRow.append(rubricRowContent);

    }


    $("#ed-content").append(contentinfoRow);
    $("#ed-content").append(rubricRow);
    $("#ed-content").append(scoringRow);

    // 
    // add score bands display
    // 
    // $("#ed-content").append(createTestBandsDisplay(cfdsdata));

}


// 
// add a download link to the csv report for this school
// 
function createCodeframeDownloadLink()
{

    $('#download-report').empty();
    $('#download-report').append("<a id='csv-download'>Download report as CSV file</a>");

}

// 
// handle the download
// 
function initCodeframeDownloadLinkHandler()
{

    var reportURL = "/naprr/downloadreport/codeframe";

    $('#csv-download').off("click");

    $('#csv-download').on("click", function(event)
    {
        window.location.href = reportURL;

    });


}
