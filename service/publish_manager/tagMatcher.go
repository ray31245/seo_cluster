package publishmanager

import (
	"strings"
)

type tagInterface interface {
	GetName() string
	GetID() int
	GetCount() int
}

type TagMatcher struct {
	tagBlackListMap map[string]bool
	siteTagsMap     map[string]tagInterface
}

func NewTagMatcher[T tagInterface](tagBlackList []string, siteTags []T) *TagMatcher {
	tagBlackListMap := make(map[string]bool)

	for _, tag := range tagBlackList {
		tagBlackListMap[strings.ToLower(tag)] = true
	}

	siteTagsMap := map[string]tagInterface{}
	for _, t := range siteTags {
		siteTagsMap[strings.ToLower(t.GetName())] = t
	}

	return &TagMatcher{
		tagBlackListMap: tagBlackListMap,
		siteTagsMap:     siteTagsMap,
	}
}

func (t *TagMatcher) IsTagBlackList(tag string) bool {
	return t.tagBlackListMap[strings.ToLower(tag)]
}

func (t *TagMatcher) IsTagInSite(tag string) (bool, tagInterface) {
	matchTag, ok := t.siteTagsMap[strings.ToLower(tag)]
	if !ok {
		return false, nil
	}

	return true, matchTag
}
