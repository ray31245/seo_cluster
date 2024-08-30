package model

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type ListSitesResponse struct {
	Sites []string `json:"sites"`
}

func (l *ListSitesResponse) FromDBSites(sites []model.Site) {
	for _, site := range sites {
		l.Sites = append(l.Sites, site.URL)
	}
}
