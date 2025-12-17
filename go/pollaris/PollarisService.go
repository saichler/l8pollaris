package pollaris

import (
	"errors"
	"github.com/saichler/l8parser/go/parser/boot"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
)

const (
	ServiceName = "Pollaris"
	ServiceArea = 0
)

type PollarisService struct {
	pollarisCenter *PollarisCenter
	serviceArea    byte
}

func Activate(vnic ifs.IVNic) error {
	initData := []interface{}{}
	for _, p := range boot.GetAllPolarisModels() {
		initData = append(initData, p)
	}
	initData = append(initData, boot.CreateK8sBootPolls())
	sla := ifs.NewServiceLevelAgreement(&PollarisService{}, ServiceName, ServiceArea, true, nil)
	sla.SetServiceItem(&l8tpollaris.L8Pollaris{})
	sla.SetInitItems(initData)
	vnic.Resources().Services().Activate(sla, vnic)
	return nil
}

func (this *PollarisService) Activate(sla *ifs.ServiceLevelAgreement, vnic ifs.IVNic) error {
	vnic.Resources().Registry().Register(&l8tpollaris.L8Pollaris{})
	this.pollarisCenter = newPollarisCenter(sla, vnic)
	this.serviceArea = sla.ServiceArea()
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
	ws := web.New(ServiceName, ServiceArea, 0)
	ws.AddEndpoint(&l8tpollaris.L8Pollaris{}, ifs.POST, &l8web.L8Empty{})
	ws.AddEndpoint(&l8tpollaris.L8Pollaris{}, ifs.PUT, &l8web.L8Empty{})
	return ws
}
