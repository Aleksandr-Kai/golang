package FilesCollection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	RootDir = "./html/img/"
	PageSize = 24
)

const (
	Current = -1
	Next = -2
	Prev = -3
)

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
	os.Mkdir(RootDir + p.Path, os.ModePerm)
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
/*
func GetAlbum(path string) (Album, error) {
	txt, err := ioutil.ReadFile(RootDir + path + "/description.txt") // пытаемся читать описание
	if err != nil {
		fmt.Printf("[GetAlbum] Не удалось открыть файл описания: %v\n", err.Error())
		return Album{}, err
	}
	al := &Album{}
	err = json.Unmarshal(txt, al)
	if err != nil {
		fmt.Printf("[GetAlbum] Ошибка распаковки json: %v\n", err.Error())
		return Album{}, err
	}
}
*/
func GetAlbums() []Album {
	files, err := ioutil.ReadDir(RootDir)	// читаем папку с альбомами
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

			images, err := ioutil.ReadDir(RootDir + file.Name() + "/s/")	// читаем папку с альбомами
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

type Collection struct {
	Files []string
	Thumbs []string
	CurrentPage uint8
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

func (p *Collection) Count() uint32{
	return uint32(len(p.Files))
}

func (p *Collection) Set(value []string){
	p.CurrentPage = 0
	p.Files = make([]string, len(value))
	copy(p.Files, value)
}

func (p *Collection) PageCount(size uint8) uint8{
	pageCount := uint8(p.Count() / uint32(size))
	if p.Count() - uint32(pageCount * size) > 0{
		pageCount++
	}
	return pageCount
}

func (p *Collection) GetPage(page int16, size uint8) []string {
	switch page {
	case Next:{
		p.CurrentPage++
		if p.CurrentPage >= p.PageCount(size){
			p.CurrentPage = 0
		}
	}
	case Prev:{
		if p.CurrentPage == 0{
			p.CurrentPage = p.PageCount(size) - 1
		}else{
			p.CurrentPage--
		}
	}
	case Current:

	default:
		if page >= 0{
			p.CurrentPage = uint8(page)
		}else {
			return nil
		}
	}

	pos := p.CurrentPage * size
	if p.CurrentPage == p.PageCount(size) - 1{
		return p.Files[pos:]
	}else{
		return p.Files[pos:pos + size]
	}
}