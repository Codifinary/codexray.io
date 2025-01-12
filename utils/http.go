package utils

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

const FRONTE_END_BRUM = "http://34.47.146.55:3000"

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

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", FRONTE_END_BRUM)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
