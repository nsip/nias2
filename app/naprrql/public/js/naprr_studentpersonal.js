// 
// get student info 1-line summary for given psi
// 
function getStudentInfoSummaryLine(psi) {

    var sp = {};
    var sp_psi = "";

    // data can be found in participation info
    $.each(studentPersonalData, function(index, studentpersonal) {
        // ei = pds.EventInfos[0];
        // sp = studentpersonal;
        $.each(studentpersonal.OtherIdList.OtherId, function(index, oid) {
        if (oid.Type == "NAPPlatformStudentId") {
            if (oid.Value == psi) {
                // var student = pds.Student;
                // sp = student;
                sp = studentpersonal;
                sp_psi = oid.Value;
                //return false;
            }
        }
        });
    });

    return $("<p>" +
        "Student: " + sp.GivenName + " " + sp.FamilyName +
        ", PSI: " + sp_psi +
        ", Homegroup: " + sp.Homegroup +
        ", Class: " + sp.ClassGroup +
        ", Offline: " + sp.OfflineDelivery +
        ", Ed. Support: " + sp.EducationSupport +
        ", Home-Schooled: " + sp.HomeSchooledStudent +
        "</p>");

}


// 
// gql query to support search
// 
function studentPersonalQuery() {

    return `
query NAPData($acaraIDs: [String]) {
  students_by_school(acaraIDs: $acaraIDs) {
    RefId
    LocalId
    StateProvinceId
    FamilyName
    GivenName
    MiddleName
    PreferredName
    IndigenousStatus
    Sex
    BirthDate
    CountryOfBirth
    StudentLOTE
    VisaCode
    LBOTE
    AddressLine1
    AddressLine2
    Locality
    StateTerritory
    Postcode
    SchoolLocalId
    YearLevel
    FTE
    Parent1LOTE
    Parent2LOTE
    Parent1Occupation
    Parent2Occupation
    Parent1SchoolEducation
    Parent2SchoolEducation
    Parent1NonSchoolEducation
    Parent2NonSchoolEducation
    LocalCampusId
    ASLSchoolId
    TestLevel
    Homegroup
    ClassGroup
    MainSchoolFlag
    FFPOS
    ReportingSchoolId
    OtherSchoolId
    EducationSupport
    HomeSchooledStudent
    Sensitive
    OfflineDelivery
    OtherIdList{
      OtherId {
        Type
        Value
      }
    }
  }
}
`

}
