set destination=test_ocr_tesseract\
set fileName=test_ocr_tesseract.exe

echo y | rmdir %destination% /S
mkdir %destination%
copy ..\deps\* %destination%
mkdir %destination%tessdata\
copy ..\tessdata\ %destination%tessdata\

go build -o %destination%%fileName%

%destination%%fileName%