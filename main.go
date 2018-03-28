package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"short/models"
	"short/pkg/genid"
	utils "short/utils"

	"github.com/avct/uasurfer"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-zoo/bone"
	"github.com/pquerna/ffjson/ffjson"
	"go.uber.org/zap"
)

var logger *zap.Logger
var session *mgo.Session

var err error

var server = utils.Config.Server
var database = utils.Config.Database

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	/*t, err := template.New("index.html").Funcs(template.FuncMap{
		"gettext": func(input string) string {
			return gotext.Get(input)
		},
	}).ParseFiles("public/index.html", "public/head.html")
	if err != nil {
		logger.Error(err.Error())
	}

	err = t.Execute(w, nil)
	if err != nil {
		logger.Error(err.Error())
	}*/
	t, err := template.New("index").Parse(`
		<!DOCTYPE html>
<html lang="en">
<body>
    <h1>It's working! Yay!</h1>
</body>
</html>
	`)
	if err != nil {
		logger.Error(err.Error())
	}

	if err = t.Execute(w, nil); err != nil {
		logger.Error(err.Error())
	}
}

func redirectFull(w http.ResponseWriter, r *http.Request) {
	newSession := session.Copy()
	defer newSession.Close()

	id := bone.GetValue(r, "id")

	// User agent
	ua := uasurfer.Parse(r.Header.Get("User-Agent"))

	var result models.URL
	err := newSession.DB("short").C("UrlCollection").Find(bson.M{"id": id}).One(&result)
	if err != nil {
		logger.Error(err.Error())
	}

	if result.LongURL != "" {
		logger.Info("Access: "+result.LongURL, zap.String("os", ua.OS.Name.String()[2:]), zap.String("browser", ua.Browser.Name.String()[7:]), zap.String("remote-ip", r.RemoteAddr))
		http.Redirect(w, r, result.LongURL, 302)

		/*go func() {
			routineSession := session.Copy()
			defer routineSession.Close()

			_, err := routineSession.DB("short").C("UrlStatsCollection").Upsert(bson.M{"id": result.ID}, bson.M{"$set": bson.M{"id": result.ID}, "$inc": bson.M{"click": 1}})
			if err != nil {
				logger.Error(err.Error())
			}
		}()*/

		return
	}
	fmt.Fprintf(w, "Not Found")
}

func delete(w http.ResponseWriter, r *http.Request) {
	newSession := session.Copy()
	defer func() {
		newSession.Close()

		var urlCopy models.URL
		urlCopy.Success = true
		json, _ := ffjson.Marshal(&urlCopy)
		w.Write(json)
	}()

	id := bone.GetValue(r, "id")

	db := newSession.DB("short")

	var urlCopy models.URL

	err := db.C("UrlCollection").Remove(bson.M{"id": id})
	if err != nil {
		logger.Error(err.Error())

		urlCopy.Success = false
		urlCopy.Err = err.Error()
		json, _ := ffjson.Marshal(&urlCopy)
		w.WriteHeader(400)
		w.Write(json)
	}
	err = db.C("UrlStatsCollection").Remove(bson.M{"id": id})
	if err != nil {
		logger.Error(err.Error())

		urlCopy.Success = false
		urlCopy.Err = err.Error()
		json, _ := ffjson.Marshal(&urlCopy)
		w.WriteHeader(400)
		w.Write(json)
	}
}

// generate handles short url endpoint
func generate(w http.ResponseWriter, r *http.Request) {
	newSession := session.Copy()
	defer newSession.Close()

	w.Header().Set("Content-Type", "application/json")

	// Unmarshal JSON request
	var urlCopy models.URL
	if err := ffjson.NewDecoder().DecodeReader(r.Body, &urlCopy); err != nil {
		logger.Error(err.Error())

		urlCopy.Success = false
		urlCopy.Err = err.Error()
		json, err := ffjson.Marshal(&urlCopy)
		if err != nil {
			logger.Error(err.Error())
		}
		w.Write(json)
		return
	}

	// Long URL empty
	if urlCopy.LongURL == "" {
		logger.Error("URL null")

		urlCopy.Success = false
		urlCopy.Err = "URL null"
		json, err := ffjson.Marshal(&urlCopy)
		if err != nil {
			logger.Error(err.Error())
		}
		w.Write(json)
		return
	}

	db := newSession.DB("short")

	var result models.URL
	db.C("UrlCollection").Find(bson.M{"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(urlCopy.LongURL)))}).One(&result)

	// URL exists
	if result.LongURL != "" {
		//result.ShortURL = "http://localhost:8080/" + result.ID
		result.ShortURL = server.Base + result.ID
		json, err := ffjson.Marshal(&result)
		if err != nil {
			logger.Error(err.Error())
		}
		w.Write(json)
		return
	}

	if err = db.C("UrlCollection").Find(bson.M{}).Sort("-id").One(&result); err != nil {
		logger.Error(err.Error())
	}

	urlCopy.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(urlCopy.LongURL)))
	id, _ := genid.GenerateNextID(result.ID)
	urlCopy.ID = id
	urlCopy.Timestamp = time.Now()

	if err = db.C("UrlCollection").Insert(&urlCopy); err != nil {
		logger.Error(err.Error())
	}

	//urlCopy.ShortURL = "http://localhost:8080/" + urlCopy.ID
	urlCopy.ShortURL = server.Base + urlCopy.ID
	urlCopy.Success = true

	json, err := ffjson.Marshal(&urlCopy)
	if err != nil {
		logger.Error(err.Error())
	}

	w.Write(json)
}

// exists return whether the given file or directory exists or not
func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return false, err
	}

	return true, nil
}

func main() {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/short.log",
		MaxSize:    25,
		MaxBackups: 2,
		MaxAge:     14,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	logger = zap.New(core)

	defer logger.Sync()

	// Connect to mongodb
	session, err = mgo.Dial(fmt.Sprintf("%s:%d", database.Host, database.Port))
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	logger.Info("Connected successfully to MongoDB", zap.String("url", fmt.Sprintf("%s:%d", database.Host, database.Port)))

	defer session.Close()

	if err := session.DB("short").C("UrlCollection").EnsureIndex(mgo.Index{
		Key:    []string{"hash", "id"},
		Unique: true,
	}); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	if err := session.DB("short").C("UrlCollection").EnsureIndex(mgo.Index{
		Key:         []string{"ttl"},
		Background:  true,
		ExpireAfter: 0,
	}); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	fmt.Printf("Listening on %s:%d\n", server.Host, server.Port)

	/*if err := session.DB("short").C("UrlStatsCollection").EnsureIndex(mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}); err != nil {
		logger.Error(err.Error())
		panic(err)
	}*/

	mux := bone.New()
	mux.GetFunc("/", index)
	mux.GetFunc("/:id", redirectFull)
	mux.PostFunc("/api/v1/generate", generate)
	mux.DeleteFunc("/api/v1/delete/:id", delete)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", server.Host, server.Port), mux).Error())
}
