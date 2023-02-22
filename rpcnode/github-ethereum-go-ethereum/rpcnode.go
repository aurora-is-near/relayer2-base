package github_ethereum_go_ethereum

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/aurora-is-near/relayer2-base/broker"
	"github.com/aurora-is-near/relayer2-base/log"
	eventbroker "github.com/aurora-is-near/relayer2-base/rpcnode/github-ethereum-go-ethereum/events"
	"golang.org/x/net/context"

	gel "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/jinzhu/copier"
)

const (
	LoggerLevelConfPath = "logger.level"
)

// GoEthereum is a container on which underlying go-ethereum services can be registered.
type GoEthereum struct {
	node.Node
	Broker broker.Broker
}

type connection struct {
	io.Reader
	io.Writer
}

func (nc connection) RemoteAddr() string {
	return ""
}

func (nc connection) Close() error {
	return nil
}

func (nc connection) SetWriteDeadline(time.Time) error { return nil }

// New creates a new node with default conf
func New() (*GoEthereum, error) {
	conf := GetConfig()
	return NewWithConf(conf)
}

// NewWithConf creates a new node with given config and the types broker if node supports websocket comm
func NewWithConf(conf *Config) (*GoEthereum, error) {
	ethConf := convertConfigurationToEthNode(conf)
	n, err := node.New(ethConf)
	if err != nil {
		return nil, err
	}
	// Disable geth p2p server operation
	n.Server().Config.NoDial = true
	n.Server().Config.NoDiscovery = true
	n.Server().Config.EnableMsgEvents = false
	configureGoEthRootLogger()

	// Start eventbroker if WS configured
	eb := eventbroker.NewEventBroker()
	if conf.WSHost != "" && conf.WSPort > 0 {
		go eb.Start()
	}

	return &GoEthereum{
		Node:   *n,
		Broker: eb,
	}, nil
}

// WithMiddleware places the given middleware at the beginning of HTTP handlers chain. A middleware is a function accepting
// http.Handler and returning http.Handler. Any request matching the path argument is first processed by this middleware.
// Middleware should either return response for the request or pass the request to next HTTP handler in the chain.
// Refer to http.HandlerFunc for middleware creation.
func (ge *GoEthereum) WithMiddleware(name string, path string, middleware func(handler http.Handler) http.Handler) {
	h, err := ge.RPCHandler()
	if err != nil {
		log.Log().Fatal().Err(err).Msg("failed to get rpc handler")
	}
	ge.RegisterHandler(name, path, middleware(h))
}

func (ge *GoEthereum) Resolve(ctx context.Context, reader io.Reader, writer io.Writer) error {
	rpcHandler, err := ge.RPCHandler()
	if err != nil {
		return err
	}
	hr := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, "POST", "", reader)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	rpcHandler.ServeHTTP(hr, req)
	_, err = writer.Write(hr.Body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// HandleConfigChange re-configures the go-eth root.Logger if needed
func (ge *GoEthereum) HandleConfigChange() {
	configureGoEthRootLogger()
}

// configureGoEthRootLogger configures the go-eth root.Logger that used by its internal packages
func configureGoEthRootLogger() {
	logConf := log.GetConfig()
	gLvl, err := gel.LvlFromString(logConf.Level)
	if err != nil {
		// go-eth doesn't support fatal and panic log levels. Therefore, LvlError is assigned when there is error
		gLvl = gel.LvlError
		log.Log().Error().Err(err).Msgf("error while setting the go-eth root.Logger Level: %s ", logConf.Level)
	}

	var consoleHandler gel.Handler
	var fileHandler gel.Handler
	if logConf.LogToConsole {
		consoleHandler = gel.LvlFilterHandler(
			gLvl,
			gel.StreamHandler(os.Stderr, gel.JSONFormatEx(false, true)))
	}
	if logConf.LogToFile {
		fileHandler = gel.LvlFilterHandler(
			gLvl,
			gel.Must.FileHandler(logConf.FilePath, gel.JSONFormatEx(false, true)))
	}

	if logConf.LogToConsole && logConf.LogToFile {
		gel.Root().SetHandler(gel.MultiHandler(consoleHandler, fileHandler))
	} else if logConf.LogToConsole {
		gel.Root().SetHandler(consoleHandler)
	} else if logConf.LogToFile {
		gel.Root().SetHandler(fileHandler)
	} else {
		gel.Root().SetHandler(gel.DiscardHandler())
	}
}

func convertConfigurationToEthNode(confAurora *Config) *node.Config {
	confEth := &node.Config{}
	copier.CopyWithOption(confEth, confAurora, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return confEth
}
