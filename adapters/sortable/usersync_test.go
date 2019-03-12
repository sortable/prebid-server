package sortable

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestSortableSyncer(t *testing.T) {
	temp := template.Must(template.New("sync-template").Parse("//c.deployads.com/prebid/cookiesync?gdpr={{.GDPR}}&gdpr_consent={{.GDPRConsent}}&redir=localhost%2Fsetuid%3Fbidder%3Dsortable%26gdpr%3D{{.GDPR}}%26gdpr_consent%3D{{.GDPRConsent}}"))
	syncer := NewSortableSyncer(temp)
	syncInfo, err := syncer.GetUsersyncInfo("1", "BOPVK28OVJoTBABABAENBs-AAAAhuAKAANAAoACwAGgAPAAxAB0AHgAQAAiABOADkA")
	assert.NoError(t, err)
	assert.Equal(t, "//c.deployads.com/prebid/cookiesync?gdpr=1&gdpr_consent=BOPVK28OVJoTBABABAENBs-AAAAhuAKAANAAoACwAGgAPAAxAB0AHgAQAAiABOADkA&redir=localhost%2Fsetuid%3Fbidder%3Dsortable%26gdpr%3D1%26gdpr_consent%3DBOPVK28OVJoTBABABAENBs-AAAAhuAKAANAAoACwAGgAPAAxAB0AHgAQAAiABOADkA", syncInfo.URL)
	assert.Equal(t, "redirect", syncInfo.Type)
	assert.EqualValues(t, 145, syncer.GDPRVendorID())
	assert.Equal(t, false, syncInfo.SupportCORS)
}
