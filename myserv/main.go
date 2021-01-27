package main

import (
	"encoding/json"
	"fmt"
	"github.com/aleksandr-kai/golang/myserv/AlbumsTools"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var rnd *renderer.Render

type ImageInfo struct {
	Name		string	`json:"-"`
	Title		string	`json:"title"`
	Description	string	`json:"description"`
}

type TmplParams struct {
	AlbumTitle	string
	Image		[]ImageInfo
}

type TmplAlbum struct {
	AlbumImg			string
	AlbumTitle			string
	AlbumTitleComment	string
	AlbumDescription	string
	AlbumPath			string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("get_content")

	switch param {
	case "":{
		tmpls := []string{"html/templates/home.html", "html/templates/templates.html"}
		err := rnd.Template(w, http.StatusOK, tmpls, nil)
		if err != nil{
			fmt.Printf("%v\n", err)
		}
	}
	case "album-list":{
		Albums := AlbumsTools.GetAlbums()
		params := make([]TmplAlbum, len(Albums))

		for i, album := range Albums {
			if len(album.Images) == 0 {
				continue
			}
			params[i].AlbumTitleComment = fmt.Sprintf(" %v фото ", len(album.Images))
			params[i].AlbumTitle = album.Title
			//fmt.Printf("img/%v/%v\n", album.Path, album.Preview)
			if album.Preview == "" {
				params[i].AlbumImg = "img/no_images.png"
			} else {
				params[i].AlbumImg = "img?album=" + album.Path + "&name=" + album.Preview + "&size=s"
			}
			params[i].AlbumDescription = album.Description
			params[i].AlbumPath = album.Path
		}
		tmpls := []string{"html/templates/album-list.html"}
		err := rnd.Template(w, http.StatusOK, tmpls, params)
		if err != nil{
			fmt.Printf("%v\n", err)
		}
	}
	default:
		fmt.Fprintf(w, "Запрос не может быть обработан!")
	}

}

func galleryHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		http.NotFound(w, r)
		fmt.Printf("[galleryHandler] Не метод GET: %v\n", r.URL)
		return
	}
	param := r.URL.Query().Get("album")
	if param == ""{
		fmt.Printf("[galleryHandler] Параметр <album> не найден: %v\n", r.URL)
		return
	}
	Albums := AlbumsTools.GetAlbums()
	for _, album := range Albums {
		if album.Path == param {
			params := TmplParams{album.Title, make([]ImageInfo, len(album.Images))}
			for i, image := range album.Images {
				params.Image[i].Name = "img?album=" + album.Path + "&name=" + image + "&size=s"
				txt, err := ioutil.ReadFile(AlbumsTools.RootDir + album.Path + "/" + strings.TrimSuffix(image, filepath.Ext(image)) + ".txt") // пытаемся читать описание
				if err == nil {
					err = json.Unmarshal(txt, &params.Image[i])
					//params.Image[i].Description = string(txt)
					//fmt.Println("[galleryHandler] Description: " + string(txt))
				}/*else{
					fmt.Println("[galleryHandler] " + err.Error())
				}*/
			}
			tmpls := []string{"html/templates/gallery.html"}
			err := rnd.Template(w, http.StatusOK, tmpls, params)
			if err != nil{
				fmt.Printf("%v\n", err)
			}
			return
		}
	}
}
func imgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		fmt.Printf("[imgHandler] Не метод GET: %v\n", r.URL)
		return
	}

	album := r.URL.Query().Get("album")
	name := r.URL.Query().Get("name")
	size := r.URL.Query().Get("size")

	path := AlbumsTools.RootDir + album + "/" + size + "/" + name
	//fmt.Printf("ServeFile: %v\n", path)
	_, err := os.Stat(path)
	if err != nil {
		//fmt.Printf("File exists: %v\n", err.Error())
		w.Header().Set("Content-Type", "image/jpeg")
		if size == "s" {
			http.ServeFile(w, r, AlbumsTools.RootDir + "no_images.png")
		}else{
			http.ServeFile(w, r, AlbumsTools.RootDir + album + "/s/" + name)
		}

	}else{
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, path)
		//fmt.Printf("%v\nAlbum: %v\nName: %v\nFormat: %v\n", r.URL, album, name, size)
	}
}

func init() {
	rnd = renderer.New()
}

func main() {
	fs := http.FileServer(http.Dir("html"))
	mux := http.NewServeMux()
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/gallery", galleryHandler)
	mux.HandleFunc("/img", imgHandler)
	mux.Handle("/img/", fs)
	mux.Handle("/css/", fs)
	mux.Handle("/js/", fs)
	port := "80"
	fmt.Println("starting server at 127.0.0.1:" + port)
	http.ListenAndServe(":" + port, mux)
}