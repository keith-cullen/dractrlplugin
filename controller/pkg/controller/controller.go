package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	v1 "k8s.io/api/core/v1"
	resourcev1alpha1 "k8s.io/api/resource/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/dynamic-resource-allocation/controller"
	"k8s.io/klog/v2"
)

const (
	controllerName = "dractrlplugin"
	defaultResync = 0
)

// Implements the Driver interface in k8s.io/dynamic-resource-allocation/controller
type Controller struct {
	name string
	clientset kubernetes.Interface
	allocated map[types.UID]string
	mutex sync.Mutex
	logger klog.Logger
}

func New(kubeconfigPath string) (*Controller, error) {
	logger := klog.LoggerWithName(klog.FromContext(context.Background()), controllerName)
	logger.Info("New", "kubeconfigPath", kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		logger.Error(err, "Failed to create controller")
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err, "Failed to create controller")
		return nil, err
	}
	allocated := make(map[types.UID]string)
	return &Controller{
		name: controllerName,
		clientset: clientset,
		allocated: allocated,
		logger: logger,
	}, nil
}

func (ctrl *Controller) Run(ctx context.Context, workers int) {
	ctrl.logger.Info("Run")
	informerFactory := informers.NewSharedInformerFactory(ctrl.clientset, defaultResync)
	dractrl := controller.New(ctx, ctrl.name, ctrl, ctrl.clientset, informerFactory) // new k8s DRA controller
	informerFactory.Start(ctx.Done())
	dractrl.Run(workers)
}

func (ctrl *Controller) GetClassParameters(ctx context.Context, class *resourcev1alpha1.ResourceClass) (interface{}, error) {
	ctrl.logger.Info("GetClassParameters")
	if class.ParametersRef == nil {
		return nil, nil
	}
	if class.ParametersRef.APIGroup != "" {
		err := fmt.Errorf("Unsupported API version: %s", class.ParametersRef.APIGroup)
		ctrl.logger.Error(err, "Failed to get class parameters")
		return nil, err
	}
	if class.ParametersRef.Kind != "ConfigMap" {
		err := fmt.Errorf("Unsupported class parameters kind: %s", class.ParametersRef.Kind)
		ctrl.logger.Error(err, "Failed to get class parameters")
		return nil, err
	}
	configMap, err := ctrl.clientset.CoreV1().ConfigMaps(class.ParametersRef.Namespace).Get(ctx, class.ParametersRef.Name, metav1.GetOptions{})
	if err != nil {
		ctrl.logger.Error(err, "Failed to get class parameters")
		return nil, err
	}
	ctrl.logger.Info("GetClassParameters", "configMap.Data", configMap.Data)
	return configMap.Data, nil
}

func (ctrl *Controller) GetClaimParameters(ctx context.Context, claim *resourcev1alpha1.ResourceClaim, class *resourcev1alpha1.ResourceClass, classParameters interface{}) (interface{}, error) {
	ctrl.logger.Info("GetClaimParameters")
	if claim.Spec.ParametersRef == nil {
		return nil, nil
	}
	if claim.Spec.ParametersRef.APIGroup != "" {
		err := fmt.Errorf("Unsupported API version: %s", claim.Spec.ParametersRef.APIGroup)
		ctrl.logger.Error(err, "Failed to get claim parameters")
		return nil, err
	}
	if claim.Spec.ParametersRef.Kind != "ConfigMap" {
		err := fmt.Errorf("Unsupported claim parameters kind: %s", claim.Spec.ParametersRef.Kind)
		ctrl.logger.Error(err, "Failed to get claim parameters")
		return nil, err
	}
	configMap, err := ctrl.clientset.CoreV1().ConfigMaps(claim.Namespace).Get(ctx, claim.Spec.ParametersRef.Name, metav1.GetOptions{})
	if err != nil {
		ctrl.logger.Error(err, "Failed to get claim parameters")
		return nil, err
	}
	ctrl.logger.Info("GetClaimParameters", "configMap.Data", configMap.Data)
	return configMap.Data, nil
}

func toEnv(src interface{}, dst map[string]string) {
	if src == nil {
		return
	}
	env, ok := src.(map[string]string)
	if !ok {
		return
	}
	for key, val := range env {
		dst[key] = val
	}
}

func (ctrl *Controller) Allocate(ctx context.Context, claim *resourcev1alpha1.ResourceClaim, claimParameters interface{}, class *resourcev1alpha1.ResourceClass, classParameters interface{}, selectedNode string) (*resourcev1alpha1.AllocationResult, error) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	immediate := selectedNode == ""
	ctrl.logger.Info("Allocate", "selectedNode", selectedNode, "immediate allocation", immediate, "UID", claim.UID)
	env := make(map[string]string)
	toEnv(classParameters, env)
	toEnv(claimParameters, env)
	data, err := json.Marshal(env)
	if err != nil {
		ctrl.logger.Error(err, "Failed to Allocate resources")
		return nil, err
	}
	ctrl.logger.Info("Allocate", "ResourceHandle", string(data))
	node, prs := ctrl.allocated[claim.UID]
	if !prs && !immediate {
		node = selectedNode
		ctrl.allocated[claim.UID] = node
	}
	var nodes []string
	nodes = append(nodes, node)
	ctrl.logger.Info("Allocate", "nodes", nodes)
	allocation := &resourcev1alpha1.AllocationResult{
		ResourceHandle: string(data),
		AvailableOnNodes: &v1.NodeSelector{
			NodeSelectorTerms: []v1.NodeSelectorTerm{
				{
					MatchExpressions: []v1.NodeSelectorRequirement{
						{
							Key:      "kubernetes.io/hostname",
							Operator: v1.NodeSelectorOpIn,
							Values:   nodes,
						},
					},
				},
			},
		},
		Shareable: false,
	}
	return allocation, nil
}

func (ctrl *Controller) Deallocate(ctx context.Context, claim *resourcev1alpha1.ResourceClaim) error {
	ctrl.logger.Info("Deallocate")
	return nil
}

func (ctrl *Controller) UnsuitableNodes(ctx context.Context, pod *v1.Pod, claims []*controller.ClaimAllocation, potentialNodes []string) error {
	ctrl.logger.Info("UnsuitableNodes")
	return nil
}
