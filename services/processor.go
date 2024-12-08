package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient     *redis.Client
	ctx             = context.Background()
	redisKeyPattern = "request_id:%s:*"
	OpenFile        = os.OpenFile //monkey patch for unit testing
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", //ideally would be passed while starting the application or from a secret store
	})

	logFile, err := OpenFile("unique_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)

	go logUniqueRequests()
}

func HandleRequest(id int, endpoint *url.URL) error {
	currentMinute := time.Now().Format("2006-01-02 15:04")
	idKey := fmt.Sprintf("request_id:%s:%d", currentMinute, id) // Creates a unique key for the id with the current minute

	isRedisKeySet, err := redisClient.SetNX(ctx, idKey, true, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to set Id in Redis: %v", err)
	}

	if isRedisKeySet {
		log.Printf("New unique ID processed: %d", id)
	} else {
		log.Printf("duplicate ID detected for the current minute: %d", id)
	}

	if endpoint != nil {
		go func(endpoint *url.URL) {
			pattern := fmt.Sprintf(redisKeyPattern, currentMinute)
			keys, err := redisClient.Keys(ctx, pattern).Result()
			if err != nil {
				log.Printf("failed to get keys from Redis: %v", err)
				return
			}
			count := len(keys)

			query := endpoint.Query()
			query.Set("count", strconv.Itoa(count))
			endpoint.RawQuery = query.Encode()

			inputMap := map[string]interface{}{
				"count": count,
			}

			payload, _ := json.Marshal(inputMap)
			contentType := "application/json"
			resp, err := http.Post(endpoint.String(), contentType, bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("failed to make POST request: %v", err)
			}
			log.Printf("POST %s - Status: %d", endpoint.String(), resp.StatusCode)
			defer resp.Body.Close()
		}(endpoint)
	}

	return nil
}

// logging count of unique id every min
func logUniqueRequests() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentMinute := time.Now().Add(-time.Minute).Format("2006-01-02 15:04")
		pattern := fmt.Sprintf(redisKeyPattern, currentMinute)

		keys, err := redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			log.Printf("failed to fetch keys from Redis: %v", err)
			continue
		}

		count := len(keys)
		log.Printf("Unique requests for %s: %d", currentMinute, count)
	}
}
