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

# dummy app
I'll use a dummy blogging app to test things. It'll be as simple as possible. See app/README.md

# todos
- [ ] verify cluster integrity https://github.com/vmware-tanzu/sonobuoy
- [ ] terraform
- [ ] everything as code for maximum reproducability
- [ ] traefik
- [ ] iptables
- [ ] namespaces
- [ ] add worker node post setup
- [ ] logging
- [ ] sysdig cloud
- [ ] fault tolerance
- [ ] bump k8s version
- [ ] persistant storage
- [ ] ci/cd
- [ ] tracing
- [ ] metrics (k8s objects, host, app, db, 4 goldens sigs)
- [x] rolling updates
- [ ] rollback
- [ ] synthetic traffic
- [ ] run custom registry https://docs.docker.com/registry/deploying/
- [x] remote kubectl
- [ ] k8s events
- [ ] drain a node https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/
- [ ] pod Disruption Budget https://kubernetes.io/docs/tasks/run-application/configure-pdb/

# license
MIT
