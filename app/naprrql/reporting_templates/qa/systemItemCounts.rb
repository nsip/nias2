require "csv"
require "pp"

# 12. Test alternate item references

item = {}
CSV.foreach("../../out/qa/itemPrinting.csv", headers: true) do |row|
  next if row["Response Correctness"] == "Not In Path"
  next unless ["P", "S"].include?(row["Participation Code"])
  item[row["Test Name"]] = {} if item[row["Test Name"]].nil?
  item[row["Test Name"]][row["Test level"]] = {} if item[row["Test Name"]][row["Test level"]].nil?
  item[row["Test Name"]][row["Test level"]][row["Test Domain"]] = {} if item[row["Test Name"]][row["Test level"]][row["Test Domain"]].nil?
  item[row["Test Name"]][row["Test level"]][row["Test Domain"]][row["Test Item Local Id"]] = {count: 0} if item[row["Test Name"]][row["Test level"]][row["Test Domain"]][row["Test Item Local Id"]].nil?
  item[row["Test Name"]][row["Test level"]][row["Test Domain"]][row["Test Item Local Id"]][:count] = item[row["Test Name"]][row["Test level"]][row["Test Domain"]][row["Test Item Local Id"]][:count] + 1
  item[row["Test Name"]][row["Test level"]][row["Test Domain"]][row["Test Item Local Id"]][:substitute] = (row["ItemSubstitutedForList"] != "{\"SubstituteItem\":[]}")
end


CSV.open("../../out/qa/systemItemCounts.rpt.csv", "wb",write_headers: true,
         headers: ["Test Name","Test level","Test Domain","Test Item Local Id","Substitute","Count"]) do |rpt|
  item.each do |k, v|
    v.each do |k1,v1|
      v1.each do |k2,v2|
        v2.each do |k3, v3|
          rpt << [k, k1, k2, k3, v3[:substitute], v3[:count] ]
        end
      end
    end
  end
end
