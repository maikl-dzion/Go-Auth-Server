package data

import (
	"database/sql"
	"fmt"
	"log"
	// "gitlab.basis-plus.ru/ibis/proto"
	"os"
	// "net/http"

	_ "github.com/lib/pq"

	"gitlab.tesonero-computers.ru/ibis/aaa/internal/model"
	// exlog "github.com/rs/zerolog/log"
)

var db *sql.DB

const AUTH_USERS_TABLE string = "public.auth_users"


func init() {

	var err error

	os.Setenv("DB_ADDRESS", "192.168.3.23")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_ROLE", "ibis")
	os.Setenv("DB_PASSWD", "ibis")
	os.Setenv("DB_DATABASE", "ibis")

	//address  := "192.168.3.23"
	//port     := "5432"
	//role     := "ibis"
	//passwd   := "ibis"
	//database := "ibis"

	address := os.Getenv("DB_ADDRESS")
	port    := os.Getenv("DB_PORT")
	role    := os.Getenv("DB_ROLE")
	passwd  := os.Getenv("DB_PASSWD")
	database := os.Getenv("DB_DATABASE")

	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		      role, passwd, address, port, database)

	db, err = sql.Open("postgres", dataSource)

	if err != nil {
		log.Fatal("db open error ...%v", err)
		return
	}

}

func RegisterUser(u *model.User) (id int, err error) {
	query := `INSERT INTO ` + AUTH_USERS_TABLE + `
		      (login, password) 
		      VALUES ($1, $2) RETURNING id
		     `
	row := db.QueryRow(query, u.Login, u.Password)
	if err = row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}


func AuthUser(login string) (password string, err error) {
	sql := `SELECT password FROM ` + AUTH_USERS_TABLE + ` WHERE login = "login"`

	row := db.QueryRow(sql)
	if err = row.Scan(&password); err != nil {
		return "-1", err
	}

	return password, nil
}



func GetUsers() ([]model.User) {

	sql := `SELECT  
               id,  
               login,
               password
            FROM 
               ` + AUTH_USERS_TABLE + `   
            ORDER BY id`


	rows, err := db.Query(sql)

	if err != nil {
		fmt.Println("users select(GetUsers)")
		log.Println(err)
	}


	defer rows.Close()

	users := []model.User{}

	// fmt.Println(rows)

	for rows.Next() {
		u := model.User{}
		err := rows.Scan(&u.Id, &u.Login, &u.Password)
		if err != nil {
			fmt.Println("scan select(GetUsers)")
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}


	return users
}



func GetUserById(id int) (*model.User, error) {

	var resp model.User
	sql := "SELECT id, login, password FROM " + AUTH_USERS_TABLE + " WHERE id=$1"
	err := db.QueryRow(sql, id).Scan(&resp.Id, &resp.Login, &resp.Password)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}



func FindLogin(login string) (*model.User, error) {
	var resp model.User
	sql := "SELECT id, login, password FROM " + AUTH_USERS_TABLE + " WHERE login=$1"
	err := db.QueryRow(sql, login).Scan(&resp.Id, &resp.Login, &resp.Password)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}



func DeleteUser(login string) (error) {

	sql := "DELETE FROM " + AUTH_USERS_TABLE + " WHERE login=$1"
	res, err := db.Exec(sql, login)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(res.RowsAffected())  // количество удаленных строк

	return err
}


func ChangeUserPasswd(ch model.ChangePasswd) (*model.User, error) {
	var resp model.User
	sql := "SELECT id, login, password FROM " + AUTH_USERS_TABLE + " WHERE login=$1"

	login := ch.Login
	oldPassword := ch.OldPassword
	newPassword := ch.NewPassword

	err := db.QueryRow(sql, login).Scan(&resp.Id, &resp.Login, &resp.Password)
	if err != nil {
		return &resp, err
	}

	if resp.Password != oldPassword {
		fmt.Println(" not password compare")
	} else {
		query := "UPDATE " + AUTH_USERS_TABLE + " SET password = $1 WHERE id = $2"
		_, err := db.Exec(query, newPassword, resp.Id)
		if err != nil {
			panic(err)
		}
	}


	u, err := GetUserById(resp.Id)
	if err != nil{
		panic(err)
	}

	return u, nil
}




// #################################
// ------  GRPC FUNC  IBIS ---------


func GetIbisAccount(authType int32, account int32) (*model.AccIbis, error) {

	var resp model.AccIbis

	//log.Println(account)
	//log.Println(authType)

	if authType == 0 {

		query := `SELECT 
                  
                   auth.client_login
				  ,auth.client_password
				  ,auth.server_login
				  ,auth.server_password

              FROM acc_ibis as ibis
			  INNER JOIN acc_auth as auth ON  ibis.account = auth.account_id
			  INNER JOIN acc_authtype as t ON  t.type_num = auth.authtype_id 
              WHERE ibis.account = $1 AND t.type_num = $2`

		row := db.QueryRow(query, account, authType)

		err := row.Scan(&resp.ClientLogin,
						&resp.ClientPassword,
						&resp.ServerLogin,
						&resp.ServerPassword)

		if err != nil {
			logPrint(err, "-####- Scan Error -####-")
			return &resp, err
		}


	} else {

		query := `SELECT 

				   auth.preshared_key

              FROM acc_ibis as ibis
			  INNER JOIN acc_auth as auth ON  ibis.account = auth.account_id
			  INNER JOIN acc_authtype as t ON  t.type_num = auth.authtype_id 
              WHERE ibis.account = $1 AND t.type_num = $2`

		row := db.QueryRow(query, account, authType)

		err := row.Scan(&resp.PresharedKey)

		if err != nil {
			logPrint(err, "-####- Scan Error -####-")
			return &resp, err
		}

	}


	// log.Println(resp)

	return &resp, nil
}




func IbisSessionAdd(uuid, currDate string, account int32) (int, error){

	id := 0

	query := `INSERT INTO acc_sessions
		      (account_id, uuid, opened_at) 
		      VALUES ($1, $2, $3) RETURNING id
		     `
	row := db.QueryRow(query, account, uuid, currDate)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil

}



func IbisSessionClose(uuid, currDate string, account int32) (int, error){

	id := 0

	query := `UPDATE acc_sessions SET closed_at = $1 WHERE uuid = $2`
	//db.QueryRow(query, currDate, uuid)
	_ , err := db.Exec(query, currDate, uuid)

	if err != nil {
        fmt.Println(err)
		return id, err
	}

	// fmt.Println(resp)

	return 1, nil

}


func checkError(err error) (bool) {
	check := true
	if err != nil {
		check = false
		panic(err)
	}
	return check
}

func logPrint(err error, message string) {
	log.Println(message)
	log.Println(err)
}



func GetAccSession(uuid string) (*model.AccSession, error) {

	var resp model.AccSession
	sql := "SELECT account_id, opened_at, uuid FROM acc_sessions WHERE uuid=$1 AND closed_at IS NULL "
	err := db.QueryRow(sql, uuid).Scan(&resp.AccountId, &resp.OpenedAt, &resp.Uuid)
	if err != nil {
		return nil, err
	}

	// fmt.Println(resp)
	return &resp, nil
}