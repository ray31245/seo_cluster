package sitemanager

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	wordpressModel "github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	wordpressInterface "github.com/ray31245/seo_cluster/pkg/wordpress_api/wordpress_interface"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

var (
	ErrSiteNotFound        = errors.New("site not found")
	ErrCategoryNumNotMatch = errors.New("category number is not match")
)

// SiteManager is a struct that contains the necessary information for the site manager service.
type SiteManager struct {
	zAPI         zInterface.ZBlogAPI
	WordpressAPI wordpressInterface.WordpressAPI
	siteDAO      dbInterface.SiteDAOInterface
}

// NewSiteManager is a constructor for SiteManager.
func NewSiteManager(zAPI zInterface.ZBlogAPI, wordpressAPI wordpressInterface.WordpressAPI, siteDAO dbInterface.SiteDAOInterface) *SiteManager {
	return &SiteManager{
		zAPI:         zAPI,
		WordpressAPI: wordpressAPI,
		siteDAO:      siteDAO,
	}
}

// AddSite is a method that adds a site to the site manager.
func (s SiteManager) AddSite(ctx context.Context, cmsType string, urlStr string, userName string, password string, expectCategoryNum uint8) error {
	var err error
	if cmsType == string(dbModel.CMSTypeWordPress) {
		err = s.addWordPressSite(ctx, urlStr, userName, password, expectCategoryNum)
	} else if cmsType == string(dbModel.CMSTypeZBlog) {
		err = s.addZBlogSite(ctx, urlStr, userName, password, expectCategoryNum)
	} else {
		err = errors.New("CMS type not support")
	}

	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	return nil
}

func (s SiteManager) addZBlogSite(ctx context.Context, urlStr string, userName string, password string, expectCategoryNum uint8) error {
	// check site is valid
	client, err := s.zAPI.NewClient(ctx, urlStr, userName, password)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// list category of site
	categories, err := client.ListCategory(ctx)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	if expectCategoryNum != 0 && uint8(len(categories)) != expectCategoryNum {
		return fmt.Errorf("AddSite: %w", ErrCategoryNumNotMatch)
	}

	// add site
	site, err := s.siteDAO.CreateSite(&dbModel.Site{URL: urlStr, UserName: userName, Password: password, CmsType: dbModel.CMSTypeZBlog})
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

		err = s.siteDAO.CreateCategory(&dbModel.Category{SiteID: site.ID, ZBlogID: uint32(cateID), Name: cate.Name})
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	if multiErr != nil {
		return fmt.Errorf("AddSite: %w", multiErr)
	}

	return nil
}

func (s SiteManager) addWordPressSite(ctx context.Context, urlStr string, userName string, password string, expectCategoryNum uint8) error {
	// check site is valid
	client, err := s.WordpressAPI.NewClient(ctx, urlStr, userName, password)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// list category of site
	categories, err := client.ListCategory(ctx, wordpressModel.ListCategoryArgs{})
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	if expectCategoryNum != 0 && uint8(len(categories)) != expectCategoryNum {
		return fmt.Errorf("AddSite: %w", ErrCategoryNumNotMatch)
	}

	// add site
	site, err := s.siteDAO.CreateSite(&dbModel.Site{URL: urlStr, UserName: userName, Password: password, CmsType: dbModel.CMSTypeWordPress})
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add category
	var multiErr error

	for _, cate := range categories {
		err = s.siteDAO.CreateCategory(&dbModel.Category{SiteID: site.ID, WordpressID: uint32(cate.ID), Name: cate.Name})
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	if multiErr != nil {
		return fmt.Errorf("AddSite: %w", multiErr)
	}

	return nil
}

// DeleteSite is a method that deletes a site from the site manager.
func (s SiteManager) DeleteSite(siteID string) error {
	site, err := s.siteDAO.GetSite(siteID)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("DeleteSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("DeleteSite: %w", err)
	}

	if site.CmsType == dbModel.CMSTypeWordPress {
		s.WordpressAPI.DeleteClient(uuid.MustParse(siteID))
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		s.zAPI.DeleteClient(uuid.MustParse(siteID))
	}

	err = s.siteDAO.DeleteSite(siteID)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("DeleteSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("DeleteSite: %w", err)
	}

	err = s.siteDAO.DeleteSiteCategories(siteID)
	if err != nil {
		return fmt.Errorf("DeleteSite: %w", err)
	}

	return nil
}

// List sites of site manager
func (s SiteManager) ListSites() ([]dbModel.Site, error) {
	res, err := s.siteDAO.ListSites()
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}

	return res, nil
}

