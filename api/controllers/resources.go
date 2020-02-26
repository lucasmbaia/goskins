package controllers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"regexp"

	"github.com/lucasmbaia/goskins/api/models/interfaces"
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
	)

	if err = r.setParams(c, fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	if data, err = m.Get(fields); err != nil {
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

	if err = r.setParams(c, data); err != nil {
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

func (r *Resources) setParams(c *gin.Context, filters interface{}) (err error) {
	var (
		rgx     = regexp.MustCompile(`\/(:[^:\/]*)`)
		matches []string
		params  = make(map[string]interface{})
		str     string
		body    []byte
	)

	matches = rgx.FindAllString(c.FullPath(), -1)

	for _, v := range matches {
		str = strings.Replace(v, "/:", "", -1)
		params[str] = c.Param(str)
	}

	if body, err = json.Marshal(params); err != nil {
		return
	}

	if err = json.Unmarshal(body, filters); err != nil {
		return
	}

	return
}
