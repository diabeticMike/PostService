package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PostService/mocks"
	"github.com/PostService/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetPosts(t *testing.T) {
	type (
		payload struct {
			mockPostSvc func(mock *mocks.MockService)
			mockLogger  func(mock *mocks.MockLogger)
			qParams     map[string]string
			path        string
		}
		expected struct {
			body       string
			statusCode int
		}
	)
	var testCases = []struct {
		name     string
		payload  payload
		expected expected
	}{
		{
			name: "GetPostsByAuthorAndName (post_name and author provided) error",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					mock.EXPECT().GetPostsByNameAndAuthor("name1", "author1").Return(nil, errors.New("custom error"))
				},
				mockLogger: func(mock *mocks.MockLogger) {
					mock.EXPECT().Error("custom error")
				},
				qParams: map[string]string{
					"post_name": "name1",
					"author":    "author1",
				},
				path: "/post",
			},
			expected: expected{
				body:       "custom error\n",
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "GetPostsByAuthorAndName only post_name provided error",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					mock.EXPECT().GetPostsByKey("name1").Return(nil, errors.New("custom error"))
				},
				mockLogger: func(mock *mocks.MockLogger) {
					mock.EXPECT().Error("custom error")
				},
				qParams: map[string]string{
					"post_name": "name1",
				},
				path: "/post",
			},
			expected: expected{
				body:       "custom error\n",
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "GetPostsByAuthorAndName only author provided error",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					mock.EXPECT().GetPostsByKey("author1").Return(nil, errors.New("custom error"))
				},
				mockLogger: func(mock *mocks.MockLogger) {
					mock.EXPECT().Error("custom error")
				},
				qParams: map[string]string{
					"author": "author1",
				},
				path: "/post",
			},
			expected: expected{
				body:       "custom error\n",
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "posts lenth is 0 error",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					mock.EXPECT().GetPostsByKey("name1").Return([]model.Post{}, nil)
				},
				mockLogger: func(mock *mocks.MockLogger) {
				},
				qParams: map[string]string{
					"post_name": "name1",
				},
				path: "/post",
			},
			expected: expected{
				body:       "No posts found",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "order is true case",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					date1 := time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)
					date2 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
					posts := []model.Post{{Name: "name2", Date: date2, Author: "author2"},
						{Name: "name1", Date: date1, Author: "author1"}}
					mock.EXPECT().GetPostsByKey("name1").Return(posts, nil)
				},
				mockLogger: func(mock *mocks.MockLogger) {
				},
				qParams: map[string]string{
					"post_name": "name1",
					"order":     "true",
				},
				path: "/post",
			},
			expected: expected{
				body:       "[{\"post_name\":\"name1\",\"author\":\"author1\",\"date\":\"01.01.20\"},{\"post_name\":\"name2\",\"author\":\"author2\",\"date\":\"01.01.00\"}]",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "no order case",
			payload: payload{
				mockPostSvc: func(mock *mocks.MockService) {
					date1 := time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)
					date2 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
					posts := []model.Post{{Name: "name2", Date: date2, Author: "author2"},
						{Name: "name1", Date: date1, Author: "author1"}}
					mock.EXPECT().GetPostsByKey("name1").Return(posts, nil)
				},
				mockLogger: func(mock *mocks.MockLogger) {
				},
				qParams: map[string]string{
					"post_name": "name1",
				},
				path: "/post",
			},
			expected: expected{
				body:       "[{\"post_name\":\"name2\",\"author\":\"author2\",\"date\":\"01.01.00\"},{\"post_name\":\"name1\",\"author\":\"author1\",\"date\":\"01.01.20\"}]",
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", setParams(tc.payload.path, tc.payload.qParams), nil)
			if err != nil {
				t.Fatal(err)
			}
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockPostSvc := mocks.NewMockService(mockCtrl)
			tc.payload.mockPostSvc(mockPostSvc)
			mockLogger := mocks.NewMockLogger(mockCtrl)
			tc.payload.mockLogger(mockLogger)

			r := mux.NewRouter()
			pc := NewPostController(mockLogger, mockPostSvc)
			rr := httptest.NewRecorder()
			r.HandleFunc("/post", pc.GetPosts).Methods("GET")
			r.ServeHTTP(rr, req)
			assert.Equal(t, tc.expected.statusCode, rr.Code)
			assert.Equal(t, tc.expected.body, rr.Body.String())
		})
	}
}

func setParams(path string, qParams map[string]string) string {
	if len(qParams) == 0 {
		return ""
	}
	resp := ""
	for k, v := range qParams {
		if resp != "" {
			resp = fmt.Sprintf("%s&%s=%s", resp, k, v)
		} else {
			resp = fmt.Sprintf("?%s=%s", k, v)
		}
	}
	return path + resp
}
