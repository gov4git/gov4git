package common

func AdsToBallotNames(ads []Advertisement) []string {
	names := make([]string, len(ads))
	for i := range ads {
		names[i] = ads[i].Name.Path()
	}
	return names
}

func FilterOpenClosedAds(closed bool, ads []Advertisement) []Advertisement {
	r := []Advertisement{}
	for _, ad := range ads {
		if ad.Closed == closed {
			r = append(r, ad)
		}
	}
	return r
}
