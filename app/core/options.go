package core

import "fmt"

type AppOptions struct {
	Port            int
	Datastore       string
	ExtraParams     string
	HideExtraParams bool
}

var Options AppOptions

func (o AppOptions) String(hideExtraParams bool) string {
	var extraParams = ""

	if !hideExtraParams {
		extraParams = o.ExtraParams
	}

	return fmt.Sprintf(
		"Options\n  Port: %d\n  Datastore: %s\n  ExtraParams: %s\n  HideExtraParams: %t\n",
		o.Port,
		o.Datastore,
		extraParams,
		o.HideExtraParams)
}
