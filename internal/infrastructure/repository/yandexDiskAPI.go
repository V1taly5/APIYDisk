package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type YandexDiskAPI struct {
	Token string
	Debug bool

	Client          HTTPClient
	shutdownChannel chan interface{} //onFuture

	CurrentPath string

	apiEndpoint string
}

func NewYandexDiskAPI(token string) (*YandexDiskAPI, error) {
	return NewYandexDiskAPIWithClient(token, APIEndpoint, &http.Client{})
}

func NewYandexDiskAPIWithClient(token, apiEndpoint string, client HTTPClient) (*YandexDiskAPI, error) {
	disk := &YandexDiskAPI{
		Token:           token,
		Client:          client,
		shutdownChannel: make(chan interface{}),
		apiEndpoint:     apiEndpoint,
	}
	return disk, nil
}

func (disk *YandexDiskAPI) SetAPIEndpoint(apiEndpoint string) {
	disk.apiEndpoint = apiEndpoint
}

// проверить существование этого пути
func (disk *YandexDiskAPI) SetCurrentPath(path string) {
	if disk.CurrentPath != path || path != "" {
		disk.CurrentPath = path
	}
}

type Params map[string]string

func buildParams(in Params) url.Values {
	if in == nil {
		return url.Values{}
	}
	out := url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}
	return out
}

type requestData struct {
	endPoint string
	method   string
	headers  map[string]string
	params   Params
}

func NewRequestParams(endPoint string, method string, headers map[string]string, params Params) requestData {
	return requestData{
		endPoint: endPoint,
		method:   method,
		headers:  headers,
		params:   params,
	}

}

func (disk *YandexDiskAPI) UploadFileLink(imgUrl, uploadPath string) (map[string]interface{}, error) {

	apiEndpoint := "resources/upload"
	method := "POST"
	headers := map[string]string{
		"Authorization": fmt.Sprintf("OAuth %s", disk.Token),
	}
	var params Params
	params = make(Params)
	params["path"] = uploadPath
	params["url"] = imgUrl

	par := NewRequestParams(apiEndpoint, method, headers, params)

	request, err := disk.MakeRequest(par)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (disk *YandexDiskAPI) MakeRequest(params requestData) (map[string]interface{}, error) {
	values := buildParams(params.params)
	url := fmt.Sprintf("%s/%v?%s", disk.apiEndpoint, params.endPoint, values.Encode())
	if disk.Debug {
		fmt.Printf("Endpoint: %s, params: %v/n", params.endPoint, values)
	}

	req, err := http.NewRequest(params.method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if params.headers != nil {
		for key, value := range params.headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := disk.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}
	responseData := make(map[string]interface{})
	json.Unmarshal(body, &responseData)

	return responseData, nil
}
