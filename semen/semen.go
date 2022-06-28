package semen

import (
    "bufio"
    "encoding/gob"
    "log"
    "math/rand"
    "os"
    "regexp"
    "strings"
)

type Dictionary struct {
    Data map[string]int
}

type Semen map[string]Dictionary

func NewFromDump(filename string) *Semen {
    s:= Semen{}
    s.LoadDump(filename)
    return &s
}

func NewFromText(filename string) *Semen {
    s:= Semen{}
    s.ReadFile(filename)
    return &s
}

func (dict Dictionary) Update(word string) {
    _, ok := dict.Data[word]
    if ok {
        dict.Data[word]++
    } else {
        dict.Data[word] = 1
    }
}

func (dict Dictionary) Count(str string) int {
    value, ok := dict.Data[str]
    if ok {
        return value
    }
    return 0
}

func (dict Dictionary) Random() string {
    if len(dict.Data) != 0 {
        i := rand.Intn(len(dict.Data))
        for key := range dict.Data {
            if i == 0 {
                return key
            }
            i--
        }
    }
    return ""
}

func (dict Dictionary) RandomWeighted() string {
    if len(dict.Data) != 0 {
        rnd := rand.Intn(len(dict.Data))
        index:= 0
        for key, value := range dict.Data {
            index += value
            if index > rnd {
                return key
            }
        }
    }
    return ""
}
func (s Semen) Random() string {
    if len(s) != 0 {
        i := rand.Intn(len(s))
        for key := range s {
            if i == 0 {
                return key
            }
            i--
        }
    }
    return ""
}

func (s Semen) Learning(words []string) {
    for i := 0; i < len(words) - 1; i++ {
        word := strings.ToLower(words[i])
        _, ok := s[word]
        if !ok {
            s[word] = Dictionary{Data:map[string]int{}}
        }
        if word != strings.ToLower(words[i + 1]) {
            s[word].Update(strings.ToLower(words[i + 1]))
        }
    }
}

func (s Semen) Talk(start string, length int) string {
    start = strings.ToLower(start)
    str := start
    if len(s) != 0 {
        _, ok := s[start]
        if start == "" || !ok {
            start = s.Random()
        }
        for i := 0; i < length; i++ {
            _, ok := s[start]
            if ok {
                start = s[start].RandomWeighted()
                str += " " + start
            }
        }
        min:=100
        last:= ""
        for key,_ := range s[start].Data {
            if len(s[key].Data) < min {
                last = key
                min = len(s[key].Data)
            }
        }
        str += " " + last
    }
    return str
}

func (s Semen) SaveDump(filename string) {
	file, err := os.Create(filename)
	if err != nil {
        log.Fatal(err)
		return
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	if err = e.Encode(s); err != nil {
        log.Fatal(err)
		return
	}
}

func (s Semen) LoadDump(filename string) {
	f, err := os.Open(filename)
	if err != nil {
        log.Fatal(err)
		return
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&s); err != nil {
        log.Fatal(err)
		return
	}
}

func (s Semen) ReadFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
        log.Fatal(err)
		return
	}
    defer f.Close()
    scanner := bufio.NewScanner(f)
    data := []string{""}
    reg, err := regexp.Compile(`\p{Cyrillic}+`)
    if err != nil {
        log.Fatal(err)
    }
    for scanner.Scan() {
        s:=scanner.Text()
        data = append(data, reg.FindAllString(s, -1)...)
    }
    s.Learning(data)
}
