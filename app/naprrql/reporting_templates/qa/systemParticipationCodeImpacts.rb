require "csv"
require "pp"

CSV.open("../../out/qa/systemParticipationCodeImpacts.rpt.csv", "wb",
        write_headers: true,
        headers: ["ACARA ID","School Name","Test Level","Test Domain","Participation Code","Family Name","Given Name","Birth Date","PSI","Path Taken For Domain","Parallel Test","Raw Score","ResponseID","Error"]) do |rpt|
  CSV.foreach("../../out/qa/systemResponses.csv", headers: true) do |row|
    #row["Error"] = ""
    if (!row["Path Taken For Domain"].nil? || !row["Parallel Test"].nil?) && !["P", "S"].include?(row["Participation Code"])
      row["Error"] = "Adaptive pathway without student writing test"
    end
    if (!row["Raw Score"].nil?) && !["P", "R"].include?(row["Participation Code"])
      row["Error"] = "Scored test with status other than P or R"
    end
    if (row["Raw Score"].nil?) && ["P", "R"].include?(row["Participation Code"])
      row["Error"] = "Unscored test with status of P or R"
    end
    if !row["Error"].nil?
      rpt << row
    end
  end
end
CSV.open("../../out/qa/systemParticipationCodeItemImpacts.rpt.csv", "wb",
        write_headers:true,
        headers: ["Test Name","Test level","Test Domain","Test Item Local Id","Test Item Name","Subdomain","Writing Genre","Birth Date","ACARA ID","PSI","Testlet Score","Item Score","Item Lapsed Time","Item Subscores","Item Response","Participation Code"]) do |rpt|
  CSV.foreach("../../out/qa/itemPrinting.csv", headers: true) do |row|
    #row["Error"] = ""
    if (!row["Item Lapsed Time"].nil? || !row["Item Response"].nil?) && !["P", "S"].include?(row["Participation Code"])
      row["Error"] = "Response captured without student writing test"
    end
    if (!row["Testlet Score"].nil? || !row["Item Score"].nil? || row["Item Subscores"] != "[]") && !["P", "R"].include?(row["Participation Code"])
      row["Error"] = "Scored test with status other than P or R"
    end
    if (row["Testlet Score"].nil? || row["Item Score"].nil? ) && ["P", "R"].include?(row["Participation Code"])
      row["Error"] = "Unscored test with status of P or R"
    end
    if (row["Item Subscores"] == "[]" && row["Test Domain"] == "Writing" ) && ["P", "R"].include?(row["Participation Code"])
      row["Error"] = "Unscored writing test with status of P or R"
    end
    if !row["Error"].nil?
      rpt << row
    end
  end
end

