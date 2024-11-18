package main

import (
	"cli-todo/db"
	"context"
	"fmt"
	"log"
)

func AddTodo(client *db.PrismaClient,title string)error{
	ctx := context.Background()

	addedTodo,err := client.Todo.CreateOne(db.Todo.Title.Set(title)).Exec(ctx)
	if err != nil{
		log.Fatal("Failed to create the todo",err);
		return err
	}

	fmt.Println("Todo was created")
	fmt.Println(addedTodo)
}