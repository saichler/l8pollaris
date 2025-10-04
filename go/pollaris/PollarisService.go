package pollaris

import (
	"errors"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
)

const (
	ServiceName = "Pollaris"
	ServiceArea = 0
	ServiceType = "PollarisService"
)

type PollarisService struct {
	pollarisCenter *PollarisCenter
	serviceArea    byte
}

func (this *PollarisService) Activate(serviceName string, serviceArea byte,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {
	r.Registry().Register(&l8tpollaris.L8Pollaris{})
	var data []interface{}
	if args != nil {
		d, ok := args[0].([]interface{})
		if ok {
			data = d
		}
	}
	this.pollarisCenter = newPollarisCenter(r, l, data)
	this.serviceArea = serviceArea
	return nil
}

func (this *PollarisService) DeActivate() error {
	this.pollarisCenter = nil
	return nil
}

func (this *PollarisService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	var err error
	for _, elem := range pb.Elements() {
		l8Pollaris, ok := elem.(*l8tpollaris.L8Pollaris)
		if ok {
			vnic.Resources().Logger().Info("Added a l8Pollaris ", l8Pollaris.Name)
			e := this.pollarisCenter.Post(l8Pollaris, pb.Notification())
			if e != nil {
				err = e
			}
		} else {
			err = errors.New("Element is not a L8Pollaris")
		}
	}
	return object.New(err, &l8web.L8Empty{})
}
func (this *PollarisService) Put(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	var err error
	for _, elem := range pb.Elements() {
		l8Pollaris, ok := elem.(*l8tpollaris.L8Pollaris)
		if ok {
			vnic.Resources().Logger().Info("Added a l8Pollaris ", l8Pollaris.Name)
			e := this.pollarisCenter.Put(l8Pollaris, pb.Notification())
			if e != nil {
				err = e
			}
		} else {
			err = errors.New("Element is not a L8Pollaris")
		}
	}
	return object.New(err, &l8web.L8Empty{})
}
func (this *PollarisService) Patch(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *PollarisService) Delete(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *PollarisService) Get(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *PollarisService) GetCopy(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *PollarisService) Failed(pb ifs.IElements, vnic ifs.IVNic, msg *ifs.Message) ifs.IElements {
	return nil
}
func (this *PollarisService) TransactionConfig() ifs.ITransactionConfig {
	return nil
}
func (this *PollarisService) WebService() ifs.IWebService {
	ws := web.New(ServiceName, this.serviceArea,
		&l8tpollaris.L8Pollaris{}, &l8web.L8Empty{},
		&l8tpollaris.L8Pollaris{}, &l8web.L8Empty{},
		nil, nil,
		nil, nil,
		nil, nil)
	return ws
}
