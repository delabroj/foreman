package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/", newPackage)
	log.Fatal(http.ListenAndServe(":8081", http.DefaultServeMux))
}

func newPackage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(newPackageFrontend))
	case "POST":
		fail := func(message string, status int) {
			w.WriteHeader(status)
			w.Write([]byte(wrapHTML("<h1>" + message + `</h1><a href="/">Back</a>`)))
		}

		const maxUploadSize = 1024 * 1024 // 1 GB

		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			fail(fmt.Sprintf("Unable to parse multipart form: %s", err), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("package")
		if err != nil {
			fail(fmt.Sprintf("Unable to find uploaded file: %s", err), http.StatusBadRequest)
			return
		}
		fileStrings := strings.Split(header.Filename, ".")
		extension := "." + fileStrings[len(fileStrings)-1]

		if extension != ".zip" {
			fail("Package must be a zip file", http.StatusBadRequest)
			return
		}

		defer file.Close()

		dir := "/tmp/forman"
		err = os.RemoveAll(dir)
		if err != nil {
			fail(fmt.Sprintf("Unable to clear temp directory: %s", err), http.StatusInternalServerError)
			return
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fail(fmt.Sprintf("Unable to create temp directory: %s", err), http.StatusInternalServerError)
				return
			}
		}

		destFile := dir + "/package.zip"

		f, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, 0700)
		if err != nil {
			fail(fmt.Sprintf("Unable to save file: %s", err), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			fail(fmt.Sprintf("Unable to save file: %s", err), http.StatusInternalServerError)
			return
		}

		err = verify(destFile, r.FormValue("hash"))
		if err != nil {
			fail(fmt.Sprintf("Error checking hash: %s", err), http.StatusBadRequest)
			return
		}

		unzipDir := dir + "/package"

		cmd := exec.Command("unzip", "-o", destFile, "-d", unzipDir)
		err = cmd.Start()
		if err != nil {
			fail(fmt.Sprintf("Unable to process package: %s", err), http.StatusInternalServerError)
			return
		}
		err = cmd.Wait()
		if err != nil {
			fail(fmt.Sprintf("Unable to process package: %s", err), http.StatusInternalServerError)
			return
		}

		files, err := ioutil.ReadDir(unzipDir)
		if err != nil {
			fail(fmt.Sprintf("Unable to process package: %s", err), http.StatusInternalServerError)
			return
		}

		installFound := false
		for _, file := range files {
			if file.Name() == "install.sh" {
				installFound = true
				break
			}
		}
		if !installFound {
			fail("Unexpected package contents: missing install script", http.StatusBadRequest)
			return
		}

		out, err := exec.Command("bash", "-x", unzipDir+"/install.sh").CombinedOutput()
		if err != nil {
			fail(fmt.Sprintf("<h1>Unable to begin install: %s</h1><h2>Install script output:</h2><pre>", err)+fmt.Sprintf("%s", out)+"</pre>", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(wrapHTML(`<h1>Package received</h1><a href="/">Back</a><h2>Install script output:</h2><pre>` + fmt.Sprintf("%s", out) + "</pre>")))
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// verify hash of file
func verify(filePath string, hash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	fileHash := fmt.Sprintf("%x", h.Sum(nil))
	if fileHash != hash {
		return errors.New("Invalid package hash")
	}

	return nil
}
