package targets

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saichler/l8orm/go/orm/persist"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
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
	p := persist.NewPostgres(db, vnic.Resources())

	sla := ifs.NewServiceLevelAgreement(&persist.OrmService{}, ServiceName, ServiceArea, true, nil)
	sla.SetServiceItem(&l8tpollaris.L8PTarget{})
	sla.SetServiceItemList(&l8tpollaris.L8PTargetList{})
	sla.SetPrimaryKeys("TargetId")
	sla.SetArgs(p)

	vnic.Resources().Services().Activate(sla, vnic)
}

func Targets(r ifs.IResources) (ifs.IServiceHandler, bool) {
	return r.Services().ServiceHandler(ServiceName, ServiceArea)
}

func openDBConection(dbname, user, pass string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, user, pass, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	fmt.Println("dbname=" + dbname + ";user=" + user + ";pass=" + pass)
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
