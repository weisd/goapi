package routers

import (
	"github.com/weisd/goapi/models"
	"github.com/weisd/goapi/modules/log"
	"github.com/weisd/goapi/modules/setting"
)

func GlobalInit() {
	setting.NewConfigContext()
	setting.NewServices()
	// mysql
	models.LoadModelsConfig()
	if err := models.NewEngine(); err != nil {
		log.Fatal(4, "Fail to initialize ORM engine: %v", err)
	}
	models.HasEngine = true
}
