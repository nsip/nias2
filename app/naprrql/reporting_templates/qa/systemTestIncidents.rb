require "csv"
require "pp"

# 16. Test incidents

CSV.open("../../out/qa/systemTestIncidents.rpt.csv", "wb",
         write_headers: true,
         headers: ["Test Disruption","ACARA ID","School Name","Test Level","Test Domain","Family Name","Given Name","Birth Date","PSI"]
        ) do |rpt|
          CSV.foreach("../../out/qa/systemTestIncidents.csv", headers: true) do |row|
            next if row["Test Disruption"] == "[]"
            row["Test Disruption"] = row["Test Disruption"].gsub(/\[{"Event":/,"").gsub(/}\]/,"").gsub(/"/,"")
            rpt << row
          end
        end

