package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/cstati/auth/internal/app/auth"
	"github.com/cstati/auth/internal/pkg/storage/db"
	"github.com/cstati/auth/internal/pkg/storage/google"
	pb "github.com/cstati/auth/pkg/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	osinit "github.com/hse-experiments-platform/library/pkg/utils/init"
	"github.com/hse-experiments-platform/library/pkg/utils/loggers"
	"github.com/hse-experiments-platform/library/pkg/utils/token"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func loadEnv() {
	file := os.Getenv("DOTENV_FILE")
	// loads values from .env into the system
	if err := godotenv.Load(file); err != nil {
		log.Fatal().Err(err).Msg("cannot load env variables")
	}
}

func initDB(ctx context.Context, dsnOSKey string, loadTypes ...string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(osinit.MustLoadEnv(dsnOSKey))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse config")
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, loadType := range loadTypes {
			t, err := conn.LoadType(context.Background(), loadType) // type
			if err != nil {
				log.Fatal().Err(err).Msg("cannot load type")
			}
			conn.TypeMap().RegisterType(t)

			t, err = conn.LoadType(context.Background(), "_"+loadType) // array of type
			if err != nil {
				log.Fatal().Err(err).Msg("cannot load type")
			}
			conn.TypeMap().RegisterType(t)
		}

		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Str("dsn", osinit.MustLoadEnv(dsnOSKey)).Msg("cannot osinit db")
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	return pool
}

func initService(ctx context.Context, maker token.Maker, storage google.Storage) pb.AuthServiceServer {
	service := auth.NewService(
		storage,
		db.New(initDB(ctx, "DB_CONNECT_STRING")),
		maker,
	)

	return service
}

func runGRPC(ctx context.Context, c context.CancelFunc, server pb.AuthServiceServer, grpcAddr string, maker token.Maker) {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		logging.WithFieldsFromContext(func(ctx context.Context) logging.Fields {
			return []any{token.UserIDContextKey, ctx.Value(token.UserIDContextKey), token.UserRolesContextKey, ctx.Value(token.UserRolesContextKey)}
		}),
	}

	s := grpc.NewServer(
		//grpc.ChainUnaryInterceptor(
		//	maker.TokenExtractorUnaryInterceptor(),
		//	logging.UnaryServerInterceptor(loggers.ZerologInterceptorLogger(log.Logger), opts...),
		//),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(loggers.ZerologInterceptorLogger(log.Logger), opts...),
		),
	)
	pb.RegisterAuthServiceServer(s, server)
	reflection.Register(s)

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get grpc net.Listener")
	}

	go func() {
		<-ctx.Done()
		log.Info().Msg("stropping grpc server")
		s.GracefulStop()
	}()

	go func() {
		log.Info().Msgf("grpc server listening on %s", grpcAddr)
		err = s.Serve(l)
		if err != nil {
			log.Error().Err(err).Msg("error in grpc.Serve")
		}
		c()
	}()
}

func runHTTP(ctx context.Context, c context.CancelFunc, grpcAddr string) {
	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, rmux, grpcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register rmux")
	}

	httpAddr := ":" + osinit.MustLoadEnv("HTTP_PORT")
	l, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get http net.Listener")
	}

	//creating swagger
	mux := http.NewServeMux()
	// mount the gRPC HTTP gateway to the root
	mux.Handle("/", rmux)
	fs := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	s := http.Server{Handler: cors.AllowAll().Handler(mux)}

	go func() {
		<-ctx.Done()
		log.Info().Msg("stropping grpc server")
		err = s.Shutdown(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot shutdown http server")
		}
	}()

	go func() {
		log.Info().Msgf("http server listening on %s", httpAddr)
		err = s.Serve(l)
		if err != nil {
			log.Error().Err(err).Msg("error in http.Serve")
		}
		c()
	}()
}

func run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	googleStorage := google.NewStorage()
	maker, err := token.NewMaker(osinit.MustLoadEnv("CIPHER_KEY"))
	if err != nil {

		log.Fatal().Err(err).Msg("cannot osinit token maker")
	}
	service := initService(ctx, maker, googleStorage)

	grpcAddr := ":" + osinit.MustLoadEnv("GRPC_PORT")

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		stop()
		cancel()
	}()

	runGRPC(ctx, cancel, service, grpcAddr, maker)
	runHTTP(ctx, cancel, grpcAddr)

	<-ctx.Done()
}

func main() {
	ctx := context.Background()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // TimeFormat: time.Un.DateTime})

	loadEnv()

	run(ctx)
}
