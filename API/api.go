package API

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/arzesh-co/arzesh-common/errors"
	"github.com/arzesh-co/arzesh-common/jwt"
	"github.com/arzesh-co/arzesh-common/tools"
	"github.com/arzesh-co/arzesh-common/tracing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"strings"
)

type InfoRequest struct {
	RequestFilter string
	PolicyFilter  string
	GroupFilter   string
	ServiceFilter []Filter
	Sort          string
	Skip          int64
	Limit         int64
	ClientToken   string
	UserToken     string
	Route         string
	Lang          string
	Tracer        trace.Tracer
	TraceShutdown func(context.Context) error
	Ctx           context.Context
}
type Filter struct {
	Condition any
	Label     string
	Operation string
}
type aggregation struct {
	GroupBy      string `json:"group_by"`
	GroupByTitle string `json:"group_by_title"`
	Aggregators  []struct {
		Aggregate string `json:"aggregate"`
		Operation string `json:"operation"`
	} `json:"aggregators"`
}
type sort struct {
	DbName string `json:"db_name"`
	Type   string `json:"type"`
}
type MongoPipeLine struct {
	Filter      map[string]any
	Sort        map[string]any
	Aggregation map[string]any
	Fields      map[string]int8
}

func createSkip(strSkip string) int64 {
	skip, err := strconv.ParseInt(strSkip, 10, 64)
	if err != nil {
		skip = 0
	}
	return skip
}
func createLimit(strLimit string) int64 {
	limit, err := strconv.ParseInt(strLimit, 10, 64)
	if err != nil {
		limit = 10
	}
	return limit
}

func New(request *http.Request, service, serviceVersion string) *InfoRequest {
	ctx := context.Background()
	req := &InfoRequest{}
	req.RequestFilter = request.URL.Query().Get("filter")
	// TODO : how can fill this part
	req.PolicyFilter = ""
	req.Lang = request.Header.Get("lang")
	req.Skip = createSkip(request.URL.Query().Get("skip"))
	req.Limit = createLimit(request.URL.Query().Get("limit"))
	req.Route = request.URL.Path
	req.ClientToken = request.Header.Get("client")
	req.GroupFilter = request.URL.Query().Get("aggregation")
	req.Sort = request.URL.Query().Get("sort")
	req.UserToken = request.Header.Get("id_token")
	shutdown, err := tracing.InitProvider(service, serviceVersion)
	if err != nil {
		return req
	}
	req.TraceShutdown = shutdown
	tracer := otel.Tracer(req.Route)
	// work begins
	req.Tracer = tracer
	ctx, span := tracer.Start(
		ctx,
		"start request presses")
	defer span.End()
	req.Ctx = ctx
	return req
}

func (r InfoRequest) ClientValidationRequest() bool {
	return jwt.Validate(r.ClientToken)
}
func (r InfoRequest) UserValidationRequest() (bool, error) {
	IsClientValid := jwt.Validate(r.ClientToken)
	if !IsClientValid {
		return false, errors2.New("client token is not valid")
	}
	IsUserTokenValid := jwt.Validate(r.UserToken)
	if !IsUserTokenValid {
		return false, errors2.New("user token is not valid")
	}
	return true, nil
}

func createFilter(cond Filter) interface{} {
	switch cond.Operation {
	case "text":
		return bson.M{"$text": bson.M{"$search": cond.Condition.(string)}}
	case "Start With":
		return primitive.Regex{Pattern: "^" + cond.Condition.(string) + ".", Options: "i"}
	case "End With":
		return primitive.Regex{Pattern: ".*" + cond.Condition.(string) + "$", Options: "i"}
	case "Equal":
		return bson.M{"$eq": cond.Condition}
	case "Include":
		return primitive.Regex{Pattern: ".*" + cond.Condition.(string) + ".*", Options: "i"}
	case "Empty":
		return bson.M{"$exists": false}
	case "not Empty":
		return bson.M{"$exists": true}
	case "=":
		if fmt.Sprintf("%T", cond.Condition) == "[]interface {}" {
			return bson.M{"$in": cond.Condition}
		}
		return bson.M{"$eq": ConvertFilterCondition(cond.Condition)}
	case ">=":
		return bson.M{"$gte": ConvertFilterCondition(cond.Condition)}
	case "<=":
		return bson.M{"$lte": ConvertFilterCondition(cond.Condition)}
	case ">":
		return bson.M{"$gt": ConvertFilterCondition(cond.Condition)}
	case "<":
		return bson.M{"$lt": ConvertFilterCondition(cond.Condition)}
	case "!=":
		if fmt.Sprintf("%T", cond.Condition) == "[]interface {}" {
			return bson.M{"$nin": cond.Condition}
		}
		return bson.M{"$ne": ConvertFilterCondition(cond.Condition)}
	}
	return bson.M{}
}
func ConvertFilterCondition(condition any) any {
	switch condition.(type) {
	case string:
		switch tools.ConvertorType(condition.(string)) {
		case "func":
			return tools.FindFunc(condition.(string))
		case "string":
			return condition
		default:
			return condition
		}
	default:
		return condition
	}
}

func convertStringFilter(strFilter string) ([]Filter, error) {
	var filters []Filter
	err := json.Unmarshal([]byte(strFilter), &filters)
	if err != nil {
		return nil, err
	}
	return filters, nil
}

