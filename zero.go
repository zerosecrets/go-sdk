package zero

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const GRAPHQL_ENDPOINT_URL = "https://core.tryzero.com/v1/graphql"

type Fetch interface {
	Fetch() (map[string]map[string]string, error)
}

type ZeroApi struct {
	token string
	pick  []string
}

func (params ZeroApi) Fetch() (map[string]map[string]string, error) {
	query := `
    query Secrets($token: String!, $pick: [String!]) {
      secrets(zeroToken: $token, pick: $pick) {
        name

        fields {
          name
          value
        }
      }
    }
  `

	variables := map[string]string{
		"token": params.token,
		"pick":  strings.Join(params.pick, ","),
	}

	graphqlBody := GraphqlRequestBody{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(graphqlBody)

	if err != nil {
		return make(map[string]map[string]string), err
	}

	response, err := http.Post(GRAPHQL_ENDPOINT_URL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return make(map[string]map[string]string), err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return make(map[string]map[string]string), errors.New("zero returned non-200 status code")
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return make(map[string]map[string]string), err
	}

	var graphqlResponseBody GraphqlResponseBody
	err = json.Unmarshal(body, &graphqlResponseBody)

	if err != nil {
		return make(map[string]map[string]string), err
	}

	if graphqlResponseBody.Errors != nil {
		return nil, errors.New(graphqlResponseBody.Errors[0].Message)
	}

	secrets := make(map[string]map[string]string)

	for _, secret := range graphqlResponseBody.Data.Secrets {
		fields := make(map[string]string)

		for _, field := range secret.Fields {
			fields[field.Name] = field.Value
		}

		secrets[secret.Name] = fields
	}

	return secrets, err
}

type GraphqlRequestBody struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

type GraphqlResponseBody struct {
	Data struct {
		Secrets []struct {
			Name string `json:"name"`

			Fields []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"secrets"`
	} `json:"data"`

	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func Zero(token string, pick []string) (*ZeroApi, error) {
	if len(token) == 0 {
		return &ZeroApi{}, errors.New("Zero token should be non-empty string")
	}

	return &ZeroApi{
		token: token,
		pick:  pick,
	}, nil
}
