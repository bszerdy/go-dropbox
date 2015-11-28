package dropbox

import (
	"encoding/json"
	"io"
	"time"
)

// Files client for files and folders.
type Files struct {
	*Client
}

// NewFiles client.
func NewFiles(config *Config) *Files {
	return &Files{
		Client: &Client{
			Config: config,
		},
	}
}

// WriteMode determines what to do if the file already exists.
type WriteMode string

// Supported write modes.
const (
	WriteModeAdd       WriteMode = "add"
	WriteModeOverwrite           = "overwrite"
)

// Metadata for a file or folder.
type Metadata struct {
	Tag            string    `json:".tag"`
	Name           string    `json:"name"`
	PathLower      string    `json:"path_lower"`
	ClientModified time.Time `json:"client_modified"`
	ServerModified time.Time `json:"server_modified"`
	Rev            string    `json:"rev"`
	Size           uint64    `json:"size"`
	ID             string    `json:"id"`
}

// GetMetadataInput request input.
type GetMetadataInput struct {
	Path             string `json:"path"`
	IncludeMediaInfo bool   `json:"include_media_info"`
}

// GetMetadataOutput request output.
type GetMetadataOutput struct {
	Metadata
}

// GetMetadata returns the metadata for a file or folder.
func (c *Files) GetMetadata(in *GetMetadataInput) (out *GetMetadataOutput, err error) {
	body, err := c.call("/files/get_metadata", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// CreateFolderInput request input.
type CreateFolderInput struct {
	Path string `json:"path"`
}

// CreateFolderOutput request output.
type CreateFolderOutput struct {
	Name      string `json:"name"`
	PathLower string `json:"path_lower"`
	ID        string `json:"id"`
}

// CreateFolder creates a folder.
func (c *Files) CreateFolder(in *CreateFolderInput) (out *CreateFolderOutput, err error) {
	body, err := c.call("/files/create_folder", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// DeleteInput request input.
type DeleteInput struct {
	Path string `json:"path"`
}

// DeleteOutput request output.
type DeleteOutput struct {
	Metadata
}

// Delete a file or folder and its contents.
func (c *Files) Delete(in *DeleteInput) (out *DeleteOutput, err error) {
	body, err := c.call("/files/delete", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// CopyInput request input.
type CopyInput struct {
	FromPath string `json:"from_path"`
	ToPath   string `json:"to_path"`
}

// CopyOutput request output.
type CopyOutput struct {
	Metadata
}

// Copy a file or folder to a different location.
func (c *Files) Copy(in *CopyInput) (out *CopyOutput, err error) {
	body, err := c.call("/files/copy", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// MoveInput request input.
type MoveInput struct {
	FromPath string `json:"from_path"`
	ToPath   string `json:"to_path"`
}

// MoveOutput request output.
type MoveOutput struct {
	Metadata
}

// Move a file or folder to a different location.
func (c *Files) Move(in *MoveInput) (out *MoveOutput, err error) {
	body, err := c.call("/files/move", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// RestoreInput request input.
type RestoreInput struct {
	Path string `json:"path"`
	Rev  string `json:"rev"`
}

// RestoreOutput request output.
type RestoreOutput struct {
	Metadata
}

// Restore a file to a specific revision.
func (c *Files) Restore(in *RestoreInput) (out *RestoreOutput, err error) {
	body, err := c.call("/files/restore", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// ListFolderInput request input.
type ListFolderInput struct {
	Path             string `json:"path"`
	Recursive        bool   `json:"recursive"`
	IncludeMediaInfo bool   `json:"include_media_info"`
	IncludeDeleted   bool   `json:"include_deleted"`
}

// ListFolderOutput request output.
type ListFolderOutput struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Entries []*Metadata
}

// ListFolder returns the metadata for a file or folder.
func (c *Files) ListFolder(in *ListFolderInput) (out *ListFolderOutput, err error) {
	body, err := c.call("/files/list_folder", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// SearchMode determines how a search is performed.
type SearchMode string

// Supported search modes.
const (
	SearchModeFilename           SearchMode = "filename"
	SearchModeFilenameAndContent            = "filename_and_content"
	SearchModeDeletedFilename               = "deleted_filename"
)

// SearchMatch represents the type of match made.
type SearchMatchType string

// Supported search match types.
const (
	SearchMatchFilename SearchMatchType = "filename"
	SearchMatchContent                  = "content"
	SearchMatchBoth                     = "both"
)

// SearchMatch represents a matched file or folder.
type SearchMatch struct {
	MatchType struct {
		Tag SearchMatchType `json:".tag"`
	} `json:"match_type"`
	Metadata *Metadata `json:"metadata"`
}

// SearchInput request input.
type SearchInput struct {
	Path       string     `json:"path"`
	Query      string     `json:"query"`
	Start      uint64     `json:"start,omitempty"`
	MaxResults uint64     `json:"max_results,omitempty"`
	Mode       SearchMode `json:"mode"`
}

// SearchOutput request output.
type SearchOutput struct {
	Matches []*SearchMatch `json:"matches"`
	More    bool           `json:"more"`
	Start   uint64         `json:"start"`
}

// Search for files and folders.
func (c *Files) Search(in *SearchInput) (out *SearchOutput, err error) {
	if in.Mode == "" {
		in.Mode = SearchModeFilename
	}

	body, err := c.call("/files/search", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// UploadInput request input.
type UploadInput struct {
	Path           string    `json:"path"`
	Mode           WriteMode `json:"mode"`
	AutoRename     bool      `json:"autorename"`
	Mute           bool      `json:"mute"`
	ClientModified time.Time `json:"client_modified,omitempty"`
	Reader         io.Reader `json:"-"`
}

// UploadOutput request output.
type UploadOutput struct {
	Metadata
}

// Upload a file smaller than 150MB.
func (c *Files) Upload(in *UploadInput) (out *UploadOutput, err error) {
	body, err := c.download("/files/upload", in, in.Reader)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// DownloadInput request input.
type DownloadInput struct {
	Path string `json:"path"`
}

// DownloadOutput request input.
type DownloadOutput struct {
	Body io.ReadCloser
}

// Download a file.
func (c *Files) Download(in *DownloadInput) (out *DownloadOutput, err error) {
	body, err := c.download("/files/download", in, nil)
	if err != nil {
		return
	}

	return &DownloadOutput{body}, nil
}
