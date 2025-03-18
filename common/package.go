package common

import "github.com/behavioral-ai/core/messaging"

const (
	WestRegion = "us-west1"
	WestZoneA  = "w-zone-a"
	WestZoneB  = "w-zone-b"

	CentralRegion = "us-central1"
	CentralZoneA  = "c-zone-a"
	CentralZoneB  = "c-zone-b"

	EastRegion = "us-east1"
	EastZoneA  = "e-zone-a"
	EastZoneB  = "e-zone-b"
)

type NewAgentFunc func(origin Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher)
