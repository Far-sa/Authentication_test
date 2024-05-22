package httpServer

import (
	"net/http"
	"user-svc/internal/service/param"
	"user-svc/ports"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type server struct {
	userSvc ports.Service
	//  ports.Validator
	logger  ports.Logger
	metrics ports.HTTPMetrics
	config  ports.Config
	Router  *echo.Echo
}

func New(config ports.Config, userSvc ports.Service, logger ports.Logger, metrics ports.HTTPMetrics,
) server {
	return server{config: config, userSvc: userSvc, logger: logger, Router: echo.New()}
}

// TODO: implement serve function
func (s server) Serve() error {

	// s.Router.Use(middleware.RequestID())
	// s.Router.Use(middleware.Recover())

	s.logger.Info("server is running")

	//TODO add group for user handler
	s.Router.POST("/user/register", s.Register)
	s.Router.GET("/metrics", s.handleMetrics)

	// port := s.config.GetHTTPConfig().Port
	// address := fmt.Sprintf(":%d", port)
	if err := s.Router.Start(":5000"); err != nil {
		s.logger.Error("Router error", zap.Error(err))
	}
	return nil
}

// RegisterUserEndpoint handles user registration
func (s server) Register(c echo.Context) error {
	s.logger.Info("Handling register request")

	//! metrics
	//start := time.Now()

	// defer func() {
	// 	duration := time.Since(start).Seconds()
	// 	s.metrics.RegisterHTTPDurationHistogram().WithLabelValues(c.Request().Method, "/register").Observe(duration)

	// 	// If an error occurred, handle it and increment error counter
	// 	if err := recover(); err != nil {
	// 		s.metrics.RegisterHTTPErrorCounter().WithLabelValues(c.Request().Method, "/register").Inc()
	// 		s.logger.Error("Recovered from panic", zap.Any("error", err))
	// 	}
	// }()

	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// fieldErrors, err := s.validator.ValidateRegisterRequest(req)
	// if err != nil {
	// 	//* richError usage
	// 	msg, code := httpmsg.Error(err)
	// 	return c.JSON(code, echo.Map{
	// 		"Message": msg,
	// 		"Errors":  fieldErrors,
	// 	})
	// }

	ctx := c.Request().Context()

	resp, err := s.userSvc.Register(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	s.logger.Info("User registered successfully", zap.String("email", req.Email))

	return c.JSON(http.StatusCreated, resp)
}

func (s server) handleMetrics(c echo.Context) error {
	// Serve Prometheus metrics using promhttp.Handler()
	promhttp.Handler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s server) Profile(e echo.Context) error {
	s.logger.Warn("unimplemented")
	panic("")
}
