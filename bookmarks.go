package linkding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ListBookmarksParams defines the parameters used when listing bookmarks.
type ListBookmarksParams struct {
	// The search query to filter bookmarks.
	Query string
	// The maximum number of bookmarks to return.
	Limit int
	// The offset for pagination.
	Offset int
	// Filter to include only unread bookmarks.
	Unread bool
	// Search for bookmarks added after this date
	AddedSince time.Time
	// Search for bookmarks modified after this date
	ModifiedSince time.Time
	// Sort order of results: added_asc, added_desc, title_asc, title_desc
	Sort string
}

// ListBookmarksResponse represents the response from the Linkding API when
// listing bookmarks.
type ListBookmarksResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Bookmark `json:"results"`
}

// Bookmark represents a bookmark object in the Linkding API.
type Bookmark struct {
	ID                    int       `json:"id"`
	URL                   string    `json:"url"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	Notes                 string    `json:"notes"`
	WebsiteTitle          string    `json:"website_title"`
	WebsiteDescription    string    `json:"website_description"`
	WebArchiveSnapshotURL string    `json:"web_archive_snapshot_url"`
	FaviconURL            string    `json:"favicon_url"`
	PreviewImageURL       string    `json:"preview_image_url"`
	IsArchived            bool      `json:"is_archived"`
	Unread                bool      `json:"unread"`
	Shared                bool      `json:"shared"`
	TagNames              []string  `json:"tag_names"`
	DateAdded             time.Time `json:"date_added"`
	DateModified          time.Time `json:"date_modified"`
}

// CreateBookmarkRequest represents the request body when creating or updating
// bookmarks.
type CreateBookmarkRequest struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
	IsArchived  bool     `json:"is_archived"`
	Unread      bool     `json:"unread"`
	Shared      bool     `json:"shared"`
	TagNames    []string `json:"tag_names"`
}

// CheckBookmarkResponse represents the response from the Linkding API when
// checking a if a URL has been bookmarked.
//
// Warning: The Bookmark field will be nil if a URL has not been bokmarked.
type CheckBookmarkResponse struct {
	Bookmark *Bookmark `json:"bookmark"`
	Metadata Metadata  `json:"metadata"`
	AutoTags []string  `json:"auto_tags"`
}

// Metadata contains metadata scraped from a website.
type Metadata struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PreviewImage string `json:"preview_image"`
}

// ListBookmarks retrieves a list of bookmarks from Linkding based on the
// provided parameters.
func (c *Client) ListBookmarks(params ListBookmarksParams) (*ListBookmarksResponse, error) {
	path := buildBookmarksQueryString("/api/bookmarks/", params)

	body, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	result := &ListBookmarksResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListArchivedBookmarks retrieves a list of archived bookmarks from Linkding.
// It also filters the list based on the provided parameters.
func (c *Client) ListArchivedBookmarks(params ListBookmarksParams) (*ListBookmarksResponse, error) {
	path := buildBookmarksQueryString("/api/bookmarks/archived/", params)

	body, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	result := &ListBookmarksResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetBookmark retrieves a single bookmark from Linkding.
func (c *Client) GetBookmark(id int) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/api/bookmarks/%d/", id), nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

// CheckBookmark checks if a URL is already bookmarked.
func (c *Client) CheckBookmark(bookmarkUrl string) (*CheckBookmarkResponse, error) {
	uri, err := url.Parse(bookmarkUrl)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("url", uri.String())

	body, err := c.makeRequest(
		http.MethodGet,
		fmt.Sprintf("/api/bookmarks/check/?%s", query.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	result := &CheckBookmarkResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateBookmark creates a new bookmark in Linkding using the provided payload.
//
// Warning: Ensure that the TagNames property in the CreateBookmarkRequest is
// initialized (even if empty) to avoid nil pointer issues.
func (c *Client) CreateBookmark(payload CreateBookmarkRequest) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodPost, "/api/bookmarks/", payload)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

// UpdateBookmark updates an existing bookmark in Linkding using the provided
// payload.
//
// Warning: Ensure that the TagNames property in the CreateBookmarkRequest is
// initialized (even if empty) to avoid nil pointer issues.
func (c *Client) UpdateBookmark(id int, payload CreateBookmarkRequest) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodPut, fmt.Sprintf("/api/bookmarks/%d/", id), payload)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

// ArchiveBookmark archives a bookmark from Linkding.
func (c *Client) ArchiveBookmark(id int) error {
	_, err := c.makeRequest(http.MethodPost, fmt.Sprintf("/api/bookmarks/%d/archive/", id), nil)

	return err
}

// UnarchiveBookmark unarchives a bookmark from Linkding.
func (c *Client) UnarchiveBookmark(id int) error {
	_, err := c.makeRequest(http.MethodPost, fmt.Sprintf("/api/bookmarks/%d/unarchive/", id), nil)

	return err
}

// DeleteBookmark deletes a bookmark from Linkding.
func (c *Client) DeleteBookmark(id int) error {
	_, err := c.makeRequest(http.MethodDelete, fmt.Sprintf("/api/bookmarks/%d/", id), nil)

	return err
}

func buildBookmarksQueryString(path string, params ListBookmarksParams) string {
	values := url.Values{}

	if params.Query != "" {
		values.Set("q", params.Query)
	}

	if params.Limit > 0 {
		values.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		values.Set("offset", strconv.Itoa(params.Offset))
	}

	if params.Unread {
		values.Set("unread", "yes")
	}

	if !params.AddedSince.IsZero() {
		values.Set("added_since", params.AddedSince.Format(time.RFC3339))
	}

	if !params.ModifiedSince.IsZero() {
		values.Set("modified_since", params.AddedSince.Format(time.RFC3339))
	}

	if params.Sort != "" {
		values.Set("sort", params.Sort)
	}

	if len(values) > 0 {
		return fmt.Sprintf("%s?%s", path, values.Encode())
	}

	return path
}
