package connections

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type TransformationResponse struct {
	Errors map[string]any `json:"errors"`
	Data   map[string]any `json:"data"`
}

func getTransformedBodyByKeyApi(account string, transformationKey string) (map[string]any, string) {
	url := os.Getenv("transformationApi") + "/api/transformation/transformations/key/" + transformationKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "REF.CANNOT_CONNECT"
	}
	req.Header.Set("account_uuid", account)
	client := &http.Client{Transport: &http.Transport{}}
	res, err := client.Do(req)
	if err != nil {
		return nil, "REF.CANNOT_CONNECT"
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "REF.INVALIDATION_ERROR"
	}
	transformationRes := &TransformationResponse{}
	err = json.Unmarshal(body, transformationRes)
	if err != nil {
		return nil, "REF.INVALIDATION_ERROR"
	}
	errorKey, ok := transformationRes.Errors["status_key"]
	if ok {
		return nil, errorKey.(string)
	}
	return transformationRes.Data, ""
}

func TransformationInfos(account string, key string) (map[string]any, string) {
	if account == "" {
		return nil, "REF.INVALIDATION_ERROR"
	}
	if key == "" {
		return nil, "REF.INVALIDATION_ERROR"
	}
	return getTransformedBodyByKeyApi(account, key)
}
