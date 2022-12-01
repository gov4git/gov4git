package common

func AdsToBallotNames(ads []Advertisement) []string {
	names := make([]string, len(ads))
	for i := range ads {
		names[i] = ads[i].Name.Path()
	}
	return names
}
