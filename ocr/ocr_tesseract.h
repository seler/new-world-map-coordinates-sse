#ifdef __cplusplus
extern "C" {
#endif

typedef void* TessBaseAPI;

TessBaseAPI TessNew();
void TessInit(TessBaseAPI api_);
char *TessGetText(TessBaseAPI api_, unsigned char* imageBytes, int size);
void TessEnd(TessBaseAPI api_);

#ifdef __cplusplus
}
#endif/* extern "C" */
