package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func TestGetHandler(c *gin.Context) {
	q, _ := url.ParseQuery(c.Request.URL.RawQuery)
	fmt.Println(q)
	h := gin.H{}
	for k, v := range q {
		h[k] = v
	}
	c.JSON(http.StatusOK, h)
}

func TestPostHandler(c *gin.Context) {
	d, _ := ioutil.ReadAll(c.Request.Body)
	defer func() {
		_ = c.Request.Body.Close()
	}()
	c.JSON(http.StatusOK, gin.H{
		"data": string(d),
	})
}
