package targets

import (
	"errors"
	"github.com/saichler/l8orm/go/orm/common"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"sync"
)

type TargetCallback struct {
	addressValidation *sync.Map
	iorm              common.IORM
}

func newTargetCallback(iorm common.IORM) *TargetCallback {
	return &TargetCallback{addressValidation: &sync.Map{}, iorm: iorm}
}

func (this *TargetCallback) Before(elem interface{}, action ifs.Action, notification bool, vnic ifs.IVNic) (interface{}, error) {
	switch action {
	case ifs.POST:
		if !notification {
			list, ok := elem.(*l8tpollaris.L8PTargetList)
			if ok {
				elems := make([]interface{}, 0)
				for _, item := range list.List {
					err := this.validateNewIP(item)
					if err != nil {
						return nil, err
					}
					elems = append(elems, item)
				}
				return elems, nil
			}
			item, ok := elem.(*l8tpollaris.L8PTarget)
			if ok {
				err := this.validateNewIP(item)
				if err != nil {
					return nil, err
				}
			}
		}
	case ifs.PATCH:
		if !notification {
			target, ok := elem.(*l8tpollaris.L8PTarget)
			if !ok {
				return nil, errors.New("invalid target")
			}
			_, err := Target(target.TargetId, vnic)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (this *TargetCallback) After(elem interface{}, action ifs.Action, notification bool, vnic ifs.IVNic) (interface{}, error) {
	if action == ifs.POST && !notification {
		target, ok := elem.(*l8tpollaris.L8PTarget)
		if !ok {
			return nil, errors.New("invalid target")
		}
		if target.State == l8tpollaris.L8PTargetState_Up {
			collectorService, collectorArea := Links.Collector(target.LinksId)
			vnic.Resources().Logger().Info("Sending target to collector:", collectorService, " area ", collectorArea)
			err := vnic.RoundRobin(collectorService, collectorArea, ifs.POST, target)
			if err != nil {
				return nil, err
			}
		}
	}
	if action == ifs.PATCH && !notification {
		target, ok := elem.(*l8tpollaris.L8PTarget)
		if !ok {
			return nil, errors.New("invalid target")
		}

		currTarget, err := Target(target.TargetId, vnic)
		if err != nil {
			return nil, err
		}
		collectorService, collectorArea := Links.Collector(currTarget.LinksId)

		switch target.State {
		case l8tpollaris.L8PTargetState_Down:
			vnic.Resources().Logger().Info("Sending stop target to collector:", collectorService, " area ", collectorArea,
				" with hosts ", currTarget.Hosts)
			err = vnic.Multicast(collectorService, collectorArea, ifs.POST, currTarget)
			if err != nil {
				return nil, err
			}
		case l8tpollaris.L8PTargetState_Up:
			vnic.Resources().Logger().Info("Sending start target to collector:", collectorService, " area ", collectorArea,
				" with hosts ", currTarget.Hosts)
			err = vnic.RoundRobin(collectorService, collectorArea, ifs.POST, currTarget)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (this *TargetCallback) validateNewIP(target *l8tpollaris.L8PTarget) error {
	if target.Hosts == nil || len(target.Hosts) == 0 {
		return errors.New("invalid target, has no hosts")
	}
	for _, host := range target.Hosts {
		if host.Configs == nil || len(host.Configs) == 0 {
			return errors.New("invalid target, host has no configs")
		}
		for _, p := range host.Configs {
			_, ok := this.addressValidation.Load(p.Addr)
			if ok {
				return errors.New("invalid target, address " + p.Addr + " already exists")
			}
		}
		for _, p := range host.Configs {
			this.addressValidation.Store(p.Addr, true)
		}
	}
	return nil
}
