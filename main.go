//go:generate go run github.com/steebchen/prisma-client-go db push

package main

import (
	"bufio"
	"cli-todo/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var IsPrismaConnected bool = false

func main() {
	client := db.NewClient()
	defer PrismaDisconnect(client)

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Use one of the following:")
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "all":
		todos, err := GetAll(client)
		if err != nil {
			log.Printf("Error fetching todos: %v\n", err)
			return
		}

		fmt.Println("ID   Title")
		fmt.Println("-----------")
		for _, todo := range todos {
			fmt.Printf("%d   %s\n", todo.ID, todo.Title)
		}
	case "add":
		// if len(os.Args) < 3 {
		// 	fmt.Println("Missing the todo title argument")
		// 	return
		// }
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the title: ")
		title,readerErr := reader.ReadString('\n')
		if readerErr != nil{
			fmt.Println("Something went wrong",readerErr)
		}

		// title := os.Args[2]
		_, err := AddTodo(client, title)
		if err != nil {
			log.Printf("Error adding todo: %v\n", err)
			return
		}
		fmt.Println("Todo added successfully!")
	case "get-todo":
		if len(os.Args) < 3 {
			fmt.Println("Missing the todo ID argument")
			return
		}

		id := os.Args[2]
		todo, err := GetTodo(client, id)
		if err != nil {
			log.Printf("Error fetching todo: %v\n", err)
			return
		}
		if todo == nil {
			fmt.Println("Todo not found")
			return
		}

		fmt.Println("Todo Found:")
		fmt.Printf("ID: %d\nTitle: %s\n", todo.ID, todo.Title)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Missing the todo ID argument")
			return
		}

		id := os.Args[2]
		delErr := DeleteTodo(client,id)
		if delErr != nil {
			log.Printf("Error deleting todo: %v\n", delErr)
			return
		}
		fmt.Printf("Task with id %s Was deleted\n",id)
		return

	default:
		fmt.Println("Invalid command.")
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  all           Get all todos")
	fmt.Println("  add <title>   Add a new todo")
	fmt.Println("  get-todo <id> Get a specific todo by ID")
}

// PrismaConnect establishes a connection to the Prisma client
func PrismaConnect(client *db.PrismaClient) {
	if IsPrismaConnected {
		return // Avoid reconnecting if already connected
	}

	if err := client.Prisma.Connect(); err != nil {
		panic(fmt.Errorf("failed to connect to Prisma: %w", err))
	}

	IsPrismaConnected = true
}

// PrismaDisconnect disconnects the Prisma client
func PrismaDisconnect(client *db.PrismaClient) {
	if !IsPrismaConnected {
		log.Println("Prisma is not connected. Skipping disconnect.")
		return
	}

	if err := client.Prisma.Disconnect(); err != nil {
		panic(fmt.Errorf("failed to disconnect from Prisma: %w", err))
	}

	IsPrismaConnected = false
}

// Run demonstrates CRUD operations (example function)
func Run(client *db.PrismaClient) error {
	PrismaConnect(client)
	ctx := context.Background()

	// Create a post
	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Prisma!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("Created post: %s\n", result)

	// Find a single post
	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(createdPost.ID),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error finding post: %w", err)
	}

	result, _ = json.MarshalIndent(post, "", "  ")
	fmt.Printf("Post: %s\n", result)

	desc, ok := post.Desc()
	if !ok {
		return fmt.Errorf("post's description is null")
	}

	fmt.Printf("The post's description is: %s\n", desc)
	return nil
}
