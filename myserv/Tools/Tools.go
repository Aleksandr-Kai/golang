package Tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	RootDir = "./html/img/"
)

type TError struct{
	Code int16
	Message string
}

type Album struct{
	Path		string		`json:"-"`
	Title		string		`json:"title""`
	Description	string		`json:"description""`
	Preview		string		`json:"image"`
	Images		[]string	`json:"-"`
	Thumbs		[]string	`json:"-"`
}

func (p *Album) CreateNew() error{
	aName := time.Now().Unix()
	p.Path = strconv.FormatInt(aName, 10)
	os.Mkdir(RootDir+ p.Path, os.ModePerm)
	desc,err := os.Create(RootDir + p.Path + "/description.txt")
	if err != nil{
		fmt.Printf("Не удалось создать файл описания альбома: %v\n", err.Error())
		return err
	}
	enc, _ := json.Marshal(p)
	_, err = desc.Write(enc)
	desc.Close()
	return err
}

func (p *Album) FullPath() string  {
	return "img/" + p.Path + "/"
}

func IsEmptyAlbum(path string) bool {
	files, err := ioutil.ReadDir(path)	// читаем папку с альбомами
	if err != nil {
		fmt.Printf("[IsEmptyAlbum]: %v\n", err.Error())
		return true
	}
	for _, file := range files {
		if !file.IsDir() && (strings.ToLower(filepath.Ext(file.Name())) == ".jpg"){
			return false
		}
	}
	return true
}

func GetAlbums() []Album {
	files, err := ioutil.ReadDir(RootDir) // читаем папку с альбомами
	if err != nil {
		fmt.Errorf("[GetAlbums]: %v", err.Error())
		return nil
	}
	albums := make([]Album, 0, 10)
	for _, file := range files {
		if file.IsDir() { // если папка, т.е. альбом
			txt, err := ioutil.ReadFile(RootDir + file.Name() + "/description.txt") // пытаемся читать описание
			al := &Album{}
			if err != nil {
				fmt.Printf("[GetAlbums] Не удалось открыть файл описания: %v\n", err.Error())
			}
			eerr := json.Unmarshal(txt, al)
			if eerr != nil {
				fmt.Printf("[GetAlbums] Ошибка распаковки json: %v\n", eerr.Error())
			}
			al.Path = file.Name()
			if (err != nil) || (eerr != nil) { // если не удалось, пытаемся исправить
				al.Title = "Album_Title"
				al.Description = "Album_Description"
				dFile, err := os.Create(RootDir + al.Path + "/description.txt")
				if err != nil{ // НЕ ВАЛИДНЫЙ АЛЬБОМ
					fmt.Printf("[GetAlbums] Не удалось создать файл описания альбома: %v\n", err.Error())
					continue
				}else {
					txt, _ = json.Marshal(al)
					_, err = dFile.Write(txt)
					dFile.Close()
					if err != nil { // НЕ ВАЛИДНЫЙ АЛЬБОМ
						continue
					}
				}
			}

			images, err := ioutil.ReadDir(RootDir + file.Name() + "/s/") // читаем папку с альбомами
			if err == nil {
				for _, img := range images {
					if !img.IsDir() && (strings.ToLower(filepath.Ext(img.Name())) == ".jpg") {
						al.Images = append(al.Images, img.Name())
					}
				}
			}
			if (al.Preview == "") && (len(al.Images) > 0) {
				al.Preview = al.Images[0]
			}

			al.Description += " \n" + file.Name()
			albums = append(albums, *al)
		}
	}
	return albums
}

func GetFilesList(path string) []string{
	files, err := ioutil.ReadDir(RootDir + path)
	if err != nil {
		fmt.Printf("Read File Sistem Fail: %v\n", err.Error())
		return nil
	}
	res := make([]string, 0)
	var dir []string
	for _, file := range files{
		if file.IsDir(){
			dir = GetFilesList(path + file.Name() + "/")
			if dir == nil{
				continue
			}
			res = append(res, dir...)
		}else{
			if strings.ToLower(filepath.Ext(file.Name())) == ".jpg" {
				res = append(res, path + file.Name())
			}
		}
	}
	return res
}

func Log(msg string, args ...interface{}){
	defer func() {
		if err := recover(); err != nil {
			log.Println("[Log] Epic fail!")
		}
	}()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	str := fmt.Sprintf("[%v] %v", funcName, msg)
	if len(args) > 0{
		str += fmt.Sprintf(":")
	}
	for _, arg := range args{
		str += fmt.Sprintf(" %v", arg)
	}

	log.Println(str)
}

func Message(msg string, args ...interface{}){
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	str := fmt.Sprintf("[%v] %v", funcName, msg)
	if len(args) > 0{
		str += fmt.Sprintf(":")
	}
	for _, arg := range args{
		str += fmt.Sprintf(" %v", arg)
	}

	fmt.Println(str)
}

func NamedMessage(prefix string, args ...interface{}){
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	if len(args) == 0{
		Log("No arguments")
		return
	}
	str := fmt.Sprintf("[%v]", prefix)
	for _, arg := range args{
		str += fmt.Sprintf(" %v", arg)
	}

	fmt.Println(str)
}