package request

import (
	"net/http"
	"tech_check/internal/model"
)

var defaultParser *Parser

func InitParser() {
	defaultParser = NewParser()
}

func ParseBody(r *http.Request, req interface{}) error {
	return defaultParser.ParseBody(r, req)
}

func GetQuerySearch(r *http.Request) *Search {
	return defaultParser.GetQuerySearch(r)
}

func GetQueryMap(r *http.Request, key string) map[string]string {
	return defaultParser.GetQueryMap(r, key)
}

func GetQueryInt(r *http.Request, key string, defaultValue int) int {
	return defaultParser.GetQueryInt(r, key, defaultValue)
}

func GetHeaderIP(r *http.Request) string {
	return defaultParser.GetHeaderIP(r)
}

func GetAuthUser(r *http.Request) (*model.User, error) {
	return defaultParser.GetAuthUser(r)
}