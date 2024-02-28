package mqtt

import (
	"errors"
	proxy "github.com/eclipse/paho.mqtt.golang"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/util"
	"time"
)

var (
	ErrPublishTimeout = errors.New("mqtt publish wait timeout")
)

type Session interface {
	Connect() error
	Publish(topic string, qos byte, retained bool, data []byte) error
	Subscribe(handler func(topic string, payload []byte), qos byte, topic string) error
	Unsubscribe(topics ...string)
}

type ClientOption func(*proxy.ClientOptions) *proxy.ClientOptions

var WithAuth = func(username, password string) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetUsername(username)
		if len(password) > 0 {
			options = options.SetPassword(password)
		}
		return options
	}
}

var WithRandomClientId = func() ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetClientID(util.UUIDString())
		log.Debug("ClientId:%s", options.ClientID)
		return options
	}
}
var WithClientId = func(clientId string) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetClientID(clientId)
		return options
	}
}

var WithBroker = func(broker string, others ...string) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.AddBroker(broker)
		for _, b := range others {
			options = options.AddBroker(b)
		}
		return options
	}
}

var WithSessionControl = func(keepAlive, pingTimeout, writeTimeout time.Duration) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.
			SetKeepAlive(keepAlive).
			SetPingTimeout(pingTimeout).
			SetWriteTimeout(writeTimeout)
		return options
	}
}

var WithOnConnectHandler = func(h proxy.OnConnectHandler) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetOnConnectHandler(h)
		return options
	}
}
var WithOnLostHandler = func(h proxy.ConnectionLostHandler) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetConnectionLostHandler(h)
		return options
	}
}
