package main 

import(
	"database/sql"
	"fmt"
	"log"
	"net"
	// "time"
	"context"
	// "net/http"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"

	"grpc_ex.com/v1/messager"
	
)


// account struct
type Accounts struct {
	Title     string `json:"title"`
	Code      string `json:"code"`
}	

//server struct 
type server struct {
	db      *sql.DB
	messager .UnimplementedMessageReceiverServer
}

func main() {
	fmt.Println("Server is started")
	var err error
	dbURL := fmt.Sprint("root:dev@tcp(127.0.0.1:3306)/dummy_account_service")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()	
	grpcServer := grpc.NewServer()
	s := &server{
		db: db,
	}
	err = s.dbPrepare()
	if err != nil {
		log.Println(err)
		return
	}
	messager.RegisterMessageReceiverServer(grpcServer, s)
	
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("Could not listen on port 9000:", err)
		return
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
		return
	}
	
}


func (s *server) dbPrepare() error {
	var err error
	_, err = s.db.Exec(`
    CREATE TABLE IF NOT EXISTS
    accounts (
        id int not null auto_increment,
        title tinytext not null,
        code tinytext not null,
        primary key (id)
    )
    `)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func (s *server) EnableAccount(ctx context.Context,account *messager.Account) (*messager.Response, error) {
	if account.Title == ""{
		return &messager.Response{ResponseMessage : "Account Title is required"}, nil
	}

	if account.Code== ""{
		return &messager.Response{ResponseMessage : "Account Code is required"}, nil
	}

	acc := &Accounts{account.Title, account.Code}
	err := s.storageAccountInsert(acc)
	if err != nil {
		log.Println(err)
		return &messager.Response{ResponseMessage : "Something went wrong"}, err
	}	
	return &messager.Response{ResponseMessage : "Account Created"}, nil
}

func (s *server) storageAccountInsert(a *Accounts) error {
	tx, err := s.db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	_, err = tx.Exec(`
    INSERT INTO accounts (
        title,
        code
    )
    VALUES (?, ?)
    `,
		a.Title, a.Code)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	tx.Commit()
	return nil
}