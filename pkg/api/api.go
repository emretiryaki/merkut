package api

import (
	"github.com/emretiryaki/merkut/pkg/api/dto"
	"github.com/go-macaron/binding"
	"github.com/emretiryaki/merkut/pkg/routing"
	m "github.com/emretiryaki/merkut/pkg/model"

)

func (hs *HTTPServer) registerRoutes() {

	bind := binding.Bind
	r := hs.RouteRegister

	r.Get("/",  Index)



	r.Group("/api", func(apiRoute routing.RouteRegister) {

		apiRoute.Get("/search/", Search)
		apiRoute.Get("/alerts/", GetAlarmList)
		apiRoute.Post("/alerts/",  bind(dto.AddAlertCommand{}),Wrap(AddAlert))

	})


}

func Index(c *m.ReqContext) {
	data, _ := setIndexViewData(c)
	c.HTML(200, "index", data)
}