package media_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/anhthii/go-echo/router"
	"github.com/xeipuuv/gojsonschema"
)

func TestGetTop100Song(t *testing.T) {
	URL := "/api/media/top100/ZWZB96AB?page=1"
	testMedia(t, URL, "top100.json")
}

func TestGetDefaultAlbum(t *testing.T) {
	URL := "/api/media/albums/default"
	testMedia(t, URL, "default_album.json")
}

func TestGetChart(t *testing.T) {
	URL := "/api/media/chart/vpop"
	testMedia(t, URL, "chart.json")
}

func testMedia(t *testing.T, URL, jsonChemaFileName string) {
	folderpath, _ := os.Getwd()
	schemaPath := fmt.Sprintf("file://%s", path.Join(folderpath, "../../jsonschema", jsonChemaFileName))
	req, _ := http.NewRequest("GET", URL, nil)
	w := httptest.NewRecorder()
	r := router.SetupRouter("prod")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	var response map[string]interface{}
	bytes, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(bytes, &response)

	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	documentLoader := gojsonschema.NewGoLoader(response)
	_result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Log("Jsonschema validation failed")
		t.Fail()
	}

	if _result.Valid() {
		t.Log("The document is valid\n")
	} else {
		t.Log("The document is not valid. see errors :\n")
		for _, desc := range _result.Errors() {
			t.Log(fmt.Sprintf("- %s\n", desc))
		}
		t.Fail()
	}
}
