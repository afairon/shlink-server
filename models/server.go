package models

// Server information about version,
// go version and platform.
//go:generate $GOPATH/bin/ffjson $GOFILE
type Server struct {
	Version    string `json:"version"`
	Platform   string `json:"platform"`
	GoVersion  string `json:"go"`
	GoPlatform string `json:"goPlatform"`
}
