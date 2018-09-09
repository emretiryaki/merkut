package model

import (
	"github.com/emretiryaki/merkut/pkg/components/simplejson"
	"time"
)

type AlertNotification struct {
	Id        int64            `json:"id"`
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	IsDefault bool             `json:"isDefault"`
	Settings  *simplejson.Json `json:"settings"`
	Created   time.Time        `json:"created"`
	Updated   time.Time        `json:"updated"`
}


type GetAlertNotificationsToSendQuery struct {
	Ids   []int64

	Result []*AlertNotification
}