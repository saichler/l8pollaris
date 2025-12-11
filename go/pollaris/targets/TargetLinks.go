package targets

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8services"
	"sync"
)

type TargetLinks struct {
	links *sync.Map
}

func newTargetLinks() *TargetLinks {
	return &TargetLinks{
		links: &sync.Map{},
	}
}

func (this *TargetLinks) AddLink(linksId,
	collectorService string, collectorArea byte,
	parserService string, parserArea byte,
	cacheService string, cacheArea byte,
	persistService string, persistArea byte) {
	links := &l8tpollaris.L8TargetLinks{}
	links.LinksId = linksId
	links.Collector = &l8services.L8ServiceLink{ZsideServiceName: collectorService, ZsideServiceArea: int32(collectorArea), Mode: int32(ifs.M_Proximity), Interval: 5}
	links.Parser = &l8services.L8ServiceLink{ZsideServiceName: parserService, ZsideServiceArea: int32(parserArea),
		AsideServiceName: collectorService, AsideServiceArea: int32(collectorArea), Mode: int32(ifs.M_Proximity), Interval: 5}
	links.Cache = &l8services.L8ServiceLink{ZsideServiceName: cacheService, ZsideServiceArea: int32(cacheArea),
		AsideServiceName: parserService, AsideServiceArea: int32(parserArea), Mode: int32(ifs.M_Proximity), Interval: 5}
	links.Persistency = &l8services.L8ServiceLink{ZsideServiceName: persistService, ZsideServiceArea: int32(persistArea),
		AsideServiceName: parserService, AsideServiceArea: int32(parserArea), Mode: int32(ifs.M_Proximity), Interval: 5}
	this.links.Store(linksId, links)
}

func (this *TargetLinks) CollectorLink(linksId string) (string, byte) {
	l, ok := this.links.Load(linksId)
	if ok {
		links := l.(*l8tpollaris.L8TargetLinks)
		return links.Collector.ZsideServiceName, byte(links.Collector.ZsideServiceArea)
	}
	return "", 0
}

func (this *TargetLinks) ParserLink(linksId string) (string, byte) {
	l, ok := this.links.Load(linksId)
	if ok {
		links := l.(*l8tpollaris.L8TargetLinks)
		return links.Parser.ZsideServiceName, byte(links.Parser.ZsideServiceArea)
	}
	return "", 0
}

func (this *TargetLinks) CacheLink(linksId string) (string, byte) {
	l, ok := this.links.Load(linksId)
	if ok {
		links := l.(*l8tpollaris.L8TargetLinks)
		return links.Cache.ZsideServiceName, byte(links.Cache.ZsideServiceArea)
	}
	return "", 0
}

func (this *TargetLinks) PersistLink(linksId string) (string, byte) {
	l, ok := this.links.Load(linksId)
	if ok {
		links := l.(*l8tpollaris.L8TargetLinks)
		return links.Persistency.ZsideServiceName, byte(links.Persistency.ZsideServiceArea)
	}
	return "", 0
}
