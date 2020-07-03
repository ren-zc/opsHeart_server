package call_http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"opsHeart/conf"
	"time"
)

var C *http.Client

func init() {
	C = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func HttpGet(host string, token string, path string, para string) (int, []byte, error) {
	urlWithPara := fmt.Sprintf("http://%s%s?%s", host, path, para)
	//fmt.Printf("%s", urlWithPara)
	req, _ := http.NewRequest("GET", urlWithPara, nil)

	// get uuid of self from conf
	u := conf.GetUUID()

	// set basic auth
	req.SetBasicAuth(u, token)

	resp, err := C.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, b, nil
}

func HttpPost(host string, token string, path string, para []byte) (int, []byte, error) {
	url := fmt.Sprintf("http://%s%s", host, path)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(para))

	// get uuid of self from conf
	u := conf.GetUUID()

	// set basic auth
	req.SetBasicAuth(u, token)

	resp, err := C.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, b, nil
}
