package modelProxy

import (
	"drip/core/serverRaiser"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

type CanaryProxy struct {
	ProxyBase
	portProd  int
	portDev   int
	startDate time.Time
	endDate   time.Time
}

func NewCanaryProxy(ms NewModelSetting) Proxy {
	modelName := ms.ModelName
	raiser := ms.Raiser
	deadline := ms.Deadline
	var err error
	var portProd, portDev int
	portProd, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	portDev, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	cmdProd, endpoint := raiser(portProd, modelName, serverRaiser.Prod)
	cmdDev, _ := raiser(portDev, modelName, serverRaiser.Dev)
	return &CanaryProxy{
		ProxyBase: ProxyBase{
			activeServers: []*exec.Cmd{cmdProd, cmdDev},
			modelName:     modelName,
			endpoint:      endpoint,
		},
		portProd:  portProd,
		portDev:   portDev,
		startDate: time.Now(),
		endDate:   deadline,
	}
}

func (c *CanaryProxy) prob() float64 {
	now := time.Now()
	if now.Before(c.startDate) {
		return 0
	}
	if now.After(c.endDate) {
		return 1
	}

	totalDuration := c.endDate.Sub(c.startDate).Seconds()
	elapsed := now.Sub(c.startDate).Seconds()

	return elapsed / totalDuration
}

func (c *CanaryProxy) forwardReq(request *http.Request, useDev bool) (gin.H, error) {
	var url, testCase string
	if useDev {
		url = fmt.Sprintf("%s:%d/%s", baseUrl, c.portDev, c.endpoint)
		testCase = "challenger"
	} else {
		url = fmt.Sprintf("%s:%d/%s", baseUrl, c.portProd, c.endpoint)
		testCase = "champion"
	}

	response, err := http.Post(url, "application/json", request.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := make(gin.H)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	result := gin.H{
		"data": data,
		"metadata": map[string]string{
			"case":  testCase,
			"reqId": genReqId(),
		},
	}

	return result, nil
}

func (c *CanaryProxy) ForwardRequest(request *http.Request) (gin.H, error) {
	useDev := rand.Float64() < c.prob()
	return c.forwardReq(request, useDev)
}

func (c *CanaryProxy) Close() {
	shutDownServer(&c.ProxyBase)
}
