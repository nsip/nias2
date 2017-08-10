// naprr_schoolinfo.js

// creates the sidebar to show school info for the selected school

// instantiate interaction listeners
$(document).ready(function()
{

    // initiate modals
    $('.modal').modal();

    // create the report when button is clicked
    $(".button-collapse").on("click", function(event)
    {
        // hideReport();
        if (currentASLId == null)
        {
            alert("Click to display a populated school info panel once you have selecetd a school.");
        }
        else
        {
            createSchoolInfoReport();
        }
    });

});


// 
// gql query to support search
// 
function schoolInfoQuery() {
    return `
query NAPData($acaraIDs: [String]) {
  school_infos_by_acaraid(acaraIDs: $acaraIDs) {
    RefId
    LocalId
    StateProvinceId
    ACARAId
    SchoolName
    LEAInfoRefId
    OtherLEA
    SchoolDistrict
    SchoolType
    StudentCount
    SchoolURL
    SessionType
    ARIA
    OperationalStatus
    FederalElectorate
    SchoolSector
    IndependentSchool
    NonGovSystemicStatus
    System
    ReligiousAffiliation
    SchoolGeographicLocation
    LocalGovernmentArea
    JurisdictionLowerHouse
    SLA
    SchoolCoEdStatus
    YearLevels{
      YearLevel {
        Code
      }
    }
    Campus {
      SchoolCampusId
      CampusType
      AdminStatus
    }
    PrincipalInfo{
      ContactName {
        Type
        Title
        FamilyName
        GivenName
        MiddleName
        Suffix
        FullName
      }
      ContactTitle
    }
    SchoolContactList{
      SchoolContact{
        ContactInfo{
          Name {
            Type
            Title
            FamilyName
            GivenName
            MiddleName
            Suffix
            FullName
          }
          PositionTitle
          Role
          Address {
            Type
            Role
            Street{
              Line1
            }
            City
            StateProvince
            Country
            PostalCode
            GridLocation {
              Latitude
              Longitude
            }
          }
          EmailList{
            Email {
              Type
              Address
            }
          }
          PhoneNumberList{
            PhoneNumber {
              Type
              Number
              Extension
              ListedStatus
            }
          }
        }
      }
    }
  }
}

`;
}

// 
// create the school info sidebar component 
// 
function createSchoolInfoReport()
{

    // console.log(schoolinfoData);

    $("#si-header").empty();
    $("#si-contact").empty();
    $("#si-establishment").empty();
    $("#si-extended").empty();

    $("#si-header").append("<h5 class=''>" +
        schoolinfoData.SchoolName + "</h5>");

    var contactExists = schoolinfoData.SchoolContactList.SchoolContact !== null ? true : false
    if (contactExists)
    {
        contactInfo = schoolinfoData.SchoolContactList.SchoolContact[0].ContactInfo;
        $("#si-contact").append("<p>" + contactInfo.Address.Street.Line1 +
            ", " + contactInfo.Address.City +
            ", " + contactInfo.Address.StateProvince +
            ", " + contactInfo.Address.PostalCode +
            ".</p>");

        $("#si-contact").append("<p>Principal: " + schoolinfoData.PrincipalInfo.ContactName.FullName +
            ".</p>");

        $("#si-contact").append("<p>Contact: " + contactInfo.Name.FullName + "<br/>" +
            "Phone: " + contactInfo.PhoneNumberList.PhoneNumber[0].Number +
            " Ext.: " + contactInfo.PhoneNumberList.PhoneNumber[0].Extension +
            ".</p>");

    }


    $("#si-establishment").append("<p>ACARA (ASL) Id:   " + schoolinfoData.ACARAId + "<br/>" +
        "Local Id:  " + schoolinfoData.LocalId + "<br/>" +
        "State/Province Id:  " + schoolinfoData.StateProvinceId + "<br/>" +
        "District:  " + schoolinfoData.SchoolDistrict + "<br/>" +
        "Campus:  " + schoolinfoData.Campus.AdminStatus + "<br/>" +
        "Independent School:  " + schoolinfoData.IndependentSchool + "<br/>" +
        "School Sector:  " + schoolinfoData.SchoolSector + "<br/>" +
        "School Type:  " + schoolinfoData.SchoolType + "<br/>" +
        "Year Levels:  " + unpackList(schoolinfoData.YearLevels.YearLevel) + "</br>" +
        "Student Count:  " + schoolinfoData.StudentCount + "<br/>" +
        ".</p>");

    $("#si-extended").append("<p>ARIA:   " + schoolinfoData.ARIA + "<br/>" +
        "NG Systemic Status:  " + schoolinfoData.NonGovSystemicStatus + "<br/>" +
        "Operational Status:  " + schoolinfoData.OperationalStatus + "<br/>" +
        "Co-Ed Status:  " + schoolinfoData.SchoolCoEdStatus + "<br/>" +
        "LGA:  " + schoolinfoData.LocalGovernmentArea + "<br/>" +
        "Federal Electorate:  " + schoolinfoData.FederalElectorate + "<br/>" +
        "Religious Affiliation:  " + schoolinfoData.ReligiousAffiliation + "<br/>" +
        "Geographic Location:  " + schoolinfoData.SchoolGeographicLocation + "<br/></p>");

    if (contactExists)
    {
        $("#si-extended").append("<p>Grid Location:  " + contactInfo.Address.GridLocation.Latitude +
            ", " + contactInfo.Address.GridLocation.Longitude + "<br/></p>");
    }




}
