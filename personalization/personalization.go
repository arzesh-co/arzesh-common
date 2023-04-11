package personalization

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type personalization struct {
	Data struct {
		UserId     string         `json:"user_uuid" bson:"user_uuid"`
		RefId      string         `json:"ref_uuid" bson:"ref_uuid"`
		RefKey     string         `json:"ref_key" bson:"ref_key"`
		Type       string         `json:"type" bson:"type"`
		CustomData map[string]any `json:"custom_data" bson:"custom_data"`
	} `json:"data"`
}

func getPersonalizationInfoApi(userToken, clientToken, personalKey string) map[string]any {
	url := os.Getenv("PersonalizationApi") + "/api/personalize/personalization/key/" + personalKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	client := &http.Client{
		Transport: &http.Transport{},
	}
	req.Header.Add("id_token", userToken)
	req.Header.Add("client", clientToken)
	req.Close = true
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	Info := &personalization{}
	err = json.Unmarshal(body, Info)
	if err != nil {
		return nil
	}
	return Info.Data.CustomData
}

func GetUserPersonalizationInfoByKey(userToken, clientToken, key string) map[string]any {
	if userToken == "" || clientToken == "" || key == "" {
		return nil
	}
	return getPersonalizationInfoApi(userToken, clientToken, key)
}
