package main

import (
	"cli-todo/db"
	"context"
	"fmt"
	"log"

)

type Todo *db.InnerTodo //NOTICE:Will fix this Later

func AddTodo(client *db.PrismaClient,title string)(*db.TodoModel,error){
	PrismaConnect(client)
	ctx := context.Background()

	addedTodo,err := client.Todo.CreateOne(db.Todo.Title.Set(title)).Exec(ctx)
	if err != nil{
		log.Fatal("Failed to create the todo",err);
		return nil,err
	}

	fmt.Println("Todo was created")
	// fmt.Println(addedTodo.InnerTodo)
	return addedTodo,nil 
}
func GetAll(client *db.PrismaClient)[]db.TodoModel{
	PrismaConnect(client)
	ctx := context.Background()

	//getting all todos from the db
	todos,err := client.Todo.FindMany().Exec(ctx)
	if err != nil{
		log.Fatal("Failed to fetch todos:",err)
		return nil
	}
	return todos
}

func UpdateTodo(){}

func DeleteTodo(){}

func GetTodo(){}

