package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"grpc_ex.com/v1/messager"
)




func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := messager.NewMessageReceiverClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.EnableAccount(ctx, &messager.Account{Title : "Title", Code : "ADF"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf(r.GetResponseMessage())
}