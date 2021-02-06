# kubernetes
this directory contains essential schema and values files to run every service expected by [delta challenge](https://github.com/hurbcom/challenge-delta)


## TL;DR;
```bash
make
make add-production
```
#### namespace
first we need somewhere to install our charts in minikube 
```bash
make add-namespace
#make rm-namespace
#"sudo rm -rf /" equivalent commented above
```
<!--
#### metrics-server
```bash
make add-metrics-server
#make rm-metrics-server
```
-->

#### database
then setup postgres internal service (this is not available from outside the cluster context)
```bash
make add-database
#make rm-database
```

#### redis
then setup redis internal service (this is not available from outside the cluster context)
```bash
make add-redis
#make rm-redis
```

#### dsock-api
then setup dsock-api internal service (this is not available from outside the cluster context)
```bash
make add-dsock-api
#make rm-dsock-api
```


#### dsock-worker
then setup dsock-worker internal service (this is not available from outside the cluster context)
```bash
make add-dsock-worker
#make rm-dsock-worker
```

#### server
once workers are.....
```bash
make add-server
# make rm-server
```
