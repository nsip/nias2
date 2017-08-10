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
function createScoreSummaryGraph(ssdata)
{

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


// domain score - cumulative bar with bubble for point value

function createDomainScoreGraph(rdsdata)
{

    var testBands = rdsdata.Test.TestContent.DomainBands;

    var data = [
        {
            'name': 'band 1',
            'lower': +testBands.Band1Lower,
            'upper': +testBands.Band1Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 2',
            'lower': +testBands.Band2Lower,
            'upper': +testBands.Band2Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 3',
            'lower': +testBands.Band3Lower,
            'upper': +testBands.Band3Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 4',
            'lower': +testBands.Band4Lower,
            'upper': +testBands.Band4Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 5',
            'lower': +testBands.Band5Lower,
            'upper': +testBands.Band5Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 6',
            'lower': +testBands.Band6Lower,
            'upper': +testBands.Band6Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 7',
            'lower': +testBands.Band7Lower,
            'upper': +testBands.Band7Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 8',
            'lower': +testBands.Band8Lower,
            'upper': +testBands.Band8Upper,
            'score': -1,
            'type': 'band'
        },
        {
            'name': 'band 9',
            'lower': +testBands.Band9Lower,
            'upper': +testBands.Band9Upper,
            'score': -1,
            'type': 'band'
        },

        {
            'name': 'band 10',
            'lower': +testBands.Band10Lower,
            'upper': +testBands.Band10Upper,
            'score': +rdsdata.Response.DomainScore.ScaledScoreValue,
            'type': 'band'
        },

    ];


    var svg = dimple.newSvg("#graphContainer", 500, 200);
    var chart = new dimple.chart(svg, data);
    chart.setBounds(50, 10, 450, 150);

    chart.defaultColors = [
        new dimple.color("#90caf9", null, 0.1),
        new dimple.color("#90caf9", null, 0.2),
        new dimple.color("#90caf9", null, 0.3),
        new dimple.color("#90caf9", null, 0.4),
        new dimple.color("#90caf9", null, 0.5),
        new dimple.color("#90caf9", null, 0.6),
        new dimple.color("#90caf9", null, 0.7),
        new dimple.color("#90caf9", null, 0.8),
        new dimple.color("#90caf9", null, 0.9),
        new dimple.color("#90caf9", null, 0.95),
        new dimple.color("#333333")
    ];


    var xAxis = chart.addCategoryAxis("x", "name");
    var orderBands = [
        'band 1',
        'band 2',
        'band 3',
        'band 4',
        'band 5',
        'band 6',
        'band 7',
        'band 8',
        'band 9',
        'band 10'
    ];

    xAxis.addOrderRule(orderBands);
    xAxis.title = null;

    var xAxis2 = chart.addMeasureAxis("x", "score")
    xAxis2.overrideMax = +testBands.Band10Upper;
    xAxis2.overrideMin = +testBands.Band1Lower;
    xAxis2.title = null;
    xAxis2.hidden = true;

    var yAxis = chart.addCategoryAxis("y", "type");
    yAxis.title = null;
    yAxis.hidden = true;

    var series = chart.addSeries("upper", dimple.plot.bar);
    series.barGap = 0.1;

    chart.addSeries("score", dimple.plot.bubble, [xAxis2, yAxis]);

    chart.draw();


}
