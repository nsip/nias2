
// globals
var txID;
var txTotal;
var display_results_once = false;


// instantiate interaction listeners
$(document).ready(function()
    {

      hideResults();
      hideProgress();
      $('.modal').modal();
      // $('.modal').modal();
      $('.modal').modal('open');


      // handler for downloading error report
      $(function()
          {
            $("#csv-download,#csv-download2").click(function(e)
                {
                  CsvDownload();
                });
          });


      // handler for validation '>' button
      $("#validate").click(function()
          {
            Validate();
          });

      // handler for csv - sifxml conversion
      $("#xml-convert,#xml-convert2").click(function()
          {
            XmlConvert();
          });


    });

// get the validation results in csv format
function CsvDownload()
{
  var filesList = $('#upload').prop('files');
  if (filesList.length < 1)
  {
    var toastContent = $('<h5>No file selected</h5>');
    Materialize.toast(toastContent, 3000, 'rounded green');
    return
  }

  if (!txID)
  {
    var toastContent = $('<h5>No results to download</h5>');
    Materialize.toast(toastContent, 3000, 'rounded green');
    return

  }

  var fname = filesList[0].name
    var url = "/naplan/reg/report/" + txID + "/" + fname;

  window.location.href = url;


}

// run the validation
function Validate()
{

  var filesList = $('#upload').prop('files');
  if (filesList.length < 1)
  {
    var toastContent = $('<h5>No file selected</h5>');
    Materialize.toast(toastContent, 4000, 'rounded green');
    return
  }

  // reset ui elements for new validation run
  hideResults();
  hideProgress(); //clear any previous data
  showProgress(filesList[0].size);
  display_results_once = false;

  Materialize.toast('Sending file to server...', 4000, 'rounded')


    $('#validateform').prop('action' , '/naplan/reg/validate');

  $('#upload').fileupload(
      {
        autoUpload: false
      });

  $('#upload').bind('fileuploaddone', function(e, data)
      {
        txTotal = data.result.Records
          txID = data.result.TxID
          console.log(txID + ": " + txTotal)
          $('#upload').fileupload('destroy');
        startStreamSocket();
        Materialize.toast('Analysis results generating...', 4000, 'rounded')

      });
  $('#upload').fileupload('send',
      {
        sequentialUploads: true,
        files: filesList[0]
      });


}

// convert the input file from csv to xml
function XmlConvert()
{
  var filesList = $('#upload').prop('files');
  if (filesList.length < 1)
  {
    var toastContent = $('<h5>No file selected</h5>');
    Materialize.toast(toastContent, 3000, 'rounded green');
    return
  }

  $('#validateform').prop('action' , '/naplan/reg/convert');
  /*
     $('#upload').fileupload(
     {   
     autoUpload: false
     });
     $('#upload').fileupload('send',
     {
     sequentialUploads: true,
     files: filesList[0]
     });
     */
  var x = $('#validateform').submit();


  /*
  // console.log(filesList);
  if(navigator.userAgent.toLowerCase().indexOf('firefox') > -1){
  // Firefox treats files as read only
  var fileArray = [filesList[0].name];
  var fc = multiFileInput = document.getElementById('fileConvert'); 
  fc.mozSetFileNameArray(fileArray, fileArray.length);
  } else {
  $('#fileConvert').prop('files', filesList);
  }
  var x = $('#convForm').submit();
  */
  // console.log(x);

}



// 
// 
// websocket handler
// 
function startStreamSocket()
{
  var loc = window.location;
  var uri = 'ws:';

  if (loc.protocol === 'https:')
  {
    uri = 'wss:';
  }
  uri += '//' + loc.host;
  uri += '/naplan/reg/stream/' + txID;

  console.log(uri)

    websocket = new WebSocket(uri);

  msgs_rcvd = 0;
  report_data = [];

  websocket.onopen = function(evt)
  {
    onOpen(evt)
  };
  websocket.onclose = function(evt)
  {
    onClose(evt)
  };
  websocket.onmessage = function(evt)
  {
    onMessage(evt)
  };
  websocket.onerror = function(evt)
  {
    onError(evt)
  };
}

function onOpen(evt)
{
  console.log('Stream ' + txID + ' Connected');
}

function onClose(evt)
{
  console.log('Stream ' + txID + ' Disconnected');
  console.log('Validation messages received: ' + msgs_rcvd);
  console.log('Mesages in report data:' + report_data.length);
}

// 
// 
// Callbacks for websocket handlers
// 
function onMessage(evt)
{

  var msg = JSON.parse(evt.data);

  var results = document.getElementById('results-message');
  var progress = document.getElementById('progress-message');

  if (msg.Type == "result")
  {
    msgs_rcvd++;
    // console.log("results message received")
    report_data.push(msg.Payload);

  }

  if (msg.Type == "progress")
  {
    progress.innerHTML = '<p>' + "Records validated: " + msg.Payload.Progress + '</p>';
    // console.log(msg)
    // console.log("messages received: " + msgs_rcvd)
    if (msg.Payload.TxComplete)
    {
      console.log("transaction " + txID + " complete");
      websocket.close();
      // console.log("socket closed");

      Materialize.toast('All records validated', 4000, 'rounded');
    }

    if (msg.Payload.TxComplete)
    {
      progress.innerHTML = '<p>' + "Records validated: " + msg.Payload.Progress + " of " + msg.Payload.Size + '</p>';
      var manifestref = "/naplan/reg/manifest/" + txID;
      var xhttp = new XMLHttpRequest();
          xhttp.open("GET", manifestref, false);
              xhttp.setRequestHeader("Content-type", "application/json");
                  xhttp.send();
                      var prgdata = JSON.parse(xhttp.responseText);
        var manifest_container = document.getElementById('manifest');
        var manifest = "<h5>Breakdown of schools</h5><ul>";
        Object.keys(prgdata.Tally).sort().forEach(function(key) {
          manifest += "<li>School " + key + ": " + prgdata.Tally[key] + " students</li>\n";
        });
        manifest += "</ul>";
        manifest_container.innerHTML = manifest;

      results.innerHTML = '<p>' + "Results are ready for review" + '</p>'
        if (msgs_rcvd < 1)
        {
          results.innerHTML = '<p>' + "No validation errors found." + '</p>'
        }
        else
        {
          renderResultsOnce(report_data);
        }

      completeProgress();
    }
    else if (msg.Payload.UIComplete)
    {
      results.innerHTML = '<p>' + "Results are ready for review" + '</p>'
        renderResultsOnce(report_data);
    }

  }

}

