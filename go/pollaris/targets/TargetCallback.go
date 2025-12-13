package targets

import (
	"errors"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
)

type TargetCallback struct{}

func (this *TargetCallback) Before(elem interface{}, action ifs.Action, notification bool, vnic ifs.IVNic) (interface{}, error) {
	if action == ifs.POST && !notification {
		list, ok := elem.(*l8tpollaris.L8PTargetList)
		if ok {
			elems := make([]interface{}, len(list.List))
			for i, item := range list.List {
				elems[i] = item
			}
			return elems, nil
		}
	}
	if action == ifs.PATCH && !notification {
		target, ok := elem.(*l8tpollaris.L8PTarget)
		if !ok {
			return nil, errors.New("invalid target")
		}
		_, err := Target(target.TargetId, vnic)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (this *TargetCallback) After(elem interface{}, action ifs.Action, notification bool, vnic ifs.IVNic) (interface{}, error) {
	if action == ifs.PATCH && !notification {
		target, ok := elem.(*l8tpollaris.L8PTarget)
		if !ok {
			return nil, errors.New("invalid target")
		}
		if target.State == l8tpollaris.L8TargetState_Up {
			realTarget, err := Target(target.TargetId, vnic)
			if err != nil {
				return nil, err
			}
			collectorService, collectorArea := Links.Collector(realTarget.LinksId)
			err = vnic.RoundRobin(collectorService, collectorArea, ifs.POST, realTarget)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
