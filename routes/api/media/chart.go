package media

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/anhthii/zing.mp3/scraper"
	"github.com/gin-gonic/gin"
)

// GetChart return data for a chart
func GetChart(c *gin.Context) {
	var URL string
	var coverIndex int
	switch c.Param("popType") {
	case "pop":
		URL = "https://beta.mp3.zing.vn/api/chart/get-chart?id=IWZ9Z0BW&week=null&year=null&ctime=1541611842&sig=1577ac38465d4983fc4830a36738c0bb5b096bf87d4314ff5de8478b93f3e87b548b596fd33a3a9956e91550d9c4ba0f9af91a56a209cf2c450a60c1c7b2f245&api_key=38e8643fb0dc04e8d65b99994d3dafff"
		coverIndex = 1
	case "kpop":
		URL = "https://beta.mp3.zing.vn/api/chart/get-chart?id=IWZ9Z0BO&week=null&year=null&ctime=1541643608&sig=deb588bb6b9b22d5f31b5f3fbd2250d9105572b35d3d780b1425d34784767dcafb8bcb3f970913c3379892b1755948370a03518af25bb7812fdc510cbeaea3a5&api_key=38e8643fb0dc04e8d65b99994d3dafff"
		coverIndex = 2
	case "vpop":
		coverIndex = 0
		URL = "https://beta.mp3.zing.vn/api/chart/get-chart?id=IWZ9Z08I&week=null&year=null&ctime=1541643715&sig=f37d92d66d506da34039c4685598ba9f1c6cdda88c91279de81307b5ea786ff43491a45f56bbb35efb3fd4269adc76eca4abde4d806ed8392fbab704e4486dd0&api_key=38e8643fb0dc04e8d65b99994d3dafff"
	}
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var result map[string]interface{}
	bytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(bytes, &result)
	img := scraper.Init().Scrape("https://mp3.zing.vn").Doc.Find(".chart-song .zthumb .fn-thumb").Eq(coverIndex)
	src, _ := img.Attr("src")
	result["data"].(map[string]interface{})["cover"] = src
	result["data"].(map[string]interface{})["items"] = result["data"].(map[string]interface{})["items"].([]interface{})[:20]
	c.JSON(http.StatusOK, result)
}
