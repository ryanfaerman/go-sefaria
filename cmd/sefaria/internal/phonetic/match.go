package phonetic

import (
	"sort"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/tilotech/go-phonetics"
)

type Match struct {
	Candidate  string
	Score      int
	LetterDist int
}

// Matches finds the closest matches to input in a slice of strings
func Matches(candidates []string, input string) []Match {
	inputNorm := normalize(input)
	inputMeta := phonetics.EncodeMetaphone(inputNorm)

	threshold := int(float64(len(inputNorm)) * 0.4)

	matches := []Match{}

	for _, c := range candidates {
		norm := normalize(c)
		meta := phonetics.EncodeMetaphone(norm)

		letterDist := levenshtein.ComputeDistance(inputNorm, norm)
		if letterDist > threshold {
			continue
		}

		phoneticDist := levenshtein.ComputeDistance(inputMeta, meta)
		score := letterDist*2 + phoneticDist

		// Bonus for common prefix
		score -= commonPrefixLength(inputNorm, norm)

		matches = append(matches, Match{
			Candidate:  c,
			Score:      score,
			LetterDist: letterDist,
		})
	}

	// Sort by score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score < matches[j].Score
	})

	// Exact match takes priority
	if len(matches) > 0 && matches[0].LetterDist == 0 {
		return matches[:1]
	}

	if len(matches) > 3 {
		return matches[:3]
	}
	return matches
}

// normalize removes non-letters and converts to lowercase
func normalize(s string) string {
	var builder []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			builder = append(builder, unicode.ToLower(r))
		}
	}
	return string(builder)
}

// commonPrefixLength returns the number of letters that match at the start
func commonPrefixLength(a, b string) int {
	l := min(len(a), len(b))

	for i := range l {
		if a[i] != b[i] {
			return i
		}
	}
	return l
}
