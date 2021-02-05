# croupier

## get started
```bash
make
```

### minikube
```bash
make minikube-start
make minikube-addons-enabled
make artifacts
```

append the following command's output in your `/etc/hosts`
```bash
echo $(minikube ip) minikube
```
e.g:
```
127.0.0.1 localhost
192.168.64.2 minikube
```

once every artifact is successfully built and the minikube is identified as a HOST in your computer, you might follow [this](https://github.com/devbytom/hurb/tree/master/kubernetes) in order to install the application charts
```bash
cd ./kubernetes
```
