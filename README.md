# k8s-lab
A repo where I keep all configs for a k8s cluster playground.

# installation
- 2 GB RAM, 2 vCPU per host
- ubuntu 18
- 1 master, 2 workers on LAN
- disable swap
- cri docker https://kubernetes.io/docs/setup/production-environment/container-runtimes
- cni flannel https://github.com/coreos/flannel
- k8s 1.18 (kubelet, kubeadm, kubectl)

# todos
- [ ] verify cluster integrity https://github.com/vmware-tanzu/sonobuoy
- [ ] terraform
- [ ] everything as code for maximum reproducability
- [ ] traefik
- [ ] iptables
- [ ] namespaces
- [ ] add worker node
- [ ] logging
- [ ] sysdig cloud
- [ ] fault tolerance
- [ ] readiness and liveness probes
- [ ] bump k8s version
- [ ] persistant storage
- [ ] ci/cd
- [ ] tracing
- [ ] metrics (k8s objects, host, app, db, 4 goldens sigs)
- [ ] rolling updates
- [ ] synthetic traffic
- [ ] run custom registry https://docs.docker.com/registry/deploying/

# license
MIT
