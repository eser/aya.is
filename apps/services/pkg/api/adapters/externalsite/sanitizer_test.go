package externalsite

import (
	"testing"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SanitizeContent(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SanitizeContent(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizeContent_Vimeo(t *testing.T) {
	t.Parallel()

	input := `{{< vimeo 146022717 >}}`
	want := `https://vimeo.com/146022717`

	got := SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_Gist(t *testing.T) {
	t.Parallel()

	input := `{{< gist spf13 7896402 >}}`
	want := `https://gist.github.com/spf13/7896402`

	got := SanitizeContent(input)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SanitizeContent(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizeContent_UnknownShortcode(t *testing.T) {
	t.Parallel()

	input := "Before\n\n{{< custom_widget foo=\"bar\" >}}\n\nAfter"
	want := "Before\n\nAfter"

	got := SanitizeContent(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSanitizeContent_ZolaShortcode(t *testing.T) {
	t.Parallel()

	input := "Before\n\n{{ youtube(id=\"dQw4w9WgXcQ\") }}\n\nAfter"
	want := "Before\n\nAfter"

	got := SanitizeContent(input)
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

	got := SanitizeContent(input)
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}
