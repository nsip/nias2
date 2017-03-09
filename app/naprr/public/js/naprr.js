// globals
var scoresummaryData = [];
var participationData = [];
var domainscoresData = [];
var schoolinfoData = {};

var currentASLId;

// set up initial ui behaviour
$(document).ready(function()
{

    // data for school chooser autocomplete
    initSchoolChooser();
    initSchoolChooserHandler();

    // clear the ui
    hideReport();

    // show the privacy warning
    $('.modal').modal();
    // $('.modal').modal('open');

    $(".button-collapse").sideNav();

});









