package zero

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestZero(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	rawMock, mock := GraphqlApiResponseMock()

	httpmock.RegisterResponder("POST", GRAPHQL_ENDPOINT_URL,
		func(request *http.Request) (*http.Response, error) {
			body := make(map[string]interface{})

			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			if body["variables"].(map[string]interface{})["pick"] == "" {
				return httpmock.NewStringResponse(200, "Failed to fetch"), nil
			}

			return httpmock.NewStringResponse(200, rawMock), nil
		},
	)

	t.Run("requires token to be non-empty string", func(t *testing.T) {
		_, err := Zero("", []string{"aws", "azure"}, nil)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("does a GraphQL request which queries the requested APIs", func(t *testing.T) {
		zeroApi, err := Zero("token", []string{"aws"}, nil)

		if err != nil {
			t.Error(err)
		}

		got, err := zeroApi.Fetch()

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, mock) {
			t.Errorf("Expected %v, got %v", mock, got)
		}
	})

	t.Run("returns an error if GraphQL API responds with error", func(t *testing.T) {
		api, err := Zero("token", []string{}, nil)

		if err != nil {
			t.Error(err)
		}

		_, err = api.Fetch()

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
