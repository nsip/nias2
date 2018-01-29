package sifxml


    type SessionInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      TimeTableCellRefId IdRefType `xml:"TimeTableCellRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      LocalId LocalIdType `xml:"LocalId"`
      TimeTableSubjectLocalId LocalIdType `xml:"TimeTableSubjectLocalId"`
      TeachingGroupLocalId LocalIdType `xml:"TeachingGroupLocalId"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      StaffPersonalLocalId LocalIdType `xml:"StaffPersonalLocalId"`
      RoomNumber string `xml:"RoomNumber"`
      DayId LocalIdType `xml:"DayId"`
      PeriodId LocalIdType `xml:"PeriodId"`
      SessionDate string `xml:"SessionDate"`
      StartTime string `xml:"StartTime"`
      FinishTime string `xml:"FinishTime"`
      RollMarked string `xml:"RollMarked"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    