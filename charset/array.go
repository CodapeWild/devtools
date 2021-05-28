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

func Intersect(set1, set2 []string) []string {
	m := make(map[string]int)
	for _, v := range set1 {
		m[v] = 0
	}

	var inter []string
	for _, v := range set2 {
		if count := m[v]; count == 0 {
			inter = append(inter, v)
			m[v]++
		}
	}

	return inter
}

func Union(set1, set2 []string) []string {
	m := make(map[string]struct{}, len(set1))
	for _, v := range set1 {
		m[v] = struct{}{}
	}

	for _, v := range set2 {
		if _, found := m[v]; !found {
			set1 = append(set1, v)
		}
	}

	return set1
}

func Differ(source, compare []string) []string {
	m := make(map[string]struct{}, len(compare))
	for _, v := range compare {
		m[v] = struct{}{}
	}

	var diff []string
	for _, v := range source {
		if _, found := m[v]; !found {
			diff = append(diff, v)
		}
	}

	return diff
}
