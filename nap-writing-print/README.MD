## nap-writing-print

nap-writing-print is a standalone tool that generates pdf files of NAPLAN writing responses.

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
1. copy the file called writing_extract.csv to the /in folder of the nap-writing-print directory.
1. invoke nap-writing-print with the following command
> ./nap-writing-print

The nap-writing-print tool will then run and create individual html files for each writing response that it finds 
in the provided writing_extract.csv file.

When the tool runs it creates two folder strucures in its /out directory

/schools

contains all writing responses organised acording to the ACARA ID of the school the response came from

/yr-level

contains all writing responses organised by the year level of the student who created the response

in each of the folders you will find the following sub-folder strucure

/audit

/script

The files are named with a prefix of the state identifier and the student participation code, followed by the anonymous id given 
to that student's script. The participation prefix is to allow users to ignore
scripts or non-scripts from students who did not complete the writing test.

writing script files end with the suffix script
audit files end with the suffix audit

The files have the same name, for cross validation, in each folder

Files in the /script directory contain
the text that the student created, with a header showing the anonymous document id of that script

Files in the /audit folder contain
the meta-data associated with the student's script; such as PSI, Local identifier, TAA identifier etc.

This means scripts can be distributed independently, but manually reconciled if needed.

As the html outout is constructed entirely from the contents of the input file, for safety at the end of a run a timestamped 
folder is created in the /backup directory of the nap-writing-print folder so that the same html files can be 
generated at any time in the future even if the working .csv file has been over-written.

The output file maintains the paragraphing of the original input from the user, bold text, underlined text, italic text, and ordered and unordered lists as available in the online naplan editor component.

It goes without saying that this tool will generate a lot of files, so is best run on 64-bit environments where constraints
on number of files in a directory or folder are not an issue if large input files are being processed.

It is advised for best performance to run on a machine with an  SSD hard disk.

Current performance on a i5 macbook creates 200,000 html files (4 files for each of 50k students) in around 1 minute.










