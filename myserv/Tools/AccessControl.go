package Tools

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type User struct {
	Name		string	`json:"name"`
	PublicName	string	`json:"public_name"`
	Password	string	`json:"password"`
	Access		int		`json:"access_lvl"`
	Active		bool	`json:"-"`
}

type TokenHeader struct {
	Algorithm		string	`json:"alg"`
	TokenType		string	`json:"typ"`
}

type TokenData struct {
	Subject		string		`json:"sub"`
	Expiration	time.Time	`json:"exp"`
}

var userList []User
var DefaultUser = "Guest"

func GetUser(name string) (User, error){
	for _, u := range userList{
		if u.Name == name{
			return u, nil
		}
	}
	return User{}, errors.New("Пользователь [" + name + "] не найден")
}

func UpdateUser(user User) error{
	for i, u := range userList{
		if u.Name == user.Name{
			userList[i] = user
			return nil
		}
	}
	return errors.New("Пользователь [" + user.Name + "] не найден")
}

func NewUser(name, publicName, password string, access int) error{
	if (len(password) < 8) && (name != DefaultUser){
		return errors.New("Пароль слишком короткий")
	}
	_, err := GetUser(name)
	if err == nil {
		return errors.New("Пользователь [" + name + "] уже существует")
	}
	for _, u := range userList{
		if (u.PublicName == publicName) || (u.Name == publicName){
			return errors.New("Имя [" + publicName + "] занято другим пользователем")
		}
	}
	userList = append(userList, User{name, publicName, password, access, false})
	return nil
}

func RemoveUser(name string) {
	for i, u := range userList{
		if u.Name == name{
			copy(userList[i:], userList[i+1:])
			userList[len(userList)-1] = User{}
			userList = userList[:len(userList)-1]
			return
		}
	}
}

func Init(path string){
	userList = make([]User, 0, 10)
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		println("[Auth Init] ", err.Error())
		return
	}
	err = json.Unmarshal(txt, &userList)
	if err != nil {
		println("[Auth Init] ", err.Error())
		return
	}
	fmt.Printf("[Auth Init] %v\n", userList)
}

func SaveUsers(path string) bool {
	txt, err := json.Marshal(userList)
	if err != nil {
		fmt.Printf("[SaveUsers] %v\n", err.Error())
		return false
	}
	err = ioutil.WriteFile(path, txt, os.ModePerm)
	if err != nil {
		fmt.Printf("[SaveUsers] %v\n", err.Error())
		return false
	}
	return true
}

func Login(name, password string) (token string, err error){
	user, err := GetUser(name)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("Неверный пароль")
	}
	user.Active = true
	UpdateUser(user)
	return GetToken(user, time.Now().Add(168 * time.Hour)), nil
}

func Logout(name string){
	user, err := GetUser(name)
	if err != nil {
		println("[Logout] ", err.Error())
		return
	}
	user.Active = false
	err = UpdateUser(user)
	if err != nil {
		println("[Logout] ", err.Error())
	}
}

func GetToken(user User, expiration time.Time) string{
	th, _ := json.Marshal(TokenHeader{"HS512", "JWT"})
	td, _ := json.Marshal(TokenData{user.Name, expiration })
	token := b64.StdEncoding.EncodeToString(th) + "." + b64.StdEncoding.EncodeToString(td)
	token += "." + b64.StdEncoding.EncodeToString([]byte(token))
	return token
}

func ParseToken(token string) (User, error){
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return User{}, errors.New("[ParseToken] Не правильный формат токена")
	}
	dec, err := b64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return User{}, err
	}
	partsDec := strings.Split(string(dec), ".")
	if (len(partsDec) != 2) || (parts[0] != partsDec[0]) || (parts[1] != partsDec[1]) {
		return User{}, errors.New("Не валидная подпись")
	}
	//th, err := b64.StdEncoding.DecodeString(parts[0])
	td, err := b64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return User{}, err
	}
	var tds TokenData
	err = json.Unmarshal(td, &tds)
	if err != nil {
		return User{}, err
	}
	if tds.Expiration.Before(time.Now()) {
		return User{}, errors.New("[ParseToken] Токен просрочен")
	}
	usr, err := GetUser(tds.Subject)
	if err != nil {
		return User{}, errors.New("[ParseToken] Токен не соответствует ни одному пользователю")
	}
	if !usr.Active {
		fmt.Printf("[ParseToken] %v\n", usr)
		return User{}, errors.New("Токен не актуален, выполнен выход")
	}
	return usr, nil
}