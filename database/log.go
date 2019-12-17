package database

import "github.com/kaspanet/kaspad/util/panics"
import "github.com/kaspanet/kasparov/logger"

var (
	log   = logger.BackendLog.Logger("DTBS")
	spawn = panics.GoroutineWrapperFunc(log)
)
