package sitemanager

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

// SiteManager is a struct that contains the necessary information for the site manager service.
type SiteManager struct {
	zAPI    zInterface.ZBlogAPI
	siteDAO dbInterface.SiteDAOInterface
}

// NewSiteManager is a constructor for SiteManager.
func NewSiteManager(zAPI zInterface.ZBlogAPI, siteDAO dbInterface.SiteDAOInterface) *SiteManager {
	return &SiteManager{
		zAPI:    zAPI,
		siteDAO: siteDAO,
	}
}

// AddSite is a method that adds a site to the site manager.
func (s SiteManager) AddSite(ctx context.Context, urlStr string, userName string, password string) error {
	// check site is valid
	client, err := s.zAPI.NewClient(ctx, urlStr, userName, password)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add site
	site, err := s.siteDAO.CreateSite(&dbModel.Site{URL: urlStr, UserName: userName, Password: password})
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// list category of site
	categories, err := client.ListCategory(ctx)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add category
	var multiErr error

	for _, cate := range categories {
		cateID, err := strconv.Atoi(cate.ID)
		if err != nil {
			multiErr = errors.Join(multiErr, err)

			continue
		}

		err = s.siteDAO.CreateCategory(&dbModel.Category{SiteID: site.ID, ZBlogID: uint32(cateID)})
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	if multiErr != nil {
		return fmt.Errorf("AddSite: %w", multiErr)
	}

	return nil
}

// List sites of publish manager
func (s SiteManager) ListSites() ([]dbModel.Site, error) {
	res, err := s.siteDAO.ListSites()
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}

	return res, nil
}