package http

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"zerologix/config"
	"zerologix/dao"
	"zerologix/http/middleware"
	"zerologix/logger"
	"zerologix/logic"
)

func newCORSAllowOrigin(allow bool) cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = allow
	config.AllowHeaders = []string{
		"x-trace-id", "csrftoken", "x-ui-request-trace", "Authorization",
		"Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding",
		"Accept-Language", "Host", "lang", "fvideo-id", "device-info"}

	if !allow {
		config.AllowOriginFunc = func(origin string) bool {
			return true
		}
	}

	return config
}

func serverRoute(ctx context.Context) *gin.Engine {
	// Init Route
	if config.Conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	defaultRouter := gin.Default()

	logic := logic.NewLogicHandler(
		dao.Clients.Mysql,
		dao.Clients.Mysql)

	router := defaultRouter.Group("")
	router.Use(middleware.TraceIdMiddleware())
	// TODO monitor middleware
	router.Use(middleware.Translations())

	router.GET("/health-check", HandlerWrapper(ctx, logic.HealthCheck))

	// TODO auth middleware
	orderRouter := router.Group("/order")
	{
		orderRouter.POST("", HandlerWrapper(ctx, logic.OrderCreate))
		// TODO
		// orderRouter.GET("", HandlerWrapper(ctx, logic.OrdersGet))
		// orderRouter.GET("/:order_id", HandlerWrapper(ctx, logic.OrderGetOne))
	}

	return defaultRouter
}

func NewServer(ctx context.Context) *http.Server {
	logger.Log(ctx).Infof("config.Conf.HttpServer.Addr: %s", config.Conf.HttpServer.Addr)
	srv := &http.Server{
		Addr:    config.Conf.HttpServer.Addr,
		Handler: serverRoute(ctx),
	}
	return srv
}

func NewMetricsServer(ctx context.Context) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	// PProf Service
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	logger.Log(ctx).Infof("config.Conf.HttpServer.MetricsAddr: %s", config.Conf.HttpServer.MetricsAddr)
	srv := &http.Server{
		Addr:    config.Conf.HttpServer.MetricsAddr,
		Handler: mux,
	}
	return srv
}

func Start(ctx context.Context, errChan chan error, httpCloseCh chan struct{}) {
	// Init Server
	srv := NewServer(ctx)
	// replace validator with custom validator
	binding.Validator = NewCustomValidator()
	// Server run
	go func() {
		logger.Log(ctx).Infof("backendsrv srv is starting on %s", config.Conf.HttpServer.Addr)
		errChan <- srv.ListenAndServe()
	}()

	// Init metrics server
	metricsSrv := NewMetricsServer(ctx)
	go func() {
		logger.Log(ctx).Infof("portal metrics srv is starting on %s", config.Conf.HttpServer.MetricsAddr)
		errChan <- metricsSrv.ListenAndServe()
	}()

	// watch the ctx exit
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Log(ctx).Infof("httpServer shutdown:%v", err)
		}
		if err := metricsSrv.Shutdown(ctx); err != nil {
			logger.Log(ctx).Infof("portal metrics srv shutdown: %v", err)
		}
		httpCloseCh <- struct{}{}
	}()
}
