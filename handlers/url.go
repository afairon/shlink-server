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

	// Client POST request
	req := models.URL{}

	// Decode req json
	// TODO: Handle error
	if err := json.NewDecoder().DecodeReader(r.Body, &req); err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify URL
	// TODO: Handle error
	if ok, err := utils.IsURL(req.LongURL); !ok || err != nil {
		if err != nil {
			utils.Error(r, err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify if URL is on the blacklist
	// TODO: Handle error
	if ban, err := utils.IsBlackList(req.LongURL); ban || err != nil {
		if err != nil {
			utils.Error(r, err)
		}
		w.WriteHeader(http.StatusBadRequest)
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
		resp.Success = true
		resp.TargetURL = utils.Conf.Server.Base + resp.ID
		js, _ := json.Marshal(&resp)
		w.Write(js)

		return
	}

	counter, _ := models.FindAndModify(bson.M{"$inc": bson.M{"sequence": 1}})

	req.ReadyToInsert(hash, &counter)

	// TODO: Handle error
	err := models.InsertURL(req)
	if err != nil {
		utils.Error(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.TargetURL = utils.Conf.Server.Base + resp.ID
	req.Success = true

	js, _ := json.Marshal(&req)

	w.Write(js)
}

// InfoURL handles responding information about url.
func InfoURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := models.InfoURL(id)

	// Err or document not found
	if err != nil {
		// TODO: Handle error
		utils.Error(r, err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if len(resp) < 1 {
		// TODO: Handle error
		w.WriteHeader(http.StatusNotFound)

		return
	}

	js, err := json.Marshal(&resp[0])

	w.Write(js)
}

// RedirectURL redirects client to the target url.
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := models.FindURL(bson.M{"id": id})

	// Err or document not found
	if err != nil {
		// TODO: Handle error
		utils.Error(r, err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	http.Redirect(w, r, resp.LongURL, 301)
	models.UpdateStats(id)
}
