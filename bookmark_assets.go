package linkding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ListBookmarkAssetsResponse represents the response from the Linkding API when
// listing bookmark assets.
type ListBookmarkAssetsResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []BookmarkAsset `json:"results"`
}

// BookmarkAsset represents a bookmark asset in the Linkding API.
type BookmarkAsset struct {
	ID          int       `json:"id"`
	Bookmark    int       `json:"bookmark"`
	AssetType   string    `json:"asset_type"`
	DateCreated time.Time `json:"date_created"`
	ContentType string    `json:"content_type"`
	DisplayName string    `json:"display_name"`
	Status      string    `json:"status"`
}

// ListBookmarkAssets retrieves a list assets for a specific bookmark.
func (c *Client) ListBookmarkAssets(bookmarkID int) (*ListBookmarkAssetsResponse, error) {
	body, err := c.makeRequest(
		http.MethodGet,
		fmt.Sprintf("/api/bookmarks/%d/assets/", bookmarkID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	result := &ListBookmarkAssetsResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetBookmarkAsset retrieves a single asset by ID for a specific bookmark.
func (c *Client) GetBookmarkAsset(bookmarkID int, id int) (*BookmarkAsset, error) {
	body, err := c.makeRequest(
		http.MethodGet,
		fmt.Sprintf("/api/bookmarks/%d/assets/%d/", bookmarkID, id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	bookmark := &BookmarkAsset{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

// TODO: Implement download and upload

// DeleteBookmarkAsset deletes an asset by ID for a specific bookmark.
func (c *Client) DeleteBookmarkAsset(bookmarkID int, id int) error {
	_, err := c.makeRequest(
		http.MethodDelete,
		fmt.Sprintf("/api/bookmarks/%d/assets/%d/", bookmarkID, id),
		nil,
	)

	return err
}
