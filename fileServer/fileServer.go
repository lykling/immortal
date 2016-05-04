package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	ip               string
	port             string
	root             string
	uploadPath       string
	simpleUploadPage string
)

func initialize() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip address")
	flag.StringVar(&port, "port", "8080", "port to listen")
	flag.StringVar(&root, "root", "./", "root directory")
	flag.StringVar(&uploadPath, "data-path", root+"./upload", "upload path")
	simpleUploadPage = `
<html>
	<head>
		<title>Upload file</title>
	</head>
	<body>
		<form enctype="multipart/form-data" action="/upload" method="post">
			<input type="file" name="file" />
			<input type="submit" value="upload" />
		</form>
	</body>
</html>
`
}

// Upload handler
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tpl := template.New("upload")
		tpl, err := tpl.Parse(simpleUploadPage)
		if err != nil {
			log.Fatal(err)
		}
		tpl.Execute(w, nil)
	case "POST":
		// Make sure the upload path exist
		os.MkdirAll(uploadPath, 0755)
		source, err := r.MultipartReader()
		if err != nil {
			log.Fatal(err)
		}
		for {
			part, err := source.NextPart()
			if err == io.EOF {
				break
			}
			defer part.Close()
			// Save every file
			if part.FormName() == "file" {
				targetPath := uploadPath + "/" + part.FileName()
				target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal(err)
					return
				}
				defer target.Close()
				io.Copy(target, part)
			}
		}
		w.Write([]byte("Upload successfuly"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	initialize()
	flag.Parse()
	fmt.Fprintf(os.Stdout,
		"Listening at %s:%s\troot: %s\n", ip, port, root)
	fileServer := http.FileServer(http.Dir(root))
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", fileServer)
	err := http.ListenAndServe(ip+":"+port, nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
