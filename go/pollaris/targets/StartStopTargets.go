package targets

import (
	"bytes"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8ql/go/gsql/interpreter"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
	"strconv"
)

func (this *TargetCallback) startStopAll(state l8tpollaris.L8PTargetState, typ l8tpollaris.L8PTargetType, vnic ifs.IVNic) {
	gsql := strings.New("select * from L8PTarget where InventoryType=")
	gsql.Add(gsql.StringOf(typ))
	gsql.Add(" and State=")
	switch state {
	case l8tpollaris.L8PTargetState_Up:
		gsql.Add(gsql.StringOf(l8tpollaris.L8PTargetState_Down))
	case l8tpollaris.L8PTargetState_Down:
		gsql.Add(gsql.StringOf(l8tpollaris.L8PTargetState_Up))
	default:
		vnic.Resources().Logger().Error("Not Supported Target State ", state.String())
		return
	}

	collectorService := ""
	collectorArea := byte(0)
	page := 0
	targets := make([]*l8tpollaris.L8PTarget, 0)
	for {
		buff := bytes.Buffer{}
		buff.Write(gsql.Bytes())
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
			patchItem := &l8tpollaris.L8PTarget{}
			patchItem.TargetId = item.TargetId
			patchItem.State = state
			targets = append(targets, patchItem)
			if collectorService == "" {
				collectorService, collectorArea = Links.Collector(item.LinksId)
			}
		}
		page++
	}

	bulk := make([]*l8tpollaris.L8PTarget, 0)
	for _, target := range targets {
		bulk = append(bulk, target)
		if len(bulk) >= 500 {
			elems := object.New(nil, bulk)
			err := this.iorm.Write(ifs.PATCH, elems, vnic.Resources())
			if err != nil {
				vnic.Resources().Logger().Error(err.Error())
			}
			bulk = make([]*l8tpollaris.L8PTarget, 0)
		}
		switch target.State {
		case l8tpollaris.L8PTargetState_Up:
			vnic.RoundRobin(collectorService, collectorArea, ifs.POST, target)
		case l8tpollaris.L8PTargetState_Down:
			vnic.Multicast(collectorService, collectorArea, ifs.POST, target)
		}
	}

	if len(bulk) >= 0 {
		elems := object.New(nil, bulk)
		err := this.iorm.Write(ifs.PATCH, elems, vnic.Resources())
		if err != nil {
			vnic.Resources().Logger().Error(err.Error())
		}
	}
}
