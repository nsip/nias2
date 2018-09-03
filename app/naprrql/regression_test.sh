rm -rf out
mv "in/sample.xml.zip" .
cp sif.minimal.xml.zip "in"
go run naprrql.go -ingest
go run naprrql.go -report
go run naprrql.go -qa
go run naprrql.go -writingextract
go run naprrql.go -xml
rm "in/sif.minimal.xml.zip"
for file in $( find ./out -name "*.csv" 2> /dev/null); do
  echo "\n\n"
  echo $file
  sort test/$file | perl -pe 's/,"\[([^]"]+)\]"/sprintf(q#,"[%s]"#, join(", ", sort(split(m#, #, $1))))/eg' > 1
  sort $file  | perl -pe 's/,"\[([^]"]+)\]"/sprintf(q#,"[%s]"#, join(", ", sort(split(m#, #, $1))))/eg' > 2
  diff 1 2
done
for file in $( find ./out -name "*.xml" 2> /dev/null); do
  echo "\n\n"
  echo $file
  diff test/$file $file
done
mv sample.xml.zip "in"

