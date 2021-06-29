package mon

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
)

type IHasura interface {
	// SetResponseKey sets response key
	// e.g. in hasura return data.user then response key is "user"
	SetResponseKey(key string) IHasura

	// SetResponseModel sets response model/schema
	// Struct with json tag
	SetResponseModel(model interface{}) IHasura

	// ReadGqlFile sets grapqhl query by reading .gql or .graphql file
	ReadGqlFile(filepath string) error

	// Exec returns nil if calling graphql API returns success
	Exec(queryType string, variables interface{}, headers map[string]string) error
}

var (
	HasuraTypeQuery    = "QUERY"
	HasuraTypeMutation = "MUTATION"
)

type Hasura struct {
	GraphqlURL string
	Secret     string
	Query      string

	ResponseKey   string
	ResponseModel interface{}
}

func (h *Hasura) SetResponseKey(key string) IHasura {
	h.ResponseKey = key
	return h
}

func (h *Hasura) SetResponseModel(model interface{}) IHasura {
	h.ResponseModel = model
	return h
}

func (h *Hasura) ReadGqlFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	h.Query = string(file)
	return nil
}

func (h *Hasura) Exec(queryType string, variables interface{}, headers map[string]string) error {
	if queryType != HasuraTypeQuery || queryType != HasuraTypeMutation {
		return errors.New("invalid hasura query type")
	}

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

	response := new(HasuraResponseSchema)
	err = json.Unmarshal(resp.Body(), response)
	if err != nil {
		return err
	}

	if response.Errors != nil {
		return errors.New("something went wrong when calling hasura")
	}

	var values interface{}
	var ok bool

	if queryType == HasuraTypeQuery {
		values, ok = response.Data.(map[string]interface{})[h.ResponseKey]
	} else if queryType == HasuraTypeMutation {
		values, ok = response.Data.(map[string]interface{})[h.ResponseKey].(map[string]interface{})["returning"]
	}

	if !ok {
		return errors.New("failed to parse response from hasura")
	}

	b, err := json.Marshal(values)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, h.ResponseModel)
}

type HasuraResponseSchema struct {
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

func NewHasura(graphqlURL string, secret string) IHasura {
	return &Hasura{
		GraphqlURL: graphqlURL,
		Secret:     secret,
	}
}
