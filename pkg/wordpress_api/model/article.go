package model

import (
	"time"

	"github.com/ray31245/seo_cluster/pkg/util"
)

type (
	ArticleStatus string
)

const (
	// status
	StatusPublish ArticleStatus = "publish"
	StatusDraft   ArticleStatus = "draft"
	StatusPending ArticleStatus = "pending"
	StatusPrivate ArticleStatus = "private"
	StatusFuture  ArticleStatus = "future"
)

// The schema defines all the fields that exist within a post record.
type ArticleSchema struct {
	// Unique identifier for the post.
	// Context: view, edit, embed
	ID int `json:"id,omitempty"`
	// The date the post was published, in the site's timezone.
	// Context: view, edit, embed
	Date string `json:"date,omitempty"`
	// The date the post was published, as GMT.
	// Context: view, edit
	DateGmt string `json:"date_gmt,omitempty"`
	// The date the post was last modified, in the site's timezone.
	// Context: view, edit
	Modified string `json:"modified,omitempty"`
	// URL to the post.
	// Context: view, edit, embed
	Link string `json:"link,omitempty"`
	// The title for the post.
	// Context: view, edit, embed
	Title ArticleTitle `json:"title"`
	// The content for the post.
	// Context: view, edit
	Content ArticleContent `json:"content"`
	// The excerpt for the post.
	// Context: view, edit, embed
	Excerpt ArticleExcerpt `json:"excerpt"`
	// The ID for the author of the post.
	// Context: view, edit, embed
	Author int `json:"author,omitempty"`
	// The featured media for the post.
	// Context: view, edit, embed
	FeaturedMedia int `json:"featured_media,omitempty"`
	// The format for the post.
	// Context: view, edit, embed
	Format string `json:"format,omitempty"`
	// Whether or not the post is sticky.
	// Context: view, edit
	Sticky bool `json:"sticky,omitempty"`
	// The post's categories.
	// Context: view, edit
	Categories []int `json:"categories,omitempty"`
	// The post's tags.
	// Context: view, edit
	Tags []int `json:"tags,omitempty"`
}

type ArticleTitle struct {
	Raw      string `json:"raw"`
	Rendered string `json:"rendered"`
}

type ArticleContent struct {
	Raw          string `json:"raw"`
	Rendered     string `json:"rendered"`
	Protected    bool   `json:"protected"`
	BlockVersion int    `json:"block_version"`
}

type ArticleExcerpt struct {
	Raw       string `json:"raw"`
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected"`
}

type ListArticleArgs struct {
	// Scope under which the request is made; determines fields present in response.
	Context ApiContext `json:"context,omitempty"`
	// Maximum number of items to be returned in result set.
	PerPage int `json:"per_page,omitempty"`
	// Current page of the collection.
	Page int `json:"page,omitempty"`
	// Limit results to those matching a string.
	Search string `json:"search,omitempty"`
	// Limit response to posts published after a given ISO8601 compliant date.
	After string `json:"after,omitempty"`
	// Limit response to posts modified after a given ISO8601 compliant date.
	ModifiedAfter string `json:"modified_after,omitempty"`
	// Limit result set to posts assigned to specific authors.
	Author int `json:"author,omitempty"`
	// Ensure result set excludes posts assigned to specific authors.
	AuthorExclude int `json:"author_exclude,omitempty"`
	// Limit response to posts published before a given ISO8601 compliant date.
	Before string `json:"before,omitempty"`
	// Limit response to posts modified before a given ISO8601 compliant date.
	ModifiedBefore string `json:"modified_before,omitempty"`
	// Ensure result set excludes specific IDs.
	Exclude []int `json:"exclude,omitempty"`
	// Limit result set to specific IDs.
	Include []int `json:"include,omitempty"`
	// Offset the result set by a specific number of items.
	Offset int `json:"offset,omitempty"`
	// Order sort attribute ascending or descending.
	Order string `json:"order,omitempty"`
	// Sort collection by post attribute.
	OrderBy string `json:"orderby,omitempty"`
	// Limit result set to posts assigned one or more statuses.
	Status []string `json:"status,omitempty"`
	// Limit result set to items with specific terms assigned in the categories taxonomy.
	Categories []int `json:"categories,omitempty"`
	// Limit result set to items except those with specific terms assigned in the categories taxonomy.
	CategoriesExclude []int `json:"categories_exclude,omitempty"`
	// Limit result set to items with specific terms assigned in the tags taxonomy.
	Tags []int `json:"tags,omitempty"`
	// Limit result set to items except those with specific terms assigned in the tags taxonomy.
	TagsExclude []int `json:"tags_exclude,omitempty"`
}

type Date struct {
	Time time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return util.EncodeFormatedTime(d.Time)
}

type CreateArticleArgs struct {
	// The date the post was published, in the site's timezone.
	Date *Date `json:"date,omitempty"`
	// The date the post was published, as GMT.
	DateGmt string `json:"date_gmt,omitempty"`
	// The title for the post.
	Title string `json:"title,omitempty"`
	// The content for the post.
	Content string `json:"content,omitempty"`
	// The excerpt for the post.
	Excerpt string `json:"excerpt,omitempty"`
	// The terms assigned to the post in the category taxonomy.
	Categories []uint32 `json:"categories,omitempty"`
	// The terms assigned to the post in the post_tag taxonomy.
	Tags []int `json:"tags,omitempty"`
	// A named status for the post.
	Status ArticleStatus `json:"status,omitempty"`
}

type UpdateArticleArgs struct {
	// Unique identifier for the post.
	ID int `json:"id,omitempty"`
	// The terms assigned to the post in the post_tag taxonomy.
	Tags []int `json:"tags,omitempty"`
}

type RetrieveArticleArgs struct {
	// Unique identifier for the post.
	ID int `json:"id,omitempty"`
	// Scope under which the request is made; determines fields present in response.
	Context ApiContext `json:"context,omitempty"`
	// The password for the post if it is password protected.
	Password string `json:"password,omitempty"`
}

type CreateArticleResponse ArticleSchema

type UpdateArticleResponse ArticleSchema

type ListArticleSchema []ArticleSchema

type ListArticleResponse ListArticleSchema

type RetrieveArticleResponse ArticleSchema
