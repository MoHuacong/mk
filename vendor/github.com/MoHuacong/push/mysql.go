package push

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	userTable = `CREATE TABLE user (
		id int AUTO_INCREMENT,
		user varchar(25) NOT NULL,
		pass varchar(25) NOT NULL,
		url varchar(50) NOT NULL,
		date TIMESTAMP default CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
)

// 用户表结构体
type User struct {
	Id int64 `db:"id"`
	User string `db:"user"`
	Pass string `db:"pass"`
	Url string `db:"url"`
	Date string `db:"DATE"`
}

type MySQL struct {
	err error
	db *sql.DB
}

func (my *MySQL) Open(args ...string) bool {
	if len(args) < 6 {
		return false
	}
	
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", args[0], args[1], args[2], args[3], args[4], args[5])
	my.db, my.err = sql.Open("mysql", dbDSN)
	
	if my.err != nil || my.db.Ping() != nil{
		return false
	}
	
	var i int = 0
	var rows *sql.Rows
	rows, my.err = my.db.Query("SELECT COUNT(*) FROM information_schema.TABLES WHERE table_name ='user';")
	
	if my.err != nil {
		return false
	}
	
	rows.Next()
	rows.Scan(&i)
	if i > 0 {
		return true
	}
	
	if _, my.err = my.db.Query(userTable); my.err != nil {
		return false
	}
	
	return true
}

func (my *MySQL) Result(rows *sql.Rows) interface{} {
	var i interface{} = 0
	rows.Next()
	rows.Scan(&i)
	return i
}

/* 读取数据 */
func (my *MySQL) Scan(rows *sql.Rows, user *User) error {
	return rows.Scan(&user.Id, &user.User, &user.Pass, &user.Url, &user.Date)
}

func (my *MySQL) Row(rows *sql.Rows, user *User) {
	for rows.Next() {
		my.Scan(rows, user)
	}
}