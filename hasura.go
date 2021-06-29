package mon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
)

type IHasura interface {
	// ReadGqlFile sets grapqhl query by reading .gql or .graphql file
	ReadGqlFile(filepath string) error

	// Exec returns nil if calling graphql API returns success
	// and it will set the response in the response object
	Exec(variables interface{}, response interface{}, headers map[string]string) error
}

type Hasura struct {
	GraphqlURL string
	Secret     string
	Query      string
}

func (h *Hasura) ReadGqlFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	h.Query = string(file)
	return nil
}

func (h *Hasura) Exec(variables interface{}, response interface{}, headers map[string]string) error {
	body := make(map[string]interface{})
	body["query"] = h.Query
	body["variables"] = variables

	jsonbody, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	client := resty.New()
	req := client.R().SetHeader("Content-Type", "application/json")
	for key, value := range headers {
		req = req.SetHeader(key, value)
	}

	resp, err := req.SetBody(jsonbody).Post(h.GraphqlURL)
	if err != nil {
		return err
	}

	if strings.Contains(string(resp.Body()), "errors") {
		return errors.New("something went wrong when calling hasura")
	}

	if response == nil {
		return nil
	}

	return json.Unmarshal(resp.Body(), response)
}

func NewHasura(graphqlURL string, secret string) IHasura {
	return &Hasura{
		GraphqlURL: graphqlURL,
		Secret:     secret,
	}
}
