var txID = "unknown";
var txTotal = 0;
var prgTotal = 0;
var fName = "";
var currentFile;

// set any initial ui state
// disable download error report until there are some results
window.onload = function() {
    $("#fetch2").prop('disabled', true);
}

// handler for downloading error report
$(function() {
    $("#fetch2").click(function(e) {
        // fname = fname.replace(/.*[\/\\]/, '');
        // console.log(fName)
        var url = "/naplan/reg/report/" + txID + "/" + fName;
        // console.log(url)
        window.location.href = url;
    });
});


// handler for drag-drop fileuploader for validation
// accepts files through browsing or drag-drop and submits for
// processing immediately
$(function() {

    var ul = $('#upload ul');

    $('#drop a').click(function() {
        // Simulate a click on the file input button
        // to show the file browser dialog
        $(this).parent().find('input').click();
    });

    // Initialize the jQuery File Upload plugin
    $('#upload').fileupload({

        sequentialUploads: true,

        // This element will accept file drag/drop uploading
        dropZone: $("#drop"),

        // don't listen for other file input forms on this ui
        fileInput: $("#vfileInput"),

        // This function is called when a file is added to the queue;
        // either via the browse button, or via drag/drop:
        add: function(e, data) {
            
            txID = "";
            txTotal = 0;
            prgTotal = 0;

            $("#processed").empty();
            $("#progress").empty();
            $("#fetch2").prop('disabled', true)

            var tpl = $('<li class="working"><input type="text" value="0" data-width="48" data-height="48"' +
                ' data-fgColor="#0788a5" data-readOnly="1" data-bgColor="#3e4043" /><p></p><span></span></li>');

            // Append the file name and file size
            tpl.find('p').text(data.files[0].name)
                .append('<i>' + formatFileSize(data.files[0].size) + '</i>');
            fName = data.files[0].name;
	    // csv? then Get first line of the file, to retrieve its header
	    if(fName.endsWith(".csv") || fName.endsWith(".CSV")) {
	    var reader = new FileReader();
	    var blob = data.files[0].slice(0, 1000);
	    reader.readAsText(blob);
	    reader.onloadend = function(evt) {
		    alert(evt.target.result);
		    $.post("/naplan/reg/sanityCSV", {file : evt.target.result}, function(data1, status){
			    if(data1["Error"]){
			            alert("Data: " + data1["Error"] );
			    }
				        });
					    }

	    }

            // Add the HTML to the UL element
            data.context = tpl.appendTo(ul);

            // Initialize the knob plugin
            tpl.find('input').knob();

            // Listen for clicks on the cancel icon
            tpl.find('span').click(function() {

                if (tpl.hasClass('working')) {
                    jqXHR.abort();
                }

                tpl.fadeOut(function() {
                    tpl.remove();
                });

            });

            currentFile = data.files[0];

            // Automatically upload the file once it is added to the queue
            var jqXHR = data.submit();
        },

        progress: function(e, data) {

            // Calculate the completion percentage of the upload
            var progress = parseInt(data.loaded / data.total * 100, 10);

            // Update the hidden input field and trigger a change
            // so that the jQuery knob plugin knows to update the dial
            data.context.find('input').val(progress).change();

            if (progress == 100) {
                data.context.removeClass('working');
            }
        },

        fail: function(e, data) {
            // Something has gone wrong!
            data.context.addClass('error');
        },

        done: function(e, data) {
            // console.log(data)
            fName = data.files[0].name;

            txTotal = data.result.Records;
            txID = data.result.TxID;
            $("#progress").html("<h5>Number of records: " + txTotal + "</h5>");
            monitorProgress();
        }

    });


    // Prevent the default action when a file is dropped on the window
    $(document).on('drop dragover', function(e) {
        e.preventDefault();
    });

    // Helper function that formats the file sizes
    function formatFileSize(bytes) {
        if (typeof bytes !== 'number') {
            return '';
        }

        if (bytes >= 1000000000) {
            return (bytes / 1000000000).toFixed(2) + ' GB';
        }

        if (bytes >= 1000000) {
            return (bytes / 1000000).toFixed(2) + ' MB';
        }

        return (bytes / 1000).toFixed(2) + ' KB';
    }

});



function monitorProgress() {
    ref = "/naplan/reg/status/" + txID
    monitorSource = new EventSource(ref);
    monitorSource.onmessage = function(event) {
        var prgdata = JSON.parse(event.data);
        // console.log(prgdata.Total)
        prgTotal = prgdata.Total
        $("#processed").html("<h5>Records processed: " + prgTotal + "</h5>")
        if (prgTotal == txTotal) {
            monitorSource.close();
            $("#fetch2").prop('disabled', false);
            renderAnalysis();
        }
    };
}

function renderAnalysis() {

    ref = "/naplan/reg/results/" + txID;
    d3.json(ref, function(error, data) {

        var errorsBarChart = dc.barChart("#errors-chart");
        var errorsByTypeChart = dc.rowChart('#errors-by-type-chart');
        var dataTable = dc.dataTable('.dc-data-table');
        var verrorsCount = dc.dataCount('.dc-data-count');
        var errorData = data;

        // normalize/parse data so dc can correctly sort & bin them
        errorData.forEach(function(d) {
            d.originalLine = +d.originalLine;
        });
        // console.log(errorData);

        var ndx = crossfilter(errorData);
        var all = ndx.groupAll();

        var lineDim = ndx.dimension(function(d) {
            return d.originalLine;
        });

        var typeDim = ndx.dimension(function(d) {
            return d.validationType;
        });
        var validationTypesGroup = typeDim.group();

        var allDim = ndx.dimension(function(d) {
            return d;
        });

        // var countPerLine = lineDim.group().reduceCount(function(d) {
        //     return d.originalLine;
        // });

        var lineGroup = lineDim.group();

        // var countPerLine = lineDim.group().reduceSum(function(d) {return d.OriginalLine;});

        errorsBarChart
            .width(350)
            .height(190)
            .margins({
                top: 20,
                right: 0,
                bottom: 0,
                left: 0
            })
            .gap(1)
            .x(d3.scale.linear().domain([0, 200000]))
            .elasticX(true)
            .elasticY(true)
            .dimension(lineDim)
            .group(lineGroup);

        errorsBarChart.dimension(lineDim);

        errorsByTypeChart
            .width(350)
            .height(190)
            .margins({
                top: 20,
                left: 10,
                right: 10,
                bottom: 20
            })
            .group(validationTypesGroup)
            .dimension(typeDim)
            .title(function(d) {
                return d.value;
            })
            .elasticX(true)
            .xAxis().ticks(4);


        verrorsCount
            .dimension(ndx)
            .group(all);


        dataTable
        // .width(900)
        // .height(800)
            .dimension(allDim)
            .group(function(d) {
                // return 'dc.js insists on putting a row here so I remove it using JS';
                // return d.originalLine;
                return 'Errors ordered by original file line number (table shows first 100 errors)'
            })
            .size(100)
            .columns([
                // function (d) { return d.txID; },
                function(d) {
                    return d.originalLine;
                },
                function(d) {
                    return d.validationType;
                },
                function(d) {
                    return d.errField;
                },
                function(d) {
                    return d.description;
                }
            ])
            .sortBy(function(d) {
                return d.originalLine;
            })
            .order(d3.ascending)
            .on('renderlet', function(table) {
                // each time table is rendered remove nasty extra row dc.js insists on adding
                // table.select('tr.dc-table-group').remove();
                table.selectAll('.dc-table-group').classed('info', true);
            });

        dc.renderAll();

        // dc.redrawAll();


    });


}
