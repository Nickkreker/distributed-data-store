package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

const bodyFilePath = "body238974234576235678"

var saveBodyFile *os.File
var mutex sync.Mutex

func main() {
	persistenceDir := flag.String("pdir", "/var/lib/ddas", "Директория, в которой хранятся персистентные данные СУБД")
	flag.Parse()

	var err error
	err = os.MkdirAll(*persistenceDir, 0777)
	if err != nil {
		panic(err)
	}

	saveBodyFile, err = os.OpenFile(path.Join(*persistenceDir, bodyFilePath), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer saveBodyFile.Close()

	mux := http.NewServeMux()

	mux.Handle("/get", http.HandlerFunc(getHandler))
	mux.Handle("/replace", http.HandlerFunc(replaceHandler))

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	saveBodyFile.Seek(0, 0)
	_, err := io.Copy(w, saveBodyFile)
	if err != nil {
		panic(err)
	}
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	mutex.Lock()
	defer mutex.Unlock()

	saveBodyFile.Truncate(0)
	saveBodyFile.Seek(0, 0)
	_, err := io.Copy(saveBodyFile, r.Body)
	if err != nil {
		panic(err)
	}
}
