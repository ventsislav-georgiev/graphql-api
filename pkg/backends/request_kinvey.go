package backends

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

type KinveyBackendProvider struct {
	Host         string
	KinveyID     string
	MasterSecret string
	HttpClient   HttpClient
}

func (p *KinveyBackendProvider) Request(endpoint string, method httpmethod.HttpMethod, responseType responsetype.ResponseType, data map[string]interface{}, filter *string, sort *string) (interface{}, error) {
	url := p.Host + endpoint
	params := map[string]string{}

	if filter != nil {
		params["query"] = *filter
	}

	if sort != nil {
		params["sort"] = *sort
	}

	resp, err := p.HttpClient.Request(RequestOptions{
		Method: method,
		Url:    url,
		Header: p.getHeader(),
		Params: params,
		Data:   data,
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return true, nil
	}

	switch responseType {
	case responsetype.JSONObject:
		return helpers.JSONDecodeObject(resp.Body)
	case responsetype.JSONArrayOfObjects:
		return helpers.JSONDecodeArrayOfObjects(resp.Body)
	case responsetype.ScalarValue:
		return io.ReadAll(resp.Body)
	}

	return nil, errors.New("unknown server response type")
}

func (p *KinveyBackendProvider) getHeader() http.Header {
	return http.Header{
		"Authorization":        []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(p.KinveyID+":"+p.MasterSecret))},
		"X-Kinvey-API-Version": []string{"5"},
	}
}
