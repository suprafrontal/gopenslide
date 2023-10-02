package gopenslide

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"testing"
)

func init() {
	// make sure testData  exists
	os.Mkdir("testData", os.ModePerm)
	// download test data
	if _, err := os.Stat("testData/JP2K-33003-1.svs"); err == nil {
		fmt.Printf("Test Data already available moving on")
	} else {
		wsi1 := "https://openslide.cs.cmu.edu/download/openslide-testdata/Aperio/JP2K-33003-1.svs"
		out, err := os.Create("testData/JP2K-33003-1.svs")
		if err != nil {
			log.Fatalf("error creating test data file %v", err)
		}
		defer out.Close()
		resp, err := http.Get(wsi1)
		if err != nil {
			log.Fatalf("error downloading test data file %v", err)
		}
		defer resp.Body.Close()
		n, err := io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("error saving downloaded test data file %v", err)
		}
		log.Printf("successfully downloaded %d bytes of test file \n", n)
	}
}

func Test_WSI_empty(t *testing.T) {
	wsi, err := OpenWSI("")
	if err == nil {
		t.Fatalf("Expected this to fail and return error but got nill as error")
	}
	if wsi.osr != nil {
		t.Fatalf("expected this to be nil but it is something %#v", wsi)
	}
}

func Test_WSI_not_found(t *testing.T) {
	wsi, err := OpenWSI("somethign that is not there")
	if err == nil {
		t.Fatalf("Expected this to fail and return error but got nill as error")
	}
	if !strings.Contains(err.Error(), "Not Found") {
		t.Fatalf("Expected a Not Found error but got %s", err.Error())
	}
	if wsi.osr != nil {
		t.Fatalf("expected this to be nil but it is something %#v", wsi)
	}
}

func Test_WSI(t *testing.T) {
	wsi, err := OpenWSI("testData/JP2K-33003-1.svs")
	if err != nil {
		t.Fatalf("What now %v", err)
	}
	if wsi.osr == nil {
		t.Fatalf("expected this not to be nil but it is ")
	}

	if wsi.GetLevelCount() != 3 {
		t.Fatalf("expected this to be 9 but it is %d", wsi.GetLevelCount())
	}

	if wsi.GetLevelDownsample(0) != 1.0 {
		t.Fatalf("expected this to be 1.0 but it is %f", wsi.GetLevelDownsample(0))
	}

	if math.Ceil(wsi.GetLevelDownsample(1)*1000000) != 4000375 {
		t.Fatalf("expected this to be 4.000375 but it is %f", wsi.GetLevelDownsample(1))
	}

	if math.Ceil(wsi.GetLevelDownsample(2)*1000000) != 8001791 {
		t.Fatalf("expected this to be 8.001790 but it is %f", math.Ceil(wsi.GetLevelDownsample(2)*1000000))
	}

	w, h, err := wsi.GetLevelDimensions(0)
	if err != nil {
		t.Fatalf("expected this to be nil but it is %v", err)
	}
	if w != 15374 {
		t.Fatalf("expected this to be 97792 but it is %d", w)
	}
	if h != 17497 {
		t.Fatalf("expected this to be 221184 but it is %d", h)
	}

	region, err := wsi.ReadRegion(0, 0, 0, 100, 100)
	if err != nil {
		t.Fatalf("expected this to be nil but it is %v", err)
	}
	if len(region) != 40000 {
		t.Fatalf("expected this to be 40000 but it is %d", len(region))
	}
	if region[0] != 255 {
		t.Fatalf("expected this to be 0 but it is %d", region[0])
	}
}
