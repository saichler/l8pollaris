package pollaris

import (
	"errors"

	"github.com/saichler/l8pollaris/go/types/l8poll"
	"github.com/saichler/l8reflect/go/reflect/introspecting"
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

func newPollarisCenter(resources ifs.IResources, listener ifs.IServiceCacheListener, initElements []interface{}) *PollarisCenter {
	pc := &PollarisCenter{}
	node, _ := resources.Introspector().Inspect(&l8poll.L8Pollaris{})
	introspecting.AddPrimaryKeyDecorator(node, "Name")
	if initElements != nil {
		resources.Logger().Info("Initializing pollarisCenter with init elements ", len(initElements))
		for _, element := range initElements {
			pc.addInit(element.(*l8poll.L8Pollaris))
		}
	} else {
		resources.Logger().Info("Initializing pollarisCenter with no init elements")
	}
	pc.name2Poll = dcache.NewDistributedCacheNoSync(ServiceName, ServiceArea, &l8poll.L8Pollaris{}, initElements,
		listener, resources)
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
	group := this.groups[name]
	if group != nil {
		result := make(map[string]string)
		for k, v := range group {
			result[k] = v
		}
		return result
	}
	return nil
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

func (this *PollarisCenter) deleteExisting(pollrs *l8poll.L8Pollaris, key string) {
	ePoll, _ := this.name2Poll.Get(pollrs)
	if ePoll == nil {
		return
	}
	existPoll, _ := ePoll.(*l8poll.L8Pollaris)
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

func (this *PollarisCenter) AddAll(pollarises []*l8poll.L8Pollaris) {
	for _, l8pollaris := range pollarises {
		this.Add(l8pollaris, false)
	}
}

func (this *PollarisCenter) Add(l8pollaris *l8poll.L8Pollaris, isNotification bool) error {
	if l8pollaris.Name == "" {
		return errors.New("Pollaris does not contain a Name")
	}
	if l8pollaris.Polling == nil {
		return errors.New("Pollaris does not contain any polling information")
	}

	for _, poll := range l8pollaris.Polling {
		if poll.What == "" {
			return errors.New("Pollaris " + l8pollaris.Name + ": poll does not contain a What value")
		}
	}

	key := this.PollarisKey(l8pollaris)
	_, ok := this.getPollName(key)
	if ok {
		this.deleteExisting(l8pollaris, key)
	}

	this.name2Poll.Put(l8pollaris, isNotification)

	this.mtx.Lock()
	defer this.mtx.Unlock()

	this.key2Name[key] = l8pollaris.Name
	if l8pollaris.Groups != nil {
		for _, gName := range l8pollaris.Groups {
			gEntry, ok := this.groups[gName]
			if !ok {
				this.groups[gName] = make(map[string]string)
				gEntry = this.groups[gName]
			}
			gEntry[key] = l8pollaris.Name
		}
	}
	return nil
}

func (this *PollarisCenter) addInit(p *l8poll.L8Pollaris) {
	key := this.PollarisKey(p)
	this.key2Name[key] = p.Name
	if p.Groups != nil {
		for _, gName := range p.Groups {
			gEntry, ok := this.groups[gName]
			if !ok {
				this.groups[gName] = make(map[string]string)
				gEntry = this.groups[gName]
			}
			gEntry[key] = p.Name
		}
	}
}

func (this *PollarisCenter) Update(l8pollaris *l8poll.L8Pollaris, isNotification bool) error {
	if l8pollaris.Name == "" {
		return errors.New("Pollaris does not contain a Name")
	}
	if l8pollaris.Polling == nil {
		return errors.New("Pollaris does not contain any polling information")
	}

	for _, poll := range l8pollaris.Polling {
		if poll.What == "" {
			return errors.New("Pollaris " + l8pollaris.Name + ": poll does not contain a What value")
		}
	}

	key := this.PollarisKey(l8pollaris)
	_, ok := this.getPollName(key)

	if ok {
		this.deleteExisting(l8pollaris, key)
	}

	this.name2Poll.Put(l8pollaris, isNotification)

	this.mtx.Lock()
	defer this.mtx.Unlock()

	this.key2Name[key] = l8pollaris.Name
	if l8pollaris.Groups != nil {
		for _, gName := range l8pollaris.Groups {
			gEntry, ok := this.groups[gName]
			if !ok {
				this.groups[gName] = make(map[string]string)
				gEntry = this.groups[gName]
			}
			gEntry[key] = l8pollaris.Name
		}
	}
	return nil
}

func (this *PollarisCenter) PollarisKey(l8pollaris *l8poll.L8Pollaris) string {
	return pollarisKey(l8pollaris.Name, l8pollaris.Vendor, l8pollaris.Series, l8pollaris.Family, l8pollaris.Software, l8pollaris.Hardware, l8pollaris.Version)
}

func (this *PollarisCenter) PollarisByName(name string) *l8poll.L8Pollaris {
	if this == nil || this.name2Poll == nil {
		return nil
	}
	filter := &l8poll.L8Pollaris{Name: name}
	p, _ := this.name2Poll.Get(filter)
	poll, _ := p.(*l8poll.L8Pollaris)
	return poll
}

func (this *PollarisCenter) PollarisByKey(args ...string) *l8poll.L8Pollaris {
	if args == nil || len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		pollName := this.key2Name[args[0]]
		filter := &l8poll.L8Pollaris{Name: pollName}
		p, _ := this.name2Poll.Get(filter)
		poll, _ := p.(*l8poll.L8Pollaris)
		return poll
	}
	buff := strings.New()
	buff.Add(args[0])
	for i := 1; i < len(args); i++ {
		addToKey(args[i], buff)
	}
	p, ok := this.getPollName(buff.String())
	if ok {
		filter := &l8poll.L8Pollaris{Name: p}
		f, _ := this.name2Poll.Get(filter)
		poll, _ := f.(*l8poll.L8Pollaris)
		return poll
	}
	return this.PollarisByKey(args[0 : len(args)-1]...)
}

func (this *PollarisCenter) Poll(pollarisName, jobName string) *l8poll.L8Poll {
	l8pollaris := this.PollarisByName(pollarisName)
	if l8pollaris == nil {
		return nil
	}
	poll, ok := l8pollaris.Polling[jobName]
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

func (this *PollarisCenter) PollsByGroup(groupName, vendor, series, family, software, hardware, version string) []*l8poll.L8Pollaris {
	names := this.Names(groupName, vendor, series, family, software, hardware, version)
	result := make([]*l8poll.L8Pollaris, 0)
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

func Poll(pollarisName, pollName string, resources ifs.IResources) (*l8poll.L8Poll, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris service")
	}
	l8pollaris := pollarisCenter.PollarisByName(pollarisName)
	if l8pollaris == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris " + pollName)
	}
	poll, ok := l8pollaris.Polling[pollName]
	if !ok {
		return nil, resources.Logger().Error("Cannot find poll " + pollName + " in Pollaris " + l8pollaris.Name)
	}
	return poll, nil
}

func PollarisByKey(resources ifs.IResources, args ...string) (*l8poll.L8Pollaris, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris service")
	}
	p := pollarisCenter.PollarisByKey(args...)
	if p == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris for keys ", args)
	}
	return p, nil
}

func PollarisByGroup(resources ifs.IResources, groupName, vendor, series, family, software, hardware, version string) ([]*l8poll.L8Pollaris, error) {
	pollarisCenter := Pollaris(resources)
	if pollarisCenter == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris service")
	}
	p := pollarisCenter.PollsByGroup(groupName, vendor, series, family, software, hardware, version)
	if p == nil {
		return nil, resources.Logger().Error("Cannot find Pollaris for group ", groupName)
	}
	return p, nil
}
