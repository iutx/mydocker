3.1
./mydocker run -ti /bin/sh

3.2
存在问题：第一次运行如果不加 -mem 可能会报错（Wait位置），是否跟Init时候空值有关系; 暂时无法复现
./mydocker run -ti -mem 10m stress --vm-bytes 200m --vm-keep -m 1

存在问题：cpu.shares 的确已经限制到了512.但是top中看到 cpu并没有限制到 50%以下
./mydocker run -ti -cpushare 512 stress --vm-bytes 200m --vm-keep -m 1 

存在问题： write /sys/fs/cgroup/cpuset/mydocker-cgroup/tasks: no space left on device
https://blog.csdn.net/xftony/article/details/80536562  问题待定
./mydocker run -ti -cpuset 2 stress --vm-bytes 200m --vm-keep -m 1 
