package handlers

import (
	"net/http"

	"github.com/afairon/shlink-server/models"

	json "github.com/pquerna/ffjson/ffjson"
)

// Status handles displaying information about
// server version, go built version and platform.
func Status(version string, platform string, goVersion string, goPlatform string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(models.Server{
			Version:    version,
			Platform:   platform,
			GoVersion:  goVersion,
			GoPlatform: goPlatform,
		})

		w.Write(js)
	}
}
