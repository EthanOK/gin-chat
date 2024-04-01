package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total interface{}
}

func Response(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	json, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(json)
}

func ResponseFail(w http.ResponseWriter, msg string) {

	Response(w, -1, nil, msg)

}

func ResponseOK(w http.ResponseWriter, data interface{}, msg string) {

	Response(w, 0, data, msg)
}

func ResponseList(w http.ResponseWriter, code int, data interface{}, total interface{}) {

	w.Header().Set("Content-Type", "application/json")
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	json, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(json)

}

func ResponseOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	ResponseList(w, 0, data, total)

}
