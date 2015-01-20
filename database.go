// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import (
    sql "database/sql"
//    "fmt"
    //Go-sqlite3 library. Documentation here: http://godoc.org/github.com/mattn/go-sqlite3
    _ "github.com/mattn/go-sqlite3"
    "log"
    "os"
    "bytes"
    "time"
)

var db *sql.DB
//
//func main() {
//    db = setupDatabase()
////    insertMessage("+12345678910",[]byte("asdfasfasdfasdfasfasasdf"),nil)
//    rows := getConversation("You")
//    defer rows.Close()
//    var msg string
//    var src string
//    for rows.Next() {
//        rows.Scan(&msg,&src)
//        log.Println(msg)
//        log.Println(src)
//    }
//}

// setupDatabase runs at start time to make sure the database exists and is ready to go
func setupDatabase() *sql.DB {
    var doesExist bool = true
    file, err := os.Open(".storage/textsecure.db")
    if err != nil {
        if os.IsNotExist(err) {
            doesExist = false
        } else {
            log.Fatal("Problem opening database file",err)
        }
    }
    file.Close()
    db, err := sql.Open("sqlite3", ".storage/textsecure.db")
    if err != nil {
        log.Fatal("Problem accessing the sqlite database:",err)
    }
    if doesExist == false {
        sqlStmt := `
        create table textsecure (id integer not null primary key, source blob, dest blob, timeRecieved timestamp, message blob, whichGroup blob);
        `
        _, err = db.Exec(sqlStmt)
        if err != nil {
            log.Fatal("Failed to intialize database", err, sqlStmt)
        }
    }
    return db
}

// Inserts a message object into the database
// INSERT INTO `textsecure`(`id`,`source`,`timeRecieved`,`message`,`whichGroup`) VALUES (1,NULL,NULL,NULL,NULL,NULL);
func insertMessage( source string, dest string, msg []byte, group []byte ) {
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    stmt, err := tx.Prepare("insert into textsecure(id, source, dest, timeRecieved, message, whichGroup) values(?, ?, ?, ?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    _, err = stmt.Exec(nil,source, dest, time.Now(),msg,group)
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
}



// getConversation gets the messages associated with a contact so they can be displayed in the message window
func getConversation(source string) *sql.Rows {
    var b bytes.Buffer
    b.WriteString("select message,source from textsecure where source = '")
    b.WriteString(source)
    b.WriteString("' or dest = '")
    b.WriteString(source)
    b.WriteString("'order by timeRecieved asc;")
    rows, err := db.Query(b.String())
    if err != nil {
        log.Fatal("Error retrieving messages from Sqlite database:", err)
    }
    return rows
}



