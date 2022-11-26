package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
)

//go:embed app
var embededFiles embed.FS

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")

		return http.FS(os.DirFS("app"))
	}

	log.Print("using embed mode")

	fsys, err := fs.Sub(embededFiles, "app")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

type UrlModel struct {
	Url string `json:"url"`
}

func main() {
	reverseProxyURL, ok := os.LookupEnv("REVERSE_PROXY_URL")
	if !ok || reverseProxyURL == "" {
		glog.Fatalf("web: environment variable not declared: reverseProxyURL")
	}

	webPort, ok := os.LookupEnv("WEB_PORT")
	if !ok || webPort == "" {
		glog.Fatalf("web: environment variable not declared: webPort")
	}

	e := echo.New()

	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/reverse-proxy-url", func(c echo.Context) error {
		return c.JSON(http.StatusOK, UrlModel{Url: reverseProxyURL})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", webPort)))
}
