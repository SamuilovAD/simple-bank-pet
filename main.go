package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SamuilovAD/simple-bank-pet/api"
	db "github.com/SamuilovAD/simple-bank-pet/db/sqlc"
	"github.com/SamuilovAD/simple-bank-pet/gapi"
	"github.com/SamuilovAD/simple-bank-pet/pb"
	"github.com/SamuilovAD/simple-bank-pet/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	config, err := util.LoadConfig(".")
	if config.ENVIRONMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}
	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
	//runGinServer(config, store)

}

func runDBMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration instance")
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	//reflection.Register(grpcServer) // allow a GRPC client to see all services and methods
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create gRPC listener")
	}
	log.Info().Msgf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener) //starts sever
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register gRPC gateway handler")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create gateway listener")
	}
	log.Info().Msgf("Start Gateway Http server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler) //starts sever
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
