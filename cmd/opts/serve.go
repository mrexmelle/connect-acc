package opts

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/grading"
	"github.com/mrexmelle/connect-emp/internal/localerror"
	"github.com/mrexmelle/connect-emp/internal/profile"
	"github.com/mrexmelle/connect-emp/internal/titling"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/dig"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func Serve(cmd *cobra.Command, args []string) {
	container := dig.New()

	container.Provide(config.NewRepository)
	container.Provide(grading.NewRepository)
	container.Provide(titling.NewRepository)

	container.Provide(config.NewService)
	container.Provide(grading.NewService)
	container.Provide(localerror.NewService)
	container.Provide(titling.NewService)

	container.Provide(grading.NewController)
	container.Provide(profile.NewController)
	container.Provide(titling.NewController)

	process := func(
		configService *config.Service,
		gradingController *grading.Controller,
		profileController *profile.Controller,
		titlingController *titling.Controller,
	) {
		r := chi.NewRouter()

		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://localhost:3000"},
			AllowedMethods:   []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		if configService.GetProfile() == "local" {
			r.Mount("/swagger", httpSwagger.Handler(
				httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", configService.GetPort())),
				httpSwagger.UIConfig(map[string]string{
					"defaultModelsExpandDepth": "-1",
				}),
			))
		}

		r.Route("/gradings", func(r chi.Router) {
			r.Post("/", gradingController.Post)
			r.Get("/{id}", gradingController.Get)
			r.Patch("/{id}", gradingController.Patch)
			r.Delete("/{id}", gradingController.Delete)
		})

		r.Route("/titlings", func(r chi.Router) {
			r.Post("/", titlingController.Post)
			r.Get("/{id}", titlingController.Get)
			r.Patch("/{id}", titlingController.Patch)
			r.Delete("/{id}", titlingController.Delete)
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Get("/{ehid}", profileController.Get)
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", configService.GetPort()), r)

		if err != nil {
			panic(err)
		}
	}

	if err := container.Invoke(process); err != nil {
		panic(err)
	}
}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start connect-emp server",
	Run:   Serve,
}
