require "csv"
require "pp"

# 17. Incomplete test attempts

CSV.open("../../out/qa/systemTestAttempts.rpt.csv", "wb", 
        write_headers: true,
        headers: ["Participation Code","ACARA ID","School Name","Test Level","Test Domain","Family Name","Given Name","Birth Date","PSI"]
) do |rpt|
  CSV.foreach("../../out/qa/systemTestAttempts.csv", headers: true) do |row|
    next if row["Participation Code"] != "S"
    rpt << row
  end
end

