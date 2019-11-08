package main

import (
  "database/sql"
  "github.com/gin-gonic/gin"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "log"
  // "os"
)

var (
  Idq int
  usernameq string
  passwordq string
  roleq string
  availBalanceq int
)

type user struct {
  Id int
  Username string
  Password string
  AvailBalance int
  Role string
}

type miniStatement struct {
  IssuedTo string
  Transactions []user
}

func main() {
  r := gin.Default()

  db, err := sql.Open("mysql", "abhi:abhi@/ebanking")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  r.LoadHTMLGlob("templates/*")
  // r.LoadHTMLFiles("pages/index.html")

  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H {
      "title" : "Ebanking",
    })
  })

  r.GET("/yourMiniStatement", func(c *gin.Context){

    rows, err := db.Query("SELECT * FROM customers")

    if err != nil {
      log.Fatal(err)
    }

    defer rows.Close()

    var arr []user

    for rows.Next() {
      if err := rows.Scan(&Idq, &usernameq, &passwordq, &roleq, &availBalanceq); err != nil {
        log.Fatal(err)
      }
      temp := user{
        Id : Idq,
        Username : usernameq,
        Password : passwordq,
        Role : roleq,
        AvailBalance : availBalanceq,
      }
      arr = append(arr,temp)
      // log.Println(Idq, usernameq, passwordq, roleq, availBalanceq)
    }

    sessionUname,_ := c.Cookie("username")

    data := miniStatement {
      IssuedTo : sessionUname,
      Transactions : arr,
    }

    fmt.Println("ram", data.Transactions[0])

    // var x interface{}

    // fmt.Println(data)

    // x = data

    // fmt.Println("guptaji/n",x)

    // c.String(http.StatusOK, "success", data)
    c.HTML(http.StatusOK, "transactions.tmpl", data)

  })

  r.GET("/admin", func(c *gin.Context) {

      cookie1, err1 := c.Cookie("isLoggedIn")
      cookie2, err2 := c.Cookie("role")

      if err1 != nil && err2 != nil {

        c.Redirect(http.StatusMovedPermanently, "/")

      } else if cookie1 == "true" && cookie2 == "admin" {

        name, _ := c.Cookie("username")
        c.HTML(http.StatusOK, "credit.tmpl", gin.H {
          "name" : name,
        })
      }

      c.Redirect(http.StatusMovedPermanently, "/")

  })

  r.POST("/signup", func(c *gin.Context) {

      role := c.PostForm("role")
      username := c.PostForm("username")
      password := c.PostForm("password")
      availBalance := c.PostForm("availBalance")

      fmt.Println(role, username, password, availBalance)

      stmt, err := db.Prepare("INSERT INTO customers (username, password, role, availBalance) VALUES (?, ?, ?, ?)")
      if err != nil {
        log.Fatal(err)
      }

      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
      if err != nil {
        log.Fatal(err)
      }

      _, err = stmt.Exec(username, string(hashedPassword), role, availBalance)
      if err != nil {
        log.Fatal(err)
      }

      c.Redirect(http.StatusMovedPermanently, "/")

  })

  r.POST("/signin", func(c *gin.Context) {

    username := c.PostForm("username")
    password := c.PostForm("password")

    q := "SELECT * FROM customers WHERE username=?"
    row := db.QueryRow(q,username)

    if err := row.Scan(&Idq, &usernameq, &passwordq, &roleq, &availBalanceq); err != nil {
      log.Fatal(err)
    }

    log.Print(Idq, usernameq, passwordq, password, roleq, availBalanceq)


    err := bcrypt.CompareHashAndPassword([]byte(passwordq), []byte(password))

    if err == nil {

      fmt.Println("passwords match")

      c.SetCookie("username", username, 120, "/", "localhost", false, true)
      c.SetCookie("isLoggedIn", "true", 120, "/", "localhost", false, true)
      c.SetCookie("role", roleq, 120, "/", "localhost", false, true)

      fmt.Println(c.Cookie("username"))
      fmt.Println(c.Cookie("isLoggedIn"))

      if roleq == "admin" {
        c.Redirect(http.StatusMovedPermanently, "/admin")
      } else if roleq == "customer" {
        c.Redirect(http.StatusMovedPermanently, "/user")
      }

    } else {

      fmt.Println("passwords don't match")

      c.JSON(http.StatusOK, gin.H {
        "Request" : "Wrong password please try again",
      })

    }

  })

  r.Run(":3000")
}
