package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/persistence"
)

const probeTimeout = 2 * time.Second

func init() {
	LivenessHandler = DefaultLiveness
	ReadinessHandler = DefaultReadiness
}

// DefaultLiveness answers whether the process is alive (no dependency checks).
func DefaultLiveness() int {
	return http.StatusOK
}

// DefaultReadiness checks optional enabled dependencies: DB, Redis, MQTT.
// Unconfigured components are skipped; any configured but unreachable dependency returns 503.
func DefaultReadiness() int {
	if err := persistence.Ping(persistence.DefaultAlias); err != nil {
		log.Debug("readiness: db", "err", err)
		return http.StatusServiceUnavailable
	}
	if err := pingRedis(); err != nil {
		log.Debug("readiness: redis", "err", err)
		return http.StatusServiceUnavailable
	}
	if err := checkMqtt(); err != nil {
		log.Debug("readiness: mqtt", "err", err)
		return http.StatusServiceUnavailable
	}
	return http.StatusOK
}

func pingRedis() error {
	client := factory.GetRedisClient()
	if client == nil || client.Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(probeCtx(), probeTimeout)
	defer cancel()
	return client.Ping(ctx).Err()
}

func checkMqtt() error {
	session := factory.GetMqtt()
	if session == nil {
		return nil
	}
	if !session.Connected() {
		return errMqttNotConnected
	}
	return nil
}

var errMqttNotConnected = errors.New("mqtt: not connected")

func probeCtx() context.Context {
	if ctx := factory.GetCtx(); ctx != nil {
		return ctx
	}
	return context.Background()
}
