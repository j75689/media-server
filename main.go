package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/video/list", showMediaList)
	apiV1.GET("/video/player/:id", viedoStream)
	fmt.Println(router.Run(":58808"))
}

var videoInfoMapping = make(map[string]interface{})            // file name index
var videoHashMapping = make(map[string]map[string]interface{}) // hash id index

func showMediaList(c *gin.Context) {
	files, err := ioutil.ReadDir("./media")
	if err != nil {
		c.JSON(200, map[string]string{"error": err.Error()})
		return
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".mp4") && videoInfoMapping[f.Name()] == nil {
			hash := StringWithCharset(10)
			videoInfoMapping[f.Name()] = hash
			videoHashMapping[hash] = map[string]interface{}{
				"name": f.Name(),
				"size": f.Size(),
				"id":   hash,
				"url":  fmt.Sprintf("/api/v1/video/player/%s", hash),
			}
		}
	}
	count := 0
	list := []map[string]interface{}{}
	for _, data := range videoHashMapping {
		count++
		list = append(list, data)
	}
	c.JSON(200, map[string]interface{}{"count": count, "list": list})
}

func viedoStream(c *gin.Context) {
	id := c.Param("id")
	videoInfo := videoHashMapping[id]
	if videoInfo == nil {
		c.JSON(200, map[string]string{"error": "file not existed"})
		return
	}
	fileName := "./media/" + videoInfo["name"].(string)
	video, err := os.Open(fileName)
	if err != nil {
		c.JSON(200, map[string]string{"error": err.Error()})
		return
	}
	defer video.Close()

	c.Header("Content-type", "video/mp4")
	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), video)

}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
