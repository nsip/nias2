require "csv"
require "pp"

codeframe = {}
$substitute = {}
$itemlocalid = {}

def normalise_testitems(testitems)
  norm_testitems = []
  testitems.each do |t|
    t1 = if $substitute.has_key?(t)
           if $itemlocalid.has_key?($substitute[t])
             $itemlocalid[$substitute[t]]
           else
             $substitute[t]
           end
         else
           $itemlocalid.has_key?(t) ? $itemlocalid[t] : t
         end
    norm_testitems << t1
  end
  norm_testitems
end


# 3. Student test results to codeframe comparison

CSV.foreach("../../out/qa/systemCodeframe.csv", headers: true) do |row|
  codeframe[row["Test"]] = {} if codeframe[row["Test"]].nil?
  codeframe[row["Test"]][row["NAPTestletLocalId"]] = [] if codeframe[row["Test"]][row["NAPTestletLocalId"]].nil?
  $itemlocalid[row["ItemRefID"]] = row["NAPTestItemLocalId"]
  if row["ItemSubstitutedForList"] == "{\"SubstituteItem\":[]}"
    codeframe[row["Test"]][row["NAPTestletLocalId"]] << row["NAPTestItemLocalId"]
  else
    sub_items = row["ItemSubstitutedForList"].gsub(/\{"SubstituteItem":\[/,"").gsub(/\]\}/,"").split(",")
    sub_items.each do |s|
      $substitute[row["ItemRefID"]] = s.gsub(/\{"SubstituteItemRefId":"/,"").gsub(/"\}/,"")
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
         headers: ["PSI","Test","Testlet Name","Participation Code","Expected Items Not Found","Found Items Not Expected"]) do |rpt|
  CSV.foreach("../../out/qa/itemPrinting.csv", headers: true) do |row|
    if ( testlet != row["NAPTestletLocalId"] || test != row["Test Name"] ) && !testitems.nil?
      norm_testitems = normalise_testitems(testitems)
      expected_not_found = codeframe[test][testlet] - norm_testitems
      found_not_expected = norm_testitems - codeframe[test][testlet]
      if !expected_not_found.empty? || !found_not_expected.empty?
        rpt << [psi, test, testlet, participationcode, expected_not_found.join(";"), found_not_expected.join(";")]
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
  norm_testitems = normalise_testitems(testitems)
  expected_not_found = codeframe[test][testlet] - norm_testitems
  found_not_expected = norm_testitems - codeframe[test][testlet]
  if !expected_not_found.empty? || !found_not_expected.empty?
    rpt << [psi, test, testlet, participationcode, expected_not_found.join(";"), found_not_expected.join(";")]
  end
end

