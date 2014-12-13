package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	length  int
	regular bool
	single  bool
	source  string
	number  string
	letter  string
	whole   string
	lower   string
	capital string
	word    string
	visible string
)

func initialize() {
	whole = ""
	for i := 0; i < 128; i++ {
		whole += string(i)
	}
	visible = whole[0x20:0x7f]
	capital = whole[0x41:0x5b]
	lower = whole[0x61:0x7b]
	letter = capital + lower
	number = whole[0x30:0x3a]
	word = number + letter
	flag.IntVar(&length, "length", 8, "length of string")
	flag.BoolVar(&regular, "regular", false, "print regular string")
	flag.BoolVar(&single, "single", false, "remove duplicate char")
	flag.StringVar(&source, "source", "", "source of random string")
}

func duplicate(str string) (ret string) {
	ret = ""
	for _, chr := range str {
		if !strings.Contains(ret, string(chr)) {
			ret += string(chr)
		}
	}
	return
}

func main() {
	initialize()
	flag.Parse()
	var src string
	if regular {
		src = word
	} else {
		src = visible
	}
	if source != "" {
		src = source
	}
	if single {
		src = duplicate(src)
		if length > len(src) {
			length = len(src)
		}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	str := ""
	perm := r.Perm(len(src))
	for i := 0; i < length; i++ {
		str += string(src[perm[i]])
	}
	fmt.Fprintf(os.Stdout, "%s\n", str)
}
