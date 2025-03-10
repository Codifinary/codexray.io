package utils

import (
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

type StaticFSWrapper struct {
	http.FileSystem
	modTime time.Time
}

func NewStaticFSWrapper(fs fs.FS) *StaticFSWrapper {
	return &StaticFSWrapper{FileSystem: http.FS(fs), modTime: time.Now()}
}

func (f *StaticFSWrapper) Open(name string) (http.File, error) {
	file, err := f.FileSystem.Open(name)
	return &StaticFileWrapper{File: file, modTime: f.modTime}, err
}

type StaticFileWrapper struct {
	http.File
	modTime time.Time
}

func (f *StaticFileWrapper) Stat() (os.FileInfo, error) {
	fileInfo, err := f.File.Stat()
	return &StaticFileInfoWrapper{FileInfo: fileInfo, modTime: f.modTime}, err
}

type StaticFileInfoWrapper struct {
	os.FileInfo
	modTime time.Time
}

func (f *StaticFileInfoWrapper) ModTime() time.Time {
	return f.modTime
}

func EnableCORS(next http.Handler, trustDomains map[string]struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			originHost := strings.TrimPrefix(origin, "http://")
			originHost = strings.TrimPrefix(originHost, "https://")

			if _, allowed := trustDomains[originHost]; allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
