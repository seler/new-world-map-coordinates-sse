#include "ocr_tesseract.h"
#include <leptonica/allheaders.h>
#include <tesseract/baseapi.h>

TessBaseAPI TessNew(){
  tesseract::TessBaseAPI *api = new tesseract::TessBaseAPI();

  return (TessBaseAPI)api;
}

void TessInit(TessBaseAPI api_) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  api->SetPageSegMode(tesseract::PageSegMode::PSM_SINGLE_LINE);
  api->SetVariable("tessedit_char_whitelist", "1234567890[],.");
  api->SetVariable("load_system_dawg", "F");
  api->SetVariable("load_freq_dawg", "F");

  if (api->Init(NULL, "eng")) {
      fprintf(stderr, "Could not initialize tesseract.\n");
      exit(1);
  }
}

void prepareImage(Pix *img) {
  img = pixScale(img, 8, 8);

  int r = 0x9a, g = 0x9c, b = 0x7b, t = 16;
  Pix *masked = pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t);

  r = 0xff, g = 0xff, b = 0xbb, t = 32;
  pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  r = 0xd7, g= 0xd7, b=0x9d, t=32;
  pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  r = 0x97, g= 0x97, b=0x70, t=16;
  pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  r = 0xb7, g= 0xdb, b=0xb5, t=32;
  pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  img = pixReduceRankBinary2(masked, 2, NULL);

  img = pixReduceRankBinary2(img, 2, NULL);
}

char* TessGetText(TessBaseAPI api_, unsigned char* imageBytes, int size) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  Pix *image = pixReadMemPng(imageBytes, (size_t)size);
  prepareImage(image);

  api->SetImage(image);
  api->SetSourceResolution(300);
  char *text = api->GetUTF8Text();
  pixDestroy(&image);

  return text;
}

void TessEnd(TessBaseAPI api_) {
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;
  api->Clear();
  api->End();
  delete api;
}