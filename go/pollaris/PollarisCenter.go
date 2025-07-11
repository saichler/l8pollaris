package pollaris

import (
	"errors"
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8services/go/services/dcache"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
	"sync"
)

type PollarisCenter struct {
	name2Poll ifs.IDistributedCache
	key2Name  map[string]string
	groups    map[string]map[string]string
	log       ifs.ILogger
	mtx       *sync.RWMutex
}

func newPollarisCenter(resources ifs.IResources, listener ifs.IServiceCacheListener) *PollarisCenter {
	pc := &PollarisCenter{}
	pc.name2Poll = dcache.NewDistributedCache(ServiceName, ServiceArea, "Pollaris",
		resources.SysConfig().LocalUuid, listener, resources)
	pc.key2Name = make(map[string]string)
	pc.groups = make(map[string]map[string]string)
	pc.log = resources.Logger()
	pc.mtx = &sync.RWMutex{}
	return pc
}

func (this *PollarisCenter) getPollName(key string) (string, bool) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	pollName, ok := this.key2Name[key]
	return pollName, ok
}

func (this *PollarisCenter) getGroup(name string) map[string]string {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	return this.groups[name]
}

func (this *PollarisCenter) deleteFromGroup(gEntry map[string]string, key string) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	delete(gEntry, key)
}

func (this *PollarisCenter) deleteFromKey2Name(key string) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	delete(this.key2Name, key)
}

func (this *PollarisCenter) deleteExisting(pollaris *types.Pollaris, key string) {
	existPoll := this.name2Poll.Get(pollaris.Name).(*types.Pollaris)
	if existPoll.Groups != nil {
		for _, gName := range existPoll.Groups {
			gEntry := this.getGroup(gName)
			if gEntry != nil {
				this.deleteFromGroup(gEntry, key)
			}
		}
	}
	this.deleteFromKey2Name(key)
}

func (this *PollarisCenter) AddAll(pollarises []*types.Pollaris) {
	for _, pollaris := range pollarises {
		this.Add(pollaris, false)
	}
}

func (this *PollarisCenter) Add(pollaris *types.Pollaris, isNotification bool) error {
	if pollaris.Name == "" {
		return errors.New("pollaris does not contain a Name")
	}
	if pollaris.Polling == nil {
		return errors.New("pollaris does not contain any polling information")
	}

	for _, poll := range pollaris.Polling {
		if poll.What == "" {
			return errors.New("pollaris " + pollaris.Name + ": poll does not contain a What value")
		}
	}

	key := this.PollarisKey(pollaris)
	_, ok := this.getPollName(key)
	if ok {
		this.deleteExisting(pollaris, key)
	}

	this.name2Poll.Put(pollaris.Name, pollaris, isNotification)

	this.mtx.Lock()
	defer this.mtx.Unlock()

	this.key2Name[key] = pollaris.Name
	if pollaris.Groups != nil {
		for _, gName := range pollaris.Groups {
			gEntry, ok := this.groups[gName]
			if !ok {
				this.groups[gName] = make(map[string]string)
				gEntry = this.groups[gName]
			}
			gEntry[key] = pollaris.Name
		}
	}
	return nil
}

func (this *PollarisCenter) Update(pollaris *types.Pollaris, isNotification bool) error {
	if pollaris.Name == "" {
		return errors.New("pollaris does not contain a Name")
	}
	if pollaris.Polling == nil {
		return errors.New("pollaris does not contain any polling information")
	}

	for _, poll := range pollaris.Polling {
		if poll.What == "" {
			return errors.New("pollaris " + pollaris.Name + ": poll does not contain a What value")
		}
	}

	key := this.PollarisKey(pollaris)
	_, ok := this.getPollName(key)

	if ok {
		this.deleteExisting(pollaris, key)
	}

	this.name2Poll.Put(pollaris.Name, pollaris, isNotification)

	this.mtx.Lock()
	defer this.mtx.Unlock()

	this.key2Name[key] = pollaris.Name
	if pollaris.Groups != nil {
		for _, gName := range pollaris.Groups {
			gEntry, ok := this.groups[gName]
			if !ok {
				this.groups[gName] = make(map[string]string)
				gEntry = this.groups[gName]
			}
			gEntry[key] = pollaris.Name
		}
	}
	return nil
}

