package services

import (
	"io"
	"log"
	"net/http"
	"os"
)

var TmpDirPath string

// 1MB
const MAX_UPLOADING_FILE_SIZE_BYTES = 1 << 20

func uploadFile(w http.ResponseWriter, r *http.Request, formFileName string) (tempFile *os.File, originalFilename string, err error) {
	r.ParseMultipartForm(MAX_UPLOADING_FILE_SIZE_BYTES)

	file, handler, err := r.FormFile(formFileName)

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	log.Printf("Uploading file: %+v\n", handler.Filename)
	log.Printf("File size: %+v\n", handler.Size)

	tempFile, err = os.CreateTemp(TmpDirPath, "upload-*")

	if err != nil {
		log.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		log.Println(err)
		return
	}

	tempFile.Write(fileBytes)
	originalFilename = handler.Filename

	log.Printf("File %+v successfully uploaded\n", handler.Filename)

	return
}
