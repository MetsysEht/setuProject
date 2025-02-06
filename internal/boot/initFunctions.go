package boot

import (
	"log"

	"github.com/MetsysEht/setuProject/pkg/config"
	"github.com/MetsysEht/setuProject/utils/osUtils"
)

func initConfig() {
	// Init config
	err := config.NewDefaultConfig().Load(osUtils.GetEnv(), &Config)
	if err != nil {
		log.Fatal(err)
	}
}
