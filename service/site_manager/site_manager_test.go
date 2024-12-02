package sitemanager

import (
	"log"
	"testing"

	"github.com/google/uuid"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	"github.com/stretchr/testify/assert"
)

func Test_syncCategories(t *testing.T) {
	gotCreate := []dbModel.Category{}
	gotUpdate := []dbModel.Category{}
	gotDelete := []dbModel.Category{}
	type args struct {
		realCategories    []dbModel.Category
		currentCategories []dbModel.Category
		getCMSID          func(dbModel.Category) uint32
	}
	tests := []struct {
		name           string
		args           args
		expectedCreate []dbModel.Category
		expectedUpdate []dbModel.Category
		expectedDelete []dbModel.Category
		wantErr        bool
	}{
		{
			name: "normal",
			args: args{
				realCategories: []dbModel.Category{
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}, Name: "category1", ZBlogID: 1},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")}, Name: "category2", ZBlogID: 3},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003")}, Name: "category3", ZBlogID: 5},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004")}, Name: "category4", ZBlogID: 7},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005")}, Name: "category5", ZBlogID: 9},
				},
				currentCategories: []dbModel.Category{
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}, Name: "category10", ZBlogID: 1},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")}, Name: "category2", ZBlogID: 2},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003")}, Name: "category3", ZBlogID: 4},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004")}, Name: "category4", ZBlogID: 6},
					{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005")}, Name: "category5", ZBlogID: 8},
				},
				getCMSID: func(category dbModel.Category) uint32 {
					return category.ZBlogID
				},
			},
			expectedCreate: []dbModel.Category{
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")}, Name: "category2", ZBlogID: 3},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003")}, Name: "category3", ZBlogID: 5},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004")}, Name: "category4", ZBlogID: 7},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005")}, Name: "category5", ZBlogID: 9},
			},
			expectedUpdate: []dbModel.Category{
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}, Name: "category10", ZBlogID: 1},
			},
			expectedDelete: []dbModel.Category{
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")}},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003")}},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004")}},
				{Base: dbModel.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005")}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			// initialize
			gotCreate = []dbModel.Category{}
			gotUpdate = []dbModel.Category{}
			gotDelete = []dbModel.Category{}
			createCategory := func(category dbModel.Category) error {
				gotCreate = append(gotCreate, category)
				return nil
			}

			updateCategory := func(category dbModel.Category) error {
				gotUpdate = append(gotUpdate, category)
				return nil
			}

			deleteCategory := func(id string) error {
				gotDelete = append(gotDelete, dbModel.Category{Base: dbModel.Base{ID: uuid.MustParse(id)}})
				return nil
			}
			err := syncCategories(tt.args.realCategories, tt.args.currentCategories, tt.args.getCMSID,
				createCategory, updateCategory, deleteCategory)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncCategoriesV2() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, v := range gotCreate {
				log.Printf("gotCreate ID: %s Name: %s CMSID: %d", v.ID, v.Name, tt.args.getCMSID(v))
			}
			for _, v := range gotDelete {
				log.Printf("gotDelete ID: %s Name: %s CMSID: %d", v.ID, v.Name, tt.args.getCMSID(v))
			}
			for _, v := range gotUpdate {
				log.Printf("gotUpdate ID: %s Name: %s CMSID: %d", v.ID, v.Name, tt.args.getCMSID(v))
			}

			assert.Equal(len(tt.expectedCreate), len(gotCreate))
			assert.ElementsMatch(tt.expectedCreate, gotCreate)

			assert.Equal(len(tt.expectedUpdate), len(gotUpdate))
			assert.ElementsMatch(tt.expectedUpdate, gotUpdate)

			assert.Equal(len(tt.expectedDelete), len(gotDelete))
			assert.ElementsMatch(tt.expectedDelete, gotDelete)
		})
	}
}
