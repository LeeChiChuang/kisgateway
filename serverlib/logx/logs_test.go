package logx

import (
	"fmt"
	"kisgateway/serverlib/conf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	setConf()
	assert.NotPanics(t, func() {
		fmt.Println("+++++++++")
		Info("hello world")
	})
}

func setConf()  {
	c := make(conf.Conf)
	c["log.mode"] = "console"
	c["log.level"] = "0"

	SetUp(c)
}