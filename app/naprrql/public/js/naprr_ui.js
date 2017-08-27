// naprr_ui.js

// helper routines for the main ui



// 
// gets schoolinfo 1-line summary for given acaraid
// 
function getSchoolInfoSummaryLine(acaraid) {

    var sinfo = schoolinfoData;

    return $("<p>" +
        "School: " + sinfo.SchoolName +
        ", Sector: " + sinfo.SchoolSector +
        ", Type: " + sinfo.SchoolType +
        ", Independent: " + sinfo.IndependentSchool +
        ", ACARA ID: " + sinfo.ACARAId +
        ", Local ID: " + sinfo.LocalId +
        "</p>");
}


// 
// builds singl line descriptor for a tetlet
// 
function getTestletSummaryLine(testlet) {

    return "Testlet: " + testlet.TestletContent.TestletName +
        ", Location in Stage: " + testlet.TestletContent.LocationInStage +
        ", Node: " + testlet.TestletContent.Node +
        ", Max. Score: " + testlet.TestletContent.TestletMaximumScore;


}

// 
// initialise the school selector that initiates the rest of the
// data downloads
// 
function initSchoolChooserQL() {

    var school_names = [];
    var sd = {};

    var query = `query NAPdata{school_details{SchoolName ACARAId}}`
    // console.log(query)
    $.post("/graphql", { Query: query }, function(response) {
        // console.log(response.data.school_details);
        $.each(response.data, function(index, school_list) {
            $.each(school_list, function(index, schooldetails) {
                var combi_name = schooldetails.ACARAId + " - " + schooldetails.SchoolName;
                sd[combi_name] = null;
            });
        });
        // console.log(sd);
        $('input.autocomplete').autocomplete({
            data: sd,
            limit: 20, // The max amount of results that can be shown at once. Default: Infinity.
        });

    });
}

// 
// listen for changes to the chosen school,
// once selected pull down all datasets
// 
function initSchoolChooserHandlerQL() {


    // 
    // background load all datasets when user selects a new school
    // 
    $("input.autocomplete").on("change", function(event) {
        // 
        // set report buttons disabled until data is downloaded
        // 
        $("#btn-scoresummary").addClass("disabled");
        $("#btn-domainscores").addClass("disabled");
        $("#btn-participation").addClass("disabled");
        $("#btn-codeframe").addClass("disabled");
        
        // get the current selected school
        var selection = $('input.autocomplete').val();
        // console.log(selection);

        var delimiter_pos = selection.indexOf(" - ");
        // console.log(delimiter_pos);

        // make sure we ignore false selections
        if (delimiter_pos < 0) {
            return
        }

        // strip out the acara id
        var acaraid = selection.slice(0, delimiter_pos);
        // console.log(">" + acaraid + "<");
        currentASLId = acaraid; // store globally for reuse

        // 
        // get the schoolinfo object for the selected school
        // 
        var query = schoolInfoQuery();
        var xhrSI = new XMLHttpRequest();
        xhrSI.responseType = 'json';
        xhrSI.open("POST", "/graphql");
        xhrSI.setRequestHeader("Content-Type", "application/json");
        xhrSI.setRequestHeader("Accept", "application/json");
        xhrSI.onload = function() {
            // console.log('school info data returned:', xhrSI.response);
            schoolinfoData = {};
            schoolinfoData = xhrSI.response.data.school_infos_by_acaraid[0];
        }
        xhrSI.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));

        // 
        // get the student information for the selected school
        // 
        var query = studentPersonalQuery();
        var xhrSP = new XMLHttpRequest();
        xhrSP.responseType = 'json';
        xhrSP.open("POST", "/graphql");
        xhrSP.setRequestHeader("Content-Type", "application/json");
        xhrSP.setRequestHeader("Accept", "application/json");
        xhrSP.onload = function() {
            // console.log('student data returned:', xhrSP.response);
            studentPersonalData = [];
            studentPersonalData = xhrSP.response.data.students_by_school;
        }
        xhrSP.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));



        // 
        // get the score summaries for the selected school
        // 
        var query = scoreSummaryQuery();
        var xhrSS = new XMLHttpRequest();
        xhrSS.responseType = 'json';
        xhrSS.open("POST", "/graphql");
        xhrSS.setRequestHeader("Content-Type", "application/json");
        xhrSS.setRequestHeader("Accept", "application/json");
        xhrSS.onload = function() {
            // console.log('score summary data returned:', xhrSS.response);
            scoresummaryData = [];
            scoresummaryData = xhrSS.response.data.score_summary_report_by_school;
            $("#btn-scoresummary").removeClass("disabled");
            hideReport();
            createScoreSummaryReport();
            showReport();
        }
        xhrSS.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));

        // 
        // get the domain scores
        // 
        var query = domainScoresQuery();
        var xhrDS = new XMLHttpRequest();
        xhrDS.responseType = 'json';
        xhrDS.open("POST", "/graphql");
        xhrDS.setRequestHeader("Content-Type", "application/json");
        xhrDS.setRequestHeader("Accept", "application/json");
        xhrDS.onload = function() {
            // console.log('domain scores data returned:', xhrDS.response);
            domainscoresData = [];
            domainscoresData = xhrDS.response.data.domain_scores_report_by_school;
            $("#btn-domainscores").removeClass("disabled");
        }
        xhrDS.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));



        // 
        // get participation data
        // 
        var query = participationQuery();
        var xhrPD = new XMLHttpRequest();
        xhrPD.responseType = 'json';
        xhrPD.open("POST", "/graphql");
        xhrPD.setRequestHeader("Content-Type", "application/json");
        xhrPD.setRequestHeader("Accept", "application/json");
        xhrPD.onload = function() {
            // console.log('participation data returned:', xhrPD.response);
            participationData = [];
            participationData = xhrPD.response.data.participation_report_by_school;
            $("#btn-participation").removeClass("disabled");
        }
        xhrPD.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));


        // 
        // get codeframe data
        // 
        var query = codeframeQuery();
        var xhrCF = new XMLHttpRequest();
        xhrCF.responseType = 'json';
        xhrCF.open("POST", "/graphql");
        xhrCF.setRequestHeader("Content-Type", "application/json");
        xhrCF.setRequestHeader("Accept", "application/json");
        xhrCF.onload = function() {
            // console.log('codeframe data returned:', xhrCF.response);
            codeframeData = [];
            codeframeData = xhrCF.response.data.codeframe_report;
            $("#btn-codeframe").removeClass("disabled");
        }
        xhrCF.send(JSON.stringify({
            query: query,
            variables: { acaraIDs: [currentASLId] },
        }));



    });

}

