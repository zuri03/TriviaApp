package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	//"github.com/elliotchance/redismock"
	"github.com/go-redis/redis/v8"
	"github.com/zuri03/user/models"
	//"github.com/stretchr/testify/assert"
)

type TestSet struct {
	Model    models.UserDetails
	Request  *http.Request
	Code     int
	Response string
}

func TestServeHtTTP(t *testing.T) {
	emptyModel := models.UserDetails{
		Username: "",
		Password: "",
	}

	bytesArr, err := json.Marshal(emptyModel)
	if err != nil {
		log.Fatalf("Error parsing buffer: %s\n", err.Error())
	}

	emptyModelRequest := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(bytesArr))

	properModel := models.UserDetails{
		Username: "testUsername",
		Password: "testPassword",
	}

	bytesArr, err = json.Marshal(properModel)
	if err != nil {
		log.Fatalf("Error parsing buffer: %s\n", err.Error())
	}

	properModelRequest := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(bytesArr))

	tests := [2]TestSet{
		TestSet{
			Model:    emptyModel,
			Request:  emptyModelRequest,
			Code:     http.StatusBadRequest,
			Response: "Missing username or password",
		},
		TestSet{
			Model:    properModel,
			Request:  properModelRequest,
			Code:     http.StatusOK,
			Response: "{\"username\": \"testUsername\", \"password\": \"testPassword\"}",
		},
	}

	miniRedis, err := miniredis.Run()

	redisClient := redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})

	createHandler := CreateHandler{
		RedisHandler: redisClient,
		Signaler:     make(chan os.Signal, 1),
		Ctx:          context.Background(),
	}

	for _, test := range tests {
		t.Run("Starting request", func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()

			createHandler.ServeHTTP(responseRecorder, test.Request)

			if responseRecorder.Code != test.Code {
				t.Errorf("Needed %d, got %d\n", test.Code, responseRecorder.Code)
			}

			if responseRecorder.Body.String() != test.Response {
				t.Errorf("Wanted %s, got %s\n", test.Response, responseRecorder.Body)
			}
		})
	}
}
