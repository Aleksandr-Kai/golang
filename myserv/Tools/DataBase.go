package Tools

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DBOpen() error{
	defer func() {
		if err := recover(); err != nil {
			Log("Panic occurred", err)
		}
	}()
	dsn := "root@tcp(localhost:3306)/akaiphoto?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil{
		Log("SQL Open DataBase", err)
		return err
	}
	db.SetConnMaxIdleTime(10)
	err = db.Ping()
	if err != nil{
		Log("SQL Ping DataBase", err)
		return err
	}
	return nil
}

func DBInit(){
	defer func() {
		if err := recover(); err != nil {
			Log("Panic occurred", err)
		}
	}()

	request := fmt.Sprintf("create table users(login char(15) not null, public_name char(15), password char(25), accesslvl int not null, active bool not null, primary key(login, public_name));")
	_, err := db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create table 'users':", err)
	}else{
		NamedMessage("SQL", "Create table 'users': Ok")
	}


	request = fmt.Sprintf("insert into users values('admin', 'Администратор', '1qaz@WSX3edc$RFV', 0, 0);")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create user 'admin':", err)
	}else{
		NamedMessage("SQL", "Create user 'admin': Ok")
	}


	request = fmt.Sprintf("insert into users values('guest', 'Гость', '', 10, 0);")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create user 'guest':", err)
		//return err
	}else{
		NamedMessage("SQL", "Create user 'guest': Ok")
	}

	request = fmt.Sprintf("create table albums(id integer not null, name text not null, description text, accesslvl int not null, primary key(id));")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create table 'albums':", err)
	}else{
		NamedMessage("SQL", "Create table 'albums': Ok")
	}
}