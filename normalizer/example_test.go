package normalizer

import (
	"strings"
	"testing"
	"time"
)

// Example data structures that might be used in real applications
type Article struct {
	ID          int
	Title       string
	Content     string
	Author      string
	Tags        []string
	Metadata    map[string]string
	PublishedAt time.Time
	Comments    []Comment
}

type Comment struct {
	ID      int
	Author  string
	Content string
	Replies []string
}

type User struct {
	Name        string
	Email       string
	Bio         string
	Preferences map[string]string
}

func TestRealWorldExamples(t *testing.T) {
	t.Run("article normalization", func(t *testing.T) {
		article := &Article{
			ID:      1,
			Title:   "The " + string(rune(0x201C)) + "Future" + string(rune(0x201D)) + " of AI &amp; Machine Learning" + string(rune(0x2026)),
			Content: "This article discusses the " + string(rune(0x201C)) + "amazing" + string(rune(0x201D)) + " advances in AI &lt; 2024.",
			Author:  "Dr. Jane Smith",
			Tags:    []string{"AI", "Machine Learning", "Future"},
			Metadata: map[string]string{
				"category": "Technology &amp; Innovation",
				"summary":  "A " + string(rune(0x201C)) + "comprehensive" + string(rune(0x201D)) + " look at AI" + string(rune(0x2026)),
			},
			PublishedAt: time.Now(),
			Comments: []Comment{
				{
					ID:      1,
					Author:  "John Doe",
					Content: "Great article! " + string(rune(0x201C)) + "Very" + string(rune(0x201D)) + " informative" + string(rune(0x2026)),
					Replies: []string{"I agree!", "Thanks for sharing"},
				},
			},
		}

		// Apply all normalizers
		normalizers := []Normalizer{
			HTMLUnescape,
			UnicodeNFC,
			Punctuation,
		}
		Apply(article, normalizers...)

		// Verify normalization
		expectedTitle := "The \"Future\" of AI & Machine Learning..."
		if article.Title != expectedTitle {
			t.Errorf("Title normalization failed: got %q, want %q", article.Title, expectedTitle)
		}

		expectedContent := "This article discusses the \"amazing\" advances in AI < 2024."
		if article.Content != expectedContent {
			t.Errorf("Content normalization failed: got %q, want %q", article.Content, expectedContent)
		}

		expectedCategory := "Technology & Innovation"
		if article.Metadata["category"] != expectedCategory {
			t.Errorf("Metadata normalization failed: got %q, want %q", article.Metadata["category"], expectedCategory)
		}

		expectedComment := "Great article! \"Very\" informative..."
		if article.Comments[0].Content != expectedComment {
			t.Errorf("Comment normalization failed: got %q, want %q", article.Comments[0].Content, expectedComment)
		}
	})

	t.Run("user profile normalization", func(t *testing.T) {
		user := &User{
			Name:  "José María",
			Email: "jose@example.com",
			Bio:   "Software developer with " + string(rune(0x201C)) + "passion" + string(rune(0x201D)) + " for AI &amp; ML" + string(rune(0x2026)),
			Preferences: map[string]string{
				"theme":         "dark",
				"language":      "English &amp; Spanish",
				"notifications": "email &lt; push",
			},
		}

		normalizers := []Normalizer{HTMLUnescape, Punctuation}
		Apply(user, normalizers...)

		expectedBio := "Software developer with \"passion\" for AI & ML..."
		if user.Bio != expectedBio {
			t.Errorf("Bio normalization failed: got %q, want %q", user.Bio, expectedBio)
		}

		expectedLanguage := "English & Spanish"
		if user.Preferences["language"] != expectedLanguage {
			t.Errorf("Language preference normalization failed: got %q, want %q", user.Preferences["language"], expectedLanguage)
		}
	})

	t.Run("mixed data types", func(t *testing.T) {
		data := map[string]interface{}{
			"title":   "Article &amp; News",
			"content": "This is " + string(rune(0x201C)) + "smart" + string(rune(0x201D)) + " text" + string(rune(0x2026)),
			"tags":    []string{"tag1", "tag2 &amp; more"},
			"meta": map[string]string{
				"author": "John &amp; Jane",
				"date":   "2024-01-01",
			},
		}

		normalizers := []Normalizer{HTMLUnescape, Punctuation}
		Apply(&data, normalizers...)

		// Verify map values
		if data["title"] != "Article & News" {
			t.Errorf("Map title normalization failed")
		}
		if data["content"] != "This is \"smart\" text..." {
			t.Errorf("Map content normalization failed")
		}

		// Verify slice values
		tags := data["tags"].([]string)
		if tags[1] != "tag2 & more" {
			t.Errorf("Slice normalization failed")
		}

		// Verify nested map values
		meta := data["meta"].(map[string]string)
		if meta["author"] != "John & Jane" {
			t.Errorf("Nested map normalization failed")
		}
	})
}

