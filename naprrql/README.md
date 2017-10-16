# NAPLAN Results Reporting Tools (GraphQL Version)

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
 
 for windos users:
 
 > naprrql.exe
 
 or you can simply double-click the executable file in explorer
 
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

### SIFQl (GraphQL) Data Explorer

NAPRRQL includes an instance of the GraphQL data explorer interface.

This allows you to make queries against the underlying data-store in accordance with the SIF Schema for NAPLAN results reporting.

The full set of queries and the data elements they return is set out in the Docs area of the data explorer ui, a key feature is that whilst queries conform to the schema you are free to request only the fields you are intrested in rather than having to receive the whole data objects.

On your first visit to the data-explorer interface it contains a set of comments which can be safely deleted.

Constructing queries within the interface is straightforward, but here are some samples to introduce you to the mechanism.

These sample queries assume naprrql is running on its default host:port combination (localhost:1329), but if you click on the queries when naprrql is running you should see the queries in the explorer interface, just click on the Play button in the toolbar to run them.

### Basic NAPLAN Data Queries

#### NAP Tests

http://localhost:1329/sifql?query=query%20NAPTests%20%7B%0A%09tests%20%7B%0A%09%20%20TestID%0A%20%20%20%20TestContent%20%7B%0A%20%20%20%20%20%20LocalId%0A%20%20%20%20%20%20TestName%0A%20%20%20%20%20%20TestLevel%0A%20%20%20%20%20%20TestDomain%0A%20%20%20%20%20%20TestYear%0A%20%20%20%20%20%20StagesCount%0A%20%20%20%20%20%20TestType%0A%20%20%20%20%7D%0A%09%7D%0A%7D%0A&operationName=NAPTests

#### NAP Testlets

http://localhost:1329/sifql?query=query%20NAPTestlets%20%7B%0A%20%20testlets%20%7B%0A%20%20%20%20TestletContent%20%7B%0A%20%20%20%20%20%20LocalId%0A%20%20%20%20%20%20NAPTestLocalId%0A%20%20%20%20%20%20TestletName%0A%20%20%20%20%20%20Node%0A%20%20%20%20%20%20LocationInStage%0A%20%20%20%20%20%20TestletMaximumScore%0A%20%20%20%20%7D%0A%20%20%20%20TestItemList%20%7B%0A%20%20%20%20%20%20TestItem%20%7B%0A%20%20%20%20%20%20%20%20TestItemLocalId%0A%20%20%20%20%20%20%20%20SequenceNumber%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D%0A&operationName=NAPTestlets

#### NAP TestItems

http://localhost:1329/sifql?query=query%20NAPTestItems%20%7B%0A%20%20testitems%20%7B%0A%20%20%20%20TestItemContent%20%7B%0A%20%20%20%20%20%20NAPTestItemLocalId%0A%20%20%20%20%20%20ItemName%0A%20%20%20%20%20%20ItemType%0A%20%20%20%20%20%20Subdomain%0A%20%20%20%20%20%20WritingGenre%0A%20%20%20%20%20%20ItemDescriptor%0A%20%20%20%20%20%20ReleasedStatus%0A%20%20%20%20%20%20MarkingType%0A%20%20%20%20%20%20MultipleChoiceOptionCount%0A%20%20%20%20%20%20CorrectAnswer%0A%20%20%20%20%20%20MaximumScore%0A%20%20%20%20%20%20ItemDifficulty%0A%20%20%20%20%20%20ItemDifficultyLogit5%0A%20%20%20%20%20%20ItemDifficultyLogit62%0A%20%20%20%20%20%20ItemDifficultyLogit5SE%0A%20%20%20%20%20%20ItemDifficultyLogit62SE%0A%20%20%20%20%20%20ItemProficiencyBand%0A%20%20%20%20%20%20ItemProficiencyLevel%0A%20%20%20%20%20%20ExemplarURL%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D%0A&operationName=NAPTestItems

#### Score Summaries

http://localhost:1329/sifql?query=query%20NAPScoreSummaries%20%7B%0A%20%20score_summaries%20%7B%0A%20%20%20%20SummaryID%0A%20%20%20%20SchoolInfoRefId%0A%20%20%20%20SchoolACARAId%0A%20%20%20%20NAPTestRefId%0A%20%20%20%20NAPTestLocalId%0A%20%20%20%20DomainNationalAverage%0A%20%20%20%20DomainSchoolAverage%0A%20%20%20%20DomainJurisdictionAverage%0A%20%20%20%20DomainTopNational60Percent%0A%20%20%20%20DomainBottomNational60Percent%0A%20%20%7D%0A%7D%0A&operationName=NAPScoreSummaries

#### Students

http://localhost:1329/sifql?query=query%20NAPStudents%20%7B%0A%20%20students%20%7B%0A%20%20%20%20RefId%0A%20%20%20%20LocalId%0A%20%20%20%20StateProvinceId%0A%20%20%20%20FamilyName%0A%20%20%20%20GivenName%0A%20%20%20%20MiddleName%0A%20%20%20%20PreferredName%0A%20%20%20%20IndigenousStatus%0A%20%20%20%20Sex%0A%20%20%20%20BirthDate%0A%20%20%20%20CountryOfBirth%0A%20%20%20%20StudentLOTE%0A%20%20%20%20VisaCode%0A%20%20%20%20LBOTE%0A%20%20%20%20AddressLine1%0A%20%20%20%20AddressLine2%0A%20%20%20%20Locality%0A%20%20%20%20StateTerritory%0A%20%20%20%20Postcode%0A%20%20%20%20SchoolLocalId%0A%20%20%20%20YearLevel%0A%20%20%20%20FTE%0A%20%20%20%20Parent1LOTE%0A%20%20%20%20Parent2LOTE%0A%20%20%20%20Parent1Occupation%0A%20%20%20%20Parent2Occupation%0A%20%20%20%20Parent1SchoolEducation%0A%20%20%20%20Parent2SchoolEducation%0A%20%20%20%20Parent1NonSchoolEducation%0A%20%20%20%20Parent2NonSchoolEducation%0A%20%20%20%20LocalCampusId%0A%20%20%20%20ASLSchoolId%0A%20%20%20%20TestLevel%0A%20%20%20%20Homegroup%0A%20%20%20%20ClassGroup%0A%20%20%20%20MainSchoolFlag%0A%20%20%20%20FFPOS%0A%20%20%20%20ReportingSchoolId%0A%20%20%20%20OtherSchoolId%0A%20%20%20%20EducationSupport%0A%20%20%20%20HomeSchooledStudent%0A%20%20%20%20Sensitive%0A%20%20%20%20OfflineDelivery%0A%20%20%7D%0A%7D%0A&operationName=NAPStudents