// Get site by id
func (s SiteManager) GetSite(siteID string) (*dbModel.Site, error) {
	site, err := s.siteDAO.GetSite(siteID)
	if dbErr.IsNotfoundErr(err) {
		return nil, fmt.Errorf("GetSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return nil, fmt.Errorf("GetSite: %w", err)
	}

	return site, nil
}

// Update site
func (s SiteManager) UpdateSite(ctx context.Context, ID string, urlStr string, userName string, password string) error {
	site, err := s.siteDAO.GetSite(ID)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("UpdateSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("UpdateSite: %w", err)
	}

	if urlStr != "" {
		site.URL = urlStr
	}

	if userName != "" {
		site.UserName = userName
	}

	if password != "" {
		site.Password = password
	}

	if site.CmsType == dbModel.CMSTypeWordPress {
		_, err = s.WordpressAPI.UpdateClient(ctx, site.ID, site.URL, site.UserName, site.Password)
		if err != nil {
			return fmt.Errorf("UpdateSite: %w", err)
		}
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		_, err = s.zAPI.UpdateClient(ctx, site.ID, site.URL, site.UserName, site.Password)
		if err != nil {
			return fmt.Errorf("UpdateSite: %w", err)
		}
	}

	err = s.siteDAO.UpdateSite(site)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("UpdateSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("UpdateSite: %w", err)
	}

	return nil
}

func (s SiteManager) SyncCategoryFromAllSite(ctx context.Context) error {
	sites, err := s.ListSites()
	if err != nil {
		return fmt.Errorf("SyncCategoryFromAllSite: %w", err)
	}

	var multiErr error

	for _, site := range sites {
		err = s.SyncCategoryFromSite(ctx, site.ID.String())
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	return err
}

func (s SiteManager) SyncCategoryFromSite(ctx context.Context, siteID string) error {
	site, err := s.siteDAO.GetSite(siteID)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("SyncCategoryFromSite: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("SyncCategoryFromSite: %w", err)
	}

	if site.CmsType == dbModel.CMSTypeWordPress {
		err = s.SyncCategoryFromSiteWordpress(ctx, site)
		if err != nil {
			return fmt.Errorf("SyncCategoryFromSite: %w", err)
		}
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		err = s.SyncCategoryFromSiteZblog(ctx, site)
		if err != nil {
			return fmt.Errorf("SyncCategoryFromSite: %w", err)
		}
	} else {
		return fmt.Errorf("SyncCategoryFromSite: CMS type not support")
	}

	return nil
}

func (s SiteManager) SyncCategoryFromSiteZblog(ctx context.Context, site *dbModel.Site) error {
	client, err := s.zAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteZblog: %w", err)
	}

	zCategories, err := client.ListCategory(ctx)
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteZblog: %w", err)
	}

	realCategories := []dbModel.Category{}

	for _, cate := range zCategories {
		catID, err := strconv.Atoi(cate.ID)
		if err != nil {
			return fmt.Errorf("SyncCategoryFromSiteZblog: %w", err)
		}

		realCategories = append(realCategories, dbModel.Category{SiteID: site.ID, ZBlogID: uint32(catID), Name: cate.Name})
	}

	currentCategories := site.Categories

	getCMSID := func(category dbModel.Category) uint32 {
		return category.ZBlogID
	}

	createCategory := func(category dbModel.Category) error {
		err := s.siteDAO.CreateCategory(&category)
		if err != nil {
			return fmt.Errorf("createCategoryCallBack: %w", err)
		}

		return nil
	}

	updateCategory := func(category dbModel.Category) error {
		err := s.siteDAO.UpdateCategory(&category)
		if err != nil {
			return fmt.Errorf("updateCategoryCallBack: %w", err)
		}

		return nil
	}

	deleteCategory := func(id string) error {
		err := s.siteDAO.DeleteCategory(id)
		if err != nil {
			return fmt.Errorf("deleteCategoryCallBack: %w", err)
		}

		return nil
	}

	err = syncCategories(realCategories, currentCategories, getCMSID, createCategory, updateCategory, deleteCategory)
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteZblog: %w", err)
	}

	return nil
}

