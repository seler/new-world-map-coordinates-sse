set logLevel=debug
set destination=build\
set fileName=map-coordinates-sse.exe

echo y | rmdir %destination% /S
mkdir %destination%
copy deps\* %destination%
mkdir %destination%tessdata\
copy tessdata\ %destination%tessdata\

for /f %%i in ('git describe --tags --dirty') do set version=%%i

go build -ldflags "-X main.version=%version% -X main.logLevel=%logLevel%" -o %destination%%fileName%