package externalsite_test

import (
	"testing"

	"github.com/eser/aya.is/services/pkg/api/adapters/externalsite"
)

const (
	testOwnerRepo = "Ardakilic/arda.pw"
	testBranch    = "master"
)

func TestSanitizeContent_Tweet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "tweet with user and id",
			input: `Check this: {{< tweet user="taylorotwell" id="414503338026102784" >}}`,
			want:  `Check this: https://twitter.com/taylorotwell/status/414503338026102784`,
		},
		{
			name:  "tweet with user and unquoted id",
			input: `{{< tweet user="golang" id=123456789 >}}`,
			want:  `https://twitter.com/golang/status/123456789`,
		},
		{
			name:  "twitter alias",
			input: `{{% twitter user="foo" id="999" %}}`,
			want:  `https://twitter.com/foo/status/999`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := externalsite.SanitizeContent(testCase.input)
			if got != testCase.want {
				t.Errorf("got %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestSanitizeContent_YouTube(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "youtube bare id",
			input: `{{< youtube dQw4w9WgXcQ >}}`,
			want:  `https://www.youtube.com/watch?v=dQw4w9WgXcQ`,
		},
		{
			name:  "youtube quoted id",
			input: `{{< youtube id="dQw4w9WgXcQ" >}}`,
			want:  `https://www.youtube.com/watch?v=dQw4w9WgXcQ`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := externalsite.SanitizeContent(testCase.input)
			if got != testCase.want {
				t.Errorf("got %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestSanitizeContent_Vimeo(t *testing.T) {
	t.Parallel()

	input := `{{< vimeo 146022717 >}}`
	want := `https://vimeo.com/146022717`

	got := externalsite.SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_Gist(t *testing.T) {
	t.Parallel()

	input := `{{< gist spf13 7896402 >}}`
	want := `https://gist.github.com/spf13/7896402`

	got := externalsite.SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_Figure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "figure with src and alt",
			input: `{{< figure src="/images/photo.jpg" alt="A photo" >}}`,
			want:  `![A photo](/images/photo.jpg)`,
		},
		{
			name:  "figure with src and caption",
			input: `{{< figure src="/img/pic.png" caption="My caption" >}}`,
			want:  `![My caption](/img/pic.png)`,
		},
		{
			name:  "figure src only",
			input: `{{< figure src="photo.jpg" >}}`,
			want:  `![](photo.jpg)`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := externalsite.SanitizeContent(testCase.input)
			if got != testCase.want {
				t.Errorf("got %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestSanitizeContent_UnknownShortcode(t *testing.T) {
	t.Parallel()

	input := "Before\n\n{{< custom_widget foo=\"bar\" >}}\n\nAfter"
	want := "Before\n\nAfter"

	got := externalsite.SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_ZolaShortcode(t *testing.T) {
	t.Parallel()

	input := "Before\n\n{{ youtube(id=\"dQw4w9WgXcQ\") }}\n\nAfter"
	want := "Before\n\nAfter"

	got := externalsite.SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_MixedContent(t *testing.T) {
	t.Parallel()

	input := `# Hello

Some text here.

{{< tweet user="taylorotwell" id="414503338026102784" >}}

More text.

{{< youtube dQw4w9WgXcQ >}}

Final paragraph.`

	want := `# Hello

Some text here.

https://twitter.com/taylorotwell/status/414503338026102784

More text.

https://www.youtube.com/watch?v=dQw4w9WgXcQ

Final paragraph.`

	got := externalsite.SanitizeContent(input)
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestResolveRelativeImages(t *testing.T) {
	t.Parallel()

	ownerRepo := testOwnerRepo
	branch := testBranch
	filePath := "content/posts/" +
		"2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/index.md"

	rawBase := "https://raw.githubusercontent.com/Ardakilic/arda.pw/master/" +
		"content/posts/2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "relative with dot-slash",
			input: `![laravel logo](./images/laravel-logo.png)`,
			want:  "![laravel logo](" + rawBase + "images/laravel-logo.png)",
		},
		{
			name:  "relative without dot-slash",
			input: `![](images/laravel-logo.png)`,
			want:  "![](" + rawBase + "images/laravel-logo.png)",
		},
		{
			name:  "parent directory reference",
			input: `![photo](../shared/photo.jpg)`,
			want: "![photo](https://raw.githubusercontent.com/" +
				"Ardakilic/arda.pw/master/content/posts/shared/photo.jpg)",
		},
		{
			name:  "absolute URL left unchanged",
			input: `![logo](https://example.com/logo.png)`,
			want:  `![logo](https://example.com/logo.png)`,
		},
		{
			name:  "protocol-relative URL left unchanged",
			input: `![logo](//cdn.example.com/logo.png)`,
			want:  `![logo](//cdn.example.com/logo.png)`,
		},
		{
			name:  "data URI left unchanged",
			input: `![pixel](data:image/png;base64,abc)`,
			want:  `![pixel](data:image/png;base64,abc)`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := externalsite.ResolveRelativeImages(
				testCase.input, ownerRepo, branch, filePath,
			)
			if got != testCase.want {
				t.Errorf("got:\n%s\n\nwant:\n%s", got, testCase.want)
			}
		})
	}
}

func TestResolveRelativeImages_HTMLImg(t *testing.T) {
	t.Parallel()

	ownerRepo := testOwnerRepo
	branch := testBranch
	filePath := "content/posts/" +
		"2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/index.md"

	rawBase := "https://raw.githubusercontent.com/Ardakilic/arda.pw/master/" +
		"content/posts/2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/"

	input := `<img src="./images/photo.jpg" alt="photo">`
	want := `<img src="` + rawBase + `images/photo.jpg" alt="photo">`

	got := externalsite.ResolveRelativeImages(input, ownerRepo, branch, filePath)
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestResolveRelativeImages_Mixed(t *testing.T) {
	t.Parallel()

	ownerRepo := testOwnerRepo
	branch := testBranch
	filePath := "content/posts/" +
		"2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/index.md"

	rawBase := "https://raw.githubusercontent.com/Ardakilic/arda.pw/master/" +
		"content/posts/2014-03-26-laravel-4-iliskiler-uzerinden-etkin-filtreleme-ve-gruplama/"

	input := "Some text\n\n![first](./images/one.png)\n\nMore text\n\n" +
		"![second](./images/two.jpg)\n\n![external](https://example.com/img.png)"
	want := "Some text\n\n![first](" + rawBase + "images/one.png)\n\nMore text\n\n" +
		"![second](" + rawBase + "images/two.jpg)\n\n" +
		"![external](https://example.com/img.png)"

	got := externalsite.ResolveRelativeImages(input, ownerRepo, branch, filePath)
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestResolveRelativeImages_StandaloneFile(t *testing.T) {
	t.Parallel()

	// Non-page-bundle file — images relative to the posts directory
	input := `![cover](./cover.png)`
	got := externalsite.ResolveRelativeImages(
		input,
		"user/repo",
		"main",
		"content/posts/my-post.md",
	)
	want := `![cover](https://raw.githubusercontent.com/user/repo/main/content/posts/cover.png)`

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
