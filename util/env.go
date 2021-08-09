package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func GetEnvIntOr(envName string, override int) int {
	v := os.Getenv(envName)
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		logrus.Infoln(envName, "=", override)
		return override
	}
	logrus.Infoln(envName, "=", n)
	return n
}
