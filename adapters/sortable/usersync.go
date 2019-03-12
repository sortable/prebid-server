package sortable

import (
	"text/template"

	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/usersync"
)

func NewSortableSyncer(temp *template.Template) usersync.Usersyncer {
	return adapters.NewSyncer("sortable", 145, temp, adapters.SyncTypeRedirect)
}
