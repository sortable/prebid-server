package sortable

import (
	"testing"

	"net/http"

	"github.com/prebid/prebid-server/adapters/adapterstest"
)

func TestJsonSamples(t *testing.T) {
	sortableAdapter := NewSortableBidder(new(http.Client), "http://c.deployads.com/openrtb2/auction")
	adapterstest.RunJSONBidderTest(t, "sortabletest", sortableAdapter)
}
