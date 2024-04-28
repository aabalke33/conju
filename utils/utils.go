package utils

import (
	"math/rand"
)

func ChooseVerb(
	verbs []map[string]string, keepPronouns map[string]bool) (
	verb map[string]string, pov, pronoun string) {

	var povs []string

	for pronoun, keepPronoun := range keepPronouns {
		if keepPronoun {
			povs = append(povs, pronoun)
		}
	}

	pronouns := map[string][]string{
		"first_single":  {"yo"},
		"first_plural":  {"nosotros"},
		"second_single": {"tu"},
		"second_plural": {"vosotros"},
		"third_single":  {"él", "ella", "usted"},
		"third_plural":  {"ellos", "ellas", "ustedes"},
	}

	idxVerb := rand.Int() % len(verbs)
	idxPov := rand.Int() % len(povs)

	verb = verbs[idxVerb]
	pov = povs[idxPov]
	idxPronoun := rand.Int() % len(pronouns[pov])
	pronoun = pronouns[pov][idxPronoun]

	return verb, pov, pronoun
}
