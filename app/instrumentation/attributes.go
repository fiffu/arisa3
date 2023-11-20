package instrumentation

import (
	"fmt"

	"github.com/fiffu/arisa3/app/log"
	"go.opentelemetry.io/otel/attribute"
)

const (
	attrTraceID            = string(log.TraceID)
	attrTraceSubID         = string(log.TraceSubID)
	attrCogName            = string(log.CogName)
	attrUser               = string(log.User)
	attrCommandName        = "command_name"
	attrEventName          = "event_name"
	attrParams             = "params"
	attrHTTPHost           = "http_host"
	attrHTTPMethod         = "http_method"
	attrHTTPPath           = "http_path"
	attrHTTPRespStatusCode = "http_resp_status"
	attrHTTPContentLength  = "http_total_content_length"
	attrDBQuery            = "db_query"
	attrDBOperation        = "db_operation"
)

type attrs struct{}

var KV = attrs{}

func (attrs) Cog(value string) attribute.KeyValue {
	return attribute.String(attrCogName, value)
}

func (attrs) CommandName(value string) attribute.KeyValue {
	return attribute.String(attrCommandName, value)
}

func (attrs) EventName(value string) attribute.KeyValue {
	return attribute.String(attrEventName, value)
}

func (attrs) User(value string) attribute.KeyValue {
	return attribute.String(attrUser, value)
}

func (attrs) Params(value map[string]any) attribute.KeyValue {
	return attribute.String(attrParams, fmt.Sprint(value))
}

func (attrs) HTTPHost(value string) attribute.KeyValue {
	return attribute.String(attrHTTPHost, value)
}
func (attrs) HTTPMethod(value string) attribute.KeyValue {
	return attribute.String(attrHTTPMethod, value)
}

func (attrs) HTTPPath(value string) attribute.KeyValue {
	return attribute.String(attrHTTPPath, value)
}

func (attrs) HTTPRespStatusCode(value int) attribute.KeyValue {
	return attribute.Int(attrHTTPRespStatusCode, value)
}

func (attrs) HTTPTotalContentLength(value int64) attribute.KeyValue {
	return attribute.Int64(attrHTTPContentLength, value)
}

func (attrs) TraceID(value string) attribute.KeyValue {
	return attribute.String(attrTraceID, value)
}

func (attrs) TraceSubID(value string) attribute.KeyValue {
	return attribute.String(attrTraceSubID, value)
}

func (attrs) DBQuery(sql string) attribute.KeyValue {
	return attribute.String(attrDBQuery, sql)
}

func (attrs) DBOperation(op string) attribute.KeyValue {
	return attribute.String(attrDBOperation, op)
}
