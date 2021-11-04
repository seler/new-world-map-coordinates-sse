#include "coordinates.h"
#include <leptonica/allheaders.h>
#include <tesseract/baseapi.h>

char* GetText(unsigned char* imageBytes, int size) {
  Pix *image = pixReadMemPng(imageBytes, (size_t)size);
  tesseract::TessBaseAPI *api = new tesseract::TessBaseAPI();

  if (api->Init(NULL, "eng")) {
      fprintf(stderr, "Could not initialize tesseract.\n");
      exit(1);
  }

  api->SetImage(image);
  char *text = api->GetUTF8Text();

  api->End();
  delete api;
  pixDestroy(&image);
  return text;
}
