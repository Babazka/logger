package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	listenAddress = flag.String("listen", ":5588", "address to listen on")
	logRoot       = flag.String("log-root", "/logs/da/logger/", "root directory for log files")
)

var (
	loggers     = make(map[string]*MyLogger)
	loggersLock sync.Mutex
)

var (
	filenameCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_\/]`)
)

type MyLogger struct {
	sync.Mutex
	f io.Writer
	w *bufio.Writer
}

func UrlToPath(url *url.URL) string {
	path := url.Path
	path = strings.Trim(path, "/")
	path = filenameCharsRegex.ReplaceAllLiteralString(path, "")
	return path
}

func NewMyLogger(logpath string) *MyLogger {
	var f io.Writer
	dir := path.Dir(logpath)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Printf("ERROR: cannot create directory %s: %s", dir, err)
	}
	f, err = os.OpenFile(logpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Printf("ERROR: cannot create logger for %s: %s", logpath, err)
		f = ioutil.Discard
	} else {
		log.Printf("opened logger for %s", logpath)
	}
	w := bufio.NewWriter(f)
	logger := &MyLogger{
		f: f,
		w: w,
	}
	return logger
}

func GetLogger(url *url.URL) *MyLogger {
	path := *logRoot + "/" + UrlToPath(url)
	loggersLock.Lock()
	logger, ok := loggers[path]
	if !ok {
		logger = NewMyLogger(path)
		loggers[path] = logger
	}
	loggersLock.Unlock()
	return logger
}

func GreatLogger(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format("2006-01-02T15:04:05.999")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read request body: %s: %s", r.RequestURI, err)
		return
	}
	r.Body.Close()
	logger := GetLogger(r.URL)
	logger.Lock()
	fmt.Fprintf(logger.w, "%s %s %s\n", t, r.Method, r.RequestURI)
	logger.w.Write(body)
	fmt.Fprintf(logger.w, "\n.\n")
	logger.w.Flush()
	logger.Unlock()
	w.Write([]byte("OK\r\n"))
}

var (
	stash []byte
)

func StashLogger(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(stash)
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("cannot read request body: %s: %s", r.RequestURI, err)
			return
		}
		r.Body.Close()
		stash = body
		w.Write([]byte("OK\r\n"))
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/stash", StashLogger)
	http.HandleFunc("/", GreatLogger)
	log.Println("Listening on ", *listenAddress)
	http.ListenAndServe(*listenAddress, nil)
}