func (r InfoRequest) MongoDbFilter() (map[string]any, error) {
	var filters []Filter
	if len(r.ServiceFilter) > 0 {
		filters = append(filters, r.ServiceFilter...)
	}
	if r.RequestFilter != "" {
		r.RequestFilter = strings.Replace(r.RequestFilter, "'", "\"", -1)
		reqFilter, err := convertStringFilter(r.RequestFilter)
		if err != nil {
			return nil, err
		}
		filters = append(filters, reqFilter...)
	}
	if r.PolicyFilter != "" {
		r.RequestFilter = strings.Replace(r.PolicyFilter, "'", "\"", -1)
		polFilter, err := convertStringFilter(r.RequestFilter)
		if err != nil {
			return nil, err
		}
		filters = append(filters, polFilter...)
	}
	clintFilterMap := make(map[string]any)
	if len(filters) > 1 {
		listCondition := make([]map[string]any, len(filters))
		for i, f := range filters {
			condition := make(map[string]any)
			condition[f.Label] = createFilter(f)
			listCondition[i] = condition
		}
		clintFilterMap["$and"] = listCondition
		return clintFilterMap, nil
	}
	clintFilterMap[filters[0].Label] = createFilter(filters[0])
	return clintFilterMap, nil
}
func (r InfoRequest) MongoDbSorting() (map[string]any, error) {
	var sorts []sort
	err := json.Unmarshal([]byte(r.Sort), &sorts)
	if err != nil {
		return nil, err
	}
	sortFilter := make(map[string]any)
	for _, s := range sorts {
		switch s.Type {
		case "asc":
			sortFilter[s.DbName] = 1
		case "des":
			sortFilter[s.DbName] = -1
		}
	}
	return sortFilter, nil
}
func (r InfoRequest) MongoDbAggregation() map[string]interface{} {
	if r.GroupFilter == "" {
		return nil
	}
	agg := &aggregation{}
	err := json.Unmarshal([]byte(r.GroupFilter), agg)
	if err != nil {
		return nil
	}
	group := make(map[string]interface{})
	group["_id"] = "$" + agg.GroupBy
	group["group_by_title"] = bson.M{"$first": "$" + agg.GroupByTitle}
	for _, aggregator := range agg.Aggregators {
		switch aggregator.Operation {
		case "avg":
			group[aggregator.Aggregate] = bson.M{"$avg": "$" + aggregator.Aggregate}
		case "sum":
			group[aggregator.Aggregate] = bson.M{"$sum": "$" + aggregator.Aggregate}
		case "count":
			group[aggregator.Aggregate] = bson.M{"$sum": 1}
		case "min":
			group[aggregator.Aggregate] = bson.M{"$min": "$" + aggregator.Aggregate}
		case "max":
			group[aggregator.Aggregate] = bson.M{"$max": "$" + aggregator.Aggregate}
		case "first":
			group[aggregator.Aggregate] = bson.M{"$first": "$" + aggregator.Aggregate}
		}
	}
	return group
}

func (r InfoRequest) PipeLineMongoDbAggregate() ([]bson.M, any) {
	line := &MongoPipeLine{}
	var err error
	line.Sort, err = r.MongoDbSorting()
	if err != nil {
		line.Sort = nil
	}
	line.Filter, err = r.MongoDbFilter()
	if err != nil {
		line.Filter = nil
	}
	line.Aggregation = r.MongoDbAggregation()
	// TODO: how can fill this part
	line.Fields = nil
	var filterCount interface{}
	if line.Filter != nil {
		filterCount = line.Filter
	}

	var skipPage int64
	if r.Skip != 0 {
		skipPage = (r.Skip - 1) * r.Limit
	} else {
		skipPage = 0
	}
	var pipe []bson.M
	if len(line.Filter) != 0 {
		pipe = append(pipe, bson.M{"$match": line.Filter})
	}
	if len(line.Fields) != 0 {
		pipe = append(pipe, bson.M{"$project": line.Fields})
	}
	if len(line.Aggregation) != 0 {
		pipe = append(pipe, bson.M{"$group": line.Aggregation})
	}
	if len(line.Sort) != 0 {
		pipe = append(pipe, bson.M{"$sort": line.Sort})
	}
	pipe = append(pipe, bson.M{"$skip": skipPage})
	pipe = append(pipe, bson.M{"$limit": r.Limit})
	return pipe, filterCount
}

type ResponseOtp struct {
	Meta map[string]any
}

func (opt ResponseOtp) SetPagination(totalCount int64, skip, limit string) {
	if opt.Meta == nil {
		opt.Meta = make(map[string]any)
	}
	current, _ := strconv.ParseInt(skip, 10, 64)
	MLimit, _ := strconv.ParseInt(limit, 10, 64)
	opt.Meta["current"] = current
	if skip == "0" {
		opt.Meta["current"] = 1
		current = 1
	}
	opt.Meta["total"] = totalCount
	if (current * MLimit) < totalCount {
		opt.Meta["next"] = current + 1
	}
	if current > 1 {
		opt.Meta["prev"] = current - 1
	}
}

type apiResponse struct {
	Error   *errors.ResponseErrors `json:"errors,omitempty"`
	Data    any                    `json:"data,omitempty"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]any         `json:"meta,omitempty"`
}

func (r InfoRequest) NewResponse(data any, messageKey string, Error *errors.ResponseErrors, opt *ResponseOtp) *apiResponse {
	res := &apiResponse{}
	//TODO : create message of template
	if messageKey != "" {
		res.Message = messageKey
	}
	if data != nil {
		res.Data = data
	}
	if Error != nil {
		res.Error = Error
	}
	if opt != nil {
		res.Meta = opt.Meta
	}
	commonAttrs := []attribute.KeyValue{
		attribute.String("messageKey", messageKey),
		attribute.String("ErrorKey", Error.StatusKey),
	}

	// work begins
	_, span := r.Tracer.Start(
		r.Ctx,
		"End request processes",
		trace.WithAttributes(commonAttrs...))
	defer span.End()
	return res
}
