package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// APIServer struct to handle all incoming requests
type APIServer struct {
	Backend interface{}
}

// StartListening for incomming request
func (s *APIServer) StartListening() {
	http.HandleFunc("/health", GenerateHandler("^/health$", HealthHandler))
	http.HandleFunc("/static/", GenerateHandler("^/(static/(js/|css/|media/)[a-zA-Z0-9._]*)$", FileHandler))
	http.HandleFunc("/api/", GenerateHandler(config.APIPattern, s.APIHandler))
	http.HandleFunc("/", GenerateHandler("^/(.*)$", IndexHandler))
	a := fmt.Sprintf("%s:%s", config.Host, config.Port)
	logger.Infof("Start listening \"%s\"...", a)
	logger.Fatale(http.ListenAndServe(a, nil), "Server crashed !")
}

// GenerateHandler handler
func GenerateHandler(p string, f func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Tracef("Catch request %s %s", p, reflect.TypeOf(f).Name())
		v := regexp.MustCompile(p)
		m := v.FindStringSubmatch(r.URL.Path)
		if m == nil {
			logger.Warningf("Invalid path \"%s\" doesn't match pattern \"%s\"", r.URL.Path, p)
			WriteJSONErrorResponse(w, "Invalid url", http.StatusInternalServerError)
			return
		}
		logger.Tracef("Pattern %s matched : %q", p, m)
		defer func(w http.ResponseWriter, r *http.Request) {
			if e := recover(); e != nil {
				logger.Recoverf("Recover from handling request : %s", e)
				WriteJSONErrorResponse(w, fmt.Sprintf("An error occurred : %s", e), http.StatusInternalServerError)
			}
		}(w, r)
		if len(m) > 1 {
			f(w, r, m[1])
		} else {
			f(w, r, m[0])
		}
	}
}

// HealthHandler handler
func HealthHandler(w http.ResponseWriter, r *http.Request, n string) {
	SetNoCacheHeaders(w)
	logger.Trace("Health check !")
	WriteJSONResponse(w, "{\"status\":\"OK\"}")
}

// IndexHandler handler
func IndexHandler(w http.ResponseWriter, r *http.Request, n string) {
	switch r.Method {
	case "GET":
		m := ""
		logger.Tracef("Index handling %s", n)
		switch n {
		case "", "host":
			logger.Tracef("Valid index path : %s", n)
			m = fmt.Sprintf("%s/index.html", config.WorkingDir)
		case "favicon.ico", "manifest.json", "logo192.png", "logo512.png", "robots.txt":
			logger.Tracef("Valid special file path : %s", n)
			m = fmt.Sprintf("%s/%s", config.WorkingDir, n)
		default:
			logger.Warningf("Invalid index path : %s ", n)
		}
		serveFile(w, r, m)
	default:
		logger.Warningf("Wrong request method \"%s\" !", r.Method)
		WriteJSONErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// FileHandler handler
func FileHandler(w http.ResponseWriter, r *http.Request, n string) {
	switch r.Method {
	case "GET":
		serveFile(w, r, fmt.Sprintf("%s/%s", config.WorkingDir, n))
	default:
		logger.Warningf("Wrong request method \"%s\" !", r.Method)
		WriteJSONErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, f string) {
	if len(f) > 0 {
		if _, e := os.Stat(f); os.IsNotExist(e) {
			logger.Debugf("File not found \"%s\"", f)
			WriteJSONErrorResponse(w, "File not found", http.StatusNotFound)
		} else {
			logger.Debugf("Serve file \"%s\"", f)
			http.ServeFile(w, r, f)
		}
	} else {
		logger.Debugf("Not found !")
		WriteJSONErrorResponse(w, "Not found", http.StatusNotFound)
	}
}

// APIHandler handler
func (s *APIServer) APIHandler(w http.ResponseWriter, r *http.Request, n string) {
	SetNoCacheHeaders(w)
	logger.Debugf("Request API received : %s", n)
	m := FormatEndpointMethod(n, r.Method)
	ExecEndpoint(s.Backend, m, w, r)
}

// FormatEndpointMethod to get method name from url
func FormatEndpointMethod(n, o string) string {
	p := strings.Split(n, "/")
	m := strings.ToLower(o)
	for _, s := range p {
		m = fmt.Sprintf("%s%s", strings.Title(m), strings.Title(s))
	}
	return m
}

// ExecEndpoint to run method from url
func ExecEndpoint(i interface{}, m string, w http.ResponseWriter, r *http.Request) {
	a := reflect.ValueOf(i)
	logger.Debugf("Method call : %s", m)
	f := a.MethodByName(m)
	if f.IsZero() {
		WriteJSONErrorResponse(w, fmt.Sprintf("Invalid method : %s", m), http.StatusInternalServerError)
	} else {
		q := []reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
		}
		f.Call(q)
	}
}

// SetNoCacheHeaders to prevent browser caching
func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// WriteTextResponse write response in plain text format
func WriteTextResponse(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, s)
}

// WriteJSONResponse write response in json format with correct header
func WriteJSONResponse(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, s)
}

// ErrorResponse struct to handle response with error message
type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// WriteJSONErrorResponse write error response in json format
func WriteJSONErrorResponse(w http.ResponseWriter, m string, s int) {
	SetNoCacheHeaders(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(s)
	e := &ErrorResponse{
		Error:   true,
		Message: m,
	}
	r, _ := json.Marshal(e)
	io.WriteString(w, string(r))
}
