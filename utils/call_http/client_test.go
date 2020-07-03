package call_http

import (
	"encoding/json"
	"testing"
)

func TestHttpGet(t *testing.T) {
	path := "/test/v1/do"
	para := "p1=testv1&p2=测试"
	_, s, err := HttpGet("localhost:8080", "", path, para)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}
	t.Logf("success: %s", s)
}

func TestHttpPost(t *testing.T) {
	type TestS struct {
		Name string
		Age  int
	}

	d := TestS{
		Name: "testName",
		Age:  20,
	}

	b, _ := json.Marshal(d)
	path := "/test/v1/do"

	_, s, err := HttpPost("localhost:8080", "", path, b)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}
	t.Logf("success: %s", s)
}
