package backends

import (
	"errors"
	"net/http"

	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

type SitefinityBackendProvider struct {
	Host       string
	Token      string
	HttpClient HttpClient
}

func (p *SitefinityBackendProvider) Request(endpoint string, method httpmethod.HttpMethod, responseType responsetype.ResponseType, data map[string]interface{}, filter *string, sort *string, expand *string) (interface{}, error) {
	url := p.Host + endpoint
	params := map[string]string{}

	if filter != nil {
		params["$filter"] = *filter
	}
	if sort != nil {
		params["$orderby"] = *sort
	}
	if expand != nil {
		params["$expand"] = *expand
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

	if responseType == responsetype.JSONArrayOfObjects {
		odataResp, err := helpers.JSONDecodeODataArrayOfObjects(resp.Body)
		if err != nil {
			return nil, err
		}

		return odataResp, nil
	}

	odataResp, err := helpers.JSONDecodeObject(resp.Body)
	if err != nil {
		return nil, err
	}

	if responseType == responsetype.JSONObject {
		return odataResp, nil
	}

	odataValue, found := odataResp["value"]
	if !found {
		bytes, err := helpers.JSONEncodeObject(odataResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("failed to parse backend response:\n" + string(bytes))
	}

	return odataValue, nil
}

func (p *SitefinityBackendProvider) getHeader() http.Header {
	return http.Header{
		"Authorization": []string{"Bearer " + p.Token},
	}
}
