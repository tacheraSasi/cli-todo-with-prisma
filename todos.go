package main

import (
	"cli-todo/db"
	"context"
	"fmt"
	"log"
	"strconv"
)

type Todo = db.TodoModel // Corrected type alias for Todo

// AddTodo adds a new todo item to the database
func AddTodo(client *db.PrismaClient, title string) (*db.TodoModel, error) {
	PrismaConnect(client)
	ctx := context.Background()

	addedTodo, err := client.Todo.CreateOne(
		db.Todo.Title.Set(title),
	).Exec(ctx)
	if err != nil {
		log.Printf("Failed to create the todo: %v\n", err)
		return nil, err
	}

	fmt.Println("Todo was successfully created:", addedTodo)
	return addedTodo, nil
}

// GetAll fetches all todos from the database
func GetAll(client *db.PrismaClient) ([]db.TodoModel, error) {
	PrismaConnect(client)
	ctx := context.Background()

	todos, err := client.Todo.FindMany().Exec(ctx)
	if err != nil {
		log.Printf("Failed to fetch todos: %v\n", err)
		return nil, err
	}

	return todos, nil
}

// UpdateTodo updates an existing todo by its ID
func UpdateTodo(client *db.PrismaClient, id string, newTitle string) error {
	PrismaConnect(client)
	ctx := context.Background()

	numericalID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID provided: %v\n", err)
		return err
	}

	_, updateErr := client.Todo.FindUnique(
		db.Todo.ID.Equals(numericalID),
	).Update(
		db.Todo.Title.Set(newTitle),
	).Exec(ctx)
	if updateErr != nil {
		log.Printf("Failed to update the todo: %v\n", updateErr)
		return updateErr
	}

	fmt.Println("Todo was successfully updated")
	return nil
}

// DeleteTodo deletes a todo by its ID
func DeleteTodo(client *db.PrismaClient, id string) error {
	PrismaConnect(client)
	ctx := context.Background()

	numericalID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID provided: %v\n", err)
		return err
	}
	_, deleteErr := client.Todo.FindUnique(
		db.Todo.ID.Equals(numericalID),
	).Delete().Exec(ctx)

	if deleteErr != nil {
		log.Printf("Failed to delete the todo: %v\n", deleteErr)
		return deleteErr
	}

	fmt.Println("Todo was successfully deleted")
	return nil
}

// GetTodo fetches a single todo by its ID
func GetTodo(client *db.PrismaClient, id string) (*db.TodoModel, error) {
	PrismaConnect(client)
	ctx := context.Background()

	numericalID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID provided: %v\n", err)
		return nil, err
	}

	todo, fetchErr := client.Todo.FindUnique(
		db.Todo.ID.Equals(numericalID),
	).Exec(ctx)
	if fetchErr != nil {
		log.Printf("Failed to fetch the todo: %v\n", fetchErr)
		return nil, fetchErr
	}

	return todo, nil
}
