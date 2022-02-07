package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/pages"
	"codeberg.org/librarian/librarian/proxy"
	"codeberg.org/librarian/librarian/static"
	"codeberg.org/librarian/librarian/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/librarian/")
	viper.AddConfigPath("$HOME/.config/librarian")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "3000")
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	viper.SetDefault("STREAMING_API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	viper.Set("AUTH_TOKEN", api.NewUser())
	viper.WriteConfig()
	if viper.GetString("HMAC_KEY") == "" {
		b := make([]byte, 36)
		rand.Read(b)
		viper.Set("HMAC_KEY", fmt.Sprintf("%x", b))
		viper.WriteConfig()
	}

	engine := handlebars.NewFileSystem(http.FS(views.GetFiles()), ".hbs")
	app := fiber.New(fiber.Config{
		Views:             engine,
		Prefork:           viper.GetBool("FIBER_PREFORK"),
		UnescapePath:      true,
		StreamRequestBody: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = ctx.Status(code).Render("error", fiber.Map{
				"err": err,
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})

	app.Use(recover.New())

	app.Use("/", etag.New())
	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(static.GetFiles()),
	}))

	app.Get("/", pages.FrontpageHandler)
	app.Get("/image", proxy.ProxyImage)
	if viper.GetBool("ENABLE_LIVE_STREAM") {
		app.Get("/live/hls/:claimId/:path", proxy.ProxyLive)
	}
	app.Get("/search", pages.SearchHandler)
	app.Get("/privacy", pages.PrivacyHandler)

	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		file, _ := static.GetFiles().ReadFile("robots.txt")
		_, err := c.Write(file)
		return err
	})
	app.Get("/sw.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		file, _ := static.GetFiles().ReadFile("js/sw.js")
		_, err := c.Write(file)
		return err
	})

	app.Get("/api/comments", api.CommentsHandler)

	app.Get("/:channel/", pages.ChannelHandler)
	app.Get("/$/invite/:channel", pages.ChannelHandler)
	app.Get("/$/invite/:channel/", pages.ChannelHandler)
	app.Get("/:channel/rss", pages.ChannelRSSHandler)
	app.Get("/embed/:channel/:claim", pages.EmbedHandler)
	app.Get("/:channel/:claim", pages.ClaimHandler)

	app.Listen(":" + viper.GetString("PORT"))
}
