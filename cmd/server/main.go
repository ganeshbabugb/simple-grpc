package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/docker/distribution/uuid"
	todov1 "github.com/ganeshbabugb/todo-grpc/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type todoServer struct {
	// todov1.UnimplementedTodoServiceServer
	mu    sync.Mutex
	store map[string]*todov1.Todo
}

func newTodoServer() *todoServer {
	return &todoServer{
		store: make(map[string]*todov1.Todo),
	}
}

func (t *todoServer) CreateTodo(context context.Context, request *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	id := uuid.Generate().String()
	todo := &todov1.Todo{
		Id:          id,
		Title:       request.Title,
		Description: request.Description,
		Completed:   false,
	}

	t.store[id] = todo

	return &todov1.CreateTodoResponse{Todo: todo}, nil
}

func (t *todoServer) GetTodo(context context.Context, request *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todo, ok := t.store[request.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo with id %s not found", request.Id)
	}

	return &todov1.GetTodoResponse{Todo: todo}, nil
}

func (t *todoServer) UpdateTodo(context context.Context, request *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	id := request.Todo.Id
	_, ok := t.store[id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo with id %s not found", request.Todo.Id)
	}

	t.store[id] = request.Todo

	return &todov1.UpdateTodoResponse{Todo: request.Todo}, nil
}

func (t *todoServer) DeleteTodo(context context.Context, request *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todo, ok := t.store[request.Id]
	if !ok {
		return &todov1.DeleteTodoResponse{Success: false}, status.Errorf(codes.NotFound, "todo with id %s not found", request.Id)
	}

	delete(t.store, todo.Id)

	return &todov1.DeleteTodoResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	todov1.RegisterTodoServiceServer(grpcServer, newTodoServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
