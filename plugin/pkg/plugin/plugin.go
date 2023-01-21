package plugin

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/dynamic-resource-allocation/kubeletplugin"
	"k8s.io/klog/v2"
	drapbv1 "k8s.io/kubelet/pkg/apis/dra/v1alpha1"
)

const (
	pluginName = "dractrlplugin"
	sockPath = "/var/lib/kubelet/plugins/dractrlplugin/dractrlplugin.sock"
	regSockPath = "/var/lib/kubelet/plugins_registry/dractrlplugin-reg.sock"
)

type Plugin struct {
	name string
	draPlugin kubeletplugin.DRAPlugin
	logger klog.Logger
}

func New() *Plugin {
	logger := klog.LoggerWithName(klog.FromContext(context.Background()), pluginName)
	logger.Info("New")
	return &Plugin{
		name: pluginName,
		logger: logger,
	}
}

func (plugin *Plugin) Run() error {
	plugin.logger.Info("Run")
	var err error
	if err = os.MkdirAll(filepath.Dir(sockPath), 0750); err != nil {
		plugin.logger.Error(err, "Failed to create DRA plugin")
		return err
	}
	plugin.draPlugin, err = kubeletplugin.Start(
	                            plugin,
	                            kubeletplugin.PluginSocketPath(sockPath),
	                            kubeletplugin.RegistrarSocketPath(regSockPath),
	                            kubeletplugin.KubeletPluginSocketPath(sockPath),
	                            kubeletplugin.DriverName(plugin.name))
	if err != nil {
		plugin.logger.Error(err, "Failed to start DRA plugin")
		return err
	}
	return nil
}

func (plugin *Plugin) NodePrepareResource(ctx context.Context, req *drapbv1.NodePrepareResourceRequest) (*drapbv1.NodePrepareResourceResponse, error) {
	plugin.logger.Info("NodePrepareResource")
        return &drapbv1.NodePrepareResourceResponse{}, nil
}

func (plugin *Plugin) NodeUnprepareResource(ctx context.Context, req *drapbv1.NodeUnprepareResourceRequest) (*drapbv1.NodeUnprepareResourceResponse, error) {
	plugin.logger.Info("NodeUnprepareResource")
	return &drapbv1.NodeUnprepareResourceResponse{}, nil
}
