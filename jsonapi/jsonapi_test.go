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

	t.Run("Error with source pointer", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR422", http.StatusUnprocessableEntity, exception.Throw(
				errors.New("validation error"),
				exception.WithDetail("The name field is required"),
				exception.WithTitle("Validation Error"),
			), jsonapi.WithSourcePointer("/data/attributes/name")),
		)

		b, _ := json.Marshal(response)
		assert.Contains(t, string(b), "\"source\":{\"pointer\":\"/data/attributes/name\"}")
		assert.Equal(t, http.StatusUnprocessableEntity, response.HTTPStatus())
	})

	t.Run("Error with source parameter", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR400", http.StatusBadRequest, exception.Throw(
				errors.New("invalid parameter"),
				exception.WithDetail("The limit parameter must be a positive integer"),
				exception.WithTitle("Invalid Parameter"),
			), jsonapi.WithSourceParameter("limit")),
		)

		b, _ := json.Marshal(response)
		assert.Contains(t, string(b), "\"source\":{\"parameter\":\"limit\"}")
		assert.Equal(t, http.StatusBadRequest, response.HTTPStatus())
	})

	t.Run("Error with source header", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR401", http.StatusUnauthorized, exception.Throw(
				errors.New("missing authorization"),
				exception.WithDetail("The Authorization header is missing"),
				exception.WithTitle("Unauthorized"),
			), jsonapi.WithSourceHeader("Authorization")),
		)

		b, _ := json.Marshal(response)
		assert.Contains(t, string(b), "\"source\":{\"header\":\"Authorization\"}")
		assert.Equal(t, http.StatusUnauthorized, response.HTTPStatus())
	})

	t.Run("Error with source and meta", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithExceptionMeta("ERR422", http.StatusUnprocessableEntity, exception.Throw(
				errors.New("validation error"),
				exception.WithDetail("The email format is invalid"),
				exception.WithTitle("Validation Error"),
			), jsonapi.Meta{
				"field":  "email",
				"format": "email",
			}, jsonapi.WithSourcePointer("/data/attributes/email")),
		)

		b, _ := json.Marshal(response)
		assert.Contains(t, string(b), "\"source\":{\"pointer\":\"/data/attributes/email\"}")
		assert.Contains(t, string(b), "\"meta\":{\"field\":\"email\",\"format\":\"email\"}")
		assert.Equal(t, http.StatusUnprocessableEntity, response.HTTPStatus())
	})

	t.Run("Multiple source options should be handled correctly", func(t *testing.T) {
		// Only the last source option should be used when multiple are provided
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR400", http.StatusBadRequest, exception.Throw(
				errors.New("invalid parameter"),
				exception.WithDetail("Invalid request"),
				exception.WithTitle("Bad Request"),
			),
				jsonapi.WithSourcePointer("/data/attributes/name"),
				jsonapi.WithSourceParameter("id"), // This should override the pointer
			),
		)

		b, _ := json.Marshal(response)
		assert.Contains(t, string(b), "\"source\":{\"pointer\":\"/data/attributes/name\",\"parameter\":\"id\"}")
	})

	t.Run("Using all source options together", func(t *testing.T) {
		response := jsonapi.BuildResponse(
			jsonapi.WithException("ERR400", http.StatusBadRequest, exception.Throw(
				errors.New("complex error"),
				exception.WithDetail("Complex error with multiple sources"),
				exception.WithTitle("Bad Request"),
			),
				jsonapi.WithSourcePointer("/data"),
				jsonapi.WithSourceParameter("sort"),
				jsonapi.WithSourceHeader("X-Custom"), // This should be the only one that remains
			),
		)

		b, _ := json.Marshal(response)
		// The source should contain only the header field since it was the last one set
		assert.Contains(t, string(b), "\"source\":{\"pointer\":\"/data\",\"parameter\":\"sort\",\"header\":\"X-Custom\"}")
	})
}
