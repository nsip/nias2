require "csv"
require "pp"

# 15. Test type impacts

CSV.open("../../out/qa/systemTestTypeImpacts.rpt.csv", "wb",write_headers: true,
        headers: ["ACARA ID","School Name","Test Level","Test Domain","Participation Code","Family Name","Given Name","Birth Date","PSI","Path Taken For Domain","Parallel Test","Raw Score","Error"]) do |rpt|
  CSV.foreach("../../out/qa/systemResponses.csv", headers: true) do |row|
    if (!row["Path Taken For Domain"].nil? || !row["Parallel Test"].nil?) && row["Test Domain"] == "Writing"
      row["Error"] = "Writing test with adaptive structure"
    end
    if (row["Path Taken For Domain"].nil? || row["Parallel Test"].nil?) && ["P", "S"].include?(row["Participation Code"]) && row["Test Domain"] != "Writing"
      row["Error"] = "Non-Writing test with non-adaptive structure"
    end
    if !row["Error"].nil?
      rpt << row
    end
  end
end
CSV.open("../../out/qa/systemTestTypeItemImpacts.rpt.csv", "wb",
        write_headers:true,
        headers: ["Test Name","Test level","Test Domain","Test Item Local Id","Test Item Name","Subdomain","Writing Genre","Birth Date","ACARA ID","PSI","Testlet Score","Item Score","Item Lapsed Time","Item Subscores","Item Response","Participation Code"]) do |rpt|
  CSV.foreach("../../out/qa/itemPrinting.csv", headers: true) do |row|
    if (row["Item Subscores"] == "[]" ) && row["Test Domain"] == "Writing"
      row["Error"] = "No subscores for Writing test"
    end
    if (row["Item Subscores"] != "[]" ) && row["Test Domain"] != "Writing"
      row["Error"] = "Subscores for Non-Writing test"
    end
    if !row["Error"].nil?
      rpt << row
    end
  end
end
