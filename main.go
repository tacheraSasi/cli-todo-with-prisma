//go:generate go run github.com/steebchen/prisma-client-go db push

package main

import (
	"cli-todo/db"
	"context"
	"encoding/json"
	"fmt"
	"os"

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

	switch os.Args[1] {
	case "all":
		todos := GetAll(client)
		for i,todo := range todos {
			fmt.Println(i+1,". ",todo.Title)
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Missing the todo title arg")
			return
		}
		//adding a todo
		_,err := AddTodo(client,os.Args[2])
		if err != nil{
			return
		}
	default:
		fmt.Println("Invalid options\n--	all (for getting all todos)\n--		add <Todo-title> (for adding a todo)")
		
	}
	
	// if err := Run(client); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Todo was added",addedTodo.InnerTodo)

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
	// fmt.Println(post)
	
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
