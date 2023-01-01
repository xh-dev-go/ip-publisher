# IP-Publisher-Kafka
The module is aimed for the publishing external ip address to kafka.

Use case like track change of external ip address and enabling the auto updating ACL list inside other programs or external services.

# Install Service
```powershell
Set-Variable -Name "serviceName" -Value "ip-publisher-service"
Set-Variable -Name "serviceDisplayName" -Value "ip publishing service"
Set-Variable -Name "topic" -Value "{topicName}"
Set-Variable -Name "device" -Value "{deviceKey}"
Set-Variable -Name "un" -Value "{username}"
Set-Variable -Name "password" -Value "{password}"
Set-Variable -Name "hosts" -Value "{hosts separted by comma}"

sc.exe delete $serviceName `

sc.exe create $serviceName `
	start=auto `
	displayname=$serviceDispalyName `
	binpath="$pwd\ip-detect-upstash-kafka.exe --topic $topic --device $device --username $un --password $password --servers $hosts"

sc.exe failure $serviceName reset= 0 actions= restart/0/restart/0/restart/0

```
