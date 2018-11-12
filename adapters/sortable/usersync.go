package sortable

import (
	"net/url"

	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/prebid/prebid-server/usersync"
)

func NewSortableSyncer(cfg *config.Configuration) usersync.Usersyncer {
	redirectURI := url.QueryEscape(cfg.ExternalURL) + "%2Fsetuid%3Fbidder%3Dsortable%26gdpr%3D{{gdpr}}%26gdpr_consent%3D{{gdpr_consent}}%26uid%3D%24UID"
	usersyncURL := cfg.Adapters[string(openrtb_ext.BidderSortable)].UserSyncURL
	return adapters.NewSyncer(
		"sortable",
		145,
		adapters.ResolveMacros(usersyncURL+"redir="+redirectURI),
		adapters.SyncTypeRedirect)
}
