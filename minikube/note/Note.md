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
