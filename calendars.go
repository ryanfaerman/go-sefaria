package sefaria

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/google/go-querystring/query"
	"github.com/ryanfaerman/go-sefaria/bidi"
)

type CalendarService service

type CalendarGetOptions struct {
	Diaspora bool   `url:"diaspora,omitempty"`
	Custom   string `url:"custom,omitempty"`
	Year     int    `url:"year,omitempty"`
	Month    int    `url:"month,omitempty"`
	Day      int    `url:"day,omitempty"`
	TimeZone string `url:"timezone,omitempty"`
}

// TODO: replace this with a proper struct
type Calendar map[string]any

func (s *CalendarService) Get(ctx context.Context, opts *CalendarGetOptions) (*Calendar, error) {
	u := s.client.BaseURL.JoinPath("/calendars")
	if opts != nil {
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

	calendar := new(Calendar)
	_, err = s.client.Do(req, calendar)
	return calendar, err
}

var Parshiot = []string{
	"Bereshit",
	"Noach",
	"Lech-Lecha",
	"Vayera",
	"Chayei Sara",
	"Toldot",
	"Vayetzei",
	"Vayishlach",
	"Vayeshev",
	"Miketz",
	"Vayigash",
	"Vayechi",
	"Shemot",
	"Vaera",
	"Bo",
	"Beshalach",
	"Yitro",
	"Mishpatim",
	"Terumah",
	"Tetzaveh",
	"Ki Tisa",
	"Vayakhel",
	"Pekudei",
	"Vayikra",
	"Tzav",
	"Shmini",
	"Tazria",
	"Metzora",
	"Achrei Mot",
	"Kedoshim",
	"Emor",
	"Behar",
	"Bechukotai",
	"Bamidbar",
	"Nasso",
	"Beha’alotcha",
	"Sh’lach",
	"Korach",
	"Chukat",
	"Balak",
	"Pinchas",
	"Matot",
	"Masei",
	"Devarim",
	"Vaetchanan",
	"Eikev",
	"Re’eh",
	"Shoftim",
	"Ki Teitzei",
	"Ki Tavo",
	"Nitzavim",
	"Vayeilech",
	"Ha’azinu",
	"Vezot Haberakhah",
	"Vayakhel-Pekudei",
	"Tazria-Metzora",
	"Achrei Mot-Kedoshim",
	"Behar-Bechukotai",
	"Chukat-Balak",
	"Matot-Masei",
	"Nitzavim-Vayeilech",
}

type ParshaReading struct {
	Parsha     Parsha          `json:"parasha"`
	Haftorah   []Haftorah      `json:"haftarah" table:"Haftorah"`
	Date       Date            `json:"date"`
	HebrewDate BilingualString `json:"he_date"`
}

type Parsha struct {
	Title        BilingualString `json:"title"`
	DisplayValue BilingualString `json:"displayValue"`
	URL          string          `json:"url"`
	Ref          string          `json:"ref"`
	HeRef        bidi.String     `json:"heRef"`
	Order        int             `json:"order" table:"-"`
	Category     string          `json:"category"`
	ExtraDetails struct {
		Aliyot []string `json:"aliyot"`
	} `json:"extraDetails"`
	Description BilingualString `json:"description"`
}

type Haftorah struct {
	Title        BilingualString `json:"title" table:"Title"`
	DisplayValue BilingualString `json:"displayValue"`
	URL          string          `json:"url"`
	Ref          string          `json:"ref"`
	Order        int             `json:"order"`
	Category     string          `json:"category"`
}

var ErrInvalidParsha = errors.New("invalid parsha")

func (s *CalendarService) NextRead(ctx context.Context, parsha string) (*ParshaReading, error) {
	if !slices.Contains(Parshiot, parsha) {
		return nil, ErrInvalidParsha
	}
	u := s.client.BaseURL.JoinPath("/calendars/next-read", parsha)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(ParshaReading)
	_, err = s.client.Do(req, out)
	return out, err
}
