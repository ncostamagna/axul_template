package main

import (
	"context"

	"github.com/joho/godotenv"

	"flag"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net"
	"net/http"

	"github.com/ncostamagna/axul_template/templates"
	"github.com/ncostamagna/axul_template/templatespb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type server struct{}

func main() {

	fmt.Println("Initial")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "postapp",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	_ = level.Info(logger).Log("msg", "service started")
	defer func() {
		_ = level.Info(logger).Log("msg", "service ended")
	}()

	fmt.Println("Env")
	err := godotenv.Load()
	if err != nil {
		_ = level.Error(logger).Log("Error loading .env file", err)
		os.Exit(-1)
	}

	var httpAddr = flag.String("http", ":"+os.Getenv("APP_PORT"), "http listen address")

	mysqlInfo := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open("mysql", mysqlInfo)
	if err != nil {
		_ = level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	db.AutoMigrate(templates.Template{})

	flag.Parse()
	ctx := context.Background()

	var srv templates.Service
	{
		repository := templates.NewRepo(db, logger)
		srv = templates.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	mux := http.NewServeMux()

	mux.Handle("/templates/", templates.NewHTTPServer(ctx, templates.MakeEndpoints(srv)))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Hello world")

	// 50051 puerto por defecto de gRPC
	lis, err := net.Listen("tcp", "0.0.0.0:50055")
	if err != nil {
		fmt.Println("Failed to listen: %v", err)
	}

	// New server
	s := grpc.NewServer()
	// le pasamos el struct server que definimos
	templatespb.RegisterTemplatesServiceServer(s, &server{})

	go func() {

		fmt.Println("listening on port:50051")
		errs <- s.Serve(lis)

	}()

	go func() {
		fmt.Println("listening on port", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)

	}()

	err = <-errs

	if err != nil {
		_ = level.Error(logger).Log("exit", err)
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
