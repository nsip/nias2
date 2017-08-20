### NAPCOMP

napcomp is a command-line utiltiy for reporting differences in student data 
between naplan results files and naplan registraiton files.

The utility is run from the command-line without a ui, and uses the following folder structure:

> napcomp/in/registration

in this folder place a naplan registration file or files. These will typically be .csv files used for 
student registration on the NAPALN online platform (SRM). 

> napcomp/in/results

in this folder place one or more archives of napalan results reporting data, either .xml or .xml.zip or .zip files,
that contain the results reporting dataset.

To run napcomp, simply navigate in the terminal to the napcomp directory and then run the following:

> ./napcomp

for linux/mac users

> napcomp.exe 

for windows users

The napcomp utility will then read all files from the input directories, will run the difference logic and then produce 
reports of the result of the difference analysis in the /out folder.

The reports created will be:

> napcomp/out/RegisteredButNotInResults.csv

> napcomp/out/ResultsButNotInRegister.csv

each file will contain the full student records of those students who were either in the registration file, but not 
found in the results file OR
students who were found in the results but did not appear in the registration file.