// show the extended data modal form
function showExtendedData() {
    $('#ed-modal').modal('open');
}

// clear the ui
function hideReport() {
    $("#report-container").addClass("hide");
}

// display the selected report
function showReport() {
    $("#report-container").removeClass("hide");
}

// 
// helper routine to render list object contents
// 
function unpackList(list) {

    if (list == null) {
        return "";
    }

    content = "";

    jQuery.each(list, function(index, item) {
        if (item != null) {
            $.each(item, function(key, val) {
                if (jQuery.isPlainObject(val)) {

                    // console.log("isObject: ", val);

                    $.each(val, function(k, v) {
                        content += v + " ";
                    });
                } else {
                    content += val + " ";
                }

            });
        }
    });


    return content;
}

// 
// helper to render item stimuli
// 
function unpackStimulusList(stimList) {
    console.log(stimList);

    var null_response = $("<p>not supplied</p>");

    if (stimList == null) {
        return null_response;
    }

    if (stimList.Stimulus == null || stimList.Stimulus.length < 1) {
        return null_response;
    }

    var response = $("<p></p>");

    jQuery.each(stimList.Stimulus, function(index, stimulus) {
        if (stimulus == null) {
            return false;
        }

        console.log(stimulus);

        response.append("Genre: " + stimulus.TextGenre + "<br/>");
        response.append("Type: " + stimulus.TextType + "<br/>");
        response.append("Words: " + stimulus.WordCount + "<br/>");
        response.append("Descriptor: " + stimulus.TextDescriptor + "<br/>");
        response.append("Content: " + stimulus.Content + "<br/>");
    });

    return response;

}

// 
// helper to manage null field display
// 
function hideNull(content) {
    var s = "";
    if (content == null) {
        return s;
    }
    return content;
}

// 
// helper to render booleans in ui
// 
function unpackBool(val) {
    if (val == null) {
        return "";
    }

    if (val === true) {
        return "yes";
    }

    if (val == "true") {
        return "yes";
    }

    return "no";
}


// 
// filter used on most reports
// 
function createYearLevelFilter() {

    var filterRow = $("<div class='row'></div>");

    // all year levels selector
    var input = $("<p>" +
        "<input name='yrlvlFilter' type='radio' id='all' level='all' checked />" +
        "<label for='all'>All</label>" +
        "</p>"
    );
    var col = $("<div class='col s2'></div>");
    col.append(input);
    filterRow.append(col);

    // year levels selector
    var yearLevels = ["3", "5", "7", "9"];

    for (i = 0; i < yearLevels.length; i++) {
        var inputCol = $("<div class='col s2'></div>");
        var inputPara = $("<p></p>");
        var inputRadio = $("<input name='yrlvlFilter' type='radio'></input>");
        inputRadio.attr("id", "yr" + yearLevels[i]);
        inputRadio.attr("level", yearLevels[i]);
        var inputLabel = $("<label></label>");
        inputLabel.text("Yr " + yearLevels[i]);
        inputLabel.attr("for", inputRadio.attr("id"));
        inputPara.append(inputRadio);
        inputPara.append(inputLabel);
        inputCol.append(inputPara);
        filterRow.append(inputCol);
    }

    return filterRow;

}

