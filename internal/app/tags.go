package app

import (
	"errors"
	"net/url"
	"strings"
)

var ErrNoTags = errors.New("no tags provided")

type Tags struct {
	tags []string
}

func ParseTags(args []string) (*Tags, error) {
	if len(args) == 0 {
		return nil, ErrNoTags
	}

	return &Tags{
		tags: args,
	}, nil
}

func (t Tags) String() string {
	return url.QueryEscape(strings.Join(t.tags, " and "))
}
