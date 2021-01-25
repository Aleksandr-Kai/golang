package FilesCollection

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetFilesList(t *testing.T) {
	fmt.Println("----------- TestGetFilesList -----------")
	files := GetFilesList("")
	fmt.Printf("Найденные файлы: %v\n", files)
	fmt.Printf("Количество файлов: %v\n", len(files))
	//fmt.Println("----------------------------------------")
}

func TestFilesCollection_Set(t *testing.T) {
	fmt.Println("-------- TestFilesCollection_Set -------")
	files := GetFilesList("")
	var imgs Collection
	imgs.Set(files)
	fmt.Printf("Количество файлов: %v\n", imgs.Count())
	if imgs.Count() != uint32(len(files)){
		t.Error("Не верное количество файлов: ", imgs.Count(),  " вместо ", len(files))
	}
	//fmt.Println("----------------------------------------")
}

func TestFilesCollection_GetPage(t *testing.T) {
	fmt.Println("------ TestFilesCollection_GetPage -----")
	var imgs Collection
	files := GetFilesList("")
	imgs.Set(files)
	page := imgs.GetPage(-200, PageSize)
	if page != nil{
		t.Error("Не обработан ошибочный параметр")
	}
	page = imgs.GetPage(Current, PageSize)
	if len(files) >= PageSize{
		if !reflect.DeepEqual(page, files[0:PageSize]){
			t.Error("Получен не верный слайс для [current]: ", page, " > ", files[0:PageSize])
		}
	}else{
		if !reflect.DeepEqual(page, files[0:]){
			t.Error("Получен не верный слайс для [current]: ", page, " > ", files[0:])
		}
	}

	page = imgs.GetPage(Next, PageSize)
	if len(files[PageSize:]) >= PageSize{
		if !reflect.DeepEqual(page, files[PageSize:PageSize*2]){
			t.Error("Получен не верный слайс для [Next]: ", page, " > ", files[PageSize:PageSize*2])
		}
	}else{
		if !reflect.DeepEqual(page, files[PageSize:]){
			t.Error("Получен не верный слайс для [Next]: ", page, " > ", files[PageSize:])
		}
	}

	page = imgs.GetPage(Prev, PageSize)
	if len(files) >= PageSize{
		if !reflect.DeepEqual(page, files[0:PageSize]){
			t.Error("Получен не верный слайс для [Prev]: ", page, " > ", files[0:PageSize])
		}
	}else{
		if !reflect.DeepEqual(page, files[0:]){
			t.Error("Получен не верный слайс для [Prev]: ", page, " > ", files[0:])
		}
	}
	for i := 0; i < int(imgs.PageCount(PageSize)) - 1; i++{
		page = imgs.GetPage(Next, PageSize)
	}
	fmt.Printf("Последняя страница: %v\n", page)
	if !reflect.DeepEqual(page, files[(imgs.PageCount(PageSize) - 1) * PageSize:]){
		t.Error("Получен не верный слайс для [tail]: ", page, " > ", files[(imgs.PageCount(PageSize) - 1) * PageSize:])
	}
	page = imgs.GetPage(Next, PageSize)
	if len(files) >= PageSize{
		if !reflect.DeepEqual(page, files[0:PageSize]){
			t.Error("Получен не верный слайс для [Next] от последней: ", page, " > ", files[0:PageSize])
		}
	}else{
		if !reflect.DeepEqual(page, files[0:]){
			t.Error("Получен не верный слайс для [Next] от последней: ", page, " > ", files[0:])
		}
	}
	page = imgs.GetPage(Prev, PageSize)
	if !reflect.DeepEqual(page, files[(imgs.PageCount(PageSize) - 1) * PageSize:]){
		t.Error("Получен не верный слайс для [Prev] от первой: ", page, " > ", files[(imgs.PageCount(PageSize) - 1) * PageSize:])
	}
}

func TestAlbum_CreateNew(t *testing.T) {
	fmt.Println("---------- TestAlbum_CreateNew ---------")
	var a Album
	a.Description = "This doesn't seem to actually change the terminal location after I run the go program. When I print the error, I see that there was no error for os.Chdir(). I see that error is <nil>"
	a.Title = "TestAlbum"
	err := a.CreateNew()
	if err != nil{
		t.Error("Ошибка при созданиий альбома: ", err.Error())
	}
}

func TestGetAlbums(t *testing.T) {
	fmt.Println("------------ TestGetAlbums -------------")
	list := GetAlbums()
	if list == nil {
		t.Error("Внутренняя ошибка")
	}
	for _, l := range list {
		fmt.Printf("\n>>>\nНазвание: %v\nОписание: %v\nПуть: %v\nИзображения: %v\nОбложка: %v\n", l.Title, l.Description, l.Path, l.Images, l.Preview)
	}
}

func TestIsEmptyAlbum(t *testing.T) {
	fmt.Println("----------- TestIsEmptyAlbum -----------")
	res := IsEmptyAlbum("..//html/img/empty/")
	if res != true {
		t.Error("Не определяет пустой альбом: ", res)
	}
	res = IsEmptyAlbum("..//html/img/not_empty/")
	if res != false {
		t.Error("Не определяет не пустой альбом: ", res)
	}
	res = IsEmptyAlbum("..//html/img/invalid/")
	if res != true {
		t.Error("Ошибка чтения не обрабатывается: ", res)
	}
}