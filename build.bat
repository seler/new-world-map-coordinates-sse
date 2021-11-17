set logLevel=info
set destination=build\
set fileName=map-coordinates-sse.exe
set collectSamples=false
set saveImages=false

echo y | rmdir %destination% /S
mkdir %destination%
copy deps\* %destination%
mkdir %destination%tessdata\
copy tessdata\ %destination%tessdata\
copy userscripts\ %destination%userscripts\
copy run.bat %destination%run.bat
mkdir %destination%img\
mkdir %destination%samples\

for /f %%i in ('git describe --tags --dirty') do set version=%%i

go build -ldflags "-X main.version=%version% -X main.logLevel=%logLevel% -X main.collectSamples=%collectSamples% -X main.saveImages=%saveImages%" -o %destination%%fileName%