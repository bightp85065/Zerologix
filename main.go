package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/exp/slog"

	"zerologix/config"
	"zerologix/dao"
	server "zerologix/http"
	"zerologix/job"
	"zerologix/logger"

	"go.uber.org/automaxprocs/maxprocs"
)

var (
	buildDate                       = "dirty"
	gitCommit                       = "dirty"
	goVersion                       = runtime.Version()
	plaform                         = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	defaultConfFile                 = "./conf.yaml"
	defaultEnvFile                  = ".env.yaml"
	confFile, envFile, debugLogPath string
)

func main() {
	maxprocs.Set()
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.WithContext(ctx, slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", "zerologix"))

	flagParse()
	initRuntime(ctx, cancel)
	defer runtimeClose(ctx)
	initServer(ctx, cancel)
}

func flagParse() {
	flag.StringVar(&confFile, "c", defaultConfFile, "config file")
	flag.StringVar(&envFile, "e", defaultEnvFile, "env file")
	flag.StringVar(&debugLogPath, "log-path", "", "debug log path")
	vf := flag.Bool("v", false, "show the version and exit")
	flag.Parse()

	if *vf {
		fmt.Fprintf(os.Stdout, "BuildDate: %s\n", buildDate)
		fmt.Fprintf(os.Stdout, "GitCommit: %s\n", gitCommit)
		fmt.Fprintf(os.Stdout, "GoVersion: %s\n", goVersion)
		fmt.Fprintf(os.Stdout, "Platform: %s\n", plaform)
		os.Exit(0)
	}
}

func initRuntime(ctx context.Context, cancel context.CancelFunc) {
	initConfig(ctx)
	initLog()
	initStore(ctx)
	runJob(ctx)
}

func runtimeClose(ctx context.Context) {
	job.JobIns.Close()
	dao.Close()
}

func initConfig(ctx context.Context) {
	if err := config.Init(ctx); err != nil {
		panic(err)
	}
}

func initLog() {
	// TODO log server/dao init
}

func initStore(ctx context.Context) {
	dao.Init(ctx)
}

func runJob(ctx context.Context) {
	job.JobIns.Run(ctx)
}

func initServer(ctx context.Context, cancel context.CancelFunc) {
	errCh := make(chan error)
	httpCloseCh := make(chan struct{})
	server.Start(ctx, errCh, httpCloseCh)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		cancel()
		logger.Log(ctx).Infof("Start gracefulshutdown")
	case err := <-errCh:
		cancel()
		logger.Log(ctx).Errorf("http err%v", err)
	}
	<-httpCloseCh
	logger.Log(ctx).Infof("exit!")
}
