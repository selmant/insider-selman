package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"insider/pkg/models/message"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
	db, err := sqlx.Connect("mysql", mysqlURI)
	if err != nil {
		log.Fatal(err)
	}
	var msg message.Message
	err = db.Get(&msg, "SELECT * FROM messages WHERE id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	msg2 := message.Message{
		Content:        "123asdas",
		RecipientPhone: "123sda",
	}
	_, err = db.NamedExec("INSERT INTO messages (content, recipient_phone) VALUES (:content, :recipient_phone)", msg2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", msg.CreatedAt)
}
