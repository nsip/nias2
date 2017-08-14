# NAPLAN Results Rpeorting Tools (GraphQL Version)

## Usage

### Download latet binaries
1. Download the relevant go-nias-[Platform].zip file from the Releases page.
2. Extract the zip file and navigate in a terminal to the /naprrql folder

### Ingest some data
Naprrql will read any .xml or zipped (.xml.zip) results reporting files that are in the /naprrql/in folder.

The binary release contains our normal results reporting sample file of a full NAPLAN test with all testlets, items etc. and dat from 4 schools each with 150 students spread across the NAPLAN year levels.

Before interacting with naprrql you must ingest the data, this only needs to be done once.

The ingest process will also generate all specified csv reports, these can be found once created in the /naprr/out folder.

To perform the ingest at the command prompt run:

  > naprrql -ingest

(for windows users naprrql will be naprrql.exe)

  > c:\> naprrql.exe -ingest  

### Launch naprrql

To run the web interfaces of naprrql, at the command prompt run:

 > naprrql
 
this will start the web/graphql servers.

Browse to:

http://localhost:1329/sifql

for the graphql data explorer

http://localhost:1329/ui

for the regular user interface.

### Reporting

The naprrql reporting engine that generates csv reports at a school & system level now uses the same queries as the web interfaces.

Use the data explorer to create new queries, and to trun them into reports save the queries as files in the /naprr/school_templates or /naprr/system_templates folders.

To regenerate csv reports at the command promt run:

> naprrql -report







