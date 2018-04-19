package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
	"golang.org/x/text/language"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	config "shlink-server/conf"
	"shlink-server/models"
	"shlink-server/pkg/genid"
	utils "shlink-server/utils"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	json "github.com/pquerna/ffjson/ffjson"
	"go.uber.org/zap"
)

var (
	conf = config.New()

	logger  *zap.Logger
	session *mgo.Session

	err error

	version    string
	goVersion  string
	goPlatform string

	debug = flag.Bool("debug", false, "Enable stdout logger")
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lg, _ := language.Parse(r.Header.Get("Accept-Language"))

	message := map[string]map[string]string{
		"en": {
			"greet": "It's rendering",
		},
		"fr": {
			"greet": "Affichage de",
		},
	}
	t, err := template.New("index.html").Funcs(template.FuncMap{
		"T": func(key string) string {
			return message[lg.String()][key]
		},
	}).ParseFiles("public/index.html", "public/head.html")
	if err != nil {
		logger.Error(err.Error())
	}

	err = t.Execute(w, nil)
	if err != nil {
		logger.Error(err.Error())
	}
}

func redirectFull(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(conf.Database.DB)

	//var result models.URL
	result := models.URL{}
	err := db.C("url").Find(bson.M{"id": id}).One(&result)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", r.Method), zap.String("path", r.RequestURI))
	}

	if result.LongURL != "" {
		logger.Info("Access", zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))
		http.Redirect(w, r, result.LongURL, 301)

		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

// generate handles short url endpoint
func generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newSession := session.Copy()
	defer newSession.Close()

	// Unmarshal JSON request
	urlCopy := models.URL{}

	if err := json.NewDecoder().DecodeReader(r.Body, &urlCopy); err != nil {
		logger.Error(err.Error(), zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))

		urlCopy.Success = false
		urlCopy.Err = err.Error()
		json, _ := json.Marshal(&urlCopy)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)
		return
	}

	// Long URL empty
	if urlCopy.LongURL == "" {
		logger.Warn("Empty URL", zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))

		urlCopy.Success = false
		urlCopy.Err = "URL null"
		json, _ := json.Marshal(&urlCopy)
		w.Write(json)
		return
	}

	// URL is invalid
	if ok, _ := utils.IsURL(urlCopy.LongURL); !ok {
		logger.Warn("Invalid URL", zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("url", urlCopy.LongURL), zap.String("client", r.RemoteAddr))

		urlCopy.Success = false
		urlCopy.Err = "URL invalid: " + urlCopy.LongURL
		json, err := json.Marshal(&urlCopy)
		if err != nil {
			logger.Error(err.Error(), zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))
		}
		w.Write(json)
		return
	}

	// URL is on the blacklist
	if blacklisted, _ := utils.IsBlackList(urlCopy.LongURL); blacklisted {
		logger.Warn("Unallowed URL", zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("url", urlCopy.LongURL), zap.String("client", r.RemoteAddr))

		urlCopy.Success = false
		urlCopy.Err = "URL is on blacklist: " + urlCopy.LongURL
		json, _ := json.Marshal(&urlCopy)
		w.Write(json)
		return
	}

	// Trim trailing slash
	if strings.HasSuffix(urlCopy.LongURL, "/") {
		urlCopy.LongURL = urlCopy.LongURL[0 : len(urlCopy.LongURL)-1]
	}

	// Reorder url query for consistency
	urlCopy.LongURL = utils.ReOrderQuery(urlCopy.LongURL)

	db := newSession.DB(conf.Database.DB)

	if !strings.HasSuffix(conf.Server.Base, "/") {
		conf.Server.Base += "/"
	}

	// Check if url exists
	result := models.URL{}
	db.C("url").Find(bson.M{"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(urlCopy.LongURL)))}).One(&result)
	if result.LongURL != "" {
		result.Success = true
		result.ShortURL = conf.Server.Base + result.ID
		json, _ := json.Marshal(&result)
		w.Write(json)
		return
	}

	// Code for FindAndModify
	doc := models.Counter{}

	changes := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"sequence": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	_, err := db.C("counter").Find(bson.M{"_id": "shlink.cc"}).Apply(changes, &doc)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))
	}

	urlCopy.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(urlCopy.LongURL)))
	id := genid.IntToBase62(doc.Sequence - 1)
	urlCopy.ID = id
	urlCopy.Timestamp = time.Now()

	if err = db.C("url").Insert(&urlCopy); err != nil {
		logger.Error(err.Error(), zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))
	}

	urlCopy.ShortURL = conf.Server.Base + urlCopy.ID
	urlCopy.Success = true

	json, _ := json.Marshal(&urlCopy)

	w.Write(json)
}

func info(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(conf.Database.DB)

	result := models.URL{}
	err := db.C("url").Find(bson.M{"id": id}).One(&result)
	if err != nil {
		logger.Error(err.Error())
	}

	if result.LongURL == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))

		return
	}

	result.Success = true

	js, err := json.Marshal(&result)
	if err != nil {
		logger.Error(err.Error())
	}

	w.Write(js)
}

// handleRobots
func handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User-agent: *\nDisallow: /"))
}

// status returns status
func status(w http.ResponseWriter, r *http.Request) {
	resp, _ := json.Marshal(struct {
		Success    bool   `json:"success"`
		Version    string `json:"version"`
		GoVersion  string `json:"go"`
		GoPlatform string `json:"platform"`
	}{
		Success:    true,
		Version:    version,
		GoVersion:  goVersion,
		GoPlatform: goPlatform,
	})

	w.Write(resp)
}

// middlewareLog handles log
func middlewareLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Access", zap.String("method", r.Method), zap.String("path", r.RequestURI), zap.String("client", r.RemoteAddr))
		next(w, r)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from http.FileSystem
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler("/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func main() {

	flag.Parse()

	conf.ReadConfig()

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.Log.Filename,
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   true,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	logger = zap.New(core)

	defer logger.Sync()

	// Connect to mongodb
	session, err = mgo.Dial(conf.Database.Host + ":" + conf.Database.Port)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	defer session.Close()

	db := session.DB(conf.Database.DB)

	if err := db.C("url").EnsureIndex(mgo.Index{
		Key:    []string{"hash", "id"},
		Unique: true,
	}); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	if err := db.C("counter").EnsureIndex(mgo.Index{
		Key:    []string{"_id", "sequence"},
		Unique: true,
	}); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	fmt.Printf("Shlink-Server %s\n", version)
	fmt.Printf("go: %s\n", goVersion)
	fmt.Printf("platform: %s\n", goPlatform)
	fmt.Printf("Listening on %s:%s\n", conf.Server.Host, conf.Server.Port)

	// Router
	r := chi.NewRouter()

	// Enable CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"Accept-Language", "Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Middlewares
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler)
	r.Use(middleware.RealIP)
	if *debug {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	// Endpoints
	r.Get("/", middlewareLog(index))
	r.Get("/robots.txt", handleRobots)
	r.Get("/{id}", redirectFull)

	r.Post("/api/generate", generate)

	r.Get("/api/status", status)

	r.Get("/api/info/{id}", info)

	// Serve files
	FileServer(r, "/public", http.Dir("./public"))

	logger.Fatal(http.ListenAndServe(conf.Server.Host+":"+conf.Server.Port, r).Error())
}
