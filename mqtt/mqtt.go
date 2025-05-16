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
	OnConnect(cb func())
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

type ClientWrapper struct {
	C proxy.Client
}

var WithOnConnectHandler = func(h func(*ClientWrapper)) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetOnConnectHandler(func(c proxy.Client) {
			h(&ClientWrapper{C: c})
		})
		return options
	}
}
var WithOnLostHandler = func(h func(*ClientWrapper, error)) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetConnectionLostHandler(func(c proxy.Client, err error) {
			h(&ClientWrapper{C: c}, err)
		})
		options = options.SetConnectRetry(true)
		return options
	}
}

var WithWill = func(topic string, payload string, qos byte, retained bool) ClientOption {
	return func(options *proxy.ClientOptions) *proxy.ClientOptions {
		options = options.SetWill(topic, payload, qos, retained)
		return options
	}
}
