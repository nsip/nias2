require "csv"
require "pp"

counts = {}
CSV.foreach("../../out/qa/systemTestAttempts.csv", headers: true) do |row|
  counts[row["ACARA ID"]] = {} if counts[row["ACARA ID"]].nil?
  counts[row["ACARA ID"]][row["Test Domain"]] = {} if counts[row["ACARA ID"]][row["Test Domain"]].nil?
  counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]] = 
    {P_attempts: [], S_attempts: [], responses: []} if counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]].nil?
  next unless ["P", "S"].include?(row["Participation Code"])
  if row["Participation Code"] == "P"
    counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]][:P_attempts] << row["PSI"]
  else
    counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]][:S_attempts] << row["PSI"]
  end
end
CSV.foreach("../../out/qa/systemResponses.csv", headers: true) do |row|
  counts[row["ACARA ID"]] = {} if counts[row["ACARA ID"]].nil?
  counts[row["ACARA ID"]][row["Test Domain"]] = {} if counts[row["ACARA ID"]][row["Test Domain"]].nil?
  counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]] = 
    {P_attempts: [], S_attempts: [], responses: [] } if counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]].nil?
  next unless ["P", "S"].include?(row["Participation Code"])
  counts[row["ACARA ID"]][row["Test Domain"]][row["Test Level"]][:responses] << row["PSI"] unless row["ResponseID"].nil?
  # the XML does not represent writing scripts distinctly
end


CSV.open("../../out/qa/systemTestCompleteness.rpt.csv", "wb",
         write_headers: true,
         headers: ["ACARA ID", "Test Domain", "Test Level", "# P attempts", "# S attempts", "# Responses", "Attempt with no response","Response with no attempt"]) do |rpt|
  counts.each do |k, v1|
    v1.each do |k1, v2|
      v2.each do |k2, v3|
        p_attempts =  v3[:P_attempts]
        s_attempts =  v3[:S_attempts]
        attempts = p_attempts + s_attempts
        responses =  v3[:responses]
        rpt << [k, k1, k2, p_attempts.size, s_attempts.size, responses.size, (attempts - responses).sort.join(";"), (responses - attempts).sort.join(";"),]
      end
    end
  end
end