func TestCustomNormalizerExample(t *testing.T) {
	// Define a custom normalizer
	TrimWhitespace := func(s string) string {
		return strings.TrimSpace(s)
	}

	// Define another custom normalizer
	ToLowercase := func(s string) string {
		return strings.ToLower(s)
	}

	type Product struct {
		Name        string
		Description string
		Category    string
	}

	product := &Product{
		Name:        "  PREMIUM WIDGET  ",
		Description: "  High-quality product with " + string(rune(0x201C)) + "excellent" + string(rune(0x201D)) + " features" + string(rune(0x2026)) + "  ",
		Category:    "  ELECTRONICS &amp; GADGETS  ",
	}

	// Apply custom normalizers along with built-in ones
	normalizers := []Normalizer{
		TrimWhitespace,
		ToLowercase,
		HTMLUnescape,
		Punctuation,
	}
	Apply(product, normalizers...)

	expectedName := "premium widget"
	if product.Name != expectedName {
		t.Errorf("Custom normalizer failed: got %q, want %q", product.Name, expectedName)
	}

	expectedDescription := "high-quality product with \"excellent\" features..."
	if product.Description != expectedDescription {
		t.Errorf("Custom normalizer failed: got %q, want %q", product.Description, expectedDescription)
	}

	expectedCategory := "electronics & gadgets"
	if product.Category != expectedCategory {
		t.Errorf("Custom normalizer failed: got %q, want %q", product.Category, expectedCategory)
	}
}

func TestPerformanceExample(t *testing.T) {
	// Test performance with a large dataset
	type Document struct {
		Title   string
		Content string
		Tags    []string
	}

	// Create a large number of documents
	documents := make([]Document, 1000)
	for i := range documents {
		documents[i] = Document{
			Title:   "Document &amp; " + string(rune('0'+i%10)),
			Content: "Content with " + string(rune(0x201C)) + "smart" + string(rune(0x201D)) + " quotes" + string(rune(0x2026)),
			Tags:    []string{"tag1", "tag2 &amp; more"},
		}
	}

	normalizers := []Normalizer{HTMLUnescape, Punctuation}

	// Apply normalization to all documents
	for i := range documents {
		Apply(&documents[i], normalizers...)
	}

	// Verify a few documents were normalized correctly
	if documents[0].Title != "Document & 0" {
		t.Errorf("Performance test normalization failed: got %q, want %q", documents[0].Title, "Document & 0")
	}
	if documents[0].Content != "Content with \"smart\" quotes..." {
		t.Errorf("Performance test normalization failed: got %q, want %q", documents[0].Content, "Content with \"smart\" quotes...")
	}
}

// ExampleApply demonstrates basic usage of the Apply function
func ExampleApply() {
	type Person struct {
		Name string
		Bio  string
	}

	p := &Person{
		Name: "John &amp; Jane",
		Bio:  "Hello " + string(rune(0x201C)) + "world" + string(rune(0x201D)) + " — nice to meet you" + string(rune(0x2026)),
	}

	normalizers := []Normalizer{
		HTMLUnescape,
		Punctuation,
	}
	Apply(p, normalizers...)

	// p.Name is now "John & Jane"
	// p.Bio is now "Hello "world" - nice to meet you..."
	_ = p
}

// ExampleApply_custom demonstrates using custom normalizers
func ExampleApply_custom() {
	// Define custom normalizers
	TrimSpace := func(s string) string {
		return strings.TrimSpace(s)
	}

	ToTitle := func(s string) string {
		return strings.Title(s)
	}

	type Article struct {
		Title   string
		Content string
	}

	article := &Article{
		Title:   "  hello world  ",
		Content: "  this is content  ",
	}

	normalizers := []Normalizer{
		TrimSpace,
		ToTitle,
	}
	Apply(article, normalizers...)

	// article.Title is now "Hello World"
	// article.Content is now "This Is Content"
	_ = article
}

// ExampleApply_complex demonstrates working with complex nested structures
func ExampleApply_complex() {
	type Comment struct {
		Author  string
		Content string
	}

	type Post struct {
		Title    string
		Content  string
		Comments []Comment
		Tags     map[string]string
	}

	post := &Post{
		Title:   "Post &amp; News",
		Content: "This is " + string(rune(0x201C)) + "smart" + string(rune(0x201D)) + " content" + string(rune(0x2026)),
		Comments: []Comment{
			{Author: "John", Content: "Great " + string(rune(0x201C)) + "post" + string(rune(0x201D)) + "!"},
			{Author: "Jane", Content: "I agree &lt; 3"},
		},
		Tags: map[string]string{
			"category": "Tech &amp; News",
			"status":   "Published",
		},
	}

	normalizers := []Normalizer{
		HTMLUnescape,
		Punctuation,
	}
	Apply(post, normalizers...)

	// All string fields in post and its nested structures are normalized
	_ = post
}
