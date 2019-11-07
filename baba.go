package main

import (
  "database/sql"
  "github.com/gin-gonic/gin"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  // "fmt"
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
  r.LoadHTMLGlob("templates/*")
  // r.LoadHTMLFiles("pages/index.html")
  r.GET("/credit", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H {"message" : "piong"})
  })
  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H {
      "title" : "Ae raand",
    })
  })
  db, err := sql.Open("mysql", "abhi:abhi@/ebanking")
  if err != nil {
    log.Fatal(err)
  }
  stmt, err := db.Prepare("INSERT INTO customers (Id, Name, Address, AvailBalance, AcType) VALUES(?, ?, ?, ?, ?)")
  if err != nil {
    log.Fatal(err)
  }
  x := customer{356,"anshik","purnea",4356,"savings"}
  _, err = stmt.Exec(x.Id, x.Name, x.Address, x.AvailBalance, x.AcType)
  if err != nil {
    log.Fatal(err)
    }
  r.Run()
}
