package Tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var tstList = []User{{"AKai","hitomi",0},{"guest", "", 10}, {"admin","admin",0}}
var tstNewUser = []User{{"AKai","hitomi",0},{"guest", "", 10}, {"admin","admin",0}, {"newuser","password",10}}


func TestInit(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestInit")
	Init("../authdata.json")
	if !reflect.DeepEqual(userList, tstList) {
		t.Error("[TestInit] Загрузка пользователей не выполнена: ", userList)
	}
	fmt.Printf("[TestInit] Загруженные данные: %v\n[TestInit] Эталонный массив: %v\n", userList, tstList)
}

func TestFindUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestFindUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	i, user := FindUser("admin")
	if (i != 2) || (user.Name != "admin") || (user.Password != "admin") || (user.Access != 0){
		t.Error("[TestFindUser] Пользователь не найден: ", i, "/", user)
	}
	fmt.Printf("[TestFindUser] Искомый пользователь: admin %v\n[TestFindUser] В массиве: %v\n", user, userList)
	i, user = FindUser("null")
	if i != -1 {
		t.Error("[TestFindUser] Функция вернула данные пользователя null: ", i, "/", user)
	}
	fmt.Printf("[TestFindUser] Искомый пользователь: null %v\n", user)
}

func TestNewUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestNewUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	NewUser("newuser", "password", 10)
	if !reflect.DeepEqual(userList, tstNewUser) {
		t.Error("[TestNewUser] Пользователь не добавлен: ", userList)
	}
	fmt.Printf("[TestNewUser] Рабочий массив: %v\n[TestNewUser] Эталонный массив: %v\n", userList, tstNewUser)
}

func TestRemoveUser(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestRemoveUser")
	userList = make([]User, 3)
	copy(userList, tstList)
	RemoveUser("guest")
	tst := []User{{"AKai","hitomi",0},{"admin","admin",0}}
	if !reflect.DeepEqual(userList, tst) {
		t.Error("[TestRemoveUser] Пользователь не удален: ", userList)
	}
	fmt.Printf("[TestRemoveUser] Рабочий массив: %v\n[TestRemoveUser] Эталонный массив: %v\n", userList, tst)
}

func TestSaveUsers(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestSaveUsers")
	tst_filename := "tst_users.json"
	userList = make([]User, 3)
	copy(userList, tstList)
	if !SaveUsers(tst_filename) {
		t.Error("[TestSaveUsers] Функция вернула false")
	}
	txt, err := ioutil.ReadFile(tst_filename)
	if err != nil {
		t.Error("[TestSaveUsers] Не удалось прочитать файл: ", err.Error())
	}
	tst := make([]User, 0, 10)
	err = json.Unmarshal(txt, &tst)
	if err != nil {
		t.Error("[TestSaveUsers] Не удалось распарсить: ", err.Error())
	}
	if !reflect.DeepEqual(tstList, tst) {
		t.Error("[TestSaveUsers] Тест не пройден")
	}
	fmt.Printf("[TestSaveUsers] Сохраняемый массив: %v\n[TestSaveUsers] Прочитанный массив: %v\n", userList, tst)
}

func TestGetToken(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestGetToken")
	exp := time.Unix(0,0)
	s := Session{User{"admin", "admin", 0}, exp}
	token := GetToken(s)
	fmt.Println("[TestGetToken] ", token)
}

func TestParseToken(t *testing.T) {
	fmt.Println("=============================================================================================================")
	fmt.Println("            TestGetToken")
	userList = make([]User, 3)
	copy(userList, tstList)
	fmt.Println("[TestParseToken] Валидный токен")
	val, err := ParseToken("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJleHAiOiIxOTcwLTAxLTAxVDAzOjAwOjAwKzAzOjAwIiwibmFtZSI6ImFkbWluIn0=.ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnpkV0lpT2lJeE1qTWlMQ0psZUhBaU9pSXhPVGN3TFRBeExUQXhWREF6T2pBd09qQXdLekF6T2pBd0lpd2libUZ0WlNJNkltRmtiV2x1SW4wPQ==")
	tst := Session{User{"admin", "admin", 0}, time.Unix(0,0)}
	if val != tst {
		fmt.Printf("[TestParseToken] Error: %v\n[TestParseToken] Получено: %v\n[TestParseToken] Ожидалось: %v\n", err, val, tst)
		t.Error("[TestParseToken] Валидный токен не расшифрован")
	}

	fmt.Println("[TestParseToken] Не валидный токен")
	val, err = ParseToken("")

	if err == nil {
		fmt.Printf("[TestParseToken] Error: %v\n[TestParseToken] Session: %v\n", err, val)
		t.Error("[TestParseToken] Не валидный токен не вернул ошибки")
	}
}