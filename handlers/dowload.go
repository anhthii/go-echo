package handlers

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
)

// DownloadSong provide clients the API to download a song with ajax request
func DownloadSong(c *gin.Context) {
	songTitle := c.Param("songTitle")
	id := c.Param("id")

	// Get song resouce URL
	pageURL := fmt.Sprintf("https://mp3.zing.vn/bai-hat/%s/%s.html", songTitle, id)
	r := regexp.MustCompile(`key=\w+`)
	s, err := utils.GetStringDataFromHTTPGet(pageURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}

	key := r.FindStringSubmatch(s)[0]
	songResourceURL := fmt.Sprintf("https://mp3.zing.vn/xhr/media/get-source?type=audio&%s", key)

	result, err := utils.GetMapDataFromHTTPGet(songResourceURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}

	songURL := result["data"].(map[string]interface{})["source"].(map[string]interface{})["128"].(string)
	fmt.Println(songURL)
	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(songURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}
	defer resp.Body.Close()

	reader := resp.Body
	contentLength := resp.ContentLength
	contentType := resp.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=%s.mp3", songTitle),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
