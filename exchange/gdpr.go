package exchange

import (
	"encoding/json"
	"strings"

	"github.com/mxmCherry/openrtb"
)

// ExtractGDPR will pull the gdpr info from an openrtb request
func extractGDPR(bidRequest *openrtb.BidRequest, usersyncIfAmbiguous bool) (gdpr int, consent string, err error) {
	var re regsExt
	if bidRequest.Regs != nil {
		err = json.Unmarshal(bidRequest.Regs.Ext, &re)
	}
	if re.GDPR == nil || err != nil {
		if usersyncIfAmbiguous {
			gdpr = 1
		} else {
			gdpr = 0
		}
	} else {
		gdpr = *re.GDPR
	}
	if err != nil {
		return
	}
	var ue userExt
	if bidRequest.User != nil {
		err = json.Unmarshal(bidRequest.User.Ext, &ue)
	}
	if err != nil {
		return
	}
	consent = ue.Consent
	return
}

type userExt struct {
	Consent string `json:"consent,omitempty"`
}

type regsExt struct {
	GDPR *int `json:"gdpr,omitempty"`
}

// cleanPI removes IP address last byte, device ID, buyer ID, and rounds off lattitude/longitude
func cleanPI(bidRequest *openrtb.BidRequest) {
	if bidRequest.User != nil {
		// Need to duplicate pointer objects
		user := *bidRequest.User
		bidRequest.User = &user
		bidRequest.User.BuyerUID = ""
		if bidRequest.User.Geo != nil {
			bidRequest.User.Geo = cleanGeo(bidRequest.User.Geo)
		}
	}
	if bidRequest.Device != nil {
		// Need to duplicate pointer objects
		device := *bidRequest.Device
		bidRequest.Device = &device
		bidRequest.Device.DIDMD5 = ""
		bidRequest.Device.DIDSHA1 = ""
		bidRequest.Device.DPIDMD5 = ""
		bidRequest.Device.DPIDSHA1 = ""
		bidRequest.Device.IP = cleanIP(bidRequest.Device.IP)
		bidRequest.Device.IPv6 = cleanIPv6(bidRequest.Device.IPv6)
		if bidRequest.Device.Geo != nil {
			bidRequest.Device.Geo = cleanGeo(bidRequest.Device.Geo)
		}
	}
}

// Zero the last byte of an IP address
func cleanIP(fullIP string) string {
	i := strings.LastIndex(fullIP, ".")
	return fullIP[0:i] + ".000"
}

// Zero the last two bytes of an IPv6 address
func cleanIPv6(fullIP string) string {
	i := strings.LastIndex(fullIP, ":")
	return fullIP[0:i] + ":0000"
}

// Return a cleaned Geo object pointer (round off the latitude/longitude)
func cleanGeo(geo *openrtb.Geo) *openrtb.Geo {
	newGeo := *geo
	newGeo.Lat = float64(int(geo.Lat*100.0+0.5)) / 100.0
	newGeo.Lon = float64(int(geo.Lon*100.0+0.5)) / 100.0
	return &newGeo
}
