# New World Map Coordinates SSE

A simple service that exposes New World's player's position as SSE stream.

This program is compatible and can be used as a replacement for [ceN's CoordinateTracker](https://github.com/mpcen/newworld-coordinate-tracker).

It works by grabbing a screenshot of upper-right-hand corner of your screen, 
reading position using Tesseract OCR and exposing it as SSE stream on localhost:5000/events


## Development

It's ment for Windows but its possible to develop and test it on other os.

go get
build.bat