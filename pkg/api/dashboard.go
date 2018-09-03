package api

import (
	"github.com/emretiryaki/merkut/pkg/bus"
	m "github.com/emretiryaki/merkut/pkg/model"
)

func GetAlarmList(c *m.ReqContext)  {

	getalertsQuery := m.GetAllAlertsQuery{}

	if err := bus.Dispatch(&getalertsQuery); err != nil {
		c.JsonApiErr(500, "Failed to get alarm list from database", err)
		return
	}
	c.JSON(200,getalertsQuery.Result)
}

func AddAlert(c *m.ReqContext){

}

