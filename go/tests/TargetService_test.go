package tests

import (
	"fmt"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"testing"
)

func TestTargetService(t *testing.T) {
	nic := topo.VnicByVnetNum(2, 2)
	targets.Activate("postgres", "probler", nic)
	tr, _ := targets.Targets(nic.Resources())
	device := &l8tpollaris.L8PTarget{TargetId: "Test"}
	resp := tr.Post(object.New(nil, device), nic)
	if resp.Error() != nil {
		fmt.Println(resp.Error().Error())
		nic.Resources().Logger().Fail(t, resp.Error().Error())
		return
	}
}
