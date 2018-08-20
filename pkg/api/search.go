package api

import (
	"strconv"
	 search "github.com/emretiryaki/merkut/pkg/services/search"
	 m "github.com/emretiryaki/merkut/pkg/model"
	"github.com/emretiryaki/merkut/pkg/bus"
)

func Search(c *m.ReqContext) {
	query := c.Query("query")
	tags := c.QueryStrings("tag")
	starred := c.Query("starred")
	limit := c.QueryInt("limit")
	dashboardType := c.Query("type")


	if limit == 0 {
		limit = 1000
	}



	dbIDs := make([]int64, 0)
	for _, id := range c.QueryStrings("dashboardIds") {
		dashboardID, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			dbIDs = append(dbIDs, dashboardID)
		}
	}

	folderIDs := make([]int64, 0)
	for _, id := range c.QueryStrings("folderIds") {
		folderID, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			folderIDs = append(folderIDs, folderID)
		}
	}

	searchQuery := search.Query{
		Title:        query,
		Tags:         tags,
		SignedInUser: c.SignedInUser,
		Limit:        limit,
		IsStarred:    starred == "true",
		OrgId:        c.OrgId,
		DashboardIds: dbIDs,
		Type:         dashboardType,
		FolderIds:    folderIDs,
	}

	err := bus.Dispatch(&searchQuery)
	if err != nil {
		c.JsonApiErr(500, "Search failed", err)
		return
	}

	c.JSON(200, searchQuery.Result)
}

