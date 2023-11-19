package instrumentation

import (
	"fmt"

	"github.com/fiffu/arisa3/app/log"
	"go.opentelemetry.io/otel/attribute"
)

const (
	attrTraceID     = string(log.TraceID)
	attrTraceSubID  = string(log.TraceSubID)
	attrCogName     = string(log.CogName)
	attrUser        = string(log.User)
	attrParams      = "params"
	attrHTTPPath    = "http_path"
	attrDBQuery     = "db_query"
	attrDBOperation = "db_operation"
)

type attrs struct{}

var KV = attrs{}

func (attrs) Cog(value string) attribute.KeyValue {
	return attribute.String(attrCogName, value)
}

func (attrs) User(value string) attribute.KeyValue {
	return attribute.String(attrUser, value)
}

func (attrs) Params(value map[string]any) attribute.KeyValue {
	return attribute.String(attrParams, fmt.Sprint(value))
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
