package Tools

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB
/*
type TItemUser struct {
	name 		string
	public_name string
	password 	string
	accesslvl 	int
	active 		bool
}
*/
func DBOpen() error{
	defer func() {
		if err := recover(); err != nil {
			log.Println("[DBOpen] panic occurred:", err)
		}
	}()
	dsn := "root@tcp(localhost:3306)/akaiphoto?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil{
		fmt.Println("[SQL] Open DB: ", err.Error())
		return err
	}
	db.SetConnMaxIdleTime(10)
	err = db.Ping()
	if err != nil{
		fmt.Println("[SQL] Ping DB: ", err.Error())
		return err
	}
	return nil
}

func DBInit() []error{
	defer func() {
		if err := recover(); err != nil {
			log.Println("[DBInit] panic occurred:", err)
		}
	}()
	errs := make([]error, 0, 4)
	request := fmt.Sprintf("create table users(login char(15) not null, public_name char(15), password char(25), accesslvl int not null, active bool not null, primary key(login, public_name));")
	_, err := db.Query(request)
	fmt.Print("[SQL] Create table 'users': ")
	if err != nil{
		fmt.Println(err.Error())
		errs = append(errs, err)
		//return err
	}else{
		fmt.Println("Ok")
	}


	request = fmt.Sprintf("insert into users values('admin', 'Администратор', '1qaz@WSX3edc$RFV', 0, 0);")
	_, err = db.Query(request)
	fmt.Print("[SQL] Create user 'admin': ")
	if err != nil{
		fmt.Println(err.Error())
		errs = append(errs, err)
		//return err
	}else{
		fmt.Println("Ok")
	}


	request = fmt.Sprintf("insert into users values('guest', 'Гость', '', 10, 0);")
	_, err = db.Query(request)
	fmt.Print("[SQL] Create user 'guest': ")
	if err != nil{
		fmt.Println(err.Error())
		errs = append(errs, err)
		//return err
	}else{
		fmt.Println("Ok")
	}

	return errs
}