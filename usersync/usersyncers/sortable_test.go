package usersyncers

import (
	"testing"
)

func TestSortableSyncer(t *testing.T) {
	sortable := NewSortableSyncer("http://localhost", "http://localhost/pbs/cookiesync?gdpr={{gdpr}}&gdpr_consent={{gdpr_consent}}&")
	syncInfo := sortable.GetUsersyncInfo("0", "")
	assertStringsMatch(t, "http://localhost/pbs/cookiesync?gdpr=0&gdpr_consent=&redir=http%3A%2F%2Flocalhost%2Fsetuid%3Fbidder%3Dsortable%26gdpr%3D0%26gdpr_consent%3D%26uid%3D%24UID", syncInfo.URL)
	assertStringsMatch(t, "redirect", syncInfo.Type)
	if sortable.GDPRVendorID() != 145 {
		t.Errorf("Wrong Sortable GDPR VendorID. Got %d", sortable.GDPRVendorID())
	}
}
