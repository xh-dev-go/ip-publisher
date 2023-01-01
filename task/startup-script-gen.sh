source ./config.sh

echo "$(pwd)/ip-publisher.exe --topic $topic --device $deviceName --username $username --password $password --servers $hosts"

