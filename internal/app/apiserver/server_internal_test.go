package apiserver

import (
	"GO-REST-API/internal/app/store/sqlStore"
	"bytes"
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)


var (
	configPath string
)

func TestServer_handleNewRequest(t *testing.T){

	flag.Parse()
	config := NewConfig()


	db, _ := newDB(config.DatabaseURL)


	defer db.Close()
	s := newServer(sqlStore.New(db))

	testCases := []struct{
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"user_id":"1",
				"message":"hello",
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			body := &bytes.Buffer{}
			json.NewEncoder(body).Encode(tc.payload)

			req,_ := http.NewRequest(http.MethodPost,"/newrequest",body)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}
