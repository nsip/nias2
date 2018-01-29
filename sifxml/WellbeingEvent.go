package sifxml


    type WellbeingEvent struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      EventId LocalId `xml:"EventId"`
      WellbeingEventNotes string `xml:"WellbeingEventNotes"`
      WellbeingEventCategoryClass string `xml:"WellbeingEventCategoryClass"`
      WellbeingEventCategoryList WellbeingEventCategoryListType `xml:"WellbeingEventCategoryList"`
      ReportingStaffRefId IdRefType `xml:"ReportingStaffRefId"`
      WellbeingEventLocationDetails WellbeingEventLocationDetailsType `xml:"WellbeingEventLocationDetails"`
      WellbeingEventCreationTimeStamp string `xml:"WellbeingEventCreationTimeStamp"`
      WellbeingEventDate string `xml:"WellbeingEventDate"`
      WellbeingEventTime string `xml:"WellbeingEventTime"`
      WellbeingEventDescription string `xml:"WellbeingEventDescription"`
      WellbeingEventTimePeriod string `xml:"WellbeingEventTimePeriod"`
      ConfidentialFlag string `xml:"ConfidentialFlag"`
      PersonInvolvementList PersonInvolvementListType `xml:"PersonInvolvementList"`
      FollowUpActionList FollowUpActionListType `xml:"FollowUpActionList"`
      Status string `xml:"Status"`
      DocumentList WellbeingDocumentListType `xml:"DocumentList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    