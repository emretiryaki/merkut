package api

import (
		"github.com/emretiryaki/merkut/pkg/routing"
	m "github.com/emretiryaki/merkut/pkg/model"
)

func (hs *HTTPServer) registerRoutes() {


	r := hs.RouteRegister

	r.Get("/",  Index)



	r.Group("/api", func(apiRoute routing.RouteRegister) {

		apiRoute.Get("/search/", Search)


	})


}

func Index(c *m.ReqContext) {
	data, _ := setIndexViewData(c)
	c.HTML(200, "index", data)
}