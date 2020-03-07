package controllers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"regexp"
	//"fmt"

	"github.com/lucasmbaia/goskins/api/models/interfaces"
	"github.com/lucasmbaia/goskins/api/repository/filter"
	"github.com/gin-gonic/gin"
)

type Resources struct {
	GetModel  func() interfaces.Models
	GetFields func() interface{}
}

func (r *Resources) Get(c *gin.Context) {
	var (
		m	= r.GetModel()
		fields	= r.GetFields()
		data	interface{}
		err	error
		filters	[]filter.Filters
		args	[]interface{}
	)

	if err = c.ShouldBindQuery(fields); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	if filters, err = r.setParams(c, fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	if data, err = m.Get(filters, args); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
	return
}

func (r *Resources) Post(c *gin.Context) {
	var (
		m     = r.GetModel()
		data  = r.GetFields()
		err   error
		v     reflect.Value
		async bool
	)

	if _, err = r.setParams(c, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	if err = c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if async, err = m.Post(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	v = reflect.ValueOf(data).Elem()

	if async {
		c.JSON(http.StatusAccepted, gin.H{"id": v.FieldByName("ID").Interface()})
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": v.FieldByName("ID").Interface()})
	}

	return
}

func (r *Resources) setParams(c *gin.Context, data interface{}) (filters []filter.Filters, err error) {
	var (
		rgx     = regexp.MustCompile(`\/(:[^:\/]*)`)
		matches []string
		params  = make(map[string]interface{})
		str     string
		body    []byte
		v	reflect.Value
		t	reflect.Type
		param	string
		name	string
		field	string
	)

	matches = rgx.FindAllString(c.FullPath(), -1)

	for _, v := range matches {
		str = strings.Replace(v, "/:", "", -1)
		params[str] = c.Param(str)
	}

	v = reflect.ValueOf(data).Elem()
	t = v.Type()

	for i := 0; i < v.NumField(); i++ {
		name = t.Field(i).Name

		if tag, ok := t.FieldByName(name); ok {
			field = tag.Tag.Get("model")
			param = tag.Tag.Get("param")

			if field == "" {
				field = name
			}

			if !v.Field(i).IsZero() {
				filters = append(filters, filter.Filters{
					Conditions: filter.Conditions{Field: field}.Eq(v.Field(i).Interface()),
				})
			}

			if param != "" {
				if _, ok := params[param]; ok {
					v.Field(i).Set(reflect.ValueOf(params[param]))


					filters = append(filters, filter.Filters{
						Conditions: filter.Conditions{Field: field}.Eq(params[param]),
					})
				}
			}
		}
	}

	if body, err = json.Marshal(params); err != nil {
		return
	}

	if err = json.Unmarshal(body, data); err != nil {
		return
	}

	return
}
