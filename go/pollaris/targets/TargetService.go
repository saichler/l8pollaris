package targets

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saichler/l8orm/go/orm/common"
	"github.com/saichler/l8orm/go/orm/persist"
	"github.com/saichler/l8orm/go/orm/plugins/postgres"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8ql/go/gsql/interpreter"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
	"strconv"
	"time"
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

	callback := newTargetCallback()

	sla := ifs.NewServiceLevelAgreement(&persist.OrmService{}, ServiceName, ServiceArea, true, callback)
	sla.SetServiceItem(&l8tpollaris.L8PTarget{})
	sla.SetServiceItemList(&l8tpollaris.L8PTargetList{})
	sla.SetPrimaryKeys("TargetId")
	sla.SetArgs(p)
	sla.SetWebService(web.New(ServiceName, ServiceArea,
		&l8tpollaris.L8PTarget{}, &l8web.L8Empty{},
		nil, nil,
		&l8tpollaris.L8PTarget{}, &l8web.L8Empty{},
		&l8api.L8Query{}, &l8web.L8Empty{},
		&l8api.L8Query{}, &l8tpollaris.L8PTargetList{}))
	vnic.Resources().Services().Activate(sla, vnic)

	InitTargets(p, vnic, callback)
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

func InitTargets(p common.IORM, vnic ifs.IVNic, callback *TargetCallback) {
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
		resp := p.Read(q, vnic.Resources())
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
			callback.validateNewIP(item)
			if item.State == l8tpollaris.L8PTargetState_Up {
				upTargets = append(upTargets, item)
			}
		}
		page++
	}

	go func() {
		time.Sleep(time.Second * 30)
		for _, item := range upTargets {
			collectorService, collectorArea := Links.Collector(item.LinksId)
			item.State = l8tpollaris.L8PTargetState_Down
			vnic.Multicast(collectorService, collectorArea, ifs.POST, item)
		}

		for _, item := range upTargets {
			collectorService, collectorArea := Links.Collector(item.LinksId)
			item.State = l8tpollaris.L8PTargetState_Up
			vnic.RoundRobin(collectorService, collectorArea, ifs.POST, item)
		}
	}()
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
