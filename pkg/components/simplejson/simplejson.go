package simplejson

import (
	"errors"
	"log"
)

type Json struct {
	data interface{}
}

func (j *Json)  CheckGet(key string) (*Json,bool) {
	m, err :=  j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}, true
		}
	}
	return nil, false
}

func (j *Json) Map() (map[string]interface{}, error) {

	if m,ok := (j.data).(map[string]interface{}); ok{
		return m,nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

func (j *Json) MustBool(args ...bool) bool {
	var def bool

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustBool() received too many arguments %d", len(args))
	}

	b, err := j.Bool()
	if err == nil {
		return b
	}

	return def
}
func (j *Json) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}