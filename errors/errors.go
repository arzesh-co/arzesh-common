package errors

import (
	"encoding/json"
	"github.com/arzesh-co/arzesh-common/tools"
	"strings"
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
type ErrorType struct {
	Id         string            `json:"_id" bson:"_id"`
	Key        string            `json:"key" bson:"key"`
	Title      map[string]string `json:"title" bson:"title"`
	Desc       map[string]string `json:"desc" bson:"desc"`
	ServiceKey string            `json:"service_key" bson:"service_key"`
	DevTeamId  string            `json:"dev_team_uuid" bson:"devTeamId"`
	Url        string            `json:"url" bson:"url"`
	Status     int8              `json:"status" bson:"status"`
}
type ErrorParam struct {
	Key          string            `json:"key" bson:"key"`
	DefaultValue map[string]string `json:"default" bson:"default"`
}
type Errors struct {
	Id         string            `json:"_id" bson:"_id"`
	ErrorType  ErrorType         `json:"error_type" bson:"error_type"`
	StatusKey  string            `json:"status_key" bson:"status_key"`
	Detail     map[string]string `json:"detail" bson:"detail"`
	ServiceKey string            `json:"service_key" bson:"service_key"`
	Title      map[string]string `json:"title" bson:"title"`
	Params     []ErrorParam      `json:"params" bson:"params"`
	HelpUrl    string            `json:"help_url" bson:"help_url"`
	MetaData   map[string]any    `json:"meta_data" bson:"meta_data"`
	Status     int8              `json:"status" bson:"status"`
}

type Entities struct {
	Id         string            `json:"_id" bson:"_id"`
	EntityName string            `json:"entity_name" bson:"entity_name"`
	Title      map[string]string `json:"title" bson:"title"`
}

func FindError(key string) *Errors {
	strErr := tools.GetValueFromShardCommonDb("Error:" + key)
	if strErr == "" {
		return nil
	}
	Err := &Errors{}
	err := json.Unmarshal([]byte(strErr), Err)
	if err != nil {
		return nil
	}
	return Err
}
func convertErrorToResponseErr(Err *Errors, lang string) *ResponseErrors {
	res := &ResponseErrors{}
	res.ErrorType.Url = Err.ErrorType.Url
	res.ErrorType.Title = Err.ErrorType.Title[lang]
	res.ErrorType.Desc = Err.ErrorType.Desc[lang]
	res.Detail = Err.Detail[lang]
	res.HelpUrl = Err.HelpUrl
	res.Title = Err.Title[lang]
	res.StatusKey = Err.StatusKey
	return res
}

func FindEntityName(entity string, lang string) string {
	strEntity := tools.GetValueFromShardCommonDb("Entity:" + entity)
	if strEntity == "" {
		return ""
	}
	en := &Entities{}
	err := json.Unmarshal([]byte(strEntity), en)
	if err != nil {
		return ""
	}
	entityTitle, ok := en.Title[lang]
	if !ok {
		return ""
	}
	return entityTitle
}

func setParamsToResponseErr(res *ResponseErrors, DefaultParams []ErrorParam, params map[string]string, entity string, lang string) *ResponseErrors {
	for _, p := range DefaultParams {
		paramVal, ok := params[p.Key]
		if !ok {
			res.Detail = strings.Replace(res.Detail, p.Key, p.DefaultValue[p.Key], -1)
			continue
		}
		if paramVal == "%entity" {
			entityTitle := FindEntityName(entity, lang)
			if entityTitle == "" {
				res.Detail = strings.Replace(res.Detail, p.Key, p.DefaultValue[p.Key], -1)
				continue
			}
			res.Detail = strings.Replace(res.Detail, p.Key, entityTitle, -1)
			continue
		}
		res.Detail = strings.Replace(res.Detail, p.Key, paramVal, -1)
	}
	return res
}

func New(key, lang, entity string, developerInfo string, params map[string]string) *ResponseErrors {
	Err := FindError(key)
	if Err == nil {
		return nil
	}
	res := convertErrorToResponseErr(Err, lang)
	res = setParamsToResponseErr(res, Err.Params, params, entity, lang)

	//TODO: this part must be remove ...
	if developerInfo != "" {
		meta := make(map[string]string)
		meta["dev_info"] = developerInfo
		res.Meta = meta
	}
	return res
}
