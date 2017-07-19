// naprr_ui.js

// helper routines for the main ui



// 
// gets schoolinfo 1-line summary for given acaraid
// 
function getSchoolInfoSummaryLine(acaraid)
{

    var sinfo = {};

    // data can be found in participation info
    $.each(participationData, function(index, pds)
    {
        var school = pds.School;

        if (school.ACARAId == acaraid)
        {
            sinfo = school;
            return false;
        }

    });

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
// get student info 1-line summary for given psi
// 
function getStudentInfoSummaryLine(psi)
{

    var sp = {};
    var ei = {};

    // data can be found in participation info
    $.each(participationData, function(index, pds)
    {
	if(!(pds == null || pds.EventInfos == null || pds.EventInfos[0] == null )) {

        ei = pds.EventInfos[0];

        if (ei.Event.PSI == psi)
        {
            var student = pds.Student;
            sp = student;
            return false;
        }
	}

    });

    if (!ei) { return $("<p>N/A</p>");}

    return $("<p>" +
        "Student: " + sp.GivenName + " " + sp.FamilyName +
        ", PSI: " + ei.Event.PSI +
        ", Homegroup: " + sp.Homegroup +
        ", Class: " + sp.ClassCode +
        ", Offline: " + sp.OfflineDelivery +
        ", Ed. Support: " + sp.EducationSupport +
        ", Home-Schooled: " + sp.HomeSchooledStudent +
        "</p>");

}

// 
// builds singl line descriptor for a tetlet
// 
function getTestletSummaryLine(testlet)
{

    return "Testlet: " + testlet.TestletContent.TestletName +
        ", Location in Stage: " + testlet.TestletContent.LocationInStage +
        ", Node: " + testlet.TestletContent.Node +
        ", Max. Score: " + testlet.TestletContent.TestletMaximumScore;


}


// 
// set up the main school selector
// 
function initSchoolChooser()
{
    var school_names = [];
    var sd = {};

    $.get("naprr/schooldetails", function(data, status)
    {
        // console.log(data);
        // console.log(status);

        $.each(data, function(index, school_list)
        {
            $.each(school_list, function(index, schooldetails)
            {
                var combi_name = schooldetails.ACARAId + " - " + schooldetails.SchoolName;
                sd[combi_name] = null;
            });
        });

        // console.log(sd);

        $('input.autocomplete').autocomplete(
        {
            data: sd,
            limit: 20, // The max amount of results that can be shown at once. Default: Infinity.
        });

    });

}

// 
// listen for changes to the chosen school,
// once selected pull down all datasets
// 
function initSchoolChooserHandler()
{

    // 
    // background load all datasets when user selects a new school
    // 
    $("input.autocomplete").on("change", function(event)
    {
        // get the current selected school
        var selection = $('input.autocomplete').val();
        // console.log(selection);

        var delimiter_pos = selection.indexOf(" - ");
        // console.log(delimiter_pos);

        // make sure we ignore false selections
        if (delimiter_pos < 0)
        {
            return
        }

        // strip out the acara id
        var acaraid = selection.slice(0, delimiter_pos);
        // console.log(">" + acaraid + "<");
        currentASLId = acaraid; // store globally for reuse

        // fetch the score summary first as is smallest dataset
        $.get("/naprr/scoresummary/" + acaraid, function(data, status)
        {
            scoresummaryData = [];
            scoresummaryData = data;
            console.log("score summary data downloaded for " + acaraid +
                " elements: " + scoresummaryData.length);

            if (debug)
            {
                console.log(scoresummaryData);
            }

            // display on screen while other reports download
            hideReport();
            createScoreSummaryReport();
            showReport();
        });

        // get domain scores
        $.get("/naprr/domainscores/" + acaraid, function(data, status)
        {
            domainscoresData = [];
            domainscoresData = data;
            console.log("domain scores data downloaded for " + acaraid +
                " elements: " + domainscoresData.length);

            if (debug)
            {
                console.log(domainscoresData);
            }

        });

        $.get("/naprr/participation/" + acaraid, function(data, status)
        {
            participationData = [];
            participationData = data;
            console.log("participation data downloaded for " + acaraid +
                " elements: " + participationData.length);

            if (debug)
            {
                console.log(participationData);
            }

        });

        $.get("/naprr/schoolinfo/" + acaraid, function(data, status)
        {
            schoolinfoData = {};
            schoolinfoData = data;
            console.log("school info data downloaded for " + acaraid);

            if (debug)
            {
                console.log(schoolinfoData);
            }


        });

        $.get("/naprr/codeframe", function(data, status)
        {
            codeframeData = {};
            codeframeData = data;
            console.log("codeframe data downloaded.");

            if (debug)
            {
                console.log(codeframeData);
            }

        });

    });

}

// show the extended data modal form
function showExtendedData()
{
    $('#ed-modal').modal('open');
}

// clear the ui
function hideReport()
{
    $("#report-container").addClass("hide");
}

// display the selected report
function showReport()
{
    $("#report-container").removeClass("hide");
}

// 
// helper routine to render list object contents
// 
function unpackList(list)
{

    if (list == null)
    {
        return "";
    }

    content = "";

    jQuery.each(list, function(index, item)
    {
        if (item != null)
        {
            $.each(item, function(key, val)
            {
                if (jQuery.isPlainObject(val))
                {

                    // console.log("isObject: ", val);

                    $.each(val, function(k, v)
                    {
                        content += v + " ";
                    });
                }
                else
                {
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
function unpackStimulusList(stimList)
{
    console.log(stimList);

    var null_response = $("<p>not supplied</p>");

    if (stimList == null)
    {
        return null_response;
    }

    if (stimList.Stimulus == null || stimList.Stimulus.length < 1)
    {
        return null_response;
    }

    var response = $("<p></p>");

    jQuery.each(stimList.Stimulus, function(index, stimulus)
    {
        if (stimulus == null)
        {
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
// helper to render booleans in ui
// 
function unpackBool(val)
{
    if (val == null)
    {
        return "";
    }

    if (val === true)
    {
        return "yes";
    }

    if (val == "true")
    {
        return "yes";
    }

    return "no";
}


// 
// filter used on most reports
// 
function createYearLevelFilter()
{

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

    for (i = 0; i < yearLevels.length; i++)
    {
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

function createDomainFilter()
{


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

    for (i = 0; i < domains.length; i++)
    {
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
function initFilterHandlers()
{

    var fltrLevel = "all";
    var fltrDomain = "all";

    // handler for yr-level selectors
    $("input[name='yrlvlFilter']").change(function()
    {
        // console.log($(this).attr('level'));
        fltrLevel = $(this).attr('level');
        applyFilters(fltrLevel, fltrDomain);
    })

    // handler for domain selectors
    $("input[name='domainFilter']").change(function()
    {
        // console.log($(this).attr('domain'));
        fltrDomain = $(this).attr('domain');
        applyFilters(fltrLevel, fltrDomain);
    })

}

// 
// sort table output according to filters
// 
function applyFilters(fltrLevel, fltrDomain)
{
    // console.log("level: ", fltrLevel);
    // console.log("domain: ", fltrDomain);

    var $rows = $('#report-table-body tr');

    $rows.show().filter(function()
    {

        if ((fltrDomain == "all") && (fltrLevel == "all"))
        {
            return false;
        }

        if (($(this).attr('yr-level') == fltrLevel) &&
            ($(this).attr('domain') == fltrDomain))
        {
            return false;
        }

        if (($(this).attr('yr-level') == fltrLevel) &&
            (fltrDomain == "all"))
        {
            return false;
        }

        if (($(this).attr('domain') == fltrDomain) &&
            (fltrLevel == "all"))
        {
            return false;
        }

        return true;

    }).hide();


}



// 
// for any report object that contains a Test member
// the domain bands can be displayed in a consistent format
// 
function createTestBandsDisplay(data)
{

    var bottomRow = $("<div class='row'></div>");
    var brTitle = $("<div class='col s2'></div>");
    brTitle.append("<h5>Score Bands:</h5>");
    bottomRow.append(brTitle);


    var brTable = $("<div class='col s10'></div>");
    var bandsTable = $("<table></table>");

    var hdr = $("<thead><tr></tr></thead>");
    // var hdr_row = $("<tr/>");
    for (i = 0; i < 10; i++)
    {
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
function openExemplarLink(url)
{
    window.open(url, '_blank');
    window.focus();
}
