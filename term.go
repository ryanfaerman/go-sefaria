package sefaria

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/ryanfaerman/go-sefaria/bidi"
)

type TermService service

type TermTitle struct {
	// The text of the title
	Text bidi.String `json:"text" table:"Title"`

	// The language of the title, either "en" or "he"
	Lang string `json:"lang" table:"Language"`

	// Whether or not the title is a primary title. For any given topics,
	// one should expect two titles with this field present and set to true,
	// an English and a Hebrew primary title. The English value with primary
	// set to true will match the string value of the primaryTitle
	// field on topic.
	Primary bool `json:"primary,omitempty" table:"Primary"`
}

type Term struct {
	// Name of the Term. Since a Term is a shared title node that
	// can be referenced and used by many different Index nodes, the
	// name field is critical as it contains the shared title.
	Name string `json:"name" table:"Name"`

	// Array of Alternative Titles for the Term in Hebrew and English.
	Titles []TermTitle `json:"titles"`

	// A shared scheme to for a group of terms.
	Scheme string `json:"scheme"`

	// Terms that share a scheme can be ordered within that scheme. So for
	// example, Parshiyot within the Parasha scheme can be ordered
	// as per the order of the Parshiyot.
	Order int `json:"order"`

	// A string representing a citation to a Jewish text. A valid Ref consists
	// of a title string followed optionally by a section string or a segment
	// string. A title string is any one of the known text titles or title
	// variants in the Sefaria Database.
	Ref string `json:"ref"`

	// The category of a specific term.
	Category string `json:"category"`
}

var ErrEmptyTerm = fmt.Errorf("term cannot be empty")

// Get the given term. A term is a shared title node. It can be referenced and used
// by many different Index nodes. Terms that use the same TermScheme can be ordered
// within that scheme. So for example, Parsha terms who all share the
// TermScheme of parsha, can be ordered within that scheme.
//
// Arguments:
//   - term: a valid english name of a sefaria term
//
// Examples of valid terms: Noah, HaChovel
func (s *TermService) Get(ctx context.Context, term string) (*Term, error) {
	if term == "" {
		return nil, ErrEmptyTerm
	}

	u := s.client.BaseURL.JoinPath("/terms", term)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Term)
	_, err = s.client.Do(req, out)
	return out, err
}

type TermNameOptions struct {
	// Number of results to return. 0 indicates no limit.
	Limit int `url:"limit,omitempty" validate:"gte=0" desc:"number of results to return. 0 = no limit."`

	// By default the Name API returns Refs, book titles, authors,
	// topics, and collections. If the type is set, the response will
	// only contain items of that type. Note: Topic includes authors,
	// topics and people without differentiation.
	//
	// If empty, the results will include all types.
	Type string `url:"type,omitempty" validate:"omitempty,oneof=ref Collection Topic TocCategory Term User" desc:"(ref, Collection, Topic, TocCategory, Term, User)"`
}

type TermCompletion struct {
	Title      string   `json:"title" table:"Title"`
	Key        string   `json:"key" table:"Key"`
	Type       string   `json:"type" table:"Type"`
	PictureURL string   `json:"pic,omitempty"`
	Primary    bool     `json:"is_primary" table:"Primary"`
	Order      int      `json:"order"`
	TopicPools []string `json:"topic_pools,omitempty"`
}

type TermCompletions struct {
	// The language of the completions, either "en" or "he"
	Lang string `json:"lang"`

	// The type of terms returned. Possible values are:
	// - ref
	// - Topic
	// - AuthorTopic
	// - PersonTopic
	// - User
	Type string `json:"type"`

	// A list of autocompletion responses for the submitted text, as the API returns it.
	CompletionTitles []string         `json:"completions"`
	Completions      []TermCompletion `json:"completion_objects"`

	// IsBook=true if the submitted text is a book level reference. e.g. (Genesis)
	IsBook bool `json:"is_book"`

	// IsSection=true if the submitted text is a section Ref (e.g. Genesis 4, as opposed
	// to a segment Ref such as Genesis 4.1).
	IsSection bool `json:"is_section"`

	// IsSegment=true if the submitted text is a segment level Ref (e.g. Genesis 43:3, as
	// opposed to a section Ref such as Genesis 43)
	IsSegment bool `json:"is_segment"`

	// IsRange=true if the submitted text is a a ranged Ref (one that spans multiple
	// sections or segments.) e.g. Genesis 4-5
	IsRange bool `json:"is_range"`

	// If type=ref, this returns the canonical ref for the submitted text.
	Ref string `json:"ref"`

	// If type=ref, this returns the URL path to link to the submitted text on Sefaria.org
	URL string `json:"url"`

	// If type=ref, this returns the canonical name of the index of the submitted text.
	Index string `json:"index"`

	// If the submitted response is a Ref, this returns the book it belongs to.
	Book string `json:"book"`

	InternalSections   []int `json:"internalSections"`
	InternalToSections []int `json:"internalToSections"`

	Sections   []string `json:"sections"`
	ToSections []string `json:"toSections"`
	Examples   []any    `json:"examples"`

	// Given a reference, this returns the names of the sections and segments at
	// each depth of that text.
	SectionNames       []string      `json:"sectionNames"`
	HebrewSectionNames []bidi.String `json:"heSectionNames"`

	// Given a partial Ref, this will return an array of strings of possible ways
	// that it might be completed.
	AddressExamples       []string      `json:"addressExamples"`
	HebrewAddressExamples []bidi.String `json:"heAddressExamples"`
}

// Name serves primarily as an autocomplete endpoint, returning potential keyword matches
// for Refs, book titles, authors, topics, and collections available on Sefaria.
//
// Arguments:
//   - Name: an arbitrary string to search for matches
//   - opts: optional parameters
func (s *TermService) Name(ctx context.Context, name string, opts *TermNameOptions) (*TermCompletions, error) {
	u := s.client.BaseURL.JoinPath("/name", name)
	if opts != nil {
		if err := s.client.validateStruct(opts); err != nil {
			return nil, err
		}

		v, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		u.RawQuery = v.Encode()
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(TermCompletions)
	_, err = s.client.Do(req, out)
	return out, err
}
