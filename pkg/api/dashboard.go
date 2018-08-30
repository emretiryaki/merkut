package api

import (
	m "github.com/emretiryaki/merkut/pkg/model"
	"time"
)

func GetAlarmList(c *m.ReqContext) {

	var alarms []*Alarm

	alarms = append(alarms, &Alarm{Name:"Test",Comment:"Deneme",Id:1,LastFired:time.Now(),LastTriggered:time.Now().Add(-4),State:"OK"})

	//if err != nil {
	//	c.JsonApiErr(500, "Failed to get tags from database", err)
	//	return
	//}
	c.JSON(200,alarms )
}


type Alarm struct {
	Name  string `json:"name"`
	Id    int    `json:"id"`
	State string  `json:"state"`
	Comment string  `json:"comment"`
	LastFired time.Time  `json:"lastFired,omitempty"`
	LastTriggered time.Time  `json:"lastTriggered,omitempty"`
}

