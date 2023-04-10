package tools

import (
	"context"
	"fmt"
	"github.com/arzesh-co/arzesh-common/date"
	"github.com/redis/go-redis/v9"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ConvertorToInt64(n any) int64 {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case string:
		num, err := strconv.ParseInt(n, 10, 64)
		if err == nil {
			return int64(0)
		}
		return num
	}
	return int64(0)
}
func GetValueFromShardCommonDb(key string) string {
	redisAddr := os.Getenv("redis")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: "common",
		Password: "", // no password set
		DB:       6,  // use default DB
	})
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
func ConvertMap(values map[string]any, params map[string]string) map[string]any {
	NewMap := make(map[string]any)
	for key, element := range params {
		switch ConvertorType(element) {
		case "func":
			NewMap[key] = FindFunc(element)
		case "string":
			NewMap[key] = values[element]
		case "map":
			NewMap[key] = FindValueOfMap(values, element)
		case "replace":
			NewMap[key] = Replacer(values, element)
		}
	}
	return NewMap
}
func Replacer(values map[string]any, param string) string {
	re := regexp.MustCompile("{{(.*?)}}")
	submatchall := re.FindAllString(param, -1)
	NewString := ""
	for _, element := range submatchall {
		repVal := element
		element = strings.Trim(element, "}")
		element = strings.Trim(element, "{")
		switch ConvertorType(element) {
		case "func":
			value := FindFunc(element)
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", value), -1)
		case "string":
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", values[element]), -1)
		case "map":
			value := FindValueOfMap(values, element)
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", value), -1)
		}
	}
	return NewString
}
func SetCommenFunc() map[string]func(days any) any {
	funcs := make(map[string]func(days any) any)
	funcs["today"] = date.StartToday
	funcs["nday_before"] = date.NDayBefore
	funcs["nday_after"] = date.NDayAfter
	funcs["begin_of_this_year"] = date.BeginOfThisYear
	funcs["begin_of_this_month"] = date.BeginOfThisMonth
	return funcs
}
func FindFunc(param string) any {
	funcs := SetCommenFunc()
	re := regexp.MustCompile(`\$(.*?)\$`)
	Findfunc := re.FindStringSubmatch(param)
	if len(Findfunc) > 0 {
		re = regexp.MustCompile(`(.*?)\(`)
		match := re.FindStringSubmatch(Findfunc[1])
		funcName := match[1]
		var params string
		re = regexp.MustCompile(`\((.*?)\)`)
		// Text between parentheses:
		submatchall := re.FindAllString(Findfunc[1], -1)
		for _, element := range submatchall {
			element = strings.Trim(element, "(")
			element = strings.Trim(element, ")")
			params = element
		}
		if params != "" {
			return funcs[funcName](params)
		} else {
			return funcs[funcName](0)
		}
	}
	return nil
}
func FindValueOfMap(value map[string]any, key string) any {
	partOfMap := strings.Split(key, ".")
	if len(partOfMap) > 0 {
		part := value[partOfMap[0]]
		if len(partOfMap) > 1 {
			for _, s := range partOfMap {
				part = part.(map[string]any)[s]
			}
		}
		return part
	}
	return nil
}
func ConvertorType(param string) string {
	re := regexp.MustCompile(`\$(.*?)\$`)
	isFunc := re.MatchString(param)
	if isFunc {
		return "func"
	}
	partOfMap := strings.Split(param, ".")
	if len(partOfMap) > 1 {
		return "map"
	}
	re = regexp.MustCompile("{{(.*?)}}")
	isReplacer := re.MatchString(param)
	if isReplacer {
		return "replace"
	}
	return "string"
}
