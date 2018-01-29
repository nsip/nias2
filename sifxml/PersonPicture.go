package sifxml


    type PersonPicture struct {
        RefId RefIdType `xml:"RefId,attr"`
      ParentObjectRefId PersonPicture_ParentObjectRefId `xml:"ParentObjectRefId"`
      SchoolYear string `xml:"SchoolYear"`
      PictureSource PersonPicture_PictureSource `xml:"PictureSource"`
      OKToPublish string `xml:"OKToPublish"`
      SIF_Metadata string `xml:"SIF_Metadata"`
      SIF_ExtendedElements string `xml:"SIF_ExtendedElements"`
      
      }
    type PersonPicture_ParentObjectRefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
type PersonPicture_PictureSource struct {
      Type string `xml:"Type,attr"`
      Value URIOrBinaryType `xml:",chardata"`
}
