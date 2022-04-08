package httpx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockHttp struct {
	postURL string
}

func (m *mockHttp) Post(url string, mapData interface{}, resp interface{}) error {
	m.postURL = url
	return nil
}

func (m *mockHttp) PostByte(url string, data []byte, resp interface{}) error {
	return nil
}

func (m *mockHttp) Get(url string, mapData interface{}, resp interface{}) error {
	return nil
}

func TestLocal(t *testing.T) {

	http := mockHttp{}
	lc := LocalClient{cf: LocalClientConfig{
		Token:  "996",
		Origin: "http://127.0.0.1:8010",
	}, http: &http}

	var resp interface{}
	err := lc.Request("/api/base/local/wx/config", map[string]interface{}{
		"name": "TT",
	}, &resp)
	assert.Nil(t, err)

	fmt.Println("URL:", http.postURL)

}
