package tgbot

import (
	"github.com/dyvdev/cybercum/utils"
	"strings"
)

func shakeSpear(str string) string {
	str = utils.CleanText(str)
	words := strings.Split(str, " ")
	if len(words) == 0 {
		return ""
	}
	str = words[len(words)-1]
	sylls := findAllSylls(str)
	if len(sylls)/2-1 < 0 {
		return ""
	}
	mid := changeSyl(sylls[len(sylls)/2-1])
	ret := []string{mid}
	ret = append(ret, sylls[len(sylls)/2:]...)

	return strings.Join(ret, "")
}

func findVow(V string) bool {
	vows := strings.Split("аеёиоуэюя", "")
	for _, v := range vows {
		if V == v {
			return true
		}
	}
	return false
}

func findAllSylls(str string) []string {
	var ret []string
	word := []rune(str)
	prevVowel := len(word)
	prevIter := 0
	for i, c := range word {
		s := string(c)
		if findVow(s) {
			prevVowel = i
			break
		}
	}
	for i, c := range word {
		s := string(c)
		if i > prevVowel && findVow(s) {
			a := prevVowel
			b := i
			if b-a == 1 {
				s := string(word[prevIter:b])
				ret = append(ret, s)
				prevIter = b
			} else if b-a == 2 {
				s := string(word[prevIter : b-1])
				ret = append(ret, s)
				prevIter = b - 1
			} else {
				s := string(word[prevIter : a+2])
				ret = append(ret, s)
				prevIter = a + 2
			}
			prevVowel = i
		}
	}
	s := string(word[prevIter:])
	ret = append(ret, s)

	return ret
}

func changeSyl(str string) string {
	rhymes := map[string]string{
		"а": "хуя",
		"е": "хуе",
		"ё": "хуё",
		"и": "хуи",
		"о": "хуё",
		"у": "хую",
		"э": "хуе",
		"ю": "хую",
		"я": "хуя",
	}
	word := []rune(str)
	if !findVow(string(word[len(word)-1])) {
		s := string(word[len(word)-1])
		return "хуе" + s
	} else {
		s := ""
		for i, c := range word {
			s = string(c)
			if findVow(s) {
				if i == len(word)-1 {
				}
				break
			}
		}
		change, _ := rhymes[s]
		return change
	}
}
