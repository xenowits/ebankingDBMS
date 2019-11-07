
      // defer row.Close()

      // for rows.Next() {
      //
        // var (
        //   Id int
        //   username string
        //   password string
        //   role string
        //   availBalance int
        // )
      //
        // if err := rows.Scan(&Id, &username, &password, &role, &availBalance); err != nil {
        //   log.Fatal(err)
        // }
      //
        // log.Print(Id, username, password, role, availBalance)
      // }
      //
      // if !rows.NextResultSet() {
      //   log.Fatalf("expected more result sets : %v", rows.Err())
      // }




      // r.POST("/signin", handleSignIn)
      // stmt, err := db.Prepare("INSERT INTO customers (Id, Name, Address, AvailBalance, AcType) VALUES(?, ?, ?, ?, ?)")
      // if err != nil {
      //   log.Fatal(err)
      // }
      // stmt2, err := db.Prepare("DELETE FROM customers WHERE Id > 0")
      // _, err = stmt2.Exec()
      // if err != nil {
      //   log.Fatal(err)
      //   }
      // x := customer{356,"anshik","purnea",4356,"savings"}
      // _, err = stmt.Exec(x.Id, x.Name, x.Address, x.AvailBalance, x.AcType)
      // if err != nil {
      //   log.Fatal(err)
      //   }
