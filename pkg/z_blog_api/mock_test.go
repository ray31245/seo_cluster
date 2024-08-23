package zblogapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	zAPI "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/origin"
)

// routeTable is a map of route path to handler function
// routeTable[mod][act] = handler
type routeTable map[string]map[string]http.HandlerFunc

func newMockServer(route routeTable) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+origin.APIPath {
			http.Error(w, `{"code":404,"message":"not found"}`, http.StatusNotFound)

			return
		}

		act := r.URL.Query().Get("act")
		mod := r.URL.Query().Get("mod")

		if route[mod] == nil {
			http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

			return
		}

		if route[mod][act] == nil {
			http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

			return
		}

		route[mod][act](w, r)
	}))
}

func fullRouteTable() routeTable {
	return routeTable{
		"member": map[string]http.HandlerFunc{
			"list":  listMemberHandler,
			"login": zAPI.LoginHandler,
		},
		"post": map[string]http.HandlerFunc{
			"post":   postArticleHandler,
			"list":   listArticleHandler,
			"delete": deleteArticleHandler,
		},
		"category": map[string]http.HandlerFunc{
			"list": ListCategoryHandler,
		},
	}
}

func newFullMockServer() *httptest.Server {
	return newMockServer(fullRouteTable())
}

func newMockFaultServer(mod, act string) *httptest.Server {
	faultHandler := func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}

	rt := fullRouteTable()
	if rt[mod] == nil {
		rt[mod] = map[string]http.HandlerFunc{}
	}

	rt[mod][act] = faultHandler

	return newMockServer(rt)
}

func validator(w http.ResponseWriter, r *http.Request, method, paramMod, paramAct string) bool {
	// assert route path is correct
	if r.URL.Path != "/"+origin.APIPath {
		http.Error(w, `{"code":404,"message":"not found"}`, http.StatusNotFound)

		return false
	}

	if r.Header.Get("Authorization") != "Bearer "+zAPI.TestToken {
		http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)

		return false
	}

	// assert method is correct
	if r.Method != method {
		http.Error(w, `{"code":405,"message":"method not allowed"}`, http.StatusMethodNotAllowed)

		return false
	}

	// assert parameter is correct
	if r.URL.Query().Get("act") != paramAct {
		http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

		return false
	}

	if r.URL.Query().Get("mod") != paramMod {
		http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

		return false
	}

	return true
}

var listMemberHandler = func(w http.ResponseWriter, r *http.Request) {
	ok := validator(w, r, http.MethodGet, origin.ModMember, origin.ActList)
	if !ok {
		return
	}

	// return member list
	data := model.ListMemberResponse{
		BasicResponse: model.BasicResponse{
			Code: 200,
		},
		Data: model.Data[model.Member]{List: mockMembers},
	}

	resBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resBytes)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

var postArticleHandler = func(w http.ResponseWriter, r *http.Request) {
	ok := validator(w, r, http.MethodPost, origin.ModPost, origin.ActPost)
	if !ok {
		return
	}

	// assert post info is correct
	var reqBody model.PostArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"code":400,"message":"bad request"}`, http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"code":200}`))
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

var listArticleHandler = func(w http.ResponseWriter, r *http.Request) {
	ok := validator(w, r, http.MethodGet, origin.ModPost, origin.ActList)
	if !ok {
		return
	}

	cateID := r.URL.Query().Get("cate_id")
	intCateID, err := strconv.Atoi(cateID)
	if err != nil {
		http.Error(w, `{"code":400,"message":"bad request"}`, http.StatusBadRequest)
	}

	req := model.ListArticleRequest{
		CateID: uint32(intCateID),
	}

	articles := []model.Article{}
	for _, v := range mockArticles {
		if req.CateID != 0 && v.CateID != strconv.Itoa(int(req.CateID)) {
			continue
		}
		articles = append(articles, v)
	}

	// return article list
	data := model.ListArticleResponse{
		BasicResponse: model.BasicResponse{
			Code: 200,
		},
		Data: model.Data[model.Article]{
			List: articles,
			PageBar: model.PageBar{
				AllCount: uint32(len(articles)),
			},
		},
	}

	resBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resBytes)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

var deleteArticleHandler = func(w http.ResponseWriter, r *http.Request) {
	ok := validator(w, r, http.MethodGet, origin.ModPost, origin.ActDelete)
	if !ok {
		return
	}

	// assert delete info is correct
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"code":400,"message":"bad request"}`, http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"code":200}`))
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

var ListCategoryHandler = func(w http.ResponseWriter, r *http.Request) {
	ok := validator(w, r, http.MethodGet, origin.ModCategory, origin.ActList)
	if !ok {
		return
	}

	// return category list
	data := model.ListCategoryResponse{
		BasicResponse: model.BasicResponse{
			Code: 200,
		},
		Data: model.Data[model.Category]{List: mockCategories},
	}

	resBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resBytes)
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

var mockMembers = []model.Member{
	{
		ID:     "1",
		Level:  "1",
		Status: "1",
		Name:   "mock1",
		Email:  "mock1@example.com",
	},
}

var mockArticles = []model.Article{
	{
		ID:     "1",
		CateID: "1",
		Title:  "mock1",
	},
	{
		ID:     "2",
		CateID: "2",
		Title:  "mock2",
	},
}

var mockCategories = []model.Category{
	{
		ID:   "1",
		Name: "mock1",
	},
	{
		ID:   "2",
		Name: "mock2",
	},
}