function createDomainFilter() {


    var filterRow = $("<div class='row'></div>");

    // all domains selector
    var input = $("<p>" +
        "<input name='domainFilter' type='radio' id='all-domains' domain='all' checked />" +
        "<label for='all-domains'>All</label>" +
        "</p>"
    );
    var col = $("<div class='col s2'></div>");
    col.append(input);
    filterRow.append(col);

    // domains selector
    var domains = ["Grammar and Punctuation", "Numeracy", "Reading", "Spelling", "Writing"];

    for (i = 0; i < domains.length; i++) {
        var inputCol = $("<div class='col s2'></div>");
        var inputPara = $("<p></p>");
        var inputRadio = $("<input name='domainFilter' type='radio'></input>");
        inputRadio.attr("id", domains[i]);
        inputRadio.attr("domain", domains[i]);
        var inputLabel = $("<label></label>");
        inputLabel.text(domains[i]);
        inputLabel.attr("for", inputRadio.attr("id"));
        inputPara.append(inputRadio);
        inputPara.append(inputLabel);
        inputCol.append(inputPara);
        filterRow.append(inputCol);
    }

    return filterRow;

}

// 
// start listeners for year level filters
// 
function initFilterHandlers() {

    var fltrLevel = "all";
    var fltrDomain = "all";

    // handler for yr-level selectors
    $("input[name='yrlvlFilter']").change(function() {
        // console.log($(this).attr('level'));
        fltrLevel = $(this).attr('level');
        applyFilters(fltrLevel, fltrDomain);
    })

    // handler for domain selectors
    $("input[name='domainFilter']").change(function() {
        // console.log($(this).attr('domain'));
        fltrDomain = $(this).attr('domain');
        applyFilters(fltrLevel, fltrDomain);
    })

}

// 
// sort table output according to filters
// 
function applyFilters(fltrLevel, fltrDomain) {
    // console.log("level: ", fltrLevel);
    // console.log("domain: ", fltrDomain);

    var $rows = $('#report-table-body tr');

    $rows.show().filter(function() {

        if ((fltrDomain == "all") && (fltrLevel == "all")) {
            return false;
        }

        if (($(this).attr('yr-level') == fltrLevel) &&
            ($(this).attr('domain') == fltrDomain)) {
            return false;
        }

        if (($(this).attr('yr-level') == fltrLevel) &&
            (fltrDomain == "all")) {
            return false;
        }

        if (($(this).attr('domain') == fltrDomain) &&
            (fltrLevel == "all")) {
            return false;
        }

        return true;

    }).hide();


}



// 
// for any report object that contains a Test member
// the domain bands can be displayed in a consistent format
// 
function createTestBandsDisplay(data) {

    var bottomRow = $("<div class='row'></div>");
    var brTitle = $("<div class='col s2'></div>");
    brTitle.append("<h5>Score Bands:</h5>");
    bottomRow.append(brTitle);


    var brTable = $("<div class='col s10'></div>");
    var bandsTable = $("<table></table>");

    var hdr = $("<thead><tr></tr></thead>");
    // var hdr_row = $("<tr/>");
    for (i = 0; i < 10; i++) {
        hdr.append("<th>" + (i + 1) + "</th>");
    };

    var body = $("<tbody/>")
    var bandsTableRow = $("<tr/>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band1Lower +
        " - " + data.Test.TestContent.DomainBands.Band1Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band2Lower +
        " - " + data.Test.TestContent.DomainBands.Band2Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band3Lower +
        " - " + data.Test.TestContent.DomainBands.Band3Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band4Lower +
        " - " + data.Test.TestContent.DomainBands.Band4Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band5Lower +
        " - " + data.Test.TestContent.DomainBands.Band5Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band6Lower +
        " - " + data.Test.TestContent.DomainBands.Band6Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band7Lower +
        " - " + data.Test.TestContent.DomainBands.Band7Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band8Lower +
        " - " + data.Test.TestContent.DomainBands.Band8Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band9Lower +
        " - " + data.Test.TestContent.DomainBands.Band9Upper + "</td>");
    bandsTableRow.append("<td>" + data.Test.TestContent.DomainBands.Band10Lower +
        " - " + data.Test.TestContent.DomainBands.Band10Upper + "</td>");
    body.append(bandsTableRow);

    bandsTable.append(hdr);
    bandsTable.append(body);
    brTable.append(bandsTable);
    bottomRow.append(brTable);

    return bottomRow;
}

// 
// open examplars in a new window
// 
function openExemplarLink(url) {
    window.open(url, '_blank');
    window.focus();
}