package jsonapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestJSONAPI(t *testing.T) {

	t.Run("Nil errors and meta won't print in json format", func(t *testing.T) {
		response := jsonapi.BuildResponse(jsonapi.WithData(map[string]interface{}{
			"name": "testing",
		}))

		b, _ := json.Marshal(response)
		assert.Equal(t, `{"data":{"name":"testing"}}`, string(b))
	})

	t.Run("Nil data won't print in json format", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR401", http.StatusNotFound, exception.Throw(
				errors.New("unexpected error"),
				exception.WithDetail("detail"),
				exception.WithTitle("title"),
			)),
			jsonapi.WithExceptionMeta("ERR401", http.StatusNotFound, exception.Throw(
				errors.New("unexpected error"),
				exception.WithDetail("detail_2"),
				exception.WithTitle("title_2"),
			), jsonapi.Meta{
				"sample": "sample",
			}),
			jsonapi.WithMeta("nice", true),
		)

		b, _ := json.Marshal(response)
		assert.Equal(t, "{\"errors\":[{\"title\":\"title\",\"detail\":\"detail\",\"code\":\"ERR401\",\"status\":404},{\"title\":\"title_2\",\"detail\":\"detail_2\",\"code\":\"ERR401\",\"status\":404,\"meta\":{\"sample\":\"sample\"}}],\"meta\":{\"nice\":true}}", string(b))
		assert.Equal(t, http.StatusNotFound, response.HTTPStatus())
		assert.Equal(t, "JSONAPI Error:\n[title] Detail: detail, Code: ERR401\n[title_2] Detail: detail_2, Code: ERR401\n", response.Errors.Error())
	})

	t.Run("Nil data won't print in json format using WithErrors option", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR401", http.StatusNotFound, exception.Throw(
				errors.New("unexpected error"),
				exception.WithDetail("detail"),
				exception.WithTitle("title"),
			)),
			jsonapi.WithExceptionMeta("ERR401", http.StatusNotFound, exception.Throw(
				errors.New("unexpected error"),
				exception.WithDetail("detail_2"),
				exception.WithTitle("title_2"),
			), jsonapi.Meta{
				"sample": "sample",
			}),
		)

		response2 := jsonapi.BuildResponse(jsonapi.WithErrors(response.Errors))

		b, _ := json.Marshal(response2)
		assert.Equal(t, "{\"errors\":[{\"title\":\"title\",\"detail\":\"detail\",\"code\":\"ERR401\",\"status\":404},{\"title\":\"title_2\",\"detail\":\"detail_2\",\"code\":\"ERR401\",\"status\":404,\"meta\":{\"sample\":\"sample\"}}]}", string(b))
		assert.Equal(t, http.StatusNotFound, response2.HTTPStatus())
		assert.Equal(t, "JSONAPI Error:\n[title] Detail: detail, Code: ERR401\n[title_2] Detail: detail_2, Code: ERR401\n", response2.Errors.Error())
	})

	t.Run("500 http status when error status not set", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR401", 0, exception.Throw(
				errors.New("unexpected error"),
				exception.WithDetail("detail"),
				exception.WithTitle("title"),
			)),
		)

		assert.Equal(t, 500, response.HTTPStatus())
	})
}
