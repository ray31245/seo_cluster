package model

type TagSchema struct {
	// Unique identifier for the tag.
	ID int `json:"id,omitempty"`
	// The number of published posts for the tag.
	Count int `json:"count,omitempty"`
	// HTML description of the term.
	Description string `json:"description,omitempty"`
	// URL of the term.
	Link string `json:"link,omitempty"`
	// HTML title for the term.
	Name string `json:"name,omitempty"`
}

func (t TagSchema) GetID() int {
	return t.ID
}

func (t TagSchema) GetName() string {
	return t.Name
}

func (t TagSchema) GetCount() int {
	return t.Count
}

type ListTagArgs struct {
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
}

type CreateTagArgs struct {
	// HTML title for the term.
	Name string `json:"name"`
	// HTML description of the tag.
	Description string `json:"description"`
}

type ListTagsSchema []TagSchema

type ListTagResponse ListTagsSchema

type CreateTagResponse struct {
	ID int `json:"id"`
}
