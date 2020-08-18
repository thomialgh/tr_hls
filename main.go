package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type response struct {
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

func main() {
	r := echo.New()
	r.Use(middleware.Logger())
	r.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, response{
			Message: "pong",
		})
	})
	r.GET("/stream", func(c echo.Context) error {
		f, err := os.Open("video/filename.m3u8")
		if err != nil {
			log.Println(err)
			return c.JSON(500, response{
				Message: "Internal Server Error",
			})
		}
		defer f.Close()
		c.Response().Header().Add("Content-Disposition", "attachment; filename=filename.m3u8")
		return c.Stream(200, "application/x-mpegURL", f)
	})
	r.GET("/:filename", func(c echo.Context) error {
		f, err := os.Open("video/" + c.Param("filename"))
		if err != nil {
			log.Println(err)
			return c.JSON(500, response{
				Message: "Internal Server Error",
			})
		}
		defer f.Close()
		return c.Stream(200, "video/MP2T", f)

	})

	r.GET("", func(c echo.Context) error {
		return c.File("html/home.html")
	})

	server := &http.Server{
		Addr:    ":12302",
		Handler: r,
	}

	panic(server.ListenAndServe())
}
