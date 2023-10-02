package gopenslide

import (
	"fmt"
	"io"
	"log"
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
	wsi, err := WSI("")
	if err == nil {
		t.Fatalf("Expected this to fail and return error but got nill as error")
	}
	if wsi != nil {
		t.Fatalf("expected this to be nil but it is something %#v", wsi)
	}
}

func Test_WSI_not_found(t *testing.T) {
	wsi, err := WSI("somethign that is not there")
	if err == nil {
		t.Fatalf("Expected this to fail and return error but got nill as error")
	}
	if !strings.Contains(err.Error(), "Not Found") {
		t.Fatalf("Expected a Not Found error but got %s", err.Error())
	}
	if wsi != nil {
		t.Fatalf("expected this to be nil but it is something %#v", wsi)
	}
}

func Test_WSI(t *testing.T) {
	wsi, err := WSI("testData/JP2K-33003-1.svs")
	if err != nil {
		t.Fatalf("What now %v", err)
	}
	if wsi == nil {
		t.Fatalf("expected this not to be nil but it is ")
	}

}
