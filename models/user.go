package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Model
	Name string `json:"name"`
	NickName string `json:"nickname"`
	Password string `json:"password"`
	Sex string `json:"sex"`
	State int `json:"state"`
	AuthToken string `json:"token"`
}

func GetUsers(pageNum int, pageSize int, maps interface {}) (users []User) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&users)

	return
}

func GetUserTotal(maps interface {}) (count int){
	db.Model(&User{}).Where(maps).Count(&count)

	return
}
func ExistUserByName(name string) bool {
	var user User
	db.Select("id").Where("name = ?", name).First(&user)
	if user.ID > 0 {
		return true
	}

	return false
}

func AddUser(name string, password string ,state int) bool{
	db.Create(&User {
		Name : name,
		Password:password,
		State : state,
	})

	return true
}
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (user *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}

func ExistUserByID(id uint) bool {
	var user User
	db.Select("id").Where("id = ?", id).First(&user)
	if user.ID > 0 {
		return true
	}

	return false
}

func DeleteUser(id uint) bool {
	db.Where("id = ?", id).Delete(&User{})

	return true
}

func EditUser(id uint, data interface {}) bool {
	db.Model(&User{}).Where("id = ?", id).Updates(data)

	return true
}
func GetUserByName(name string) (user User) {
	db.Where("name = ?", name).First(&user)

	return
}

func GetUserById(id uint) (user User) {
	db.Where("id = ?", id).First(&user)

	return
}