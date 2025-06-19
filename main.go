package main

import (
	"database/sql"
	"github.com/SamuilovAD/simple-bank-pet/api"
	db "github.com/SamuilovAD/simple-bank-pet/db/sqlc"
	"github.com/SamuilovAD/simple-bank-pet/gapi"
	"github.com/SamuilovAD/simple-bank-pet/pb"
	"github.com/SamuilovAD/simple-bank-pet/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	//reflection.Register(grpcServer) // allow a GRPC client to see all services and methods
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create gRPC listener")
	}
	log.Printf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener) //starts sever
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
