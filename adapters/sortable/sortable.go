package sortable

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type SortableAdapter struct {
	http *adapters.HTTPAdapter
	URI  string
}

type impExts struct {
	Bidder map[string]interface{} `json:"bidder"`
}

func NewSortableBidder(client *http.Client, endpoint string) *SortableAdapter {
	a := &adapters.HTTPAdapter{Client: client}

	return &SortableAdapter{
		http: a,
		URI:  endpoint,
	}
}

func isValidRequest(request *openrtb.BidRequest) bool {
	return !((request.Site == nil || request.Site.Publisher == nil || request.Site.Publisher.ID == "") &&
		(request.App == nil || request.App.Bundle == "")) &&
		// These fields should always be here, unless PBS couldn't infer them from the request
		!(request.Device == nil || request.Device.IP == "" || request.Device.UA == "")
}

func (s *SortableAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {
	errs := make([]error, 0, len(request.Imp))
	if !isValidRequest(request) {
		errs = append(errs, errors.New("Sortable requires site.publisher.id or app.bundle to be set, and that both device.ip and device.ua are set"))
		return nil, errs
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	// Hoist the contents of ext.bidder up one level
	for i, imp := range request.Imp {
		var extStuff impExts
		err := json.Unmarshal(imp.Ext, &extStuff)
		if err != nil {
			errs = append(errs, err)
		}
		// Further hoist certain fields up one level further
		if request.Imp[i].BidFloor == 0 {
			maybeBidfloor := extStuff.Bidder["bidfloor"]
			bidfloor, ok := maybeBidfloor.(float64)
			if ok {
				request.Imp[i].BidFloor = bidfloor
				delete(extStuff.Bidder, "bidfloor")
			}
		}

		if request.Imp[i].TagID == "" {
			maybeTagid := extStuff.Bidder["tagid"]
			tagid, ok := maybeTagid.(string)
			if ok {
				request.Imp[i].TagID = tagid
				delete(extStuff.Bidder, "tagid")
			}
		}
		marshalled, err := json.Marshal(extStuff.Bidder)
		if err != nil {
			errs = append(errs, err)
		} else {
			request.Imp[i].Ext = marshalled
		}

	}

	reqJSON, err := json.Marshal(request)
	if err != nil {
		errs = append(errs, err)
	}

	if request.User != nil && request.User.BuyerUID != "" {
		var cookies map[string]string
		err := json.Unmarshal([]byte(request.User.BuyerUID), &cookies)
		if err != nil {
			errs = append(errs, err)
		} else {
			for key, value := range cookies {
				headers.Add("Cookie", "d7s_"+key+"="+value)
			}
		}
	}
	return []*adapters.RequestData{{
		Method:  "POST",
		Uri:     s.URI,
		Body:    reqJSON,
		Headers: headers,
	}}, errs
}

func (s *SortableAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(5)

	var bidResp openrtb.BidResponse
	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{&errortypes.BadServerResponse{
			Message: err.Error(),
		}}
	}

	for _, sb := range bidResp.SeatBid {
		for i := 0; i < len(sb.Bid); i++ {
			bid := sb.Bid[i]
			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &bid,
				BidType: openrtb_ext.BidTypeBanner,
			})
		}
	}
	return bidResponse, nil
}
