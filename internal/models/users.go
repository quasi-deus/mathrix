package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	UserID int
	Name string
	URN string
	Phone int64
	College string
	Dept string
	Year int
	Degree string
	Email string
	HashedPassword []byte
	Authority bool
	Created time.Time
}

type UserModel struct{
	DB *sql.DB
}

func (m *UserModel) Insert(name, urn string, phone int64, college, dept string, year int, degree string, email, password string, authority bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, urn, phone, college, dept, year, degree, email, hashed_password, authority, created) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW()::timestamp(0))`
	_,err=m.DB.Exec(stmt, name,urn, phone, college, dept, year, degree,  email, string(hashedPassword),authority)
	if err!=nil{
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError){
			if pSQLError.Code == "23505" && strings.Contains(pSQLError.Message, "users_uc_email"){
			return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(urn, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt:="SELECT userid, hashed_password FROM users where urn=$1"

	err:=m.DB.QueryRow(stmt,urn).Scan(&id, &hashedPassword)
	if err!=nil{
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrInvalidCredentials
		}else {
			return 0, err
		}
	}
	
	err=bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err!=nil{
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
		return 0, ErrInvalidCredentials
	}else{
		return 0,err
		}
	}
	return id, nil
}

func (m *UserModel) Authorized(id int) (bool, error){
	var authority bool

	stmt:="SELECT authority FROM users WHERE userid=$1"

	err:=m.DB.QueryRow(stmt, id).Scan(&authority)

	return authority, err
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt:="SELECT EXISTS(SELECT true FROM users WHERE userid=$1)"

	err:=m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (m *UserModel) ViewAll() ([]*User, error) {
	users:= []*User{}
	rows, err:= m.DB.Query("SELECT  userid,name, urn, phone, college, dept, year, degree, email, authority, created FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		s:=&User{}
		err=rows.Scan(&s.UserID, &s.Name, &s.URN, &s.Phone, &s.College, &s.Dept, &s.Year, &s.Degree, &s.Email, &s.Authority, &s.Created)
	if err != nil {
		return nil, err
	}
	users=append(users,s)
	}
	if err=rows.Err();err!=nil{
		return nil, err
	}

	return users, nil
}

func (m *UserModel) Get(eventid int) (*User, error) {
	s := &User{}

	err:= m.DB.QueryRow("SELECT  userid,name, urn, phone, college, dept, year, degree, email, authority, created FROM users WHERE userid=$1",eventid).Scan(&s.UserID, &s.Name, &s.URN, &s.Phone, &s.College,&s.Dept,&s.Year,&s. Degree,&s.Email,&s.Authority,&s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *UserModel) Update(userid int, name, urn string, phone int64, college, dept string, year int, degree string, email, password string, authority bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `UPDATE users SET name=$2, urn=$3, phone=$4, college=$5, dept=$6, year=$7, degree=$8, email=$9, hashed_password=$10, authority=$11 WHERE userid=$1 RETURNING userid`
	_,err=m.DB.Exec(stmt,userid,name,urn, phone, college, dept, year, degree,  email, string(hashedPassword),authority)
	if err!=nil{
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError){
			if pSQLError.Code == "23505" && strings.Contains(pSQLError.Message, "users_uc_email"){
			return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Delete(userid int)error{
	stmt:="DELETE FROM users WHERE userid=$1"
	_,err:=m.DB.Exec(stmt,userid)
	if err!=nil{
		return err
	}
	return nil
}

func (m *UserModel) RegisterEvent(userid,eventid int)error{
	stmt:=`INSERT INTO eventlist (userid, eventid) VALUES($1,$2)`
	_,err:=m.DB.Exec(stmt,userid,eventid)
	if err!=nil{
		return err
	}
	return nil
}

//SELECT eventlist.userid,events.eventname FROM eventlist,events WHERE eventlist.userid=14 AND eventlist.eventid=events.eventid GROUP BY eventlist.userid,events.eventid;
