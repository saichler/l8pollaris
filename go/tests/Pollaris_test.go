package tests

import (
	"testing"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/boot"
	"github.com/saichler/l8pollaris/go/pollaris"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestPollaris(t *testing.T) {
	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	vnic.Resources().Services().Activate(pollaris.ServiceType, pollaris.ServiceName, 0, vnic.Resources(), vnic)
	p := pollaris.Pollaris(vnic.Resources())
	pollrs := boot.CreateBoot01()
	err := p.Post(pollrs, false)
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}

	byName := p.PollarisByName(pollrs.Name)
	if byName == nil {
		vnic.Resources().Logger().Fail(t, "No such pollaris")
		return
	}

	byGroup, err := pollaris.PollarisByGroup(vnic.Resources(), common.BOOT_STAGE_01,
		"", "", "", "", "", "")
	if byGroup == nil {
		vnic.Resources().Logger().Fail(t, "No such group")
		return
	}
}
