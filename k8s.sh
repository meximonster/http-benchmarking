#!/bin/bash

export TERM="xterm"
red=`tput setaf 1`

source .env
app=${app}
port_name=${port_name}
namespace=${namespace}
ip=${ip}
port=${port}

if [ -z "${ip}" ];
then
    ip="localhost"
fi

if [ -z "${port}" ];
then
    port="8080"
fi

if [ -z "${app}" ];
then
    printf "${red}app name not set. Please fill .env file\n"
    exit 1;
fi

if [ -z "${port_name}" ];
then
    printf "${red}port name not set. Please fill .env file\n"
    exit 1;
fi

if [ -z "${namespace}" ];
then
    printf "${red}namespace not set. Please fill .env file\n"
    exit 1;
fi

mkdir ./k8s/deploy
for i in endpoint service serviceMonitor; do
    cp k8s/$i.yml k8s/deploy/$i-$app.yml
    sed -i "s/APP_NAME/$app/g" k8s/deploy/$i-$app.yml
    sed -i "s/PORT_NAME/$port_name/g" k8s/deploy/$i-$app.yml
    sed -i "s/HTTP_PORT/$port/g" k8s/deploy/$i-$app.yml
    sed -i "s/IP_ADDRESS/$ip/g" k8s/deploy/$i-$app.yml
done
sed -i "s/NAMESPACE_NAME/$namespace/g" k8s/deploy/serviceMonitor-$app.yml
