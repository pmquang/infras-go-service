package main

import (
    "net/http"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.GET("/process", process)
    e.Logger.Fatal(e.Start(":9000"))
}

func process(c echo.Context) error {
  sum := 0
  for i := 0; i < 100000000; i++ {
    sum += 1
  }
  return c.JSON(http.StatusCreated, sum)
}
