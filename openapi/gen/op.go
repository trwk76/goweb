package gen

import (
	"fmt"
	"net/http"

	"github.com/trwk76/goweb/openapi/spec"
)

func (p *PathItem) GET(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().GET = p.op(operationID, http.MethodGet, false, setup)
}

func (p *PathItem) PUT(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().PUT = p.op(operationID, http.MethodPut, true, setup)
}

func (p *PathItem) POST(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().POST = p.op(operationID, http.MethodPost, true, setup)
}

func (p *PathItem) DELETE(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().DELETE = p.op(operationID, http.MethodDelete, false, setup)
}

func (p *PathItem) OPTIONS(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().OPTIONS = p.op(operationID, http.MethodOptions, false, setup)
}

func (p *PathItem) HEAD(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().HEAD = p.op(operationID, http.MethodHead, false, setup)
}

func (p *PathItem) PATCH(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().PATCH = p.op(operationID, http.MethodPatch, true, setup)
}

func (p *PathItem) TRACE(operationID string, setup SetupFunc[spec.Operation]) {
	p.ensureSpec().TRACE = p.op(operationID, http.MethodTrace, false, setup)
}

func (p *PathItem) op(operationID string, method string, acceptBody bool, setup SetupFunc[spec.Operation]) *spec.Operation {
	if _, ex := p.api.ops[operationID]; ex {
		panic(fmt.Errorf("operationId '%s' redefined", operationID))
	}

	p.api.ops[operationID] = fmt.Sprintf("%s %s", method, p.pth)

	item := &spec.Operation{
		OperationID: operationID,
		Responses:   make(spec.Responses),
		Parameters:  p.params(),
	}

	if setup != nil {
		setup(p.api, item)
	}

	if item.RequestBody != nil && !acceptBody {
		panic(fmt.Errorf("method '%s' does not accept a request body", method))
	}

	return item
}
