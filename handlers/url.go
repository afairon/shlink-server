package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/globalsign/mgo/bson"
	json "github.com/pquerna/ffjson/ffjson"

	"github.com/go-chi/chi"

	"shlink-server/models"
	"shlink-server/utils"
)

// Generate handles the creation of shorten url.
func Generate(w http.ResponseWriter, r *http.Request) {

	var js []byte

	defer func() {
		w.Write(js)
	}()

	// Client POST request
	req := models.URL{}

	if err := json.NewDecoder().DecodeReader(r.Body, &req); err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusBadRequest)
		js, _ = json.Marshal(httpError(http.StatusBadRequest, "Couldn't decode json."))

		return
	}

	// Verify URL
	if ok, err := utils.IsURL(req.LongURL); !ok || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err != nil {
			utils.Error(r, err)
			js, _ = json.Marshal(httpError(http.StatusBadRequest, err.Error()))

			return
		}

		js, _ = json.Marshal(httpError(http.StatusBadRequest, fmt.Sprintf("%s is invalid.", req.LongURL)))

		return
	}

	// Verify if URL is on the blacklist
	if ban, err := utils.IsOnBlackList(req.LongURL); ban || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err != nil {
			utils.Error(r, err)
			js, _ = json.Marshal(httpError(http.StatusBadRequest, err.Error()))

			return
		}

		js, _ = json.Marshal(httpError(http.StatusBadRequest, fmt.Sprintf("%s is on the blacklist.", req.LongURL)))

		return
	}

	// Trim trailing slash
	if strings.HasSuffix(req.LongURL, "/") {
		req.LongURL = req.LongURL[0 : len(req.LongURL)-1]
	}

	// Reorder url query for consistency
	req.LongURL = utils.ReOrderQuery(req.LongURL)

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.LongURL)))

	resp, _ := models.FindURL(bson.M{"hash": hash})
	if resp.LongURL != "" {
		resp.TargetURL = utils.Conf.Server.Base + resp.ID
		js, _ = json.Marshal(&resp)

		return
	}

	counter, _ := models.FindAndModify(bson.M{"$inc": bson.M{"sequence": 1}})

	req.ReadyToInsert(hash, &counter)

	err := models.InsertURL(req)
	if err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		js, _ = json.Marshal(httpError(http.StatusInternalServerError, err.Error()))

		return
	}

	req.TargetURL = utils.Conf.Server.Base + req.ID

	js, _ = json.Marshal(&req)
}

// InfoURL handles responding information about url.
func InfoURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var js []byte

	defer func() {
		w.Write(js)
	}()

	resp, err := models.InfoURL(id)

	// Err or document not found
	if err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusNotFound)
		js, _ = json.Marshal(httpError(http.StatusNotFound, err.Error()))

		return
	}

	// Document empty
	if len(resp) < 1 {
		w.WriteHeader(http.StatusNotFound)
		js, _ = json.Marshal(httpError(http.StatusNotFound, fmt.Sprintf("No info found for %s.", id)))

		return
	}

	js, err = json.Marshal(&resp[0])
}

// RedirectURL redirects client to the target url.
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := models.FindURL(bson.M{"id": id})

	// Err or document not found
	if err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))

		return
	}

	http.Redirect(w, r, resp.LongURL, 301)
	models.UpdateStats(id)
}
