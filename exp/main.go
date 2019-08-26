package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/samueldaviddelacruz/go-job-board/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlinfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()
	u := models.User{
		Name:  "A name here",
		Email: "email@email.com",
	}
	if err := us.Create(&u); err != nil {
		panic(err)
	}
	u.Email = "email2@email.com"
	if err := us.Update(&u); err != nil {
		panic(err)
	}
	ru, err := us.ByEmail(u.Email)
	if err != nil {
		panic(err)
	}
	fmt.Println(ru)
	/*
		us.DestructiveReset()

		u := models.User{
			Name:  "A name here",
			Email: "email2@email.com",
		}

		fmt.Println(u)
	*/
	/*
		u, err := us.ByID(1)
		fmt.Println(u)
	*/
	//us.db.LogMode(true)
	//us.db.AutoMigrate(&User{}, &Order{})
	/*
		var u User

		if err := db.Preload("Orders").First(&u).Error; err != nil {
			panic(err)
		}
		fmt.Println(u.Orders)
	*/
	/*
		createOrder(db, u, 1001, "some desc1")
		createOrder(db, u, 9999, "some desc2")
		createOrder(db, u, 100, "some desc3")
	*/
	//db.DropTableIfExists(&User{})

}
