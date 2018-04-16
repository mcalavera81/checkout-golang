package test

import (
	"net/http/httptest"
	"net/http"
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"checkout-service/server/models"
)

type apiTestCase struct {
	tag      string
	method   string
	url      string
	body     string
	status   int
	response string
}

func newRouter() *mux.Router {
	router := new(mux.Router).StrictSlash(true)
	return router
}

func testAPI(router *mux.Router, method, URL, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, URL, bytes.NewBufferString(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, router *mux.Router, tests []apiTestCase) {
	for _, test := range tests {
		res := testAPI(router, test.method, test.url, test.body)
		assert.Equal(t, test.status, res.Code, test.tag)
		if test.response != "" {
			assert.JSONEq(t, test.response, res.Body.String(), test.tag)
		}
	}
}

func matchItemType(code models.Code) interface{}{
	return  mock.MatchedBy(func(item models.CatalogItem) bool { return item.Code == code})
}
