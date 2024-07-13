package http

import (
	"github.com/gofiber/fiber/v2"
	middleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/illiafox/spy-cat-test-assignment/app/config"
	"go.uber.org/zap"
	"net"
	"strconv"
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	logger  *zap.Logger
	service Service
}

func NewServer(cfg *config.Config, logger *zap.Logger, service Service) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "spy-cats-api",
		DisableStartupMessage: !cfg.Debug,
	})

	server := &Server{
		service: service,
		cfg:     cfg,
		logger:  logger,
		app:     app,
	}
	server.setupRoutes(logger)

	return server
}

func (s *Server) Start() error {
	addr := net.JoinHostPort(s.cfg.ServerHost, strconv.Itoa(s.cfg.ServerPort))
	s.logger.Info("HTTP server started", zap.String("addr", addr))

	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) setupRoutes(logger *zap.Logger) {
	s.app.Use(middleware.New())
	s.app.Use(ErrorMiddleware(logger.With(
		zap.String("server", "fiber"),
	)))

	handler := Handler{s.service}

	s.app.Route("/cats", func(router fiber.Router) {
		router.Get("/", handler.GetCats)
		router.Post("/", handler.CreateCat)

		router.Route("/:cat_id", func(router fiber.Router) {
			router.Get("/", handler.GetCatByID)
			router.Patch("/", handler.UpdateCatByID)
			router.Delete("/", handler.DeleteCatByID)
		})
	})

	s.app.Route("/missions", func(router fiber.Router) {
		router.Get("/", handler.GetMissions)
		router.Post("/", handler.CreateMission)

		router.Route("/:mission_id", func(router fiber.Router) {
			router.Get("/", handler.GetMissionByID)
			router.Patch("/", handler.UpdateMissionByID)
			router.Post("/complete", handler.CompleteMissionByID)
			router.Delete("/", handler.DeleteMissionByID)

			router.Route("/targets", func(router fiber.Router) {
				router.Get("/", handler.GetMissionTargets)
				router.Post("/add", handler.AddMissionTargets)

				router.Route("/:target_id", func(router fiber.Router) {
					router.Get("/", handler.GetTargetByID)
					router.Post("/complete", handler.CompleteTargetByID)
					router.Delete("/", handler.DeleteTargetByID)

					router.Post("/notes/add", handler.AddTargetNote)
				})
			})
		})
	})
}
