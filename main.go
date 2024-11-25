//go:generate go run github.com/steebchen/prisma-client-go db push

package main

import (
	"bufio"
	"cli-todo/db"
	"fmt"
	"log"
	"os"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
		cmdAll(client)
	case "add":
		cmdAdd(client)
	case "get-todo":
		cmdGetTodo(client)
	case "update":
		cmdUpdate(client)
	case "delete":
		cmdDel(client)
	default:
		fmt.Println("Invalid command.")
		printUsage()
	}
}

func cmdAll(client *db.PrismaClient){
	todos, err := GetAll(client)
	if err != nil {
		log.Printf("Error fetching todos: %v\n", err)
		return
	}
	printTodosTable(todos)
}

func cmdDel(client *db.PrismaClient){
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
}

func cmdUpdate(client *db.PrismaClient){
	if len(os.Args) < 3 {
		fmt.Println("Missing the todo ID argument")
		return
	}
	id := os.Args[2]
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the new title: ")
	newTitle,readerErr := reader.ReadString('\n')
	if readerErr != nil{
		fmt.Println("Something went wrong",readerErr)
	}
	err := UpdateTodo(client,id,newTitle)
	if err != nil {
		log.Printf("Error updating todo: %v\n", err)
		return
	}
}

func cmdGetTodo(client *db.PrismaClient){
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

	// fmt.Println("Todo Found:")
	// fmt.Printf("ID: %d\nTitle: %s\n", todo.ID, todo.Title)
	printSingleTodo(todo)
}
func cmdAdd(client *db.PrismaClient){
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

func printTodosTable(todos []db.TodoModel) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetTitle("Todos")
	t.AppendHeader(table.Row{"ID", "Title", "UUID"})

	for _, todo := range todos {
		t.AppendRow(table.Row{todo.ID, todo.Title, todo.UID})
	}

	// Applying dark style
	style := t.Style()
	style.Color.Header = text.Colors{text.BgHiBlack, text.FgWhite}
	style.Color.Row = text.Colors{text.BgBlack, text.FgHiWhite}
	t.SetStyle(*style)

	t.Render()
}

func printSingleTodo(todo *db.TodoModel){
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetTitle("Todo Found")
	t.AppendHeader(table.Row{"ID", "Title", "UUID"})

	t.AppendRow(table.Row{todo.ID, todo.Title, todo.UID})
	

	// Applying dark style
	style := t.Style()
	style.Color.Header = text.Colors{text.BgHiBlack, text.FgWhite}
	style.Color.Row = text.Colors{text.BgBlack, text.FgHiWhite}
	t.SetStyle(*style)

	t.Render()
}
