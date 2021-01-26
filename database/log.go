package database

import (
	"github.com/kaspanet/faucet/logger"
	"github.com/kaspanet/kaspad/util/panics"
)

var (
	log   = logger.BackendLog.Logger("DTBS")
	spawn = panics.GoroutineWrapperFunc(log)
)
