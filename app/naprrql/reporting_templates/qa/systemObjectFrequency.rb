require "csv"
require "pp"

summaries = {}
CSV.foreach("../../out/qa/systemScoreSummaries.csv", headers: true) do |row|
  summaries[row["School ACARA ID"]] = [] if summaries[row["School ACARA ID"]].nil?
  summaries[row["School ACARA ID"]] << "#{row["Test Level"]}:#{row["Test Domain"]}"
end
CSV.open("../../out/qa/systemTestSummaryCount.rpt.csv", "wb",write_headers: true,
         headers: ["School ACARA ID","Summary Count","Summaries"]) do |rpt|
  summaries.each do |k, v|
    rpt << [k, v.size, v.sort.join(";")]
  end
end
