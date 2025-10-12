package sefaria

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ryanfaerman/go-sefaria/types"
)

type TextService service

type Text struct {
	Ref                         string    `json:"ref"`
	HeRef                       string    `json:"heRef"`
	IsComplex                   bool      `json:"isComplex"`
	Text                        []string  `json:"text"`
	He                          []string  `json:"he"`
	Versions                    []Version `json:"versions"`
	TextDepth                   int       `json:"textDepth"`
	SectionNames                []string  `json:"sectionNames"`
	AddressTypes                []string  `json:"addressTypes"`
	Lengths                     []int     `json:"lengths"`
	Length                      int       `json:"length"`
	HeTitle                     string    `json:"heTitle"`
	TitleVariants               []string  `json:"titleVariants"`
	HeTitleVariants             []string  `json:"heTitleVariants"`
	Type                        string    `json:"type"`
	PrimaryCategory             string    `json:"primary_category"`
	Book                        string    `json:"book"`
	Categories                  []string  `json:"categories"`
	Order                       []int     `json:"order"`
	Sections                    []any     `json:"sections"`
	ToSections                  []any     `json:"toSections"`
	Commentator                 string    `json:"commentator"`
	HeCommentator               string    `json:"heCommentator"`
	CollectiveTitle             string    `json:"collectiveTitle"`
	HeCollectiveTitle           string    `json:"heCollectiveTitle"`
	BaseTexTitles               []string  `json:"baseTexTitles"`
	IsDependant                 bool      `json:"isDependant"`
	IndexTitle                  string    `json:"indexTitle"`
	HeIndexTitle                string    `json:"heIndexTitle"`
	SectionRef                  string    `json:"sectionRef"`
	FirstAvailableSectionRef    string    `json:"firstAvailableSectionRef"`
	HeSectionRef                string    `json:"heSectionRef"`
	IsSpanning                  bool      `json:"isSpanning"`
	HeVersionTitle              string    `json:"heVersionTitle"`
	HeVersionTitleInHebrew      string    `json:"heVersionTitleInHebrew"`
	HeShortVersionTitle         string    `json:"heShortVersionTitle"`
	HeShortVersionTitleInHebrew string    `json:"heShortVersionTitleInHebrew"`
	HeVersionSource             string    `json:"heVersionSource"`
	HeVersionStatus             string    `json:"heVersionStatus"`
	HeVersionNotes              string    `json:"heVersionNotes"`
	HeExtendedNotes             string    `json:"heExtendedNotes"`
	HeExtendedNotesHebrew       string    `json:"heExtendedNotesHebrew"`
	HeVersionNotesInHebrew      string    `json:"heVersionNotesInHebrew"`
	HeDigitizedBySefaria        bool      `json:"heDigitizedBySefaria"`
	HeLicense                   string    `json:"heLicense"`
	FormatHeAsPoetry            bool      `json:"formatHeAsPoetry"`
	Title                       string    `json:"title"`
	HeBook                      string    `json:"heBook"`
	Alts                        []any     `json:"alts"`
	IndexOffsetsByDepth         struct{}  `json:"index_offsets_by_depth"`
	Next                        string    `json:"next"`
	Prev                        string    `json:"prev"`
	Commentary                  []any     `json:"commentary"`
	Sheets                      []any     `json:"sheets"`
	Layer                       []any     `json:"layer"`
}

type Version struct {
	Title                     string                  `json:"title"`
	VersionTitle              string                  `json:"versionTitle"`
	VersionSource             string                  `json:"versionSource"`
	Language                  string                  `json:"language"`
	Status                    string                  `json:"status"`
	License                   string                  `json:"license"`
	VersionNotes              string                  `json:"versionNotes"`
	DigitizedBySefaria        types.StringOr[bool]    `json:"digitizedBySefaria"`
	Priority                  types.StringOr[float32] `json:"priority"`
	VersionTitleInHebrew      string                  `json:"versionTitleInHebrew"`
	VersionNotesInHebrew      string                  `json:"versionNotesInHebrew"`
	ExtendedNotes             string                  `json:"extendedNotes"`
	ExtendedNotesHebrew       string                  `json:"extendedNotesHebrew"`
	PurchaseInformationImage  string                  `json:"purchaseInformationImage"`
	PurchaseInformationURL    string                  `json:"purchaseInformationURL"`
	ShortVersionTitle         string                  `json:"shortVersionTitle"`
	ShortVersionTitleInHebrew string                  `json:"shortVersionTitleInHebrew"`
	FirstSectionRef           string                  `json:"firstSectionRef"`

	FormatAsPoetry         string `json:"formatAsPoetry"`
	Method                 string `json:"method"`
	HeversionSource        string `json:"heversionSource"`
	VersionURL             string `json:"versionUrl"`
	HasManuallyWrappedRefs string `json:"hasManuallyWrappedRefs"`
	ActualLanguage         string `json:"actualLanguage"`
	LanguageFamilyName     string `json:"languageFamilyName"`
	IsSource               bool   `json:"isSource"`
	IsPrimary              bool   `json:"isPrimary"`
	Direction              string `json:"direction"`
}

type TextFormat string

const (
	FormatDefault            TextFormat = "default"
	FormatTextOnly           TextFormat = "text_only"
	FormatStripOnlyFootnotes TextFormat = "strip_only_footnotes"
	FormatWrapAllEntities    TextFormat = "wrap_all_entities"
)

type TextVersion struct {
	Language string
	Title    string
}

func (t TextVersion) String() string {
	if t.Title != "" {
		return fmt.Sprintf("%s|%s", t.Language, t.Title)
	}
	return t.Language
}

type TextOptions struct {
	Versions            []TextVersion
	FillMissingSegments bool
	Format              TextFormat
}

func (t TextOptions) Query() string {
	v := url.Values{}

	for _, ver := range t.Versions {
		v.Add("version", ver.String())
	}

	if t.Format == "" {
		t.Format = FormatDefault
	}
	v.Set("return_format", string(t.Format))

	if t.FillMissingSegments {
		v.Set("fill_missing_segments", "1")
	} else {
		v.Set("fill_missing_segments", "0")
	}

	return v.Encode()
}

func (s *TextService) Get(ctx context.Context, tref string, opts *TextOptions) (*Text, error) {
	u := s.client.BaseURL.JoinPath("/v3/texts/", tref)
	u.RawQuery = opts.Query()

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	text := new(Text)
	_, err = s.client.Do(req, text)
	return text, err
}
