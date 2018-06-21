## nap-writing-sanitiser

nap-writing-sanitiser is a standalone tool that sanitises the writing extract from NAPLAN, to strip extraneous HTML styling introduced through copy-paste, and to wrap responses in HTML where responses lack any HTML markup.

It is designed to work with the .csv file that is output by the naprrql data analysis tool when naprrql is run with the 
--writingextract flag.

The full process to go from receiveing a NAPLAN results data file is:

1. Place the results data file in the /in directory of naprrql
1. Navigate in a command/terminal shell to the naprrql directory
1. Ingest the data file into naprrql by running the following command
> ./naprrql --ingest
1. wait for naprrql ingest to complete
1. Run the writing extract report of naprrql by running the following command
> ./naprrql --writingextract
1. wait for the report to complete.
1. Go to the /out/writing_extract folder of naprrql
1. copy the file called writing_extract.csv to the /in folder of the nap-writing-sanitiser directory.
1. invoke nap-writing-sanitiser with the following command
> ./nap-writing-sanitiser

The nap-writing-sanitiser tool will then generate two files:

* out/writing_extract_sanitised.csv: a file with the same structure as writing_extract.csv, but with the HTML content sanitised
* out/sanitise_report.csv: a before-and-after list of all sanitised responses

The following HTML markup in responses is preserved during sanitisation:

* The elements "strong", "em", "span", "p", "ol", "ul", "li", "br", "u", "font", "h1", "h2", "h3", "h4", "h5", "h6"
* The attribute "size" on the element "font", with a numeric value.
* The following values of the attribute "style" on any element, singly or in combination:
  * text-decoration:underline;
  * text-decoration-line:underline;
  * font-size:16px;
  * font-size:18px;
  * font-size:large;
  * text-align:left;
  * text-align:center;
  * text-align:start;
  * background-color:rgba(255, 255, 255, 0);

All other attributes are stripped, including CSS classes, and other values of the "style" attribute. All other elements are stripped; this includes the bulk of elements in the Writing prompt.

If any responses lack any HTML markup, the tool inserts a `<p>` wrapper around the response, and in any instances of double carriage return.

To make sense of the sanitise_report.csv report, we suggest a visual diff tool, that highlights the changes to markup introduced by the sanitiser. https://text-compare.com is one example of such a tool.

// Current performance on a i5 macbook creates 200,000 html files (4 files for each of 50k students) in around 1 minute.
