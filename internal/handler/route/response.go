package route

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, code int, errMsg, msg string) {
	log.Println(errMsg)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Error{Code: int32(code), Message: msg})
}
