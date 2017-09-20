package util

func ElementsNotIn(array []string, knownElements []string) []string {
	extraElements := []string{}

	for _, element := range array {
		if !Contains(knownElements, element) {
			extraElements = append(extraElements, element)
		}
	}

	return extraElements
}

func Contains(array []string, element string) bool {
	for _, e := range array {
		if element == e {
			return true;
		}
	}
	return false
}