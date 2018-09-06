package notifiers

import "github.com/emretiryaki/merkut/pkg/components/simplejson"

type NotifierBase struct {
	Name        string
	Type        string
	Id          int64
	IsDeault    bool
	UploadImage bool
}

func NewNotifierBase(id int64, isDefault bool, name, notifierType string, model *simplejson.Json) NotifierBase {
	uploadImage := true
	value, exist := model.CheckGet("uploadImage")
	if exist {
		uploadImage = value.MustBool()
	}

	return NotifierBase{
		Id:          id,
		Name:        name,
		IsDeault:    isDefault,
		Type:        notifierType,
		UploadImage: uploadImage,
	}
}

