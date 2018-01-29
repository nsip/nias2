package sifxml


    type Identity struct {
        RefId RefIdType `xml:"RefId,attr"`
      SIF_RefId Identity_SIF_RefId `xml:"SIF_RefId"`
      AuthenticationSource string `xml:"AuthenticationSource"`
      IdentityAssertions IdentityAssertionsType `xml:"IdentityAssertions"`
      PasswordList PasswordListType `xml:"PasswordList"`
      AuthenticationSourceGlobalUID IdRefType `xml:"AuthenticationSourceGlobalUID"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type Identity_SIF_RefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
