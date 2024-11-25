package model

type CommentSchema struct {
	ID              int            `json:"id,omitempty"`
	Author          int            `json:"author,omitempty"`
	AuthorName      string         `json:"author_name,omitempty"`
	AuthorURL       string         `json:"author_url,omitempty"`
	AuthorIP        string         `json:"author_ip,omitempty"`
	AuthorEmail     string         `json:"author_email,omitempty"`
	AuthorUserAgent string         `json:"author_user_agent,omitempty"`
	Content         CommentContent `json:"content,omitempty"`
	Date            string         `json:"date,omitempty"`
	DateGmt         string         `json:"date_gmt,omitempty"`
	Link            string         `json:"link,omitempty"`
	Parent          int            `json:"parent,omitempty"`
	Post            int            `json:"post,omitempty"`
	Status          string         `json:"status,omitempty"`
	Type            string         `json:"type,omitempty"`
}

type CommentContent struct {
	Raw      string `json:"raw,omitempty"`
	Rendered string `json:"rendered,omitempty"`
}

type ListCommentArgs struct {
	// Scope under which the request is made; determines fields present in response.
	Context ApiContext `json:"context,omitempty"`
	// Current page of the collection.
	Page int `json:"page,omitempty"`
	// Maximum number of items to be returned in result set.
	PerPage int `json:"per_page,omitempty"`
	// Limit results to those matching a string.
	Search string `json:"search,omitempty"`
	// Limit response to comments published after a given date.
	After string `json:"after,omitempty"`
	// Limit response to comments published before a given date.
	Before string `json:"before,omitempty"`
	// Ensure result set excludes specific parent IDs.
	Exclude []int `json:"exclude,omitempty"`
	// Limit result set to specific parent IDs.
	Include []int `json:"include,omitempty"`
	// Offset the result set by a specific number of items.
	Offset int `json:"offset,omitempty"`
	// Order sort attribute ascending or descending.
	Order ApiOrder `json:"order,omitempty"`
	// Sort collection by object attribute.
	OrderBy string `json:"orderby,omitempty"`
	// Limit result set to comments assigned to specific post IDs.
	Post int `json:"post,omitempty"`
	// Limit result set to comments assigned a specific status.
	Status string `json:"status,omitempty"`
	// Limit result set to comments of specific type.
	Type string `json:"type,omitempty"`
}

type CreateCommentArgs struct {
	// The ID of the post object the comment belongs to.
	Post int `json:"post"`
	// The ID of the parent comment.
	Parent int `json:"parent,omitempty"`
	// The name of the author of the comment.
	Author string `json:"author,omitempty"`
	// The email of the author of the comment.
	AuthorEmail string `json:"author_email,omitempty"`
	// The URL of the author of the comment.
	AuthorURL string `json:"author_url,omitempty"`
	// The IP address for the comment author.
	AuthorIP string `json:"author_ip,omitempty"`
	// The date the comment was posted in the site's timezone.
	Date string `json:"date,omitempty"`
	// The date the comment was posted in GMT timezone.
	DateGmt string `json:"date_gmt,omitempty"`
	// The content of the comment.
	Content string `json:"content,omitempty"`
	// The status of the comment.
	Status string `json:"status,omitempty"`
}

type ListCommentsSchema []CommentSchema

type ListCommentResponse ListCommentsSchema

type CreateCommentResponse struct {
	ID int `json:"id"`
}
