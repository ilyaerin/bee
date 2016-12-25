package main

import (
	"time"
	"math/rand"
	"fmt"
	"os/exec"
	"bufio"
	"os"
	"strings"
	"strconv"
	"github.com/fatih/color"
)

const COUNT = 75
const BASES_COUNT = 3
const WIDTH = 120
const HEIGHT = 30
const TURN_TIME = 250

type Point struct {
	x int
	y int
}

type Bee struct{
	Point
	base *Point
	live bool
}

var bases [BASES_COUNT]Point

func main() {
	rand.Seed(time.Now().Unix())

	var bees [COUNT*BASES_COUNT]*Bee

	for i := 0; i < BASES_COUNT; i++ {
		bases[i] = Point{rand.Intn(WIDTH), rand.Intn(HEIGHT)}
		for j := 0; j < COUNT; j++ {
			bees[i * COUNT + j] = Born(&bases[i])
			go bees[i * COUNT + j].Living()
		}

	}

	go Monitoring(bees)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		loop(bees, scanner.Text())
	}
}

func loop(bees [COUNT*BASES_COUNT]*Bee, input string) {
	switch {
	case strings.HasPrefix(input, "k "):
		data := strings.Split(input, "k ")
		// Killing all bees from base
		for _, bee := range bees {
			id, _ := strconv.Atoi(data[1])
			if bases[id] == *bee.base { bee.Kill() }
		}
	case input == "q":
		os.Exit(0)
	}
}

func Monitoring(bees [COUNT*BASES_COUNT]*Bee)  {
	t := 0
	y := color.New(color.FgHiYellow).SprintfFunc()
	b := color.New(color.FgBlue).SprintfFunc()
	w := color.New(color.FgWhite).SprintfFunc()

	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		t += 1
		//out := fmt.Sprintf("Turn: %d\n", t)
		out := w("Turn: %d\n", t)
		for j := 0; j < HEIGHT; j++ {
			for i := 0; i < WIDTH; i++ {
				s := " "
				for _, bee := range bees {
					if bee.x == i && bee.y == j && bee.live {
						s = y("*")
					}
				}
				for _, base := range bases {
					if base.x == i && base.y == j {
						s = b("@")
					}
				}
				//fmt.Print(s)
				out += s
			}
			//fmt.Print("\n")
			out += "\n"
		}
		exec.Command("sh", "-c", "clear")
		fmt.Println(out)
	}
}

func (bee *Bee) Move(x int, y int) {
	newX := bee.x + x
	newY := bee.y + y
	if newX >= 0 && newX <= WIDTH { bee.x = newX }
	if newY >= 0 && newY <= HEIGHT { bee.y = newY }
}

func (bee *Bee) Living() {
	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		if bee.live {
			bee.Move(rand.Intn(3)-1, rand.Intn(3)-1)
		} else {
			bee.Revival(1 + rand.Intn(5))
		}
	}
}

func (bee *Bee) Kill() {
	bee.live = false
}

func (bee *Bee) Revival(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
	bee.x = bee.base.x
	bee.y = bee.base.y
	bee.live = true
}

func Born(base *Point) (bee *Bee) {
	return &Bee{Point: *base, base: base, live: true}
}