func (s SiteManager) SyncCategoryFromSiteWordpress(ctx context.Context, site *dbModel.Site) error {
	client, err := s.WordpressAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteWordpress: %w", err)
	}

	wordpressCategories, err := client.ListCategory(ctx, wordpressModel.ListCategoryArgs{PerPage: int(time.Now().UnixMicro() % 100)})
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteWordpress: %w", err)
	}

	categories := []dbModel.Category{}
	for _, cate := range wordpressCategories {
		categories = append(categories, dbModel.Category{SiteID: site.ID, WordpressID: uint32(cate.ID), Name: cate.Name})
	}

	currentCategories := site.Categories

	getCMSID := func(category dbModel.Category) uint32 {
		return category.WordpressID
	}

	createCategory := func(category dbModel.Category) error {
		err := s.siteDAO.CreateCategory(&category)
		if err != nil {
			return fmt.Errorf("createCategoryCallBack: %w", err)
		}

		return nil
	}

	updateCategory := func(category dbModel.Category) error {
		err := s.siteDAO.UpdateCategory(&category)
		if err != nil {
			return fmt.Errorf("updateCategoryCallBack: %w", err)
		}

		return nil
	}

	deleteCategory := func(id string) error {
		err := s.siteDAO.DeleteCategory(id)
		if err != nil {
			return fmt.Errorf("deleteCategoryCallBack: %w", err)
		}

		return nil
	}

	err = syncCategories(categories, currentCategories, getCMSID, createCategory, updateCategory, deleteCategory)
	if err != nil {
		return fmt.Errorf("SyncCategoryFromSiteWordpress: %w", err)
	}

	return nil
}

func syncCategories(
	realCategories []dbModel.Category, currentCategories []dbModel.Category, getCMSID func(dbModel.Category) uint32,
	createCategory func(dbModel.Category) error, updateCategory func(dbModel.Category) error, deleteCategory func(string) error,
) error {
	expectedCategories := map[uint32]dbModel.Category{}
	for _, cate := range realCategories {
		expectedCategories[getCMSID(cate)] = cate
	}

	currentCategoriesMap := map[uint32]dbModel.Category{}
	for _, cate := range currentCategories {
		currentCategoriesMap[getCMSID(cate)] = cate
	}

	var multiErr error

	for cmsID, cate := range expectedCategories {
		if currentCate, ok := currentCategoriesMap[cmsID]; ok {
			if cate.Name != currentCate.Name {
				currentCate.Name = cate.Name
				err := updateCategory(currentCate)
				if err != nil {
					multiErr = errors.Join(multiErr, err)
				}
			}
		} else {
			err := createCategory(cate)
			if err != nil {
				multiErr = errors.Join(multiErr, err)
			}
		}
	}

	for cmsID, cate := range currentCategoriesMap {
		if _, ok := expectedCategories[cmsID]; !ok {
			err := deleteCategory(cate.ID.String())
			if err != nil {
				multiErr = errors.Join(multiErr, err)
			}
		}
	}

	return nil
}

// Increase lack count of site
func (s SiteManager) IncreaseLackCount(siteID string, count int) error {
	err := s.siteDAO.IncreaseLackCount(siteID, count)
	if dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("IncreaseLackCount: %w", errors.Join(ErrSiteNotFound, err))
	} else if err != nil {
		return fmt.Errorf("IncreaseLackCount: %w", err)
	}

	return nil
}
