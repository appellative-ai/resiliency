package common

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	RegionZoneHostFmt         = "%v:%v.%v.%v"
	RegionZoneSubZoneHostFmt  = "%v:%v.%v.%v.%v"
	uriFmt                    = "%v:%v"
	RegionZoneSubZoneHostFmt2 = "%v.%v.%v.%v"
)

// Origin - location
type Origin struct {
	Region     string `json:"region"`
	Zone       string `json:"zone"`
	SubZone    string `json:"sub-zone"`
	Host       string `json:"host"`
	Route      string `json:"routing"`
	InstanceId string `json:"instance-id"`
}

func (o Origin) Tag2() string {
	tag := o.Region
	if o.Zone != "" {
		tag += ":" + o.Zone
	}
	if o.SubZone != "" {
		tag += ":" + o.SubZone
	}
	if o.Host != "" {
		tag += ":" + o.Host
	}
	return tag
}

func (o Origin) Uri(class string) string {
	return fmt.Sprintf(uriFmt, class, o)
}

func (o Origin) String() string {
	var uri = o.Region

	if o.Zone != "" {
		uri += "." + o.Zone
	}
	if o.SubZone != "" {
		uri += "." + o.SubZone
	}
	if o.Host != "" {
		uri += "." + o.Host
	}
	if o.Route != "" {
		uri += "." + o.Route
	}
	return uri
}

func NewValues(o Origin) url.Values {
	values := make(url.Values)
	if o.Region != "" {
		values.Add(regionKey, o.Region)
	}
	if o.Zone != "" {
		values.Add(zoneKey, o.Zone)
	}
	if o.SubZone != "" {
		values.Add(subZoneKey, o.SubZone)
	}
	if o.Host != "" {
		values.Add(hostKey, o.Host)
	}
	if o.Route != "" {
		values.Add(routeKey, o.Route)
	}
	return values
}

func NewOrigin(values url.Values) Origin {
	o := Origin{}
	if values != nil {
		o.Region = values.Get(regionKey)
		o.Zone = values.Get(zoneKey)
		o.SubZone = values.Get(subZoneKey)
		o.Host = values.Get(hostKey)
		o.Route = values.Get(routeKey)
	}
	return o
}

func OriginMatch(target Origin, filter Origin) bool {
	isFilter := false
	if filter.Region != "" {
		if filter.Region == "*" {
			return true
		}
		isFilter = true
		if !StringMatch(target.Region, filter.Region) {
			return false
		}
	}
	if filter.Zone != "" {
		isFilter = true
		if !StringMatch(target.Zone, filter.Zone) {
			return false
		}
	}
	if filter.SubZone != "" {
		isFilter = true
		if !StringMatch(target.SubZone, filter.SubZone) {
			return false
		}
	}
	if filter.Host != "" {
		isFilter = true
		if !StringMatch(target.Host, filter.Host) {
			return false
		}
	}
	if filter.Route != "" {
		isFilter = true
		if !StringMatch(target.Route, filter.Route) {
			return false
		}
	}
	return isFilter
}

func StringMatch(target, filter string) bool {
	//if filter == "" {
	//	return true
	//}
	return strings.ToLower(target) == strings.ToLower(filter)
}
