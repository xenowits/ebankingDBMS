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

  r.POST("/signup", func(c *gin.Context) {

      role := c.PostForm("role")
      username := c.PostForm("username")
      password := c.PostForm("password")
      availBalance := c.PostForm("availBalance")

      fmt.Println(role, username, password, availBalance)

      // stmt, _:= db.Prepare("SELECT * FROM customers")
      // p, err := stmt.Exec()

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

      // c.HTML(http.StatusOK, "credit.html", gin.H{
      //   // "badhiya" : "ha sb badhiya"
      // })
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


    err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordq))
    if err != nil {
      fmt.Println("passwords match")
      c.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
    } else {
      fmt.Println("no matching bosses")
      c.JSON(http.StatusOK, gin.H {
        "health" : "nice",
      })
      //c.Redirect(200, "https://www.netflix.com/in/")
    }

  })

  r.Run(":3000")
}
