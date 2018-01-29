package sifxml


    type SystemRole struct {
        RefId RefIdType `xml:"RefId,attr"`
      SIF_RefId SystemRole_SIF_RefId `xml:"SIF_RefId"`
      SystemContextList SystemRole_SystemContextList `xml:"SystemContextList"`
      SIF_Metadata string `xml:"SIF_Metadata"`
      SIF_ExtendedElements string `xml:"SIF_ExtendedElements"`
      
      }
    type SystemRole_SIF_RefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
type SystemRole_SystemContextList struct {
      SystemContext []SystemRole_SystemContext `xml:"SystemContext"`
}
type SystemRole_SystemContext struct {
      SystemId string `xml:"SystemId,attr"`
      RoleList SystemRole_RoleList `xml:"RoleList"`
}
type SystemRole_RoleList struct {
      Role []SystemRole_Role `xml:"Role"`
}
type SystemRole_Role struct {
      RoleId string `xml:"RoleId,attr"`
      RoleScopeList SystemRole_RoleScopeList `xml:"RoleScopeList"`
}
type SystemRole_RoleScopeList struct {
      RoleScope []SystemRole_RoleScope `xml:"RoleScope"`
}
type SystemRole_RoleScope struct {
       RoleScopeName string `xml:"RoleScopeName"`
      RoleScopeRefId SystemRole_RoleScopeRefId `xml:"RoleScopeRefId"`
}
type SystemRole_RoleScopeRefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value string `xml:",chardata"`
}
