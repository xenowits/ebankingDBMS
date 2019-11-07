package main

import (
  "database/sql"
  "github.com/gin-gonic/gin"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  // "crypto/sha256"
  "log"
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
  r.GET("/credit", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H {"message" : "piong"})
  })
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

      var (
        Idq int
        usernameq string
        passwordq string
        roleq string
        availBalanceq int
      )

      q := "SELECT * FROM customers"
      row := db.QueryRow(q)

      if err := row.Scan(&Idq, &usernameq, &passwordq, &roleq, &availBalanceq); err != nil {
        log.Fatal(err)
      }

      log.Print(Idq, usernameq, passwordq, roleq, availBalanceq)

      c.JSON(http.StatusOK, gin.H{
        "badhiya" : "ha sb badhiya",
      })
  })

  r.Run()
}
