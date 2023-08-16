package code

import (
	"net/http"

	"github.com/novalagung/gubrak"
	"github.com/rshulabs/HgCache/pkg/errorx"
)

type ErrCode struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

var (
	_ errorx.Coder = &ErrCode{}
)

func (coder ErrCode) Code() int {
	return coder.C
}

func (coder ErrCode) String() string {
	return coder.Ext
}

func (coder ErrCode) Reference() string {
	return coder.Ref
}

func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 201, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 201, 400, 401, 403, 404, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	errorx.MustRegister(coder)
}
