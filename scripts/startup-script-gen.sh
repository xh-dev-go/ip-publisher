source ./config.sh

sourcePath=$(realpath $0)

echo "$sourcePath/ip-publisher.exe --topic $topic --device $deviceName --username $username --password $password --servers $hosts"

