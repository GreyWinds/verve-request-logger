package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jarcoal/httpmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	defer os.Remove("unique_requests.log")
	mockRedis := miniredis.RunT(t)
	redisClient = redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockEndpoint := "http://example.com/callback"
	mockURL, _ := url.Parse(mockEndpoint)

	httpmock.RegisterResponder("POST", mockEndpoint, func(req *http.Request) (*http.Response, error) {
		var payload map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), payload["count"])
		return httpmock.NewStringResponse(200, `{"success": true}`), nil
	})

	id := 123
	err := HandleRequest(id, mockURL)
	assert.NoError(t, err)

	currentMinute := time.Now().Format("2006-01-02 15:04")
	idKey := fmt.Sprintf("request_id:%s:%d", currentMinute, id)
	exists := mockRedis.Exists(idKey)
	assert.True(t, exists)
}

func TestDuplicateRequest(t *testing.T) {
	defer os.Remove("unique_requests.log")
	mockRedis := miniredis.RunT(t)
	redisClient = redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockEndpoint := "http://example.com/callback"
	mockURL, _ := url.Parse(mockEndpoint)

	httpmock.RegisterResponder("POST", mockEndpoint, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("POST request should not be triggered for duplicate IDs")
		return nil, nil
	})

	id := 123

	err := HandleRequest(id, mockURL)
	assert.NoError(t, err)

	currentMinute := time.Now().Format("2006-01-02 15:04")
	idKey := fmt.Sprintf("request_id:%s:%d", currentMinute, id)
	exists := mockRedis.Exists(idKey)
	assert.True(t, exists)

	//second call with same id
	err = HandleRequest(id, mockURL)
	assert.NoError(t, err)

	//make sure the original key exists and is not set again
	exists = mockRedis.Exists(idKey)
	assert.True(t, exists)
}
