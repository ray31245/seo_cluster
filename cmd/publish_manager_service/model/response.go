package model

import (
	"github.com/google/uuid"
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type site struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Lack int       `json:"lack"`
}

type ListSitesResponse struct {
	Sites []site `json:"sites"`
}

func (l *ListSitesResponse) FromDBSites(sites []model.Site) {
	for _, s := range sites {
		l.Sites = append(l.Sites, site{ID: s.ID, URL: s.URL, Lack: s.LackCount})
	}
}
