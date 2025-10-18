package routes

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"gofr.dev/pkg/gofr"
)

type Route struct {
	app *gofr.App
}

func NewRoute(app *gofr.App) *Route {
	return &Route{app: app}
}

func (route *Route) Register() {

	log.Debug().Msg("registering the  routes....")
	
	route.app.GET("/", func(c *gofr.Context) (any, error) {
		infoMap := map[string]string{
			"Server-Name": "Stock Agent",
			"Time":        time.Now().String(),
		}
		marshal, err := json.Marshal(infoMap)
		if err != nil {
			log.Err(err).Msg("marshal failed")
			return nil, err
		}
		return marshal, nil
	})
}
