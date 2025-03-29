package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/watchakorn-18k/scalar-go"
)

type Config struct {
	PathURL   string
	SpecURL   string
	PageTitle string
	Theme     string
	Layout    string
	DarkMode  bool
}

// NewConfig สร้าง config พร้อมค่าดีฟอลต์
func NewConfig(pathURL, specURL, pageTitle string, theme, layout string, darkMode *bool) Config {
	// ตั้งค่าดีฟอลต์
	if theme == "" {
		theme = "purple"
	}
	if layout == "" {
		layout = "modern"
	}
	if darkMode == nil {
		defaultDarkMode := true
		darkMode = &defaultDarkMode
	}

	return Config{
		PathURL:   pathURL,
		SpecURL:   specURL,
		PageTitle: pageTitle,
		Theme:     theme,
		Layout:    layout,
		DarkMode:  *darkMode,
	}
}

// ScalarMiddleware สร้าง middleware จาก config
func ScalarMiddleware(config Config) fiber.Handler {
	if config.PageTitle == "" {
		name, _ := GetModuleName()
		config.PageTitle = fmt.Sprintf("%v API documentation", name)
	}
	return func(c *fiber.Ctx) error {
		if c.Path() == config.PathURL {
			htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL: config.SpecURL,
				CustomOptions: scalar.CustomOptions{
					PageTitle: config.PageTitle,
				},
				Theme:    scalar.ThemeId(config.Theme),
				Layout:   scalar.ReferenceLayoutType(config.Layout),
				DarkMode: config.DarkMode,
			})
			if err != nil {
				return err
			}
			c.Type("html")
			return c.SendString(htmlContent)
		}
		return c.Next()
	}
}
