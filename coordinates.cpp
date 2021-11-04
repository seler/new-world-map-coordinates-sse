#include "coordinates.h"
#include <leptonica/allheaders.h>
#include <tesseract/baseapi.h>

TessBaseAPI TessNew(){
  tesseract::TessBaseAPI *api = new tesseract::TessBaseAPI();

  return (TessBaseAPI)api;
}

void TessInit(TessBaseAPI api_) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  if (api->Init(NULL, "eng")) {
      fprintf(stderr, "Could not initialize tesseract.\n");
      exit(1);
  }
}

char* TessGetText(TessBaseAPI api_, unsigned char* imageBytes, int size) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  Pix *image = pixReadMemPng(imageBytes, (size_t)size);

  api->SetImage(image);
  char *text = api->GetUTF8Text();
  pixDestroy(&image);

  return text;
}

void TessEnd(TessBaseAPI api_) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  api->End();
  delete api;
}
