package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	// "net/url"
	// "os"
)

var (
	Idq           int //for the id
	usernameq     string
	passwordq     string
	roleq         string
	availBalanceq int
)

type tranStruct struct {
	Id                        int
	T_datetime                string
	User_credited             string
	User_debited              string
	CreditedUser_finalBalance int64
	DebitedUser_finalBalance  int64
	T_amount                  int64
}

var (
	Idt                        int
	T_datetimet                string
	User_creditedt             string
	User_debitedt              string
	CreditedUser_finalBalancet int64
	DebitedUser_finalBalancet  int64
	T_amountt                  int64
)

type miniStatement struct {
	IssuedTo     string
	Transactions []tranStruct
}

func main() {
	r := gin.Default()

	db, err := sql.Open("mysql", "abhishek:abhishek@/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r.LoadHTMLGlob("templates/*")
	// r.LoadHTMLFiles("pages/index.html")

	r.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Ebanking",
		})

	})

	r.GET("/tStatement", func(c *gin.Context) {

		cookie1, err1 := c.Cookie("isLoggedIn")
		cookie2, err2 := c.Cookie("role")
		cookie3, err3 := c.Cookie("username")

		if err1 != nil || err2 != nil || err3 != nil || cookie1 != "true" || cookie2 != "user" {
			fmt.Println(err1, err2, err3)
			c.Redirect(307, "/")
			return
		}

		fmt.Println(cookie1, cookie2, cookie3)

		stmt, err := db.PrepareContext(c, "SELECT * FROM transactions WHERE user_credited=? OR user_debited=?")
		if err != nil {
			c.String(http.StatusBadRequest, "MiniStatement cannot be processed right now!!")
			return
		}

		rows, err := stmt.QueryContext(c, cookie3, cookie3)

		if err != nil {
			c.String(http.StatusBadRequest, "Sorry for inconvenience!!")
			return
		}
		defer stmt.Close()
		defer rows.Close()

		var trans []tranStruct

		for rows.Next() {
			if err := rows.Scan(&Idt, &T_datetimet, &User_creditedt, &User_debitedt, &CreditedUser_finalBalancet, &DebitedUser_finalBalancet, &T_amountt); err != nil {
				log.Fatal(err)
			}
			temp := tranStruct{
				Id:                        Idt,
				T_datetime:                T_datetimet,
				User_credited:             User_creditedt,
				User_debited:              User_debitedt,
				CreditedUser_finalBalance: CreditedUser_finalBalancet,
				DebitedUser_finalBalance:  DebitedUser_finalBalancet,
				T_amount:                  T_amountt,
			}
			trans = append(trans, temp)

		}
		data := miniStatement{
			IssuedTo:     cookie3,
			Transactions: trans,
		}

		c.HTML(http.StatusOK, "transactions.tmpl", data)

	})

	r.POST("/logout", func(c *gin.Context) {

		c.SetCookie("username", "", -8, "/", "localhost", false, true)
		c.SetCookie("isLoggedIn", "", -8, "/", "localhost", false, true)
		c.SetCookie("role", "", -8, "/", "localhost", false, true)

		c.Redirect(302, "/")

	})

	r.GET("/admin", func(c *gin.Context) {

		cookie1, err1 := c.Cookie("isLoggedIn")
		cookie2, err2 := c.Cookie("role")

		if err1 != nil && err2 != nil {

			c.Redirect(http.StatusMovedPermanently, "/")

		} else if cookie1 == "true" && cookie2 == "admin" {

			name, _ := c.Cookie("username")
			c.HTML(http.StatusOK, "credit.tmpl", gin.H{
				"name": name,
			})
		}

		c.Redirect(http.StatusMovedPermanently, "/")

	})

	r.GET("/user", func(c *gin.Context) {

		cookie, _ := c.Cookie("isLoggedIn")
		cookie2, _ := c.Cookie("username")
		cookie3, _ := c.Cookie("role")

		if cookie == "true" && cookie3 == "user" {
			c.HTML(http.StatusOK, "user.tmpl", gin.H{
				"name": cookie2,
			})
		} else {
			fmt.Println("some error occured")
			c.Redirect(http.StatusMovedPermanently, "/")
		}

	})

	r.POST("/transact", func(c *gin.Context) {

		cookie1, err1 := c.Cookie("isLoggedIn")
		cookie2, err2 := c.Cookie("username")

		if err1 != nil || err2 != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "sorry the transaction could not be processed",
			})

		}

		type Trn struct {
			To  string `json:"name" binding:"required"`
			Amt int64  `json:"amount" binding:"required"`
		}

		var json Trn

		if err := c.ShouldBindJSON(&json); err != nil {

			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})

			return
		}

		fmt.Println(reflect.TypeOf(json.To), json.To, cookie1, cookie2)
		fmt.Println(reflect.TypeOf(json.Amt), json.Amt)

		// check if destination user exists, if exists , save his availBalance value in a variable

		q := "SELECT username, availBalance, role FROM customers WHERE username=?"
		row := db.QueryRow(q, json.To)

		err := row.Scan(&usernameq, &availBalanceq, &roleq)

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username doesn't exists",
			})

			return
		}

		if cookie2 == json.To {

			//😨 errors

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "whoa there, u want to send yourself this money? We don't this here!!",
			})

			return

		}

		if roleq == "admin" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ohoh!!U can't send money to an admin user",
			})

			return

		}

		var MyAvailBalance int64 = 0

		q = "SELECT availBalance FROM customers WHERE username=?"

		row = db.QueryRow(q, cookie2)

		err = row.Scan(&MyAvailBalance)

		if MyAvailBalance < json.Amt {
			//😨 errors

			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "insufficient funds",
				"whattodo": "go to the bank and deposit some money",
			})

			return

		}

		//real code for transaction

		t := time.Now().String()

		stmt, err := db.Prepare("INSERT INTO transactions (t_datetime, user_credited, user_debited, creditedUser_finalBalance, debitedUser_finalBalance, t_amount) VALUES (?, ?, ?, ?, ?, ?)")

		if err != nil {
			log.Print(err)
			return
		}
		//prepared statements take up server resources and should
		//be closed after use
		defer stmt.Close()

		MyNewBalance := MyAvailBalance - json.Amt
		HisNewBalance := int64(availBalanceq) + json.Amt

		_, err = stmt.Exec(t, usernameq, cookie2, HisNewBalance, MyNewBalance, json.Amt)

		if err != nil {

			log.Print(err)
			c.JSON(http.StatusBadRequest, gin.H{

				"data": "Some error occured",
			})
			return

		}

		//Now update availBalance of debiter

		stmt, _ = db.Prepare("UPDATE customers SET availBalance = ? WHERE username = ? ")
		_, err1 = stmt.Exec(MyNewBalance, cookie2)

		stmt, _ = db.Prepare("UPDATE customers SET availBalance = ? WHERE username = ? ")
		_, err2 = stmt.Exec(HisNewBalance, json.To)

		if err1 != nil || err2 != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unable to proceed with the transaction, try again later",
			})

			return
		}

		defer stmt.Close()

		//😀 done
		c.JSON(http.StatusOK, gin.H{

			"data": "congratulations@TRANSACTION was successful",
		})

	})

	r.POST("/signup", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")
		availBalance := c.PostForm("availBalance")

		fmt.Println(username, password, availBalance)

		stmt, err := db.Prepare("INSERT INTO customers (username, password, role, availBalance) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
			return
		}

		defer stmt.Close()

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(username, string(hashedPassword), "user", availBalance)
		if err != nil {
			log.Fatal(err)
		}

		c.Redirect(http.StatusMovedPermanently, "/")

	})

	r.POST("/signin", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		q := "SELECT * FROM customers WHERE username=?"
		row := db.QueryRow(q, username)

		err := row.Scan(&Idq, &usernameq, &passwordq, &roleq, &availBalanceq)

		log.Print(Idq, usernameq, passwordq, password, roleq, availBalanceq)

		log.Print(err)

		if err == nil {

			err1 := bcrypt.CompareHashAndPassword([]byte(passwordq), []byte(password))

			if err1 == nil {

				fmt.Println("passwords match")

				c.SetCookie("username", username, 600, "/", "localhost", false, true)
				c.SetCookie("isLoggedIn", "true", 600, "/", "localhost", false, true)
				c.SetCookie("role", roleq, 600, "/", "localhost", false, true)

				fmt.Println(c.Cookie("username"))
				fmt.Println(c.Cookie("isLoggedIn"))

				if roleq == "admin" {

					c.Redirect(302, "/admin")

				} else if roleq == "user" {

					c.Redirect(302, "/user")

				}
			} else {

				c.JSON(401, gin.H{
					"Request":    "Passwords don't match",
					"what to do": "try again",
				})

			}
		} else {

			fmt.Println("username doesn't exist or some other error")

			c.JSON(http.StatusOK, gin.H{
				"Request": "Wrong username or password or some other error",
			})
		}
	})

	r.Run(":3000")
}
