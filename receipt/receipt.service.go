package receipt

import (
	"encoding/json"
	"fmt"
	"gowebservice/cors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const receiptPath = "receipts"

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptsList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptsList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) //5Mb
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666) //creating a file on disk
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()
		io.Copy(f, file) //copy the file data into the file which is cretaed on disk
		w.WriteHeader(http.StatusCreated)

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath)) //checking thr productID
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)

}

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))
}
