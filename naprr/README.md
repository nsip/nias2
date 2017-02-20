The naprr executable runs as a single binary, and launches its own streaming server instance within the application.

The unzipped binary files should have the `naprr` executable, an `/in` directory and a `/templates` directory. The `/templates` directory contains templates for the various report formats, in [Golang's template language](https://golang.org/pkg/text/template/). The `/in` directory contains zipped XML files to be processed by the `naprr` executable; the zip distribution comes packaged with `/in/master_nap.xml.zip` as a test xml results file. This will get processed automatically when the executable is run. 

To process the results file, just use a terminal to `cd `to the `naprr` folder and run the `naprr` executable.

All reports will be written to the `/out` directory; this contains summary reports at the top level, and then a folder for each school with relevant reports inside.

The current report structure in `/out` is:
* `codeframe.csv`: Flattened CSV file containing codeframe information for test items, along with the testlets and tests that contain them (given redundantly)		
* `codeframe_writing.csv`: Flattened CSV file containing codeframe information for test items in writing, along with the testlets and tests that contain them (given redundantly). Given separately from `codeframe.csv`, because of the higher level of detail in writing test items.
* `domain_scores.csv`: Domain scores for each student in the file; students are indicated by their Platform Student Identifier.	
* `participation.csv`: Which students participated in which NAPLAN tests. Demographic information and Platform Student Identifier are given for each student.
* `score_summary.csv`: Overall domain scores for each school and each test in the file. Schools identified by ACARA ID.
* _SCHOOL (ACARA ID)_
  * `domain_scores.csv`: 	Domain scores for each student in the school; students are indicated by their Platform Student Identifier.
  * `participation.csv`:	Which students in the school participated in which NAPLAN tests. Demographic information and Platform Student Identifier are given for each student.
  * `score_summary.csv`: Overall domain scores for the school for each test in the file. Schools identified by ACARA ID.

All reports have been generated as requested, except for `codeframe_writing.csv`. For the writing codeframe, we still need to understand the structure more clearly; but sensible real defaults have been put into the file from the results and reporting dataset, so the reports are structurally correct.

To test the reporting engine only, as opposed to data parsing, delete the `/out` directory then run the `naprr` executable with the flag `-rewrite`. This will simply regenerate the reports without reloading the data.

Report formats are all governed by templates so they can be modified, but will be better once a UI is in place.
