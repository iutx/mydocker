### 3.1
```shell script
./mydocker run -ti /bin/sh
```

### 3.2
存在问题：第一次运行如果不加 -mem 可能会报错（Wait位置），是否跟Init时候空值有关系; 暂时无法复现
```shell script
/mydocker run -ti -mem 10m stress --vm-bytes 200m --vm-keep -m 1
```

存在问题：cpu.shares 的确已经限制到了512.但是top中看到 cpu并没有限制到 50%以下
```shell script
./mydocker run -ti -cpushare 512 stress --vm-bytes 200m --vm-keep -m 1 
``` 

存在问题： write /sys/fs/cgroup/cpuset/mydocker-cgroup/tasks: no space left on device
https://blog.csdn.net/xftony/article/details/80536562  问题待定
```shell script
./mydocker run -ti -cpuset 2 stress --vm-bytes 200m --vm-keep -m 1
``` 

### 3.3
顺序应该在3.2前面
```shell script
./mydocker run -ti sh
```

### 4.1
为了方便接下来任务继续进行，暂时注释掉 ./subsystems/subsystem.go/&CpuSetSubSystem{} 19 lines
该问题具体解决进度：https://github.com/xianlubird/mydocker/issues/74
需要提前解包 busybox 到 /root/busybox。 在 ./container/container_proccess.go 中 cmd.dir 设置
```shell script
./mydocker run -ti sh
```

### 4.2-overlay
Linux内核 4.2.0-27-generic
/opt/busybox 为镜像层， /opt/merged 为镜像内层
 
```shell script
./mydocker run -ti sh
```
overlay 挂载
```shell script
mount -t overlay overlay -o lowerdir=/opt/busybox,upperdir=/opt/container_layer,workdir=/opt/index /opt/merged
```

### 4.3-overlay
Linux内核 4.2.0-27-generic

```shell script
./mydocker run -ti -v /opt/mountVolume:/tmp sh
```

### 4.4-overlay
Linux内核 4.2.0-27-generic

```shell script
./mydocker run -ti sh
```
另外运行，打包镜像为 tar，简单版本
```shell script
./mydocker commit imagesName
```

### 5.1
如果没有-ti，程序退出，子进程成为孤儿进程；但是并不是像书中说的被1进程托管，而是 init --user 进程托管（测试环境为 Ubuntu 14.04 Desktop）
如果未登录Desktop，SSH接入，孤儿进程由SSH连接进程托管，SSH退出以后，孤儿进程也退出了。此处书中代码并不是基于4.4继续添加的，此处存在疑问；
```shell script
./mydocker run -d top
```

### 5.2
容器信息存储路径  /var/run/mydocker/$PID/config.json  
```shell script
./mydocker ps
```

### 5.3
容器日志存储路径 /var/run/mydocker/$PID/container.log  
```shell script
./mydocker log containerName 
```

### 5.4
进入容器，需要调用C代码，一定注意 exec 要调用 nsenter.go
```shell script
./mydocker ps
./mydocker exec container_name sh
```

### 5.5
```shell script
./mydocker ps 
./mydocker stop container_name
ps -ef | grep pid 
./mydocker ps 
```

### 5.6
目前实现：删除不在运行状态的容器，需要增加强制删除，只需要判断加上停止容器即可。
```shell script
./mydocker rm container_name
./mydocker ps 
```

### 5.7-overlay
容器镜像的打包基于简单的去tar merged层去做，实际docker并不是这么做的

```shell script
./mydocker run -d --name container1 -v /opt/volume:/tmp busybox bin/top
./mydocker run -d --name container2 -v /opt/volume2:/tmp busybox bin/top
```
Container1中
```shell script
./mydocker exec container1 sh
cd /tmp
vi container1.txt   // Hello container1 in /tmp
mkdir /workspace
vi /workspace/container1 // Hello container1 in /workpsace

./mydocker commit container1 image1
./mydocker stop container1
./mydocker rm container1

./mydocker run -d --name image1 -v /opt/volume:/tmp image1 bin/top
cat /tmp/container1.txt
cat /workspace/container1
```

### 5.8
```shell script
./mydocker run -ti --name test -e viper=hello -e world=viper  busybox sh
env | grep viper
```

```shell script
./mydocker run -d --name test -e viper=hello -e world=viper  busybox bin/top
./mydocker exec test sh
env | grep viper
```





