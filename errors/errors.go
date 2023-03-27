package errors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type ResponseErrors struct {
	ErrorType struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Url   string `json:"url"`
	} `json:"error_type"`
	StatusKey string            `json:"status_key"`
	Detail    string            `json:"detail"`
	Title     string            `json:"title"`
	Meta      map[string]string `json:"meta"`
	HelpUrl   string            `json:"help_url"`
}

func getError(key string, account string, lang string, params map[string]string) *ResponseErrors {
	req, err := http.NewRequest("GET", os.Getenv("coreApi")+"/api/core/errors/key/"+key, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("account_uuid", account)
	q := req.URL.Query()
	q.Add("lang", lang)
	paramsS, _ := json.Marshal(params)
	q.Add("params", string(paramsS))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{
		Transport: &http.Transport{},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	Info := &ResponseErrors{}
	err = json.Unmarshal(body, Info)
	if err != nil {
		return nil
	}
	return Info
}

//TODO: find error by key, replace params with key params check if entities full .

func New(key, entity string, developerInfo string, params map[string]string) *ResponseErrors {
	return nil
}
