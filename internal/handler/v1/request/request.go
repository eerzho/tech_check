package request

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"tech_check/internal/def"
	"tech_check/internal/model"

	"github.com/go-playground/validator/v10"
)

type Parser struct {
	validate *validator.Validate
}

func NewParser() *Parser {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	return &Parser{
		validate: validate,
	}
}

func (p *Parser) ParseBody(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return def.ErrInvalidBody
		}
		return err
	}

	err = p.validate.StructCtx(r.Context(), req)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) GetQuerySearch(r *http.Request) *Search {
	return &Search{
		Pagination: Pagination{
			Page:  p.GetQueryInt(r, "pagination[page]", 1),
			Count: p.GetQueryInt(r, "pagination[count]", 10),
		},
		Filters: p.GetQueryMap(r, "filters"),
		Sorts:   p.GetQueryMap(r, "sorts"),
	}
}

func (p *Parser) GetQueryMap(r *http.Request, key string) map[string]string {
	valuesMap := make(map[string]string)

	for queryKey, values := range r.URL.Query() {
		if strings.HasPrefix(queryKey, key+"[") && strings.HasSuffix(queryKey, "]") {
			fieldName := strings.TrimSuffix(strings.TrimPrefix(queryKey, key+"["), "]")
			if len(values) > 0 {
				valuesMap[fieldName] = values[0]
			}
		}
	}

	return valuesMap
}

func (p *Parser) GetQueryInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return valueInt
}

func (p *Parser) GetHeaderIP(r *http.Request) string {
	forwardedFor := r.Header.Get(def.HeaderForwardedFor.String())
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")

		return strings.TrimSpace(ips[0])
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

func (p *Parser) GetAuthUser(r *http.Request) (*model.User, error) {
	user, ok := r.Context().Value(def.ContextAuthUser).(*model.User)
	if !ok {
		return nil, def.ErrInvalidUserType
	}

	return user, nil
}
