set logLevel="debug"

echo y | rmdir release\ /S
mkdir release\
copy deps\* release\
mkdir release\tessdata\
copy tessdata\ release\tessdata\

for /f %%i in ('git describe --tags --dirty') do set version=%%i

go build -ldflags "-X main.version=%version% -X main.logLevel=%logLevel%" -o release\