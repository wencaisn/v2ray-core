package core

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
)

// CreateObject creates a new object based on the given V2Ray instance and config. The V2Ray instance may be nil.
func CreateObject(v *Instance, config interface{}) (interface{}, error) {
	ctx := context.Background()
	if v != nil {
		ctx = context.WithValue(ctx, v2rayKey, v)
	}
	return common.CreateObject(ctx, config)
}

// StartInstance starts a new V2Ray instance with given serialized config.
func StartInstance(configFormat string, configBytes []byte) (*Instance, error) {
	var mb buf.MultiBuffer
	common.Must2(mb.Write(configBytes))
	config, err := LoadConfig(configFormat, "", &mb)
	if err != nil {
		return nil, err
	}
	instance, err := New(config)
	if err != nil {
		return nil, err
	}
	if err := instance.Start(); err != nil {
		return nil, err
	}
	return instance, nil
}

// Dial provides an easy way for upstream caller to create net.Conn through V2Ray.
// It dispatches the request to the given destination by the given V2Ray instance.
// Since it is under a proxy context, the LocalAddr() and RemoteAddr() in returned net.Conn
// will not show real addresses being used for communication.
func Dial(ctx context.Context, v *Instance, dest net.Destination) (net.Conn, error) {
	r, err := v.Dispatcher().Dispatch(ctx, dest)
	if err != nil {
		return nil, err
	}
	return net.NewConnection(net.ConnectionInputMulti(r.Writer), net.ConnectionOutputMulti(r.Reader)), nil
}
