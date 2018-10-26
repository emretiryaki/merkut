package sqlstore

import (
	"github.com/emretiryaki/merkut/pkg/bus"
	m"github.com/emretiryaki/merkut/pkg/model"
)

func init(){

	bus.AddHandler("sql",GetAlerts)
}

func GetAlerts(query *m.GetAlertsQuery) error{

	var alerts = make([]*m.Alert, 0)
	err :=x.Sql("select * from alarms").Find(&alerts)

	query.Result = alerts
	return err

}



