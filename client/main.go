package main

import (
	"context"
	"log"

	todov1 "github.com/ganeshbabugb/todo-grpc/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := todov1.NewTodoServiceClient(conn)

	// CREATE
	ct, err := client.CreateTodo(context.Background(), &todov1.CreateTodoRequest{
		Title:       "Buy groceries",
		Description: "Milk, Bread, Eggs",
	})
	if err != nil {
		log.Fatalf("CreateTodo failed: %v", err)
	}
	log.Printf("Created Todo: ID=%s, Title=%s, Status=%t", ct.Todo.Id, ct.Todo.Title, ct.Todo.Completed)

	// GET
	gt, err := client.GetTodo(context.Background(), &todov1.GetTodoRequest{
		Id: ct.Todo.Id,
	})
	if err != nil {
		log.Fatalf("CreateTodo failed: %v", err)
	}
	log.Printf("GET Todo: ID=%s, Title=%s, Status=%t", gt.Todo.Id, gt.Todo.Title, gt.Todo.Completed)

	updateTodo := gt.Todo
	updateTodo.Completed = true

	// UPDATE
	ut, err := client.UpdateTodo(context.Background(), &todov1.UpdateTodoRequest{
		Todo: updateTodo,
	})
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update Todo: ID=%s, Title=%s, Status=%t", ut.Todo.Id, ut.Todo.Title, ut.Todo.Completed)

	// VERIFY UPDATE
	vt, err := client.GetTodo(context.Background(), &todov1.GetTodoRequest{
		Id: ut.Todo.Id,
	})
	if err != nil {
		log.Fatalf("Verify todo failed: %v", err)
	}
	log.Printf("Verify Todo: ID=%s, Title=%s, Status=%t", ut.Todo.Id, ut.Todo.Title, ut.Todo.Completed)

	// DELETE
	dt, err := client.DeleteTodo(context.Background(), &todov1.DeleteTodoRequest{
		Id: vt.Todo.Id,
	})
	if err != nil {
		log.Fatalf("Verify todo failed: %v", err)
	}
	log.Printf("Delete Todo: ID=%s, Status=%t", vt.Todo.Id, dt.Success)
}
