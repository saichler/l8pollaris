package targets

import (
	"bytes"
	"github.com/saichler/l8bus/go/overlay/health"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8ql/go/gsql/interpreter"
	"github.com/saichler/l8types/go/ifs"
	"strconv"
	"time"
)

func (this *TargetCallback) InitTargets(vnic ifs.IVNic) {
	leader := vnic.Resources().Services().GetLeader(ServiceName, ServiceArea)
	if leader != vnic.Resources().SysConfig().LocalUuid {
		return
	}
	gsql := "select * from L8PTarget limit 500 page "
	page := 0
	upTargets := make([]*l8tpollaris.L8PTarget, 0)
	for {
		buff := bytes.Buffer{}
		buff.WriteString(gsql)
		buff.WriteString(strconv.Itoa(page))
		q, e := interpreter.NewQuery(buff.String(), vnic.Resources())
		if e != nil {
			panic(e)
		}
		resp := this.iorm.Read(q, vnic.Resources())
		if resp.Error() != nil {
			break
		}
		if resp.Elements() == nil || len(resp.Elements()) == 0 {
			break
		}
		if resp.Element() == nil {
			break
		}
		for _, elem := range resp.Elements() {
			item := elem.(*l8tpollaris.L8PTarget)
			this.validateNewIP(item)
			if item.State == l8tpollaris.L8PTargetState_Up {
				upTargets = append(upTargets, item)
			}
		}
		page++
	}

	go func() {
		time.Sleep(time.Second * 30)
		cService := ""
		cArea := byte(0)
		for _, item := range upTargets {
			if cService == "" {
				cService, cArea = Links.Collector(item.LinksId)
			}
			item.State = l8tpollaris.L8PTargetState_Down
			vnic.Multicast(cService, cArea, ifs.POST, item)
		}

		roundRobin := health.NewRoundRobin(cService, cArea, vnic.Resources())
		for _, item := range upTargets {
			item.State = l8tpollaris.L8PTargetState_Up
			next := roundRobin.Next()
			vnic.Unicast(next, cService, cArea, ifs.POST, item)
		}
	}()
}
