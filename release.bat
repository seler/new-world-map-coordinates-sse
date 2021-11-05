ECHO "Release"
mkdir release\
go build -o release\
copy deps\* release\
mkdir release\tessdata\
copy tessdata\ release\tessdata\
ECHO "Done"
PAUSE