func (this *PollarisCenter) PollarisKey(pollaris *types.Pollaris) string {
	return pollarisKey(pollaris.Name, pollaris.Vendor, pollaris.Series, pollaris.Family, pollaris.Software, pollaris.Hardware, pollaris.Version)
}

func (this *PollarisCenter) PollarisByName(name string) *types.Pollaris {
	if this == nil || this.name2Poll == nil {
		return nil
	}
	poll, _ := this.name2Poll.Get(name).(*types.Pollaris)
	return poll
}

func (this *PollarisCenter) PollarisByKey(args ...string) *types.Pollaris {
	if args == nil || len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		pollName := this.key2Name[args[0]]
		poll, _ := this.name2Poll.Get(pollName).(*types.Pollaris)
		return poll
	}
	buff := strings.New()
	buff.Add(args[0])
	for i := 1; i < len(args); i++ {
		addToKey(args[i], buff)
	}
	p, ok := this.getPollName(buff.String())
	if ok {
		poll, _ := this.name2Poll.Get(p).(*types.Pollaris)
		return poll
	}
	return this.PollarisByKey(args[0 : len(args)-1]...)
}

func (this *PollarisCenter) Poll(pollarisName, jobName string) *types.Poll {
	pollaris := this.PollarisByName(pollarisName)
	if pollaris == nil {
		return nil
	}
	poll, ok := pollaris.Polling[jobName]
	if !ok {
		return nil
	}
	return poll
}

func (this *PollarisCenter) Names(groupName, vendor, series, family, software, hardware, version string) []string {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]string, 0)
	group, ok := this.groups[groupName]
	if !ok {
		return result
	}
	for _, name := range group {
		result = append(result, name)
	}
	return result
}

func (this *PollarisCenter) PollsByGroup(groupName, vendor, series, family, software, hardware, version string) []*types.Pollaris {
	names := this.Names(groupName, vendor, series, family, software, hardware, version)
	result := make([]*types.Pollaris, 0)
	for _, name := range names {
		poll := this.PollarisByKey(name, vendor, series, family, software, hardware, version)
		if poll != nil {
			result = append(result, poll)
		}
	}
	return result
}

func Pollaris(resource ifs.IResources) *PollarisCenter {
	sp, ok := resource.Services().ServiceHandler(ServiceName, ServiceArea)
	if !ok {
		return nil
	}
	return (sp.(*PollarisService)).pollarisCenter
}

func Poll(pollarisName, pollName string, resources ifs.IResources) (*types.Poll, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find pollaris service")
	}
	pollaris := pollarisCenter.PollarisByName(pollarisName)
	if pollaris == nil {
		return nil, resources.Logger().Error("Cannot find pollaris " + pollName)
	}
	poll, ok := pollaris.Polling[pollName]
	if !ok {
		return nil, resources.Logger().Error("Cannot find poll " + pollName + " in pollaris " + pollaris.Name)
	}
	return poll, nil
}

func PollarisByKey(resources ifs.IResources, args ...string) (*types.Pollaris, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find pollaris service")
	}
	p := pollarisCenter.PollarisByKey(args...)
	if p == nil {
		return nil, resources.Logger().Error("Cannot find pollaris for keys ", args)
	}
	return p, nil
}

func PollarisByGroup(resources ifs.IResources, groupName, vendor, series, family, software, hardware, version string) ([]*types.Pollaris, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find pollaris service")
	}
	p := pollarisCenter.PollsByGroup(groupName, vendor, series, family, software, hardware, version)
	if p == nil {
		return nil, resources.Logger().Error("Cannot find pollaris for group ", groupName)
	}
	return p, nil
}
