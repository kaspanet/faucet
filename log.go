package main

import (
	"github.com/kaspanet/faucet/logger"
	"github.com/kaspanet/kaspad/util/panics"
)

var (
	log   = logger.Logger("FAUC")
	spawn = panics.GoroutineWrapperFunc(log)
)
