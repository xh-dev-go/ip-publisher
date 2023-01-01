Set-Variable -Name "serviceName" -Value "ip-publisher-service"
Set-Variable -Name "serviceDisplayName" -Value "ip publishing service"
Set-Variable -Name "device" -Value "{device}"
Set-Variable -Name "hosts" -Value "{kafka-hosts}"
Set-Variable -Name "un" -Value "{username}"
Set-Variable -Name "password" -Value "{password}"
Set-Variable -Name "topic" -Value "{topic}"

sc.exe delete $serviceName

sc.exe create $serviceName `
	start=auto `
	displayname=$serviceDisplayName `
	binpath="$pwd\ip-publisher.exe --topic $topic --device $device --username $un --password $password --servers $hosts"

sc.exe failure $serviceName reset= 0 actions= restart/0/restart/0/restart/0
