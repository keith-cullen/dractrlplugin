# DraControllerPlugin

A skeleton k8s DRA controller/plugin project.

## Instructions

1. Modify kubelet configuration

        edit /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

            Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --feature-gates=DynamicResourceAllocation=true"

        $ systemctl daemon-reload
        $ systemctl restart kubelet

2. Modify kube-apiserver configuration

        edit /etc/kubernetes/manifests/kube-apiserver.yaml

            spec:
              containers:
              - command:
                - kube-apiserver
                - --feature-gates=DynamicResourceAllocation=true
                - --runtime-config=resource.k8s.io/v1alpha1

3. Modify kube-controller-manager configuration

        edit /etc/kubernetes/manifests/kube-controller-manager.yaml

            spec:
              containers:
              - command:
                - kube-controller-manager
                - --feature-gates=DynamicResourceAllocation=true

4. Modify kube-scheduler configuration

        edit /etc/kubernetes/manifests/kube-scheduler.yaml

            spec:
              containers:
              - command:
                - kube-scheduler
                - --feature-gates=DynamicResourceAllocation=true

5. Build and run the DRA controller

        $ cd controller

        edit run.sh

            KUBECONFIG=
            KUBERNETES_MASTER=
            KUBERNETES_SERVICE_HOST=
            KUBERNETES_SERVICE_PORT=

        $ make
        $ ./run.sh

6. Build and run the DRA plugin

        $ cd plugin

        edit run.sh

            KUBECONFIG=
            KUBERNETES_MASTER=
            KUBERNETES_SERVICE_HOST=
            KUBERNETES_SERVICE_PORT=

        $ make
        $ ./run.sh

6. Create a ResourceClass, ConfigMap, ResourceClaimTemplate and Pod

        $ kubectl apply -f manifest/spec.yaml
