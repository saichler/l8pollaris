package tests

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
)

func CreateDeviceTarget(ip, linksId, crId string) *l8tpollaris.L8PTarget {
	device := &l8tpollaris.L8PTarget{}
	device.TargetId = ifs.NewUuid()
	device.LinksId = linksId
	device.Hosts = make(map[string]*l8tpollaris.L8PHost)
	host := &l8tpollaris.L8PHost{}
	host.TargetId = ifs.NewUuid()

	host.Configs = make(map[int32]*l8tpollaris.L8PHostProtocol)
	device.Hosts[host.TargetId] = host

	sshConfig := &l8tpollaris.L8PHostProtocol{}
	sshConfig.Protocol = l8tpollaris.L8PProtocol_L8PSSH
	sshConfig.Port = 22
	sshConfig.Addr = ip
	sshConfig.CredId = crId
	sshConfig.Terminal = "vt100"
	sshConfig.Timeout = 60

	host.Configs[int32(sshConfig.Protocol)] = sshConfig

	snmpConfig := &l8tpollaris.L8PHostProtocol{}
	snmpConfig.Protocol = l8tpollaris.L8PProtocol_L8PPSNMPV2
	snmpConfig.Addr = ip
	snmpConfig.Port = 161
	snmpConfig.Timeout = 60
	snmpConfig.CredId = crId

	host.Configs[int32(snmpConfig.Protocol)] = snmpConfig

	return device
}
