// go build -o libopenslide.so -buildmode=c-shared openslide.go
//go:build linux || darwin || windows
// +build linux darwin windows

package gopenslide

//#include <stdio.h>
//#include <stdlib.h>
//#include <stdint.h>

/*
#cgo CFLAGS: -g -Wall -I${SRCDIR}/include/openslide
#cgo LDFLAGS: -L. -lopenslide
#include "openslide.h"
#include "openslide-features.h"
*/
import "C"

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"unsafe"
)

var NOT_FOUND_ERROR = fmt.Errorf("OpenSlideErr: Not Found")
var NOT_SUPPORTED_ERROR = fmt.Errorf("OpenSlideErr: Format Not Supported")
var NOT_FOUND_OR_SUPPORTED_ERROR = fmt.Errorf("OpenSlideErr: Format or File Not Supported or Found!")

type WSI struct {
	osr *C.openslide_t
}

type KVPair map[string]string

func OpenWSI(pathToFile string) (WSI, error) {
	osr, err := openOpenSlide(pathToFile)
	if err != nil {
		return WSI{}, err
	}
	return WSI{osr}, nil
}

func (wsi *WSI) GetLevelCount() int {
	return int(openslide_get_level_count(wsi.osr))
}

func (wsi *WSI) GetLevelDimensions(level int) (int64, int64, error) {
	return openslide_get_level_dimensions(wsi.osr, level)
}

func (wsi *WSI) GetLevelDownsample(level int) float64 {
	return float64(openslide_get_level_downsample(wsi.osr, level))
}

func (wsi *WSI) ReadRegion(level int, x int64, y int64, w int64, h int64) ([]byte, error) {
	return openslide_read_region(wsi.osr, level, x, y, w, h)
}

func (wsi *WSI) GetPropertyNames() []string {
	return openslide_get_property_names(wsi.osr)
}

func (wsi *WSI) GetPropertyValue(name string) string {
	return openslide_get_property_value(wsi.osr, name)
}

func (wsi *WSI) GetPropertyKVPairs() []KVPair {
	kvPairs := []KVPair{}
	for _, name := range wsi.GetPropertyNames() {
		kvPairs = append(kvPairs, KVPair{name: wsi.GetPropertyValue(name)})
	}
	return kvPairs
}

func (wsi *WSI) Close() {
	C.openslide_close(wsi.osr)
}

//-------------------------------------------------

func openOpenSlide(path string) (*C.openslide_t, error) {
	if _, err := os.Stat(path); err == nil {
		pathToSVS := C.CString(path)
		vendor := C.openslide_detect_vendor(pathToSVS)
		if vendor == nil {
			return nil, NOT_SUPPORTED_ERROR
		}
		osr := C.openslide_open(pathToSVS)
		if osr != nil {
			e := C.openslide_get_error(osr)
			if e != nil {
				err = fmt.Errorf("OpenSlideErr: %s", C.GoString(e))
			}
		}
		return osr, err
	} else {
		return nil, NOT_FOUND_ERROR
	}
}

func openslide_get_level_count(osr *C.openslide_t) C.int {
	return C.openslide_get_level_count(osr)
}

func openslide_get_level_dimensions(osr *C.openslide_t, level int) (int64, int64, error) {
	var w, h C.int64_t
	levelC := C.int(level)
	C.openslide_get_level_dimensions(osr, levelC, &w, &h)
	return int64(w), int64(h), nil
}

func openslide_get_level_downsample(osr *C.openslide_t, level int) float64 {
	levelC := C.int(level)
	return float64(C.openslide_get_level_downsample(osr, levelC))
}

func openslide_read_region(osr *C.openslide_t, level int, x int64, y int64, w int64, h int64) ([]byte, error) {
	buf := make([]byte, w*h*4)
	levelC := C.int(level)
	xC := C.int64_t(x)
	yC := C.int64_t(y)
	wC := C.int64_t(w)
	hC := C.int64_t(h)
	C.openslide_read_region(osr, (*C.uint32_t)(unsafe.Pointer(&buf[0])), xC, yC, levelC, wC, hC)
	return buf, nil
}

func openslide_get_property_names(osr *C.openslide_t) []string {
	cPropNames := C.openslide_get_property_names(osr)
	names := []string{}
	for _, s := range (*[1 << 28]*C.char)(unsafe.Pointer(cPropNames))[:2:2] {
		names = append(names, C.GoString(s))
	}
	return names
}

func openslide_get_property_value(osr *C.openslide_t, name string) string {
	nameC := C.CString(name)
	valueC := C.openslide_get_property_value(osr, nameC)
	return C.GoString(valueC)
}

//-------------------------------------------------
// Some helper functions to make life easier

func NewRGBAImageFromRegionInLevel(osr *C.openslide_t, level int, x, y, width, height int64) (*image.RGBA, error) {
	var err error
	imgBuf, _ := openslide_read_region(osr, level, x, y, width, height)
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	img.Pix, err = openslidePreMultipliedARGB2PixUINT8RGBA(imgBuf)
	return img, err
}

func openslidePreMultipliedARGB2PixUINT8RGBA(imgBuff []byte) ([]uint8, error) {
	pix := []uint8{}
	var err error
	buff := bytes.NewBuffer(imgBuff)
	for err == nil {
		pixel := buff.Next(4)
		if len(pixel) < 4 {
			err = io.EOF
			break
		}
		//                R         G         B         A
		a := pixel[3]
		r := pixel[2]
		g := pixel[1]
		b := pixel[0]
		if a != 0 && a != 255 {
			r = r * 255 / a
			g = g * 255 / a
			b = b * 255 / a
		}
		pix = append(pix, r, g, b, a)
	}
	if err == io.EOF {
		err = nil
	}
	return pix, err
}

// Pixel struct example
// type Pixel struct {
// 	R int
// 	G int
// 	B int
// 	A int
// }

// // img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
// func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
// 	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
// }

// type argb2rgba struct {
// 	buff bytes.Buffer
// }

// func (a *argb2rgba) Copy2PixUINT8() ([]uint8, error) {
// 	pix := []uint8{}
// 	var err error
// 	for err == nil {
// 		pixel := a.buff.Next(4)
// 		if len(pixel) < 4 {
// 			err = io.EOF
// 			break
// 		}
// 		//                R         G         B         A
// 		a := pixel[3]
// 		r := pixel[2]
// 		g := pixel[1]
// 		b := pixel[0]
// 		if a != 0 && a != 255 {
// 			r = r * 255 / a
// 			g = g * 255 / a
// 			b = b * 255 / a
// 		}
// 		pix = append(pix, r, g, b, a)
// 	}
// 	if err == io.EOF {
// 		err = nil
// 	}
// 	return pix, err
// }

// // Get the bi-dimensional pixel array
// func getPixels(file io.Reader) ([][]Pixel, error) {
// 	img, _, err := image.Decode(file)

// 	if err != nil {
// 		return nil, err
// 	}

// 	bounds := img.Bounds()
// 	width, height := bounds.Max.X, bounds.Max.Y

// 	var pixels [][]Pixel
// 	for y := 0; y < height; y++ {
// 		var row []Pixel
// 		for x := 0; x < width; x++ {
// 			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
// 		}
// 		pixels = append(pixels, row)
// 	}

// 	return pixels, nil
// }
