package frodo

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
)

// FileUploadsPath for now, declares the path to upload the files
var FileUploadsPath = "./assets/uploads/"

// UploadedFile struct/type is the data that makes up an uploaded file
// once it is recieved and parsed eg. using request.FormFile()
type UploadedFile struct {
	multipart.File
	*multipart.FileHeader
	/*
	   type FileHeader struct {
	       Filename string
	       Header   textproto.MIMEHeader
	   }
	*/
}

// Name returns the name of the file when it was uploaded
func (file *UploadedFile) Name() string {
	// found in *multipart.FileHeader
	return file.Filename
}

// Size returns the size of the file in question
func (file *UploadedFile) Size() int64 {
	defer file.Close()
	return file.Size()
}

// Extension returns the extension of the file uploaded
func (file *UploadedFile) Extension() string {
	// _, header, error := r.FormFile(name)
	ext := filepath.Ext(file.Filename)
	return ext
}

// Move basically moves/transfers the uploaded file to the upload folder provided
//
// Using ...interface{} because I want the user to only pass more than one argument
// when changing upload dir and filename, if none is changed then defaults  are used
//
//    eg. file.Move(true)
//        ----- or -----
//        file.Move("../new_upload_path/", "newfilename.png")
//
func (file *UploadedFile) Move(args ...interface{}) bool {
	file.Open()
	defer file.Close()
	name := args[0]
	val := reflect.ValueOf(name)

	// If a string was give, then treat is a the FileUploadsPath
	if val.Kind().String() == "string" {
		FileUploadsPath = name.(string)
	}

	var FileName string
	// Check to see if a file name was given, 2nd argument
	if len(args) > 1 {
		FileName = args[1].(string)
	} else {
		FileName = file.Name()
	}

	savedFile, err := os.OpenFile(FileUploadsPath+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, ioerr := io.Copy(savedFile, file)
	if ioerr != nil {
		fmt.Println(ioerr)
		return false
	}

	return true
}

// MimeType returns the mime/type of the file uploaded
func (file *UploadedFile) MimeType() string {
	mimetype := file.Header.Get("Content-Type")
	return mimetype
}

// IsValid checks if the file is alright by opening it up
// if errors come up while opening it is an invalid upload
func (file *UploadedFile) IsValid() bool {
	_, err := file.Open()
	defer file.Close()
	if err != nil {
		return false
	}
	return true
}
