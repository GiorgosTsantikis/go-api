package util

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func CleanHTML() {

}

// Extract a code from a URL. Return the default code if code
// is missing or code is not a valid number.
func GetCode(r *http.Request, index int) (int, error) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) <= index {
		return 0, errors.New("missing ammount of tables")
	} else {
		val, err := strconv.Atoi(p[index])
		if err != nil {
			return 0, errors.New("path param isn't a number")
		}
		return val, nil
	}
}
