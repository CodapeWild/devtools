package charset

func Contains(set []string, elem string) bool {
	for _, v := range set {
		if v == elem {
			return true
		}
	}

	return false
}

func Unique(input []string) []string {
	output := make([]string, len(input))
	j := 0
NEXT:
	for i := 0; i < len(input); i++ {
		for _, v := range output {
			if v == input[i] {
				continue NEXT
			}
		}
		output[j] = input[i]
		j++
	}

	return output[:j]
}
