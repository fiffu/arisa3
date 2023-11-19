package instrumentation

import (
	"fmt"
	"net/url"

	"github.com/fiffu/arisa3/app/log"
	"go.opentelemetry.io/otel/attribute"
)

const (
	attrTraceID     = string(log.TraceID)
	attrTraceSubID  = string(log.TraceSubID)
	attrHTTPPath    = "http_path"
	attrDBQuery     = "db_query"
	attrDBOperation = "db_operation"
)

type attrs struct{}

var KV = attrs{}

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

func (attrs) HTTPRequestPath(method string, u *url.URL) attribute.KeyValue {
	formatted := fmt.Sprintf("%s %s//%s%s", method, u.Scheme, u.Host, u.EscapedPath())
	return attribute.String(attrHTTPPath, formatted)
}