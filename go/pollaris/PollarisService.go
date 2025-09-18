package pollaris

import (
	"github.com/saichler/l8pollaris/go/types/l8poll"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
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
	r.Registry().Register(&l8poll.L8Pollaris{})
	this.pollarisCenter = newPollarisCenter(r, l)
	this.serviceArea = serviceArea
	return nil
}

func (this *PollarisService) DeActivate() error {
	this.pollarisCenter = nil
	return nil
}

func (this *PollarisService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	l8Pollaris := pb.Element().(*l8poll.L8Pollaris)
	vnic.Resources().Logger().Info("Added a l8Pollaris ", l8Pollaris.Name)
	return object.New(this.pollarisCenter.Add(l8Pollaris, pb.Notification()), &l8poll.L8Pollaris{})
}
func (this *PollarisService) Put(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
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
	ws := web.New(ServiceName, this.serviceArea, &l8poll.L8Pollaris{},
		&l8poll.L8Pollaris{}, nil, nil, nil, nil, nil, nil, nil, nil)
	return ws
}
