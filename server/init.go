package server

import (
	"dns-check/config"
	"dns-check/server/router"
	"fmt"
)

func StartGinServer() error {
	r := router.SetupRouter()
	if err := r.Run(fmt.Sprintf("0.0.0.0:%s", config.Conf.Http.Port)); err != nil {
		return err
	} else {
		return nil
	}
}
