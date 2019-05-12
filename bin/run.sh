name=scrawl-verify
rm $name
go build ../src/$name
./$name  -beautify  -verbose -protocol protocol.txt -output protocol.out.txt
