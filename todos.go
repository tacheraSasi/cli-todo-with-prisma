package main

import (
	"cli-todo/db"
	"context"
)

func AddTodo(client *db.PrismaClient,title string)error{
	ctx := context.Background()

	addedTodo, err := client.Todo.CreateOne(
		db
	)
}