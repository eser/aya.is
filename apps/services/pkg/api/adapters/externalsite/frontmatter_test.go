package externalsite //nolint:testpackage // white-box test uses unexported helpers

import (
	"testing"
	"time"
)

func TestParseMarkdownFile_TOMLFrontmatter(t *testing.T) {
	t.Parallel()

	content := `+++
title = "FP Roadmap"
date = 2024-03-15T10:00:00Z
slug = "fp-roadmap"
description = "A functional programming roadmap"

[taxonomies]
tags = ["functional", "programming"]
+++

# Hello World

This is the body.`

	frontMatter, body, err := ParseMarkdownFile(content, "content/posts/fp-roadmap.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frontMatter.Title != "FP Roadmap" {
		t.Errorf("expected title %q, got %q", "FP Roadmap", frontMatter.Title)
	}

	if frontMatter.Slug != "fp-roadmap" {
		t.Errorf("expected slug %q, got %q", "fp-roadmap", frontMatter.Slug)
	}

	if frontMatter.Description != "A functional programming roadmap" {
		t.Errorf(
			"expected description %q, got %q",
			"A functional programming roadmap",
			frontMatter.Description,
		)
	}

	if frontMatter.Date == nil {
		t.Fatal("expected date to be set")
	}

	expectedDate := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	if !frontMatter.Date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, *frontMatter.Date)
	}

	if body != "# Hello World\n\nThis is the body." {
		t.Errorf("unexpected body: %q", body)
	}

	// Check extra fields (taxonomies)
	if frontMatter.Extra["taxonomies"] == nil {
		t.Error("expected taxonomies in extra")
	}
}

//nolint:cyclop // thorough frontmatter field assertions
func TestParseMarkdownFile_YAMLFrontmatter(t *testing.T) {
	t.Parallel()

	content := `---
title: "My Hugo Post"
date: "2024-06-01"
slug: "my-hugo-post"
tags:
  - go
  - web
description: "A Hugo blog post"
---

Some markdown content here.`

	frontMatter, body, err := ParseMarkdownFile(content, "content/posts/my-hugo-post.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frontMatter.Title != "My Hugo Post" {
		t.Errorf("expected title %q, got %q", "My Hugo Post", frontMatter.Title)
	}

	if frontMatter.Slug != "my-hugo-post" {
		t.Errorf("expected slug %q, got %q", "my-hugo-post", frontMatter.Slug)
	}

	if frontMatter.Description != "A Hugo blog post" {
		t.Errorf("expected description %q, got %q", "A Hugo blog post", frontMatter.Description)
	}

	expectedDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if frontMatter.Date == nil || !frontMatter.Date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, frontMatter.Date)
	}

	if len(frontMatter.Tags) != 2 || frontMatter.Tags[0] != "go" || frontMatter.Tags[1] != "web" {
		t.Errorf("expected tags [go, web], got %v", frontMatter.Tags)
	}

	if body != "Some markdown content here." {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestParseMarkdownFile_NoFrontmatter(t *testing.T) {
	t.Parallel()

	content := `# Just a markdown file

No frontmatter here.`

	frontMatter, body, err := ParseMarkdownFile(content, "content/posts/plain.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frontMatter.Slug != "plain" {
		t.Errorf("expected slug %q, got %q", "plain", frontMatter.Slug)
	}

	if body != content {
		t.Errorf("expected full content as body")
	}
}

func TestParseMarkdownFile_LanguageFromPath(t *testing.T) {
	t.Parallel()

	content := `---
title: "Turkish Post"
---

İçerik burada.`

	frontMatter, _, err := ParseMarkdownFile(content, "content/posts/my-post.tr.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frontMatter.Language != "tr" {
		t.Errorf("expected language %q, got %q", "tr", frontMatter.Language)
	}

	if frontMatter.Slug != "my-post.tr" {
		t.Errorf("expected slug %q, got %q", "my-post.tr", frontMatter.Slug)
	}
}

func TestParseMarkdownFile_LanguageInFrontmatter(t *testing.T) {
	t.Parallel()

	content := `+++
title = "Explicit Language"
language = "de"
+++

Inhalt hier.`

	frontMatter, _, err := ParseMarkdownFile(content, "content/posts/explicit-lang.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frontMatter.Language != "de" {
		t.Errorf("expected language %q, got %q", "de", frontMatter.Language)
	}
}

func TestSlugFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want string
	}{
		{"content/posts/fp-roadmap.md", "fp-roadmap"},
		{"_posts/2024-03-15-my-post.md", "2024-03-15-my-post"},
		{"post.md", "post"},
		{"content/talks/my-talk.tr.md", "my-talk.tr"},
		// Hugo page bundles — index.md should return parent directory name
		{"content/posts/2024-01-02-hello-hugo/index.md", "2024-01-02-hello-hugo"},
		{"content/posts/my-post/index.md", "my-post"},
		// Other markdown extensions
		{"content/posts/next-post.mdx", "next-post"},
		{"content/posts/markdoc-post.mdoc", "markdoc-post"},
		{"content/posts/long-ext.markdown", "long-ext"},
		// Page bundle with .mdx
		{"content/posts/mdx-bundle/index.mdx", "mdx-bundle"},
	}

	for _, testCase := range tests {
		got := slugFromPath(testCase.path)
		if got != testCase.want {
			t.Errorf("slugFromPath(%q) = %q, want %q", testCase.path, got, testCase.want)
		}
	}
}

func TestFilterMarkdownFiles(t *testing.T) {
	t.Parallel()

	tree := []treeEntry{
		{Path: "content", Type: "tree"},
		{Path: "content/posts", Type: "tree"},
		{Path: "content/posts/plain-post.md", Type: "blob"},
		{Path: "content/posts/mdx-post.mdx", Type: "blob"},
		{Path: "content/posts/markdoc-post.mdoc", Type: "blob"},
		{Path: "content/posts/long-ext.markdown", Type: "blob"},
		{Path: "content/posts/page-bundle/index.md", Type: "blob"},
		{Path: "content/posts/mdx-bundle/index.mdx", Type: "blob"},
		// Should be skipped:
		{Path: "content/_index.md", Type: "blob"},
		{Path: "content/posts/_index.md", Type: "blob"},
		{Path: "content/index.md", Type: "blob"},
		{Path: "README.md", Type: "blob"},
		{Path: "content/README.markdown", Type: "blob"},
		{Path: "content/posts/image.png", Type: "blob"},
	}

	files := filterMarkdownFiles(tree)

	expected := map[string]bool{
		"content/posts/plain-post.md":        true,
		"content/posts/mdx-post.mdx":         true,
		"content/posts/markdoc-post.mdoc":    true,
		"content/posts/long-ext.markdown":    true,
		"content/posts/page-bundle/index.md": true,
		"content/posts/mdx-bundle/index.mdx": true,
	}

	if len(files) != len(expected) {
		t.Errorf("expected %d files, got %d: %v", len(expected), len(files), files)
	}

	for _, f := range files {
		if !expected[f] {
			t.Errorf("unexpected file in result: %q", f)
		}
	}
}

func TestIsMarkdownFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want bool
	}{
		{"post.md", true},
		{"post.mdx", true},
		{"post.mdoc", true},
		{"post.markdown", true},
		{"POST.MD", true},
		{"post.MDX", true},
		{"image.png", false},
		{"data.json", false},
		{"script.js", false},
		{"noext", false},
	}

	for _, testCase := range tests {
		got := isMarkdownFile(testCase.path)
		if got != testCase.want {
			t.Errorf("isMarkdownFile(%q) = %v, want %v", testCase.path, got, testCase.want)
		}
	}
}
