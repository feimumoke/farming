###PRACTICE

```shell
minikube delete
```
```shell
adduser minikube
sudo usermod -aG docker minikube  && newgrp docker
```
```shell
minikube start --driver=docker
```
```shell
kubectl get no
NAME       STATUS   ROLES                  AGE     VERSION
minikube   Ready    control-plane,master   2m56s   v1.20.2

minikube node add -p minikube
minikube node add -p minikube

kubectl get node -o wide
NAME           STATUS   ROLES                  AGE    VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
minikube       Ready    control-plane,master   11m    v1.20.2   192.168.49.2   <none>        Ubuntu 20.04.1 LTS   5.10.47-linuxkit   docker://20.10.3
minikube-m02   Ready    <none>                 116s   v1.20.2   192.168.49.3   <none>        Ubuntu 20.04.1 LTS   5.10.47-linuxkit   docker://20.10.3
minikube-m03   Ready    <none>                 75s    v1.20.2   192.168.49.4   <none>        Ubuntu 20.04.1 LTS   5.10.47-linuxkit   docker://20.10.3

```

```shell
minikube dashboard
ss -ntulp | grep 34057
kubectl proxy --address='0.0.0.0' --disable-filter=true
```

```shell
systemctl start kube-proxy
#失败重启coredns
[minikube@VM-0-2-centos ~]$ kubectl exec -i -t dnsutil-rs-dk5fx -- nslookup kubernetes.default
;; connection timed out; no servers could be reached

command terminated with exit code 1
[minikube@VM-0-2-centos ~]$ kubectl -n kube-system rollout restart deployment coredns
deployment.apps/coredns restarted
[minikube@VM-0-2-centos ~]$ kubectl exec -i -t dnsutil-rs-dk5fx -- nslookup kubernetes.default
Server:		10.96.0.10
Address:	10.96.0.10#53

Name:	kubernetes.default.svc.cluster.local
Address: 10.96.0.1
```

```shell
kubectl exec -it server-ground-5487d5f88f-6t7fb -- /bin/sh
#
# nslookup server-farmer
Server:		10.96.0.10
Address:	10.96.0.10#53

Name:	server-farmer.default.svc.cluster.local
Address: 10.107.178.182

# curl -X POST server-farmer:28088/FarmerService/SelectFarmer
{"Result":"Farmer from Server::{\"HostName\":\"server-farmer-74d5ffb5f7-6l2sk\",\"HostIp\":\"172.18.0.4\"}"}#
#
```