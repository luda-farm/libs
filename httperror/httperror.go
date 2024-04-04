package httperror

import (
	"fmt"
	"net/http"

	"github.com/luda-farm/libs/std"
)

type httpError struct {
	Status  int
	Message string
}

func ExitWithError(status int, message string) {
	panic(httpError{
		Status:  status,
		Message: message,
	})
}

func (err httpError) Write(w http.ResponseWriter) {
	http.Error(w, err.Message, err.Status)
}

func PanicHandler(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func handlePanic(w http.ResponseWriter) {
	switch err := recover().(type) {
	case nil:
	case httpError:
		err.Write(w)
	case error:
		fmt.Printf("[ERROR] %s\n", err.Error())
		std.PrintLocalStackTrace()
		w.WriteHeader(http.StatusInternalServerError)
	default:
		fmt.Printf("[ERROR] %v\n", err)
		std.PrintLocalStackTrace()
		w.WriteHeader(http.StatusInternalServerError)
	}
}
