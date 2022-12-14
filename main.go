package main

import (
	"fmt"
	"html/template"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nfnt/resize"
)

func main() {

	r := httprouter.New()
	r.GET("/", Anasayfa)
	r.POST("/upload", Upload)
	r.ServeFiles("/admin/assets/*filepath", http.Dir("admin/assets"))
	log.Println("listening 7000")
	if err := http.ListenAndServe(":7000", r); err != nil {
		log.Println(err)
	}

}
func Anasayfa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, _ := template.ParseFiles("index.html")
	view.Execute(w, r)

}

func Upload(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	os.Remove("./resized.png")
	view, _ := template.ParseFiles("upload.html")
	view.Execute(w, r)
	// r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	header.Filename = "test.png"
	d, err := os.Create(header.Filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err1 := io.Copy(d, file)
	if err1 != nil {
		log.Fatal(err1)
	}
	old, err2 := os.Open("test.png")
	if err2 != nil {
		log.Fatal(err)
	}

	// f, err := os.OpenFile(header.Filename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0777)

	img, err3 := png.Decode(old)
	if err3 != nil {
		log.Fatal(err3)
	}
	old.Close()
	if err != nil {
		fmt.Println("hey")
		log.Fatal(err)
	}
	new, err := os.Create("admin/assets/resized.png")
	width := r.FormValue("width")

	u64, err := strconv.ParseUint(width, 10, 32)
	wd := uint(u64)
	m := resize.Resize(wd, 0, img, resize.Lanczos3)
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(new, m)

}
