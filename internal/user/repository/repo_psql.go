package repository

import (
	"database/sql"
	"fmt"
	"searcher/internal/user/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type psqlUserRepo struct {
	db *sqlx.DB
}

func NewPsqlUserRepo(db *sqlx.DB) IPsqlUserRepo {
	r := &psqlUserRepo{db: db}
	if err := r.initTables(); err != nil {
		panic(err)
	}
	return r
}

func (r *psqlUserRepo) GetCMSUsers() ([]entity.GetUser, error) {
	var users []entity.GetUser
	var ids []int

	str := `select id_user from user_info where proffession = 'cms'`

	if err := r.db.Select(&ids, str); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	for _, value := range ids {
		var user entity.GetUser

		str = `select * from users where id = $1`

		if err := r.db.Get(&user, str, value); err != nil {
			return nil, newm_helper.Trace(err, str)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *psqlUserRepo) GetUsersByGroupName(group string)([]entity.GetUser, error){
	var users []entity.GetUser
	var id []int

	str := `select id_user from user_info where group_name = $1`

	if err := r.db.Select(&id, str, group); err != nil{
		return nil, newm_helper.Trace(err, str)
	}

	str = `select * from users where id = $1`

	stmt, err := r.db.Preparex(str)
	if err != nil{
		return nil, newm_helper.Trace(err)
	}
	defer stmt.Close()

	for _, value := range id{
		var user entity.GetUser

		if err := stmt.Select(&user, value); err != nil{
			return nil, newm_helper.Trace(err, str)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *psqlUserRepo) GetAllAdmins()([]entity.GetUser, error) {
	var users []entity.GetUser

	str := `select * from users where role = 'admin'`

	if err := r.db.Select(&users, str); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	return users, nil
}

func (r *psqlUserRepo) UpdateUserAvatar(userId int, avatar string) error {
	str := `update users set avatar = $1 where id = $2`

	_, err := r.db.Exec(str, avatar, userId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) GetUserInfo(userId int) (entity.GetUserInfo, error) {
	var info entity.GetUserInfo

	str := `select * from user_info where id_user = $1`

	if err := r.db.Get(&info, str, userId); err != nil {
		if err != sql.ErrNoRows {
			return entity.GetUserInfo{}, newm_helper.Trace(err, str)
		}
	}

	return info, nil
}

func (r *psqlUserRepo) UpdateUserInfo(info entity.CreateUserInfo) error {
	str := `update user_info set name = $1, middle_name = $2, last_name = $3, course_number = $4, group_name = $5, 
	proffession = $6, proffession_number = $7 where id_user = $8`

	_, err := r.db.Exec(str, info.Name, info.MiddleName, info.LastName, info.CourseNumber, info.GroupName,
		info.Proffession, info.ProffessionNumber, info.IdUser)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) CreateUserInfo(info entity.CreateUserInfo) error {
	var id int

	str := "select id from user_info where id_user = $1"

	if err := r.db.QueryRow(str, info.IdUser).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return newm_helper.Trace(err, str)
		}
	}

	if id != 0 {
		return fmt.Errorf("info exists")
	}

	str = `insert into user_info(id_user, name, middle_name, last_name, course_number, group_name, proffession, proffession_number) 
	values($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(str, info.IdUser, info.Name, info.MiddleName, info.LastName,
		info.CourseNumber, info.GroupName, info.Proffession, info.ProffessionNumber)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) GetUserById(id int) (entity.GetUser, error) {
	var user entity.GetUser

	str := "select * from users where id = $1"

	if err := r.db.Get(&user, str, id); err != nil {
		return entity.GetUser{}, newm_helper.Trace(err, str)
	}

	return user, nil
}

func (r *psqlUserRepo) UpdateUserRole(role string, userId int) error {
	str := `update users set role = $1 where id = $2`

	_, err := r.db.Exec(str, role, userId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) UpdatePassword(user entity.UpdateUserPassword) error {
	str := `update users set password = $1 where id = $2`

	_, err := r.db.Exec(str, user.Password, user.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) UpdateUser(user entity.UpdateUser) error {
	str := `update users set email = $1, phone = $2 where id = $3`

	_, err := r.db.Exec(str, user.Email, user.Phone, user.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) GetUser(login, password string) (entity.GetUser, error) {
	var user entity.GetUser

	str := `select * from users where login = $1 and password = $2`

	if err := r.db.Get(&user, str, login, password); err != nil {
		if err == sql.ErrNoRows {
			return entity.GetUser{}, nil
		}
		return entity.GetUser{}, newm_helper.Trace(err, str)
	}

	return user, nil
}

func (r *psqlUserRepo) CreateUser(user entity.CreateUser) error {
	var id int

	str := "select id from users where login = $1"

	if err := r.db.QueryRow(str, user.Login).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return newm_helper.Trace(err, str)
		}
	}

	if id != 0 {
		return fmt.Errorf(fmt.Sprintf("user with login %s already exists", user.Login))
	}

	str = `insert into users(login, password, email, role, avatar, date_create) values($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(str, user.Login, user.Password, user.Email, user.Role, user.Avatar, user.DateCreate)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}
