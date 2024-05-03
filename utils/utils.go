package utils

import (
	"math/rand"
)

func ChoosePronouns(
	DefaultPronouns map[string]bool, db Database, tense string) (
	povs []string, pronouns map[string][]string) {

	var userSelectedPronouns []string

	for userPronoun, keep := range DefaultPronouns {
		if keep {
			userSelectedPronouns = append(userSelectedPronouns, userPronoun)
		}
	}

	pronouns = db.GetPronouns(tense, userSelectedPronouns)

	return povs, pronouns
}

func ChooseVerb(
	verbs []map[string]string, pronouns map[string][]string) (
	verb map[string]string, pov, pronoun string) {

	getRandomPOV := func(pronouns map[string][]string) string {
		povs := make([]string, 0, len(pronouns))
		for k := range pronouns {
			povs = append(povs, k)
		}
		idxPov := rand.Int() % len(povs)
		return povs[idxPov]
	}

	idxVerb := rand.Int() % len(verbs)

	verb = verbs[idxVerb]
	pov = getRandomPOV(pronouns)

	idxPronoun := rand.Int() % len(pronouns[pov])
	pronoun = pronouns[pov][idxPronoun]

	return verb, pov, pronoun
}
