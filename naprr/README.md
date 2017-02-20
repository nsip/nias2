Single binary, launches streaming server within app.

Unzip files shold have naprr executable, an /in directory and a /templates directory.

To process results file just use a terminal to cd to the naprr folder and run the naprr executable.

Zip distribution includes a test xml zip file in the /in directory - this will get processed automatically when the executable is run.

All reports will be written to the /out directory; contains summary reports at the top level and then a folder for each school with relevant reports inside.

All reports as requested except code frame_writing where we still need to understand the structure more clearly, but sensible real defaults are put into the file from the r/r dataset so structurally correct.

To test reporting engine only, as opposed to data parsing, delete the /out directory then run the naprr executable with the flag -rewrite this will simply regenerate the reports without reloading the data.

Report formats all governed by template so can be modified, but will be better once a ui is in place.
