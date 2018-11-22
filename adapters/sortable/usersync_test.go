package sortable

import (
	"testing"

	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
)

func TestSortableSyncer(t *testing.T) {
	config := config.Configuration{ExternalURL: "http://localhost", Adapters: map[string]config.Adapter{
		string(openrtb_ext.BidderSortable): {
			UserSyncURL: "http://localhost/prebid/cookiesync?gdpr={{gdpr}}&gdpr_consent={{gdpr_consent}}&",
		},
	}}
	sortable := NewSortableSyncer(&config)
	syncInfo := sortable.GetUsersyncInfo("0", "")
	assert.Equal(t, "http://localhost/prebid/cookiesync?gdpr=0&gdpr_consent=&redir=http%3A%2F%2Flocalhost%2Fsetuid%3Fbidder%3Dsortable%26gdpr%3D0%26gdpr_consent%3D%26uid%3D%24UID", syncInfo.URL)
	assert.Equal(t, "redirect", syncInfo.Type)
	if sortable.GDPRVendorID() != 145 {
		t.Errorf("Wrong Sortable GDPR VendorID. Got %d", sortable.GDPRVendorID())
	}
}
