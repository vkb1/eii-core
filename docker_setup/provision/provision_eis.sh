#!/bin/bash -e
# Provision EIS
# Usage: sudo ./provision_eis <path-of-docker-compose-file>

hostIP=`hostname -I | awk '{print $1}'`
export HOST_IP=$hostIP
echo 'System IP Address is:' $HOST_IP

set -a
source ../.env
set +a

if [ $ETCD_NAME = 'node1' ]; then
	export ETCD_INITIAL_CLUSTER_STATE=new
else
    	if [ -f dep/$ETCD_NAME.env ]; then
		set -a
    		source dep/$ETCD_NAME.env
    		set +a
		else
			echo "Required file dep/$ETCD_NAME.env does not exits, can not join ETCD cluster."
			exit -1
		fi
fi


check_ETCD_port() {
	echo "Checking if ETCD ports are already up..."
	ports=($ETCD_CLIENT_PORT)
	for port in "${ports[@]}"
	do
		set +e
		fuser $port/tcp
		if [ $? -eq 0 ]; then
			echo "$port is already being used, so please kill that process and re-run the script."
			exit -1
		fi
		set -e
	done
}

echo "1 Generating required certificates"

 if [ $DEV_MODE = 'true' ]; then
 	echo "EIS is not running in Secure mode. Generating certificates is not required.. "
 else
 	 pip3 install -r cert_requirements.txt
	 rm -rf Certificates
	 python3 gen_certs.py --f $1
 fi

echo "2 Bringing down existing EIS containers"
docker-compose -f dep/docker-compose-provision.yml down

cd ..
docker-compose down 
cd provision
echo "2.2 Checking ETCD port"

check_ETCD_port


echo "3. Create $EIS_USER_NAME if it doesn't exists. Update UID from env if already exits with different UID"

# EIS containers will be executed as eisuser
if ! id $EIS_USER_NAME >/dev/null 2>&1; then
    groupadd $EIS_USER_NAME -g $EIS_UID
    useradd -r -u $EIS_UID -g $EIS_USER_NAME $EIS_USER_NAME
else
    if ! [ $(id -u $EIS_USER_NAME) = $EIS_UID ]; then
        usermod -u $EIS_UID $EIS_USER_NAME
        groupmod -g $EIS_UID $EIS_USER_NAME
    fi
fi


echo "4. Creating Required Directories"


if [ $ETCD_RESET = 'true' ]; then
    rm -rf $EIS_INSTALL_PATH/data/etcd
fi

mkdir -p $EIS_INSTALL_PATH/data/influxdata
mkdir -p $EIS_INSTALL_PATH/data/etcd
mkdir -p $EIS_INSTALL_PATH/sockets/
chown -R $EIS_USER_NAME:$EIS_USER_NAME $EIS_INSTALL_PATH

if [ $DEV_MODE = 'false' ]; then
	chown -R $EIS_USER_NAME:$EIS_USER_NAME Certificates 
	chmod -R 760 $EIS_INSTALL_PATH/data/
else
	chmod -R 777 $EIS_INSTALL_PATH/data/
fi


echo "5. Copying docker compose yaml file which is provided as argument."
# This file will be volume mounted inside the provisioning container and deleted once privisioning it done

cp $1 ./docker-compose.yml

echo "5. Starting and provisioning ETCD ..."
# ia_etcd_provision will run only for first node on the cluster. Not required for any other node joining

if [ $DEV_MODE = 'true' ]; then
		docker-compose -f dep/docker-compose-provision.yml up --build -d
else
		docker-compose -f dep/docker-compose-provision.yml -f dep/docker-compose-provision.override.prod.yml up --build -d
fi

rm ./docker-compose.yml
