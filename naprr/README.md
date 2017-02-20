The naprr executable runs as a single binary, and launches its own streaming server instance within the application.

The unzipped binary files should have the `naprr` executable, an `/in` directory and a `/templates` directory. The `/templates` directory contains templates for the various report formats, in [Golang's template language](https://golang.org/pkg/text/template/). The `/in` directory contains zipped XML files to be processed by the `naprr` executable; the zip distribution comes packaged with `/in/master_nap.xml.zip` as a test xml results file. This will get processed automatically when the executable is run. 

To process the results file, just use a terminal to `cd `to the `naprr` folder and run the `naprr` executable.

All reports will be written to the `/out` directory; this contains summary reports at the top level, and then a folder for each school with relevant reports inside.

All reports have been generated as requested, except for code frame_writing. For the code frame, we still need to understand the structure more clearly; but sensible real defaults have been put into the file from the results and reporting dataset, so the reports are structurally correct.

To test the reporting engine only, as opposed to data parsing, delete the `/out` directory then run the `naprr` executable with the flag `-rewrite`. This will simply regenerate the reports without reloading the data.

Report formats are all governed by templates so they can be modified, but will be better once a UI is in place.
