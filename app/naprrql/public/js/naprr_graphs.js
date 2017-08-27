// naprr_graphs.js

// 
// helper routines to generate graphs used in ui
// 

// NOTE upgrade these!!!
// <script src="/lib/d3.v4.3.0.js"></script>
// <script src="http://dimplejs.org/dist/dimple.v2.3.0.min.js"></script>

// 
// generate score summary graph
// 
function createScoreSummaryGraph(ssdata) {

    var data = [];

    var attributes = {};
    attributes['Measure'] = 'School Average';
    attributes['Value'] = +ssdata.Summ.DomainSchoolAverage;
    data.push(attributes);

    var attributes = {};
    attributes['Measure'] = 'Jurisdiction Average';
    attributes['Value'] = +ssdata.Summ.DomainJurisdictionAverage;
    data.push(attributes);

    var attributes = {};
    attributes['Measure'] = 'National Average';
    attributes['Value'] = +ssdata.Summ.DomainNationalAverage;
    data.push(attributes);

    var attributes = {};
    attributes['Measure'] = 'Top National 60%';
    attributes['Value'] = +ssdata.Summ.DomainTopNational60Percent;
    data.push(attributes);

    var attributes = {};
    attributes['Measure'] = 'Bottom National 60%';
    attributes['Value'] = +ssdata.Summ.DomainBottomNational60Percent;
    data.push(attributes);

    // console.log(data);

    var orderMeasures = [
        'Bottom National 60%',
        'School Average',
        'Jurisdiction Average',
        'National Average',
        'Top National 60%'
    ];

    var svg = dimple.newSvg("#graphContainer", 800, 200);
    var chart = new dimple.chart(svg, data);
    chart.setBounds(50, 10, 750, 150);
    var xAxis = chart.addCategoryAxis("x", "Measure");
    xAxis.addOrderRule(orderMeasures);
    xAxis.title = null;
    var yAxis = chart.addMeasureAxis("y", "Value");
    yAxis.title = null;
    yAxis.overrideMax = 999;
    // yAxis.useLog = true;
    chart.addSeries("Measure", dimple.plot.bubble);
    // myChart.addLegend(200, 10, 360, 20, "right");
    chart.draw();


}

// 
// linear ds graph with display of student score on bands
// 
function createDomainScoreGraph(rdsdata) {

    var margin = { top: 5, right: 5, bottom: 25, left: 80 },
        width = 960 - margin.left - margin.right,
        height = 80 - margin.top - margin.bottom;

    var chart = d3.bullet()
        .width(width)
        .height(height);

    var testBands = rdsdata.Test.TestContent.DomainBands;

    var dataset = [];
    dataset[0] = {};
    dataset[0].title = " Jrsd. Avg.";
    dataset[0].subtitle = "School Avg.";
    dataset[0].ranges = [+testBands.Band1Upper, +testBands.Band2Upper, +testBands.Band3Upper, +testBands.Band4Upper, +testBands.Band5Upper, +testBands.Band6Upper, +testBands.Band7Upper, +testBands.Band8Upper, +testBands.Band9Upper, +testBands.Band10Upper];
    dataset[0].markers = [+rdsdata.Response.DomainScore.ScaledScoreValue];
    // get jur & school scores for comparison
    var tlevel = rdsdata.Test.TestContent.TestLevel;
    var tdomain = rdsdata.Test.TestContent.TestDomain;
    dataset[0].measures = getScoreSummaryAverages(tlevel, tdomain);
    
    // console.log(data);

    var svg = d3.select("#graphContainer").selectAll("svg")
        .data(dataset)
        .enter().append("svg")
        .attr("class", "bullet")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
        .call(chart);

    var title = svg.append("g")
        .style("text-anchor", "end")
        .attr("transform", "translate(-6," + height / 2 + ")");

    title.append("text")
        .attr("class", "title")
        .text(function(d) { return d.title; });

    title.append("text")
        .attr("class", "subtitle")
        .attr("dy", "1em")
        .text(function(d) { return d.subtitle; });


}