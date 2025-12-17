package targets

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saichler/l8orm/go/orm/persist"
	"github.com/saichler/l8orm/go/orm/plugins/postgres"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
)

const (
	ServiceName = "Targets"
	ServiceArea = byte(91)
)

var Links TargetLinks

func Activate(creds, dbname string, vnic ifs.IVNic) {
	_, user, pass, _, err := vnic.Resources().Security().Credential(creds, dbname, vnic.Resources())
	if err != nil {
		panic(err)
	}
	db := openDBConection(dbname, user, pass)
	p := postgres.NewPostgres(db, vnic.Resources())

	callback := newTargetCallback(p)

	sla := ifs.NewServiceLevelAgreement(&persist.OrmService{}, ServiceName, ServiceArea, true, callback)
	sla.SetServiceItem(&l8tpollaris.L8PTarget{})
	sla.SetServiceItemList(&l8tpollaris.L8PTargetList{})
	sla.SetPrimaryKeys("TargetId")
	sla.SetArgs(p)

	ws := web.New(ServiceName, ServiceArea, 0)
	ws.AddEndpoint(&l8tpollaris.L8PTarget{}, ifs.POST, &l8web.L8Empty{})
	ws.AddEndpoint(&l8tpollaris.L8PTarget{}, ifs.PUT, &l8web.L8Empty{})
	ws.AddEndpoint(&l8tpollaris.L8PTarget{}, ifs.PATCH, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Query{}, ifs.DELETE, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Query{}, ifs.GET, &l8tpollaris.L8PTargetList{})
	sla.SetWebService(ws)

	vnic.Resources().Services().Activate(sla, vnic)

	callback.InitTargets(vnic)
}

func Targets(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return vnic.Resources().Services().ServiceHandler(ServiceName, ServiceArea)
}

func Target(targetId string, vnic ifs.IVNic) (*l8tpollaris.L8PTarget, error) {
	this, ok := Targets(vnic)
	if !ok {
		return nil, errors.New("No Targets Service Found")
	}
	filter := &l8tpollaris.L8PTarget{TargetId: targetId}
	resp := this.Get(object.New(nil, filter), vnic)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return resp.Element().(*l8tpollaris.L8PTarget), nil
}

func openDBConection(dbname, user, pass string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, user, pass, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	return db
}

/*
   func (this *TargetService) startDevice(device *l8tpollaris.L8PTarget, vnic ifs.IVNic, isNotificaton bool) {
   	vnic.Resources().Logger().Info("TargetService.startDevice: ", device.TargetId)
   	if !isNotificaton {
   		err := vnic.RoundRobin(common.CollectorService, this.serviceArea, ifs.POST, device)
   		if err != nil {
   			vnic.Resources().Logger().Error("Device Service:", err.Error())
   		}
   	}
   }

   func (this *TargetService) updateDevice(device *l8tpollaris.L8PTarget, vnic ifs.IVNic, isNotificaton bool) {
   	vnic.Resources().Logger().Info("TargetService.startDevice: ", device.TargetId)
   	if !isNotificaton {
   		err := vnic.Multicast(common.CollectorService, this.serviceArea, ifs.PUT, device)
   		if err != nil {
   			vnic.Resources().Logger().Error("Device Service:", " ", err.Error())
   		}
   	}
   }

   func (this *TargetService) stopDevice(device *l8tpollaris.L8PTarget, vnic ifs.IVNic, isNotificaton bool) {
   	if !isNotificaton {
   		err := vnic.Multicast(common.CollectorService, this.serviceArea, ifs.DELETE, device)
   		if err != nil {
   			vnic.Resources().Logger().Error("Device Service:", " ", err.Error())
   		}
   	}
   }*/
