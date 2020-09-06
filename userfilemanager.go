package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileData struct {
	File string
	Name string
}

func sendfile(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		f, h, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(filepath.Join("./recieved/", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		err = ioutil.WriteFile(filepath.Join("./recieved/", h.Filename), bs, os.ModeAppend)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if req.FormValue("sendback") == "on" {
			http.ServeFile(w, req, filepath.Join("./recieved/", h.Filename))
			return
		}

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpls.ExecuteTemplate(w, "sendfile.gohtml", nil)

}

func returnRecieved(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/recieved/" {
		var files []string

		err := filepath.Walk("./recieved", func(path string, info os.FileInfo, err error) error {
			files = append(files, info.Name())
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files = escapeBads(files[1:])
		fileNames := unescapeBads(files)

		var fileDatas []fileData

		for i, file := range files {
			fileDatas = append(fileDatas, fileData{file, fileNames[i]})
		}

		tpls.ExecuteTemplate(w, "recieved.gohtml", fileDatas)
		return

	}

	http.ServeFile(w, req, filepath.Join(".", req.URL.Path))
}

func escapeBads(slice []string) []string {
	strs := make([]string, len(slice))
	copy(strs, slice)

	for key, _ := range badChars {
		for strI, str := range strs {
			strs[strI] = strings.Replace(str, key, badChars[key], -1)
		}
	}

	return strs
}

func unescapeBads(slice []string) []string {
	strs := make([]string, len(slice))
	copy(strs, slice)

	for key, val := range badChars {
		for strI, str := range strs {
			strs[strI] = strings.Replace(str, val, key, -1)
		}
	}

	return strs
}
