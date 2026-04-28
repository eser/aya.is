package devto

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "log/slog"
  "net/http"
  neturl "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/eser/aya.is/services/pkg/ajan/httpclient"
  "github.com/eser/aya.is/services/pkg/api/business/siteimporter"
)

var (
  ErrUserNotFound     = errors.New("dev.to user not found")
  ErrUnexpectedStatus = errors.New("unexpected status")
)

const (
  devtoAPIBase   = "https://dev.to/api"
  devtoMaxPerPage = 100
)

type devtoUser struct {
  Name     string `json:"name"`
  Username string `json:"username"`
}

type devtoArticle struct {
  ID           int     `json:"id"`
  Title        string  `json:"title"`
  Description  string  `json:"description"`
  PublishedAt  string  `json:"published_at"`
  CanonicalURL string  `json:"canonical_url"`
  URL          string  `json:"url"`
  CoverImage   string  `json:"cover_image"`
  SocialImage  string  `json:"social_image"`
  BodyMarkdown string  `json:"body_markdown"`
  TagList      tagList `json:"tag_list"`
  User         devtoUser `json:"user"`
}

type tagList []string

func (t *tagList) UnmarshalJSON(data []byte) error {
  if string(data) == "null" {
    return nil
  }

  var list []string
  if err := json.Unmarshal(data, &list); err == nil {
    *t = list
    return nil
  }

  var raw string
  if err := json.Unmarshal(data, &raw); err == nil {
    *t = splitTags(raw)
    return nil
  }

  return nil
}

// FetchAll fetches all articles for a Dev.to user.
func (p *Provider) FetchAll(
  ctx context.Context,
  username string,
) ([]*siteimporter.ImportItem, error) {
  p.logger.DebugContext(ctx, "Fetching Dev.to articles",
    slog.String("username", username))

  var items []*siteimporter.ImportItem

  for page := 1; ; page++ {
    articles, err := fetchArticlesPage(ctx, p.httpClient, username, page, devtoMaxPerPage)
    if err != nil {
      return nil, fmt.Errorf("failed to fetch Dev.to articles: %w", err)
    }

    if len(articles) == 0 {
      break
    }

    for _, article := range articles {
      // The list endpoint does not return the full body_markdown, so we need to
      // fetch the full article details.
      fullArticle, err := fetchArticleDetails(ctx, p.httpClient, article.ID)
      if err == nil && fullArticle != nil {
        article.BodyMarkdown = fullArticle.BodyMarkdown
      } else {
        p.logger.WarnContext(ctx, "Failed to fetch full article details, using summary",
          slog.Int("article_id", article.ID),
          slog.Any("error", err))
      }

      item := articleToImportItem(article)
      if item != nil {
        items = append(items, item)
      }
    }

    p.logger.DebugContext(ctx, "Fetched Dev.to page",
      slog.Int("page", page),
      slog.Int("items", len(articles)))
  }

  return items, nil
}

func fetchArticlesPage(
  ctx context.Context,
  client *httpclient.Client,
  username string,
  page int,
  perPage int,
) ([]devtoArticle, error) {
  if perPage <= 0 {
    perPage = devtoMaxPerPage
  }

  params := neturl.Values{}
  params.Set("username", username)
  params.Set("page", strconv.Itoa(page))
  params.Set("per_page", strconv.Itoa(perPage))

  url := devtoAPIBase + "/articles?" + params.Encode()

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Accept", "application/json")

  resp, err := client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch dev.to articles: %w", err)
  }

  defer func() { _ = resp.Body.Close() }()

  if resp.StatusCode == http.StatusNotFound {
    return nil, fmt.Errorf("%w: %s", ErrUserNotFound, username)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
  }

  var articles []devtoArticle

  err = json.NewDecoder(resp.Body).Decode(&articles)
  if err != nil {
    return nil, fmt.Errorf("failed to decode response: %w", err)
  }

  return articles, nil
}

func fetchArticleDetails(
  ctx context.Context,
  client *httpclient.Client,
  id int,
) (*devtoArticle, error) {
  url := fmt.Sprintf("%s/articles/%d", devtoAPIBase, id)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Accept", "application/json")

  resp, err := client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch dev.to article details: %w", err)
  }

  defer func() { _ = resp.Body.Close() }()

  if resp.StatusCode == http.StatusNotFound {
    return nil, fmt.Errorf("article not found: %d", id)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
  }

  var article devtoArticle

  err = json.NewDecoder(resp.Body).Decode(&article)
  if err != nil {
    return nil, fmt.Errorf("failed to decode response: %w", err)
  }

  return &article, nil
}

func articleToImportItem(article devtoArticle) *siteimporter.ImportItem {
  remoteID := ""
  if article.ID != 0 {
    remoteID = strconv.Itoa(article.ID)
  } else if article.URL != "" {
    remoteID = article.URL
  }

  if remoteID == "" {
    return nil
  }

  publishedAt := parsePublishedAt(article.PublishedAt)

  link := article.CanonicalURL
  if link == "" {
    link = article.URL
  }

  thumbnailURL := article.SocialImage
  if thumbnailURL == "" {
    thumbnailURL = article.CoverImage
  }

  props := make(map[string]any)

  if article.BodyMarkdown != "" {
    props["content"] = article.BodyMarkdown
  }

  if len(article.TagList) > 0 {
    props["tags"] = []string(article.TagList)
  }

  if article.URL != "" {
    props["devto_url"] = article.URL
  }

  if article.CanonicalURL != "" {
    props["canonical_url"] = article.CanonicalURL
  }

  return &siteimporter.ImportItem{
    RemoteID:     remoteID,
    Title:        article.Title,
    Description:  article.Description,
    PublishedAt:  publishedAt,
    Link:         link,
    ThumbnailURL: thumbnailURL,
    StoryKind:    "article",
    Properties:   props,
  }
}

func parsePublishedAt(value string) time.Time {
  if value == "" {
    return time.Now()
  }

  parsed, err := time.Parse(time.RFC3339, value)
  if err == nil {
    return parsed
  }

  return time.Now()
}

func splitTags(raw string) []string {
  raw = strings.TrimSpace(raw)
  if raw == "" {
    return nil
  }

  parts := strings.Split(raw, ",")
  tags := make([]string, 0, len(parts))

  for _, part := range parts {
    tag := strings.TrimSpace(part)
    if tag != "" {
      tags = append(tags, tag)
    }
  }

  return tags
}
