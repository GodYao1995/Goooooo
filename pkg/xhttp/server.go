package xhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/fx"
)

const (
	_Production      = "production"
	_ShortProduction = "prod"
)

var Module = fx.Provide(New)

func New(lifecycle fx.Lifecycle, vp *viper.Viper) *gin.Engine {
	engine := gin.New()

	pprof.Register(engine)

	// Common Middlewares
	engine.Use(gin.Logger(), gin.Recovery(), NoCache(), Cors(), Secure(), RequestID(), Translations())

	if mode := vp.GetString("APP.MODE"); mode == _Production || mode == _ShortProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.GET("/docs/*any", swagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:        vp.GetString("APP.SERVER_HOST"),
		Handler:     engine,
		ReadTimeout: 500 * time.Millisecond,
		// WriteTimeout: 500 * time.Millisecond,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("\033[1;32;32m=========== Server     Running: [ %s ] \033[0m", srv.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return engine
}
