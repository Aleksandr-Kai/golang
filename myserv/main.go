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
	UserName	string	`json:"user_name"`
}

func CheckSession(r *http.Request) (user Tools.User) {
	session, err := r.Cookie("session_id")

	if err != nil {
		println("[CheckAccess] ", err.Error())
		user, err = Tools.GetUser(Tools.DefaultUser)
		if err != nil{
			println("[CheckAccess] DefaultUser: ", err.Error())
			user = Tools.User{Name: Tools.DefaultUser, PublicName: "Гость", Password: "", Access: 10, Active: true}
		}
	}else{
		user, err = Tools.ParseToken(session.Value)
		if err != nil {
			println("[CheckAccess] ", err.Error())
			user, err = Tools.GetUser(Tools.DefaultUser)
			if err != nil{
				println("[CheckAccess] DefaultUser", err.Error())
				user = Tools.User{Name: Tools.DefaultUser, PublicName: "Гость", Password: "", Access: 10, Active: true}
			}
		}
	}
	return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := CheckSession(r)

		switch r.URL.Path {
		//-------------------------------------------------------------------------------------------------
		case "/":
				fallthrough
		//-------------------------------------------------------------------------------------------------
		/*         /home         */
		case "/home":{
			get_content := r.URL.Query().Get("get_content")
			switch get_content {
			//=============================================================================================================
			case "":{
				tmpls := []string{"html/templates/home.html", "html/templates/templates.html"}
				name := user.PublicName
				if name == ""{
					name = user.Name
				}
				prms := struct{UserName string; UserMenu interface{}}{name, GetUserMenu(user.Access)}
				err := rnd.Template(w, http.StatusOK, tmpls, prms)
				if err != nil{
					println("[rootHandler] ", err.Error())
				}
			}
			//=============================================================================================================
			case "album-list":{
				Albums := Tools.GetAlbums()
				params := make([]TmplAlbum, len(Albums))

				for i, album := range Albums {
					if len(album.Images) == 0 {
						continue
					}
					params[i].AlbumTitleComment = fmt.Sprintf(" %v фото ", len(album.Images))
					params[i].AlbumTitle = album.Title
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
			//=============================================================================================================
			case "config":{
				err := rnd.Template(w, http.StatusOK, []string{"html/templates/config.html"}, nil)
				if err != nil{
					println("[homeHandler] ", err.Error())
				}
			}
			//=============================================================================================================
			default:
				fmt.Fprintf(w, "Запрос не может быть обработан!")
			}
		}
		/*         /home         */
		//-------------------------------------------------------------------------------------------------
		/*         default (404)         */
		default:
			rnd.Template(w, http.StatusOK, []string{"html/templates/404.html"}, nil)
		}
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rnd.Template(w, http.StatusOK, []string{"html/templates/sign_in.html", "html/templates/templates.html"}, nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	inputName := r.FormValue("name")
	inputLogin := r.FormValue("login")
	inputPassword := r.FormValue("password")

	resp := LoginResponse{false, "Unknown error", ""}
	if (inputName == "") || (inputLogin == "") || (inputPassword == ""){
		fmt.Printf("[registerHandler] Не все поля заполнены: name:%v  login:%v  password:%v\n", inputName, inputLogin, inputPassword)
		resp.Success = false
		resp.Message = fmt.Sprintf("Не все поля заполнены: name:%v  login:%v  password:%v\n", inputName, inputLogin, inputPassword)
	}else{
		err := Tools.NewUser(inputLogin, inputName, inputPassword, 1)
		if err != nil{
			println("[registerHandler] Ошибка регистрации пользователя: ", err.Error())
			resp.Success = false
			resp.Message = err.Error()
		}else{
			Tools.SaveUsers("authdata.json")
			resp.Success = true
			resp.Message = "Добро пожаловать!"
			resp.UserName = inputName
		}
	}

	answer, _ := json.Marshal(resp)
	w.Write(answer)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	inputLogin := r.FormValue("login")
	inputPassword := r.FormValue("password")

	token, err := Tools.Login(inputLogin, inputPassword)
	resp := LoginResponse{false, "Unknown error", "Гость"}
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
		resp.UserName = inputLogin
	}
	answer, _ := json.Marshal(resp)
	w.Write(answer)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == nil{
		//fmt.Println("[logoutHandler] Куки получено")
		session.Expires = time.Now().AddDate(0, 0, -1)
		user, err := Tools.ParseToken(session.Value)
		if err != nil{
			println("[logoutHandler] ", err.Error())
		}else{
			Tools.Logout(user.Name)
			fmt.Println("[logoutHandler] Выход выполнен")
			Tools.Login(Tools.DefaultUser, "")
		}
		session.Value = ""
		http.SetCookie(w, session)
	}else{
		fmt.Println("[logoutHandler] ", err.Error())
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

func GetUserMenu(lvl int) template.HTML{
	var mLogin = template.HTML(`
		<li id="login"><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#modal-login">Вход</a></li>
	`)
	var mLogout = template.HTML(`
		<li id="logout"><a class="dropdown-item" href="/logout">Выход</a></li>
	`)
	var mConfig = template.HTML(`
		<li><a class="dropdown-item" href="javascript:void(0);" onclick="LoadConfig();">Настройки</a></li>
	`)
	var mLine = template.HTML(`
		<li id="divider"><hr class="dropdown-divider"></li>
	`)

	switch lvl{
	case 0:{
		return mConfig + mLine + mLogout
	}
	case 1:
		return mConfig + mLine + mLogout
	default:
		return mLogin
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
	err := Tools.NewUser(Tools.DefaultUser, "Гость", "", 10)
	if err != nil{
		println("[Main Add Guest] ", err.Error())
		err = Tools.UpdateUser(Tools.User{Name: Tools.DefaultUser, PublicName: "Гость", Password: "", Access: 10, Active: false})
		if err != nil {
			println("[Main Update Guest] ", err.Error())
		}
	}
	println("[Main] Сохранение пользователе: ", Tools.SaveUsers("authdata.json"))
	fs := http.FileServer(http.Dir("html"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/sign-in", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/gallery", galleryHandler)
	mux.HandleFunc("/img", imgHandler)
	mux.Handle("/img/", fs)
	mux.Handle("/css/", fs)
	mux.Handle("/js/", fs)
	port := "80"
	fmt.Println("starting server at 127.0.0.1:" + port)
	http.ListenAndServe(":" + port, mux)
}