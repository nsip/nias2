// naprr_codeframe.js

// instantiate interaction listeners
$(document).ready(function()
{

    // initiate modals
    $('.modal').modal();

    // create the report when button is clicked
    $("#btn-codeframe").on("click", function(event)
    {
        hideReport();
        alert("Sorry, this report still under construction.");
        // createCodeframeReport();
        // showReport();
    });

});
