package api

import (
	"github.com/emretiryaki/merkut/pkg/middleware"
	"github.com/emretiryaki/merkut/pkg/routing"
	m "github.com/emretiryaki/merkut/pkg/model"
)

func (hs *HTTPServer) registerRoutes() {
	reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})


	r := hs.RouteRegister

	r.Get("/", reqSignedIn, Index)
	r.Get("/logout", Logout)
	r.Get("/login", LoginView)


	r.Group("/api", func(apiRoute routing.RouteRegister) {
		// Search
		apiRoute.Get("/search/", Search)


	}, reqSignedIn)


}

func Index(c *m.ReqContext) {
	data, _ := setIndexViewData(c)
	c.HTML(200, "index", data)
}