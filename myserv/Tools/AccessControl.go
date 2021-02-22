package Tools

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

var DefaultUser = "Guest"
var db *sql.DB

func GetUser(name string) (User, error){
	request := fmt.Sprintf("select * from users where login = '%v'", name)
	rows, err := db.Query(request)
	if err != nil{
		println(err.Error())
		return User{}, err
	}
	defer rows.Close()
	post := &User{}
	rows.Next()
	err = rows.Scan(&post.Name, &post.PublicName, &post.Password, &post.Access, &post.Active)
	if err != nil{
		println(err.Error())
	}
	return *post, err
}

func UpdateUser(user User) error{
	request := "update users set "
	if user.PublicName != ""{
		request += fmt.Sprintf("public_name='%v', ", user.PublicName)
	}
	if user.Password != ""{
		request += fmt.Sprintf("password='%v', ", user.Password)
	}
	request += fmt.Sprintf("active = %v  where login = '%v'", user.Active, user.Name)
	_, err := db.Query(request)
	if err != nil{
		println("[SQL] Update user: ", err.Error())
		return err
	}

	return nil
}

func NewUser(name, publicName, password string, access int) error{
	if publicName == ""{
		publicName = name
	}
	request := fmt.Sprintf("insert into users values ('%v', '%v', '%v', %v, 0);", name, publicName, password, access)
	_, err := db.Query(request)
	if err != nil{
		println("[SQL] Add user: ", err.Error())
		return err
	}
	return nil
}

func RemoveUser(name string) {
	request := fmt.Sprintf("delete from users where login='%v';", name)
	_, err := db.Query(request)
	if err != nil{
		println("[SQL] Add user: ", err.Error())
		return
	}
	fmt.Println("[SQL] User removed: ", name)
}

func Init(){
	dsn := "root@tcp(localhost:3306)/akaiphoto?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	db.SetConnMaxIdleTime(10)
	err = db.Ping()
	if err != nil{
		println("[SQL] Open DB: ", err.Error())
		return
	}
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