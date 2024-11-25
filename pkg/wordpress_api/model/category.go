package model

type CategorySchema struct {
	// Unique identifier for the category.
	// Context: view, edit, embed
	ID int `json:"id,omitempty"`
	// The number of published posts for the category.
	// Context: view, edit
	Count int `json:"count,omitempty"`
	// HTML description of the term.
	// Context: view, edit
	Description string `json:"description,omitempty"`
	// URL of the term.
	// Context: view, edit, embed
	Link string `json:"link,omitempty"`
	// HTML title for the term.
	// Context: view, edit, embed
	Name string `json:"name,omitempty"`
	// The parent term ID.
	// Context: view, edit
	Parent int `json:"parent,omitempty"`
}

type ListCategoryArgs struct {
	// Scope under which the request is made; determines fields present in response.
	Context ApiContext `json:"context,omitempty"`
	// Current page of the collection.
	Page int `json:"page,omitempty"`
	// Maximum number of items to be returned in result set.
	PerPage int `json:"per_page,omitempty"`
	// Limit results to those matching a string.
	Search string `json:"search,omitempty"`
	// Order sort attribute ascending or descending.
	Order ApiOrder `json:"order,omitempty"`
	// Sort collection by object attribute.
	OrderBy string `json:"orderby,omitempty"`
	// Whether to hide terms not assigned to any posts.
	HideEmpty bool `json:"hide_empty,omitempty"`
	// Limit result set to terms assigned to a specific parent.
	Parent int `json:"parent,omitempty"`
}

type ListCategoriesSchema []CategorySchema

type ListCategoryResponse ListCategoriesSchema
