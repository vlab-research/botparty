package botparty

import (
	"encoding/json"
	"net/http"
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

type TestTransport func(req *http.Request) (*http.Response, error)

func (r TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}


func TestSendNilOn200(t *testing.T) {
	expectedEvent := `{"user":"foo","page":"bar","event":{"type":"baz","value":"hello"}}` 
 	testFn := func(req *http.Request) (*http.Response, error) {
		data, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, expectedEvent, string(data))

		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	bp := &BotParty{Botserver: "foo", Client: &http.Client{Transport: TestTransport(testFn)}}

	b, _ := json.Marshal("hello")
	rm := json.RawMessage(b)
	err := bp.Send(NewExternalEvent("foo", "bar", "baz", &rm))
	
	assert.Nil(t, err)
}


func TestSendErrorsOnNon200Code(t *testing.T) {
 	testFn := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound}, nil
	}

	bp := &BotParty{Botserver: "foo", Client: &http.Client{Transport: TestTransport(testFn)}}

	b, _ := json.Marshal("hello")
	rm := json.RawMessage(b)
	err := bp.Send(NewExternalEvent("foo", "bar", "baz", &rm))
	
	assert.NotNil(t, err)
}
