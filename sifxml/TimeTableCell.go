package sifxml


    type TimeTableCell struct {
        RefId RefIdType `xml:"RefId,attr"`
      TimeTableRefId IdRefType `xml:"TimeTableRefId"`
      TimeTableSubjectRefId IdRefType `xml:"TimeTableSubjectRefId"`
      TeachingGroupRefId IdRefType `xml:"TeachingGroupRefId"`
      RoomInfoRefId IdRefType `xml:"RoomInfoRefId"`
      StaffPersonalRefId IdRefType `xml:"StaffPersonalRefId"`
      TimeTableLocalId LocalIdType `xml:"TimeTableLocalId"`
      SubjectLocalId LocalIdType `xml:"SubjectLocalId"`
      TeachingGroupLocalId LocalIdType `xml:"TeachingGroupLocalId"`
      RoomNumber HomeroomNumberType `xml:"RoomNumber"`
      StaffLocalId LocalIdType `xml:"StaffLocalId"`
      DayId LocalIdType `xml:"DayId"`
      PeriodId LocalIdType `xml:"PeriodId"`
      CellType string `xml:"CellType"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      TeacherList ScheduledTeacherListType `xml:"TeacherList"`
      RoomList RoomListType `xml:"RoomList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    