# NAPLAN Result and Reporting code

## NAPLAN Results and Reporting processing

### Note from a user perspective this code has been replaced by naprrql which now provides the interfaces and data access, most of the data classes here, however, are still in active use.

The `naprr` executable runs as a single binary, and launches its own streaming server instance within the application. Unlike other NIAS executables, it does not require a batch file or shell script to be executed.

The unzipped binary files released for `naprr` should contain the `naprr` executable, an `/in` directory and a `/templates` directory. The `/templates` directory contains templates for the various report formats, in [Golang's template language](https://golang.org/pkg/text/template/). The `/in` directory contains zipped XML files to be processed by the `naprr` executable; the zip distribution comes packaged with `/in/master_nap.xml.zip` as a test xml results file. This will get processed automatically when the executable is run. 

To process the results file, just use a terminal to `cd `to the `naprr` folder and run the `naprr` executable:

    cd naprr
    ./naprr 

    cd naprr
    ./naprr.exe

All reports will be written to the `/out` directory; this contains summary reports at the top level, and then a folder for each school with relevant reports inside. The school folders are named by their ACARA IDs.

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
  * `napevents.xml`: The NAP Events objects specific to that school, in SIF/XML format.
  * `napresponses.xml`: The NAP Reponse objects specific to that school, in SIF/XML format.
  * `students.xml`: The Student objects specific to that school, in SIF/XML format.

All reports have been generated as requested, except for `codeframe_writing.csv`. For the writing codeframe, we still need to understand the structure more clearly; but sensible real defaults have been put into the file from the results and reporting dataset, so the reports are structurally correct.

To test the reporting engine only, as opposed to data parsing, delete the `/out` directory then run the `naprr` executable with the flag `-rewrite`. This will simply regenerate the reports without reloading the data.

Report formats are all governed by templates so they can be modified, but will be better once a UI is in place.

## NAPLAN Year 3 Writing preprocessing code

The `napyr3w` executable runs as a single binary. It processes both Pearson format and FujiXerox format files for Year 3 Writing results, and converts them into SIF/XML files consistent with those received from the NAPLAN Online Assessment Platform.

* Any Pearson format files are expected to be in the `in/Pearson` directory, as fixed format txt files.
* Any FujiXerox format files are expected to be in the `in/FujiXerox` directory, as csv format files.
* The program generates a report of which students are matched between the Pearson or FujiXerox files, and any XML files received from the platform. For that reason, the `/in` directory is still expected to contain a zipped XML file of results.

* The SIF/XML output is generated to the file `yr3w/codeframe_writing.xml`. This includes a dummy Year 3 Writing codeframe, to which test results are attached.
* A report of matches and mismatches between the students in the Pearson/FujiXerox files and the SIF/XML results file is generated to the file `yr3w/codeframe_report.txt`.

## Comparing Registration CSV files with received results

The `auditdiff` executable runs as a single binary. It compares a CSV file of Student Registration records (located in the `in/Reg` directory) with the student records in the NAPLAN Online results XML file in the `in` directory, and detects which students appear only in one or the other file.

* The comparison runs in two passes. 
  * First, records in the two files which have the same Platform Identifier (PSI) are eliminated. (For the comparison to run efficiently, users should endeavour  to download from the Student Registration Management system a CSV file containing all student records, and including their allocated PSIs.) 
  * Second, all remaining records from the two files are compared according to fields they have in common. The fields are specified in the `naprr.toml` file, under the key `MatchAttributes`, which contains a list of field names from the `xml.RegistrationRecord` struct. So `MatchAttributes = ["FamilyName", "GivenName"]`, for instance, will compare the  remaiining records according to their family name and given name.
* The mismatches are output to the file `reg_rep_mismatches.txt`. This contains:
  * A count of the students found to be unique to the Registration record
  * A listing of the students found to be unique to the Registration record (PSI and user-defined key)
  * A count of the students found to be unique to the Results & Reporting record
  * A listing of the students found to be unique to the Results & Reporting record (PSI, user-defined key, and RefId)
  * A listing of all the student records unique to the Registration record, in CSV
  * A listing of all the student records unique to the Results & Reporting record, in SIF/XML

