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

type customer struct {
  Id int
  Name string
  Address string
  AvailBalance int
  AcType string
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

  r.GET("/credit", func(c *gin.Context) {

      cookie, err := c.Cookie("isLoggedIn")

      if err != nil {

        c.Redirect(http.StatusMovedPermanently, "/")

      } else if cookie == "true" {
        name, _ := c.Cookie("username")
        c.HTML(http.StatusOK, "credit.tmpl", gin.H {
          "name" : name,
        })
      }

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

    var (
      Idq int
      usernameq string
      passwordq string
      roleq string
      availBalanceq int
    )

    q := "SELECT * FROM customers WHERE username=?"
    row := db.QueryRow(q,username)

    if err := row.Scan(&Idq, &usernameq, &passwordq, &roleq, &availBalanceq); err != nil {
      log.Fatal(err)
    }

    log.Print(Idq, usernameq, passwordq, password, roleq, availBalanceq)


    err := bcrypt.CompareHashAndPassword([]byte(passwordq), []byte(password))

    if err == nil {

      fmt.Println("passwords match")

      c.SetCookie("username", username, 3600, "/", "localhost", false, true)
      c.SetCookie("isLoggedIn", "true", 3600, "/", "localhost", false, true)

      fmt.Println(c.Cookie("username"))
      fmt.Println(c.Cookie("isLoggedIn"))

      c.Redirect(http.StatusMovedPermanently, "/credit")

    } else {

      fmt.Println("passwords don't match")

      c.JSON(http.StatusOK, gin.H {
        "Request" : "Wrong password please try again",
      })

    }

  })

  r.Run(":3000")
}
