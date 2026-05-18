package http

import (
	"io"
	"decent/internal/infrastructure/auth"
	"decent/internal/infrastructure/http/handler"
	"decent/internal/infrastructure/http/middleware"

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

	e.Use(echomiddleware.Logger())
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
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	r.echo.POST("/echo", func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "failed to read body"})
		}
		return c.JSON(200, map[string]interface{}{"body": string(body)})
	})

	r.echo.POST("/auth/token", authHandler.GenerateToken)

	books := r.echo.Group("/books")

	books.POST("", bookHandler.CreateBook)

	books.GET("/:id", bookHandler.GetBook)

	books.GET("", func(c echo.Context) error {
		if c.QueryParam("author") != "" {
			return bookHandler.GetBooksByAuthor(c)
		}
		if c.QueryParam("page") != "" || c.QueryParam("limit") != "" {
			return bookHandler.GetBooksPaginated(c)
		}
		return bookHandler.GetAllBooks(c)
	})

	booksProtected := books.Group("")
	booksProtected.Use(r.auth.RequireAuth)
	booksProtected.PUT("/:id", bookHandler.UpdateBook)
	booksProtected.DELETE("/:id", bookHandler.DeleteBook)
}

func (r *Router) Run(port string) error {
	return r.echo.Start(":" + port)
}
