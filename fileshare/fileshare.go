package main

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "log"
    "io/ioutil"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
    files, _ := ioutil.ReadDir("uploads")
    data["files"] = files
    templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    //GET displays the upload form.
    case "GET":
        data := make(map[string]interface{})
        display(w, "upload", data)

    //POST takes the uploaded file(s) and saves it to disk.
    case "POST":
        //get the multipart reader for the request.
        reader, err := r.MultipartReader()

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        //copy each part to destination.
        for {
            part, err := reader.NextPart()
            if err == io.EOF {
                break
            }

            //if part.FileName() is empty, skip this iteration.
            if part.FileName() == "" {
                continue
            }
            dst, err := os.Create("uploads/" + part.FileName())
            defer dst.Close()

            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if _, err := io.Copy(dst, part); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
        //display success message.
        res := make(map[string]interface{})
        res["info"] = "Upload successful."
        display(w, "upload", res)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
     http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

    //static file handler.
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

    //Listen on port 8080
    err := http.ListenAndServeTLS(":443", "fileshare.crt", "fileshare.key", nil)
    if err != nil {
        log.Fatal(err)
    }
}
