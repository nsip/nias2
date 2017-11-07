require "csv"
require "pp"

# 13. Test object frequency

summaries = {}
r = {}
CSV.foreach("../../out/qa/systemScoreSummaries.csv", headers: true) do |row|
  summaries[row["School ACARA ID"]] = [] if summaries[row["School ACARA ID"]].nil?
  summaries[row["School ACARA ID"]] << "#{row["Test Level"]}:#{row["Test Domain"]}"
end
=begin
CSV.foreach("../../out/qa/systemTestAttempts.csv", headers: true) do |row|
  r[row["PSI"]] = {events: [], responseful_events: [], responses: []} if r[row["PSI"]].nil?
  next if row["EventID"].nil?
  r[row["PSI"]][:events] << "#{row["EventID"]}:#{row["Test Level"]}:#{row["Test Domain"]}"
  if ["P", "R", "S"].include?(row["Participation Code"])
    r[row["PSI"]][:responseful_events] << "#{row["EventID"]}:#{row["Test Level"]}:#{row["Test Domain"]}"
  end
end
CSV.foreach("../../out/qa/systemResponses.csv", headers: true) do |row|
  r[row["PSI"]] = {events: [], responseful_events: [], responses: []} if r[row["PSI"]].nil?
  next if row["ResponseID"].nil?
  r[row["PSI"]][:responses] << "#{row["ResponseID"]}:#{row["Test Level"]}:#{row["Test Domain"]}"
end
=end
CSV.open("../../out/qa/systemObjectSummaryFrequency.rpt.csv", "wb",write_headers: true,
         headers: ["School ACARA ID","Object Type","Object Count","Objects"]) do |rpt|
  summaries.each do |k, v|
    rpt << [k, "NAPSchoolScoreSummary", v.size, v.sort.join(";")]
  end
end
=begin
CSV.open("../../out/qa/systemObjectFrequency.rpt.csv", "wb", write_headers: true,
         headers: ["PSI","# Event Counts","# Event Counts (PRS)","# Responses","Event Counts","Event Counts (PRS)","Responses"]) do |rpt|


  r.each do |k, v|
    rpt << [k, v[:events].size,  v[:responseful_events].size, v[:responses].size, v[:events].sort.join(";"), v[:responseful_events].sort.join(";"), v[:responses].sort.join(";"),]

  end
end
=end
