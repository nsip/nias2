// globals
var scoresummaryData = [];
var participationData = [];
var domainscoresData = [];
var codeframeData = [];
var studentPersonalData = [];
var schoolinfoData = {};
var debug =true;

var currentASLId;

// set up initial ui behaviour
$(document).ready(function()
{

    // data for school chooser autocomplete
    initSchoolChooserQL();
    initSchoolChooserHandlerQL();

    // clear the ui
    hideReport();

    // show the privacy warning
    $('.modal').modal();
    $('#modal1').modal('open');

    $(".button-collapse").sideNav();

});









