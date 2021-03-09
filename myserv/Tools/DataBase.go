package Tools

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type DBAlbum struct{
	Name string
	Description string
	AccessLvl int
	Images []DBImage
}

type DBImage struct {
	Name string
	Description string
	AccessLvl int
}

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

	request := fmt.Sprintf("create table if not exists users(login char(15) not null, public_name char(15), password char(25), accesslvl int not null, active bool not null, primary key(login, public_name));")
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

	request = fmt.Sprintf("create table if not exists albums(name nvarchar(128) not null, description nvarchar(256), accesslvl int not null, primary key(name));")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create table 'albums':", err)
	}else{
		NamedMessage("SQL", "Create table 'albums': Ok")
	}

	request = fmt.Sprintf("create table if not exists photos(name nvarchar(128) not null, description nvarchar(256), accesslvl int not null, primary key(name));")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create table 'photos':", err)
	}else{
		NamedMessage("SQL", "Create table 'photos': Ok")
	}

	request = fmt.Sprintf("create table if not exists album_photo(album nvarchar(128) not null, photo nvarchar(128) not null, " +
		"foreign key(album) references albums(name), foreign key(photo) references photos(name));")
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create table 'album_photo':", err)
	}else{
		NamedMessage("SQL", "Create table 'album_photo': Ok")
	}
}

func DBCreateAlbum(name string, description string, accesslvl int){
	request := fmt.Sprintf("insert into albums values(N'%v', N'%v', %v);", name, description, accesslvl)
	_, err := db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Create album:", err)
		//return err
	}else{
		NamedMessage("SQL", "Create album: Ok")
	}
}

func DBGetAlbums(accesslvl int) []string{
	/*
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	request := fmt.Sprintf("select * from albums where accesslvl >= %v", )
	rows, err := db.Query(request)
	if err != nil{
		fmt.Println(err.Error())
		return User{}, err
	}
	defer rows.Close()
	post := &User{}
	rows.Next()
	err = rows.Scan(&post.Name, &post.PublicName, &post.Password, &post.Access, &post.Active)
	if err != nil{
		fmt.Println(err.Error())
	}
	return *post, err*/
	return nil
}

func DBPlaceImage(image DBImage, album string) error{
	request := fmt.Sprintf("insert into photos values('%v', '%v', %v);", image.Name, image.Description, image.AccessLvl)
	_, err := db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Не удалось добавить изображение ", image.Name, " в базу: ", err)
		return err
	}

	request = fmt.Sprintf("insert into album_photo values('%v', '%v');", album, image.Name)
	_, err = db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Не удалось добавить изображение ", image.Name, " в базу: ", err)
		return err
	}

	return nil
}

func DBDeleteImage(image DBImage) error{
	request := fmt.Sprintf("delete from photos where name='%v';", image.Name)
	_, err := db.Query(request)
	if err != nil{
		NamedMessage("SQL", "Не удалось удалить изображение ", image.Name, err)
		return err
	}

	DeleteImage(image.Name)
	return nil
}
/*
func DBGetImages(album string, accesslvl int) []DBImage{
	request := fmt.Sprintf("select * from albums al, images im,  where name='%v';", image.Name)
}
*/