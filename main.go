package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var TYPES = []string{
	"Normal",
	"Fighting",
	"Flying",
	"Poison",
	"Ground",
	"Rock",
	"Bug",
	"Ghost",
	"Steel",
	"Fire",
	"Water",
	"Grass",
	"Electric",
	"Psychic",
	"Ice",
	"Dragon",
	"Dark",
	"Fairy",
}

func parseArgs() (dex, chart string) {
	flag.StringVar(&dex, "d", "dex.txt", "pokedex file")
	flag.StringVar(&chart, "c", "chart.txt", "type chart file")
	flag.Parse()
	return
}

type Pokemon struct {
	Number int
	Name   string
	Types  []string
}

type Effectiveness float32

func ParseEffectiveness(effs string) Effectiveness {
	switch effs {
	case "½×":
		return Weak
	case "1×":
		return Normal
	case "2×":
		return Strong
	case "0×":
		return Immune
	default:
		panic(effs)
	}
}

func (e Effectiveness) String() string {
	switch e {
	case Immune:
		return "0"
	case Weak:
		return "½"
	case Normal:
		return "1"
	case Strong:
		return "2"
	default:
		panic(e)
	}
}

const (
	Immune Effectiveness = 0
	Weak   Effectiveness = 0.5
	Normal Effectiveness = 1
	Strong Effectiveness = 2
)

type Pokedex []Pokemon

type Chart [][]Effectiveness

func (c Chart) String() string {
	s := "     |"
	for i := range TYPES {
		if i > 0 {
			s += " |"
		}
		s += " " + TYPES[i]
	}
	s += "\n"
	for i, t := range c {
		if i%3 == 0 {
			s += strings.Repeat("-", 4+7*len(TYPES)) + "\n"
		}
		s += TYPES[i] + " |"
		for j, e := range t {
			if j > 0 {
				s += "  |"
			}
			s += "   " + e.String()
		}
		s += "\n"
	}
	return s
}

func parsePokedex(dexfile io.ReadCloser) *Pokedex {
	defer dexfile.Close()
	f := bufio.NewReader(dexfile)
	mons := make(Pokedex, 0, 800)
	for {
		bytes, _, err := f.ReadLine()
		if err != nil {
			if err == io.EOF {
				return &mons
			} else {
				panic(err)
			}
		}
		line := string(bytes)
		mon := &Pokemon{}
		mon.Number, _ = strconv.Atoi(strings.TrimLeft(line, "#0"))
		name, _, _ := f.ReadLine()
		mon.Name = string(name)
		types, _, _ := f.ReadLine()
		mon.Types = strings.Split(string(types), " · ")
		mons = append(mons, *mon)
	}
}

func parseChart(chartfile io.ReadCloser) *Chart {
	defer chartfile.Close()
	f := bufio.NewReader(chartfile)
	chart := make(Chart, len(TYPES))
	for i := 0; i < len(TYPES); i++ {
		chart[i] = make([]Effectiveness, len(TYPES))

		bytes, _, err := f.ReadLine()
		if err != nil {
			panic(err)
		}
		types := strings.Split(string(bytes), "  ")
		for j, t := range types {
			chart[i][j] = ParseEffectiveness(t)
		}
	}
	return &chart
}

type Soln []int

func (s Soln) Includes(t int) bool {
	for _, i := range s {
		if i == t {
			return true
		}
	}
	return false
}

func (s Soln) Append(t int) Soln {
	return append(s, t)
}

func (s Soln) Max() int {
	m := -1
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}

func (s Soln) String() string {
	r := ""
	for i, t := range s {
		if i > 0 {
			r += " "
		}
		r += TYPES[t]
	}
	return r
}

func (c *Chart) AllCoveredOffense(soln Soln) bool {
	covered := make(map[int]bool)
	for i := range TYPES {
		covered[i] = false
	}
	for _, t := range soln {
		for j, e := range (*c)[t] {
			if e == Strong {
				covered[j] = true
			}
		}
	}
	for _, b := range covered {
		if !b {
			return false
		}
	}
	return true
}

func (c *Chart) AllCoveredDefense(soln Soln) bool {
	covered := make(map[int]bool)
	for i := range TYPES {
		covered[i] = false
	}
	for _, t := range soln {
		for i := range TYPES {
			e := (*c)[i][t]
			if e == Immune || e == Weak {
				covered[t] = true
			}
		}
	}
	for _, b := range covered {
		if !b {
			return false
		}
	}
	return true
}

func (c *Chart) FindTypeCoverage(inuse Soln, offensive bool) Solutions {
	if inuse == nil {
		if offensive {
			inuse = Soln{1} // Skip Normal since it sucks ofensively
		} else {
			inuse = Soln{}
		}
	}
	covered := false
	if offensive {
		covered = c.AllCoveredOffense(inuse)
	} else {
		covered = c.AllCoveredDefense(inuse)
	}
	if covered {
		// copy by value
		v := make(Soln, len(inuse))
		for i, t := range inuse {
			v[i] = t
		}
		return Solutions{v}
	}
	solutions := make(Solutions, 0, len(TYPES))
	for i := inuse.Max() + 1; i < len(TYPES); i++ {
		if !inuse.Includes(i) {
			solns := c.FindTypeCoverage(inuse.Append(i), offensive)
			if len(solns) > 0 {
				solutions = append(solutions, solns...)
			}
		}
	}
	return solutions
}

type Solutions []Soln

func (s Solutions) Len() int {
	return len(s)
}

func (s Solutions) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func (s Solutions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_, chartfn := parseArgs()

	chartf, _ := os.Open(chartfn)
	chart := parseChart(chartf)

	solutions := chart.FindTypeCoverage(nil, true)
	sort.Sort(solutions)

	for i := 0; i < len(solutions); i++ {
		fmt.Fprintf(os.Stdout, "%v\n", solutions[i])
	}
}
