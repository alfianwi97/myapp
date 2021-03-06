package uploads

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alfianwi97/myapp/pkg/router"
	"github.com/alfianwi97/myapp/pkg/server"
	"github.com/alfianwi97/myapp/pkg/store"
)

// UploadFile Function to Upload a File
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse Multipart Form Data
	err := r.ParseMultipartForm(server.Config.GetInt64("SERVER_UPLOAD_LIMIT"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	// Get File Content from Multipart Data
	mpFile, mpHeader, err := r.FormFile("file")
	if err != nil {
		router.ResponseBadRequest(w, err.Error())
		return
	}
	defer mpFile.Close()

	// Get File Metadata
	metaFileName := mpHeader.Filename
	metaFileSize := mpHeader.Size
	metaFileType := mpHeader.Header.Get("Content-Type")

	// Upload to Cloud Storage If Storage Driver Defined Else Save it to Local Storage
	switch strings.ToLower(server.Config.GetString("STORAGE_DRIVER")) {
	case "aws", "minio":
		err := store.S3UploadFile(metaFileName, metaFileSize, metaFileType, mpFile)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

		router.ResponseSuccess(w, "")
	default:
		// Default Save Uploaded File to Local Storage
		wrFile, err := os.OpenFile(server.Config.GetString("SERVER_UPLOAD_PATH")+"/"+metaFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
		defer wrFile.Close()

		// Copy Uploaded File Data from Multipart Data
		io.Copy(wrFile, mpFile)

		router.ResponseSuccess(w, "")
	}
}