function onError(evt)
{
  console.log('<span style="color: red;">ERROR:</span> ' + evt.data);
}

function doSend(message)
{
  console.log("SENT: " + message);
  websocket.send(message);
}
// 
// end of results stream websocket handlers 
// 

// ensure reults display does not flicker as data updates
function renderResultsOnce(data)
{

  if (display_results_once == true)
  {
    return
  }
  renderAnalysis(data);
  display_results_once = true;
}



// 
// 
// create results graphs
// 
// 
// 
function renderAnalysis(data)
{

  // $("#results_container").toggleClass("hide");

  var errorsBarChart = dc.barChart("#errors-chart");
  var errorsByTypeChart = dc.rowChart('#errors-by-type-chart');
  var dataTable = dc.dataTable('.dc-data-table');
  var verrorsCount = dc.dataCount('.dc-data-count');
  var errorData = data;

  // normalize/parse data so dc can correctly sort & bin them
  errorData.forEach(function(d)
      {
        d.originalLine = +d.originalLine;
      });
  // console.log(errorData);

  var ndx = crossfilter(errorData);
  var all = ndx.groupAll();

  var lineDim = ndx.dimension(function(d)
      {
        return d.originalLine;
      });

  var typeDim = ndx.dimension(function(d)
      {
        return d.validationType;
      });
  var validationTypesGroup = typeDim.group();

  var allDim = ndx.dimension(function(d)
      {
        return d;
      });

  // var countPerLine = lineDim.group().reduceCount(function(d) {
  //     return d.originalLine;
  // });

  var lineGroup = lineDim.group();

  // var countPerLine = lineDim.group().reduceSum(function(d) {return d.OriginalLine;});

  errorsBarChart
    .width(350)
    .height(180)
    .margins(
        {
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
    .colors(["#4caf50"])
    .group(lineGroup);

  errorsBarChart.dimension(lineDim);

  errorsByTypeChart
    .width(350)
    .height(180)
    .margins(
        {
          top: 20,
          left: 10,
          right: 10,
          bottom: 20
        })
  .group(validationTypesGroup)
    .dimension(typeDim)
    .title(function(d)
        {
          return d.value;
        })
  .colors(["#4caf50"])
    .elasticX(true)
    .xAxis().ticks(4);


  verrorsCount
    .dimension(ndx)
    .group(all);


  dataTable
    // .width(900)
    // .height(800)
    .dimension(allDim)
    .group(function(d)
        {
          // return 'dc.js insists on putting a row here so I remove it using JS';
          // return d.originalLine;
          // return 'Errors ordered by original file line number (table shows first 100 errors)'
          return ''
        })
  .size(100)
    .columns([
        // function (d) { return d.txID; },
        function(d)
        {
          return d.originalLine;
        },
        function(d)
        {
          return d.validationType;
        },
        function(d)
        {
          return d.errField;
        },
        function(d)
        {
          return d.description;
        },
        function(d)
        {
          return d.severity.substring(0,1).toUpperCase()
        }
  ])
    .sortBy(function(d)
        {
          return d.originalLine;
        })
  .order(d3.ascending)
    .on('renderlet', function(table)
        {
          // each time table is rendered remove nasty extra row dc.js insists on adding
          // table.select('tr.dc-table-group').remove();
          table.selectAll('.dc-table-group').classed('info', true);
        });

  dc.renderAll();

  // dc.redrawAll();

  showResults();


}

// 
// 
// Collection of utility methods for ui
// 
function hideResults()
{
  $("#results_container").addClass("hide");
}

function showResults()
{
  $("#results_container").removeClass("hide");
}

function hideProgress()
{
  $("#progress").addClass("hide");
  $("#upload-message").empty();
  $("#results-message").empty();
  $("#progress-message").empty();
  $("#manifest_accordion").addClass("hide");

}

function completeProgress()
{
  $("#progress").addClass("hide");
  $("#upload-message").empty();
  $("#manifest_accordion").removeClass("hide");
}


function showProgress(fileSizeBytes)
{
  var prgETA = calculateETA(fileSizeBytes);
  $("#progress").removeClass("hide");
  $("#upload-message").text("Validating input file..." + prgETA);
}

function calculateETA(fileSizeBytes)
{
  var ttc_seconds = ((fileSizeBytes / 1024000) * 4);
  if (ttc_seconds <= 60)
  {
    return "estimated analysis time: " + ttc_seconds.toFixed(0) + " seconds."
  }
  var ttc_minutes = (ttc_seconds / 60);
  return "estimated analysis time: " + ttc_minutes.toFixed(0) + " minutes."
}
