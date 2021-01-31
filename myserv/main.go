package main

import (
	"encoding/json"
	"fmt"
	"github.com/aleksandr-kai/golang/myserv/Tools"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rnd.Template(w, http.StatusOK, []string{"html/templates/404.html"}, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	inputLogin := r.FormValue("login")
	inputPassword := r.FormValue("password")

	fmt.Println("Login: ", inputLogin)
	fmt.Println("Password: ", inputPassword)

	coockie := http.Cookie{
		Name: "session_id",
		Value: inputLogin,
		Expires: time.Now().Add(5 * time.Minute),
	}
	http.SetCookie(w, &coockie)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != http.ErrNoCookie{
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}
	http.Redirect(w, r, "/home", http.StatusFound)
	fmt.Println("[logoutHandler] Session closed")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("get_content")

	switch param {
	case "":{
		session, err := r.Cookie("session_id")
		if err == nil {
			fmt.Printf("[homeHandler] Session: %v  %v  %v\n", session.Path, session.Name, session.Value)
		}
		tmpls := []string{"html/templates/home.html", "html/templates/templates.html"}
		err = rnd.Template(w, http.StatusOK, tmpls, nil)
		if err != nil{
			fmt.Printf("%v\n", err)
		}
	}
	case "album-list":{
		Albums := Tools.GetAlbums()
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
	case "test":{
		session, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println("[homeHandler] ", err.Error())
			rnd.Template(w, http.StatusOK, []string{"html/templates/null.html"}, nil)
			return
		}
		if session.Value != "AKai" {
			fmt.Println("[homeHandler] В доступе отказано")
			rnd.Template(w, http.StatusOK, []string{"html/templates/null.html"}, nil)
			return
		}
	}
	case "header":{

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
	Albums := Tools.GetAlbums()
	for _, album := range Albums {
		if album.Path == param {
			params := TmplParams{album.Title, make([]ImageInfo, len(album.Images))}
			for i, image := range album.Images {
				params.Image[i].Name = "img?album=" + album.Path + "&name=" + image + "&size=s"
				txt, err := ioutil.ReadFile(Tools.RootDir + album.Path + "/" + strings.TrimSuffix(image, filepath.Ext(image)) + ".txt") // пытаемся читать описание
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

	path := Tools.RootDir + album + "/" + size + "/" + name
	//fmt.Printf("ServeFile: %v\n", path)
	_, err := os.Stat(path)
	if err != nil {
		//fmt.Printf("File exists: %v\n", err.Error())
		w.Header().Set("Content-Type", "image/jpeg")
		if size == "s" {
			http.ServeFile(w, r, Tools.RootDir + "no_images.png")
		}else{
			http.ServeFile(w, r, Tools.RootDir + album + "/s/" + name)
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
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
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