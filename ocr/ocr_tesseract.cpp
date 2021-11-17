#include "ocr_tesseract.h"
#include <leptonica/allheaders.h>
#include <tesseract/baseapi.h>

TessBaseAPI TessNew()
{
  tesseract::TessBaseAPI *api = new tesseract::TessBaseAPI();

  return (TessBaseAPI)api;
}

void TessInit(TessBaseAPI api_)
{
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  // api->SetPageSegMode(tesseract::PageSegMode::PSM_SINGLE_LINE);
  // api->SetVariable("tessedit_char_whitelist", "on1234567890[,.");
  // api->SetVariable("tessedit_char_blacklist", "L");
  // api->SetVariable("load_system_dawg", "F");
  // api->SetVariable("load_freq_dawg", "F");

  // if (api->Init("./", "eng")) {
  //     fprintf(stderr, "Could not initialize tesseract.\n");
  //     exit(1);
  // }

  char *configs[] = {"tessdata/eng.config"};
  int configs_size = 1;
  if (api->Init(NULL, "eng", tesseract::OEM_LSTM_ONLY, configs, configs_size, NULL, NULL, false))
  {
    fprintf(stderr, "Could not initialize tesseract.\n");
    exit(1);
  }

  // api->SetPageSegMode(tesseract::PageSegMode::PSM_SINGLE_LINE);
  // api->SetVariable("tessedit_char_whitelist", "on1234567890[,.");
  // api->SetVariable("tessedit_char_blacklist", "L");
}

// red start, red end, green start, green tend, blue start, blue end
int colorsToSelect[][6] = {
    {225, 255, 225, 255, 168, 187},
};

void save(Pix *img, int i)
{
  char filename[10];
  sprintf(filename, "img/%d.png", i);
  pixWritePng(filename, img, 0);
}

Pix *deduplicateOffColumns(Pix *pix)
{
  l_int32 i, j, k, w, h, wpl;
  l_uint32 *line, *data, pixel;
  bool columnHasOnPixel;
  l_int32 columnsWithoutOnPixels = 0;
  l_int32 columnsBetweenOnPixels = 8;

  pixGetDimensions(pix, &w, &h, NULL);

  Pix *newPix = pixCreate(w, h, 1);

  for (k = 0, i = 0; i < w; i++)
  {
    columnHasOnPixel = false;
    for (j = 0; j < h; j++)
    {
      pixGetPixel(pix, i, j, &pixel);
      if (pixel != 0)
      {
        columnHasOnPixel = true;
      }
      pixSetPixel(newPix, k, j, pixel);
    }
    if (!columnHasOnPixel)
    {
      columnsWithoutOnPixels++;
    }
    k++;

    if (columnsWithoutOnPixels >= columnsBetweenOnPixels && columnHasOnPixel)
    {
      i--; // read from same column again
      k -= columnsWithoutOnPixels - columnsBetweenOnPixels - 1; // write to proper column
      if (k < 0){
        k == 0;
      }
      columnsWithoutOnPixels = 0;
    }

  }

  return newPix;
}

Pix *prepareImage(Pix *img, int saveImages)
{
  int i = 0;
  if (saveImages) save(img, i++);

  img = pixScale(img, 4, 4);
  if (saveImages) save(img, i++);

  Pix *maskColor = pixMaskOverColorRange(img, 120, 255, 120, 255, 100, 200);
  if (saveImages) save(img, i++);

  maskColor = pixScale(maskColor, 2, 2);
  if (saveImages) save(img, i++);

  img = pixScale(img, 2, 2);
  if (saveImages) save(img, i++);

  Pix *hue = pixConvertRGBToHue(img);
  if (saveImages) save(img, i++);

  img = pixGenerateMaskByBand(hue, 30, 65, 1, 1);
  if (saveImages) save(img, i++);

  img = pixAnd(NULL, maskColor, img);
  if (saveImages) save(img, i++);

  pixOpenBrick(img, img, 3, 3);
  if (saveImages) save(img, i++);

  img = pixReduceRankBinary2(img, 2, NULL);
  if (saveImages) save(img, i++);

  img = deduplicateOffColumns(img);
  if (saveImages) save(img, i++);

  return img;

  // pixWritePng("hue.png", hue, 0);

  // img = pixMaskOverColorRange(img, 207, 255, 210, 255, 155, 200);
  // pixWritePng("color.png", img, 0);

  // img = pixOr(NULL, hue, img);
  // pixWritePng("final.png", img, 0);

  //img = pixUnsharpMasking(img, 1, 0.5);
  // int r = 0x9a, g = 0x9c, b = 0x7b, t = 16;
  // Pix *masked = pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t);

  // r = 0xff, g = 0xff, b = 0xbb, t = 32;
  // pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  // r = 0xd7, g= 0xd7, b=0x9d, t=32;
  // pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  // r = 0x97, g= 0x97, b=0x70, t=16;
  // pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  // r = 0xb7, g= 0xdb, b=0xb5, t=32;
  // pixOr(masked, masked, pixMaskOverColorRange(img, r-t, r+t, g-t, g+t, b-t, b+t));

  // img = pixReduceRankBinary2(masked, 2, NULL);

  // img = pixReduceRankBinary2(img, 2, NULL);
}

char *TessGetText(TessBaseAPI api_, unsigned char *imageBytes, int size, int saveImages)
{
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;

  Pix *image = pixReadMemPng(imageBytes, (size_t)size);
  image = prepareImage(image, saveImages);

  api->SetImage(image);
  api->SetSourceResolution(300);
  char *text = api->GetUTF8Text();
  pixDestroy(&image);

  return text;
}

void TessEnd(TessBaseAPI api_)
{
  tesseract::TessBaseAPI *api = (tesseract::TessBaseAPI *)api_;
  api->Clear();
  api->End();
  delete api;
}