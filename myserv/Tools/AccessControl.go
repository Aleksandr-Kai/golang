package Tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	b64 "encoding/base64"
)

type User struct {
	Name		string	`json:"name"`
	Password	string	`json:"password"`
	Access		int		`json:"access_lvl"`
}

type Session struct {
	Account		User
	Expires		time.Time
}

type TokenHeader struct {
	Algorithm		string	`json:"alg"`
	TokenType		string	`json:"typ"`
}

type TokenData struct {
	Subject		string		`json:"sub"`
	Expiration	time.Time	`json:"exp"`
	UserName	string		`json:"name"`
}

var userList []User

var activeSessions map[string]Session

func FindUser(name string) (int, User){
	for index, u := range userList{
		if u.Name == name{
			return index, u
		}
	}
	return -1, User{}
}

func NewUser(login, password string, access int) bool{
	i, user := FindUser(login)
	if i != -1 {
		return false
	}
	user.Name = login
	user.Password = password
	user.Access = access
	userList = append(userList, user)
	return true
}

func RemoveUser(login string) bool{
	i, user := FindUser(login)
	if user.Name == "" {
		return false
	}
	copy(userList[i:], userList[i+1:])
	userList[len(userList)-1] = User{}
	userList = userList[:len(userList)-1]
	return true
}

func Init(path string){
	userList = make([]User, 0, 10)
	txt, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(txt, &userList)
		if err != nil {
			fmt.Println("[Auth Init] ", err.Error())
			return
		}
	}
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

func Login(name, password string) (token string){
	index, user := FindUser(name)
	if (index < 0) || (user.Password != password) {
		return ""
	}

	s := Session{user, time.Now().Add(48 * time.Hour)}
	activeSessions[token] = s
	return token
}

func GetToken(session Session) string{
	th, _ := json.Marshal(TokenHeader{"HS512", "JWT"})
	td, _ := json.Marshal(TokenData{"123", time.Unix(0, 0), session.Account.Name })
	token := b64.StdEncoding.EncodeToString(th) + "." + b64.StdEncoding.EncodeToString(td)
	token += "." + b64.StdEncoding.EncodeToString([]byte(token))
	return token
}

func ParseToken(token string) (Session, error){
	//fmt.Printf("[ParseToken] Token: %v\n", token)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		//fmt.Printf("[ParseToken] len(parts) != 3: %v\n", parts)
		return Session{}, errors.New("[ParseToken] Не правильный формат токена")
	}
	dec, err := b64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		//fmt.Printf("[ParseToken] Decode fail: %v\n", err.Error())
		return Session{}, err
	}
	partsDec := strings.Split(string(dec), ".")
	if len(partsDec) != 2 {
		//fmt.Printf("[ParseToken] len(partsDec) != 2: %v\n", partsDec)
		return Session{}, errors.New("[ParseToken] Не валидная подпись")
	}
	if (parts[0] != partsDec[0]) || (parts[1] != partsDec[1]) {
		//fmt.Printf("[ParseToken] parts != partsDec: %v != %v\n", parts, partsDec)
		return Session{}, errors.New("[ParseToken] Не валидная подпись")
	}
	//th, err := b64.StdEncoding.DecodeString(parts[0])
	td, err := b64.StdEncoding.DecodeString(parts[1])
	//fmt.Printf("[ParseToken] TokenData: %v\n", string(td))
	var tds TokenData
	err = json.Unmarshal(td, &tds)
	_, usr := FindUser(tds.UserName)
	s := Session{usr, tds.Expiration}

	return s, err
}