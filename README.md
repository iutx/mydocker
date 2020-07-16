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

### 4.2-aufs
Linux 内核 3.13.0-85-generic 
使用AUFS包装busybox, 提前放置 busybox.tar 到 /opt

```shell script
./mydocker run -ti sh
```
如果存在删除不了/opt/mnt情况 ,手动
```shell script
umount -f /opt/mnt
```

### 4.3-aufs
基于aufs挂载映射数据卷

```shell script
./mydocker run -ti -v /opt/volume1/:/work sh
```



