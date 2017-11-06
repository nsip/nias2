require "csv"
require "pp"

# 13. Test object frequency

CSV.open("../../out/qa/orphan.rpt.csv", "wb", write_headers: true, headers: ["ACARA ID","School RefID", "Object Type","Object RefID"]) do |rpt|
  CSV.foreach("../../out/qa/orphanEvents.csv", headers: true) do |row|
    next if row["EventID"].nil?
    rpt << [row["SchoolID"], row["SchoolRefID"], "NAPEventStudentLink", row["EventID"]]
  end
  CSV.foreach("../../out/qa/orphanScoreSummaries.csv", headers: true) do |row|
    next if row["SummaryID"].nil?
    rpt << [row["SchoolACARAId"], row["SchoolInfoRefID"], "NAPTestSchoolSummary", row["SummaryID"]]
  end
  CSV.foreach("../../out/qa/orphanStudents.csv", headers: true) do |row|
    next if row["RefId"].nil?
    rpt << [row["ASLSchoolId"], "", "StudentPersonal", row["RefId"]]
  end
end
