source ./config.sh

sed "s/{device}/$deviceName/g" template-ip-publisher.ps1 | \
  sed "s/{kafka-hosts}/$hosts/g" | \
  sed "s/{username}/$username/g" | \
  sed "s/{password}/$password/g" | \
  sed "s/{topic}/$topic/g"

