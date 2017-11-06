require "csv"
require "pp"

# 3. Student test results to codeframe comparison

codeframe = {}
substitute = {}
itemlocalid = {}
CSV.foreach("../../out/qa/systemCodeframe.csv", headers: true) do |row|
  codeframe[row["Test"]] = {} if codeframe[row["Test"]].nil?
  codeframe[row["Test"]][row["NAPTestletLocalId"]] = [] if codeframe[row["Test"]][row["NAPTestletLocalId"]].nil?
  codeframe[row["Test"]][row["NAPTestletLocalId"]] << row["NAPTestItemLocalId"]
  itemlocalid[row["ItemRefID"]] = row["NAPTestletLocalId"]
  if row["ItemSubstitutedForList"] != "{\"SubstituteItem\":[]}"
    sub_items = row["ItemSubstitutedForList"].gsub(/\{"SubstituteItem":\[/,"").gsub(/\]\}/,"").split(",")
    sub_items.each do |s|
      substitute[s.gsub(/\{"SubstituteItemRefId":"/,"").gsub(/"\}/,"")] = row["NAPTestItemLocalId"]
    end
  end
end

codeframe.each do |k, v|
  v.each do |k1, v1|
    codeframe[k][k1] = codeframe[k][k1].sort
  end
end

testlet = ""
test = ""
psi = ""
participationcode = ""
testitems = nil
CSV.open("../../out/qa/systemCodeframeMap.rpt.csv", "wb",write_headers: true,
         headers: ["PSI","Test","Testlet Name","Participation Code","Expected Items","Found Items"]) do |rpt|
  CSV.foreach("../../out/qa/itemPrinting.csv", headers: true) do |row|
    if ( testlet != row["NAPTestletLocalId"] || test != row["Test Name"] ) && !testitems.nil?
      norm_testitems = []
      testitems.each do |t|
        norm_testitems << (substitute.has_key?(t) ? substitute[t] : itemlocalid[t])
      end
      norm_testitems.sort! 
      if codeframe[test][testlet] != norm_testitems
        rpt << [psi, test, testlet, participationcode, norm_testitems.join(";"), codeframe[test][testlet].join(";")]
      end
    end
    if ( testlet != row["NAPTestletLocalId"] || test != row["Test Name"] ) 
      testitems = []
    end
    testlet = row["NAPTestletLocalId"]
    test = row["Test Name"]
    psi = row["PSI"]
    participationcode = row["Participation Code"]
    testitems << row["ItemRefID"]
  end
  testitems.sort!
  if !testitems.nil? && codeframe[test][testlet] != testitems
    rpt << [psi, test, testlet, participationcode, testitems.join(";"), codeframe[test][testlet].join(";")]
  end
end

