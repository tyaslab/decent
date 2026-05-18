package http

import (
	"decent/internal/infrastructure/auth"
	"decent/internal/infrastructure/http/handler"
	"decent/internal/infrastructure/http/middleware"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	echo       *echo.Echo
	auth       *middleware.AuthMiddleware
	jwtService *auth.JWTService
}

func NewRouter(jwtService *auth.JWTService) *Router {
	e := echo.New()

	loggingMiddleware := middleware.NewLoggingMiddleware()
	e.Use(loggingMiddleware.Logger)
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	return &Router{
		echo:       e,
		auth:       authMiddleware,
		jwtService: jwtService,
	}
}

func (r *Router) SetupRoutes(bookHandler *handler.BookHandler, authHandler *handler.AuthHandler) {
	r.echo.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]any{"success": true})
	})

	r.echo.POST("/echo", func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "failed to read body"})
		}

		if len(body) == 0 {
			body = []byte("{}")
		}

		result := map[string]any{}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "failed to unmarshal body"})
		}

		return c.JSONBlob(http.StatusOK, body)
	})

	r.echo.POST("/auth/token", authHandler.GenerateToken)

	books := r.echo.Group("/books")
	books.POST("", bookHandler.CreateBook)
	books.GET("/:id", bookHandler.GetBook)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	booksProtected := books.Group("")
	booksProtected.Use(r.auth.RequireAuth)
	booksProtected.GET("", func(c echo.Context) error {
		if c.QueryParam("author") != "" {
			return bookHandler.GetBooksByAuthor(c)
		}
		if c.QueryParam("page") != "" || c.QueryParam("limit") != "" {
			return bookHandler.GetBooksPaginated(c)
		}
		return bookHandler.GetAllBooks(c)
	})
}

func (r *Router) Run(port string) error {
	return r.echo.Start(":" + port)
}
