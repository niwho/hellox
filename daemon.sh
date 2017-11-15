#! /bin/sh
#进程名字可修改
echo $1
if [ "X$1" == "X" ]; then
    echo "请设置守护的进程可执行文件名称"
    exit -1
fi
PRO_NAME=$1

while true ; do

#    用ps获取$PRO_NAME进程数量
  NUM=`ps aux | grep ${PRO_NAME} | grep -v grep |wc -l`
#  echo $NUM
#    少于1，重启进程
  if [ "${NUM}" -lt "1" ];then
    echo "${PRO_NAME} was killed"
    ./output/bootstrap.sh
#    大于1，杀掉所有进程，重启
  elif [ "${NUM}" -gt "1" ];then
    echo "more than 1 ${PRO_NAME},killall ${PRO_NAME}"
    killall -9 $PRO_NAME
    ./output/bootstrap.sh
  fi
#    kill僵尸进程
  NUM_STAT=`ps aux | grep ${PRO_NAME} | grep T | grep -v grep | wc -l`

  if [ "${NUM_STAT}" -gt "0" ];then
    killall -9 ${PRO_NAME}
    ./output/bootstrap.sh
  fi
  sleep 1
done

exit 0m
