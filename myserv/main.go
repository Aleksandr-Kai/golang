package main

import (
	"encoding/json"
	"fmt"
	"github.com/aleksandr-kai/golang/myserv/Tools"
	"github.com/thedevsaddam/renderer"
	"html/template"
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

type LoginResponse struct {
	Success		bool	`json:"success"`
	Message		string	`json:"message"`
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

	token, err := Tools.Login(inputLogin, inputPassword)
	resp := LoginResponse{false, "Unknown error"}
	if err != nil {
		println("[loginHandler] Ошибка получения токена: ", err.Error())
		resp.Success = false
		resp.Message = err.Error()
	}else{
		coockie := http.Cookie{
			Name: "session_id",
			Value: token,
			Expires: time.Now().Add(512 * time.Hour),
			Path: "/",
		}
		http.SetCookie(w, &coockie)
		fmt.Println("[loginHandler] Выполнен вход под именем [", inputLogin, "]")
		resp.Success = true
		resp.Message = ""
	}
	answer, _ := json.Marshal(resp)
	w.Write(answer)
	//http.Redirect(w, r, "/home", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != http.ErrNoCookie{
		//fmt.Println("[logoutHandler] Куки получено")
		session.Expires = time.Now().AddDate(0, 0, -1)
		user, err := Tools.ParseToken(session.Value)
		if err != nil{
			println("[logoutHandler] ", err.Error())
		}else{
			Tools.Logout(user.Name)
			fmt.Println("[logoutHandler] Выход выполнен")
		}
		session.Value = ""
		http.SetCookie(w, session)
	}else{
		fmt.Println("[logoutHandler] ", err.Error())
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

var menuGuest = template.HTML(`
	<li id="login"><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#modal-login">Вход</a></li>
`)
var menuAdmin = template.HTML(`
	<li><a class="dropdown-item" href="#" aria-disabled="true">Настройки</a></li>
	<li id="divider"><hr class="dropdown-divider"></li>
	<li id="logout"><a class="dropdown-item" href="/logout">Выход</a></li>
`)
var menuUser = template.HTML(`
	<li><a class="dropdown-item disabled" href="#" aria-disabled="true">Настройки</a></li>
	<li id="divider"><hr class="dropdown-divider"></li>
	<li id="logout"><a class="dropdown-item" href="/logout">Выход</a></li>
`)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("get_content")

	switch param {
	case "":{
		session, err := r.Cookie("session_id")
		var user Tools.User
		if err != nil {
			println("[homeHandler] ", err.Error())
			user = Tools.User{Name: "Guest", Password: "", Access: 10, Active: true}
		}else{
			user, err = Tools.ParseToken(session.Value)
			if err != nil {
				println("[homeHandler] ", err.Error())
				user = Tools.User{Name: "Guest", Password: "", Access: 10, Active: true}
			}
		}
		fmt.Printf("[homeHandler] Запрос от пользователя: %v\n", user)
		/*
		tmplFuncs := template.FuncMap{"UserMenu" : UserMenu}
		tmpl, err := template.New("").Funcs(tmplFuncs).ParseFiles("html/templates/home.html", "html/templates/templates.html")
		if err == nil {
			err = tmpl.ExecuteTemplate(w, "home.html", struct{UserName string}{user.Name})
			if err != nil {
				fmt.Println("[homeHandler] ", err.Error())
			}
			return
		}
		fmt.Println("[homeHandler] ", err.Error())*/

		tmpls := []string{"html/templates/home.html", "html/templates/templates.html"}
		var menu template.HTML
		switch user.Access {
		case 0: menu = menuAdmin
		case 1: menu = menuUser
		default:
			menu = menuGuest
		}
		prms := struct{UserName string; UserMenu interface{}}{user.Name, menu}
		err = rnd.Template(w, http.StatusOK, tmpls, prms)
		if err != nil{
			println("[homeHandler] %v\n", err.Error())
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
	Tools.Init("authdata.json")
	Tools.NewUser("Guest", "", 10)
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