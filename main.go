//go:generate go run github.com/steebchen/prisma-client-go db push

package main

import (
	"context"
	"cli-todo/db"
	"encoding/json"
	"fmt"
	// "github.com/tacheraSasi/prisma-go-demo.git/db"
)
//AddTodo
//getAllTodos
//getOneTodo
//updateTodo
//deleteTodo


func main() {
	client := db.NewClient()
	defer PrismaDisconnect(client)
	
	if err := Run(client); err != nil {
		panic(err)
	}

	//adding a todo
	err := AddTodo(client,"Create all the other CRUD functions")
	if err != nil{
	    return
	}

}

func PrismaConnect(client *db.PrismaClient){
	if err := client.Prisma.Connect(); err != nil{
		panic(err)
	}
}

func PrismaDisconnect(client *db.PrismaClient){
	if err := client.Prisma.Disconnect();err!=nil{ 
		panic(err)
	}
}

func Run(client *db.PrismaClient) error {
	PrismaConnect(client)
	ctx := context.Background()

	// create a post
	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Prisma!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("created post: %s\n", result)

	// find a single post

	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(createdPost.ID),
	).Exec(ctx)
	
	if err != nil {
		return err
	}

	result, _ = json.MarshalIndent(post, "", "  ")
	fmt.Printf("post: %s\n", result)

	desc, ok := post.Desc()
	if !ok {
		return fmt.Errorf("post's description is null")
	}

	fmt.Printf("The posts's description is: %s\n", desc)

	return nil
}
