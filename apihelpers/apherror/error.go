package apherror

import (
	"net/http"

	"google.golang.org/api/googleapi"

	"github.com/go-chi/render"
)

// ErrResponse is to manage the http error response
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render implements https://godoc.org/github.com/go-chi/render#Renderer interface
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrServer reports internal server errors
func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "server error",
		ErrorText:      err.Error(),
	}
}

// ErrGdrive reports errors from gdrive api
func ErrGdrive(err error) render.Renderer {
	gerr, _ := err.(*googleapi.Error)
	return &ErrResponse{
		Err:            gerr,
		HTTPStatusCode: gerr.Code,
		StatusText:     gerr.Message,
		ErrorText:      gerr.Error(),
	}
}
