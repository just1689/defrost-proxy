package bin

import (
	"github.com/just1689/scale-aware-proxy-operator/v1/client"
	"github.com/just1689/scale-aware-proxy-operator/v1/netio"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("~~~ Starting Cold Start ~~~")
	netio.NewOrchestration(func() {
		client.Scaler.HandleInstruction()
	})
}