#### Schools

http://localhost:1329/sifql?query=query%20NAPSchools%20%7B%0A%20%20schools%20%7B%0A%20%20%20%20RefId%0A%20%20%20%20LocalId%0A%20%20%20%20StateProvinceId%0A%20%20%20%20ACARAId%0A%20%20%20%20SchoolName%0A%20%20%20%20LEAInfoRefId%0A%20%20%20%20OtherLEA%0A%20%20%20%20SchoolDistrict%0A%20%20%20%20SchoolType%0A%20%20%20%20StudentCount%0A%20%20%20%20SchoolURL%0A%20%20%20%20SessionType%0A%20%20%20%20ARIA%0A%20%20%20%20OperationalStatus%0A%20%20%20%20FederalElectorate%0A%20%20%20%20SchoolSector%0A%20%20%20%20IndependentSchool%0A%20%20%20%20NonGovSystemicStatus%0A%20%20%20%20System%0A%20%20%20%20ReligiousAffiliation%0A%20%20%20%20SchoolGeographicLocation%0A%20%20%20%20LocalGovernmentArea%0A%20%20%20%20JurisdictionLowerHouse%0A%20%20%20%20SLA%0A%20%20%20%20SchoolCoEdStatus%0A%20%20%7D%0A%7D%0A&operationName=NAPSchools

### School by School Queries

These queries return rich datasets for a single school, or for mutliple schools. 
These queries make use of the Varaibles area of the ui. 
Schools to report on are selected by passing the ACAR Id (ASL Id) of the school(s) to report on.
The variable acaraIDs is an array that can contain one or more acaraids to identify the selected schools.

#### Domain Scores

http://localhost:1329/sifql?query=query%20schoolDomianScores(%24acaraIDs%3A%20%5BString%5D)%20%7B%0A%20%20domain_scores_report_by_school(acaraIDs%3A%20%24acaraIDs)%20%7B%0A%20%20%20%20Test%20%7B%0A%20%20%20%20%20%20TestContent%20%7B%0A%20%20%20%20%20%20%20%20TestName%0A%20%20%20%20%20%20%20%20TestYear%0A%20%20%20%20%20%20%20%20TestLevel%0A%20%20%20%20%20%20%20%20TestDomain%0A%20%20%20%20%20%20%20%20StagesCount%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%20%20Response%20%7B%0A%20%20%20%20%20%20ReportExclusionFlag%0A%20%20%20%20%20%20CalibrationSampleFlag%0A%20%20%20%20%20%20EquatingSampleFlag%0A%20%20%20%20%20%20PathTakenForDomain%0A%20%20%20%20%20%20ParallelTest%0A%20%20%20%20%20%20PSI%0A%20%20%20%20%20%20DomainScore%20%7B%0A%20%20%20%20%20%20%20%20RawScore%0A%20%20%20%20%20%20%20%20ScaledScoreValue%0A%20%20%20%20%20%20%20%20ScaledScoreLogitValue%0A%20%20%20%20%20%20%20%20ScaledScoreStandardError%0A%20%20%20%20%20%20%20%20ScaledScoreLogitStandardError%0A%20%20%20%20%20%20%20%20StudentDomainBand%0A%20%20%20%20%20%20%20%20StudentProficiency%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D%0A&operationName=schoolDomianScores

#### Participation Report (queries two schools)

http://localhost:1329/sifql?query=query%20schoolParticipation(%24acaraIDs%3A%20%5BString%5D)%20%7B%0A%20%20participation_report_by_school(acaraIDs%3A%20%24acaraIDs)%20%7B%0A%20%20%20%20School%20%7B%0A%20%20%20%20%20%20LocalId%0A%20%20%20%20%20%20ACARAId%0A%20%20%20%20%20%20SchoolName%0A%20%20%20%20%20%20SchoolDistrict%0A%20%20%20%20%7D%0A%20%20%20%20Student%20%7B%0A%20%20%20%20%20%20YearLevel%0A%20%20%20%20%20%20GivenName%0A%20%20%20%20%20%20FamilyName%0A%20%20%20%20%7D%0A%20%20%20%20Summary%20%7B%0A%20%20%20%20%20%20Domain%0A%20%20%20%20%20%20ParticipationCode%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D%0A&variables=%7B%0A%20%20%22acaraIDs%22%3A%20%5B%0A%20%20%20%20%2249360%22%2C%20%2249453%22%0A%20%20%5D%0A%7D&operationName=schoolParticipation


Note that if you find a query helpful and want to use it to produce csv reports simply save the query into a text file with a '.gql' extension in one of the template directories; /school_templates for queries that search by school, and /system_templates for generic data queries.

Once templates are saved in these folders re-running napprql with the -report flag will use those queries to generate report .csv files based on the queries in the /out folder.






