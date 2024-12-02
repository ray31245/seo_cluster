package model

import (
	"github.com/google/uuid"
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type site struct {
	ID      uuid.UUID `json:"id"`
	URL     string    `json:"url"`
	CMSType string    `json:"cms_type"`
	Lack    int       `json:"lack"`
}

type category struct {
	Name        string `json:"name"`
	ZblogID     uint32 `json:"z_blog_id"`
	WordpressID uint32 `json:"wordpress_id"`
}

type ListSitesResponse struct {
	Sites []site `json:"sites"`
}

func (l *ListSitesResponse) FromDBSites(sites []model.Site) {
	for _, s := range sites {
		l.Sites = append(l.Sites, site{ID: s.ID, URL: s.URL, Lack: s.LackCount, CMSType: string(s.CmsType)})
	}
}

type GetSiteResponse struct {
	Site       site       `json:"site"`
	Categories []category `json:"categories"`
}

func (g *GetSiteResponse) FromDBSite(s model.Site) {
	g.Site = site{ID: s.ID, URL: s.URL, Lack: s.LackCount, CMSType: string(s.CmsType)}
	for _, c := range s.Categories {
		g.Categories = append(g.Categories, category{Name: c.Name, ZblogID: c.ZBlogID, WordpressID: c.WordpressID})
	}
}
