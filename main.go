package main

import (
	"time"
	"math/rand"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"github.com/fatih/color"
)

const COUNT = 50
const BASES_COUNT = 4
const WIDTH = 180
const HEIGHT = 46
const TURN_TIME = 50

type Point struct {
	x int
	y int
}

type Bee struct{
	Point
	base *Base
	live bool
}

type Base struct {
	Point
	color func(format string, a ...interface{}) string
}

var colors = []color.Attribute{
	color.FgHiRed,
	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
}

var bases [BASES_COUNT]*Base

func main() {
	rand.Seed(time.Now().Unix())

	var bees [COUNT*BASES_COUNT]*Bee

	for i := 0; i < BASES_COUNT; i++ {
		bases[i] = MakeBase()
		for j := 0; j < COUNT; j++ {
			bees[i * COUNT + j] = Born(bases[i])
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
	switch input {
	case "1", "2", "3", "4":
		// Killing all bees from base
		for _, bee := range bees {
			id, _ := strconv.Atoi(input)
			bases[id].color = randomColor() // Change color of base
			if bases[id] == bee.base { bee.Kill() }
		}
	case "q":
		os.Exit(0)
	}
}

func Monitoring(bees [COUNT*BASES_COUNT]*Bee)  {
	t := 0
	b := color.New(color.FgWhite).SprintfFunc()
	w := color.New(color.FgWhite).SprintfFunc()

	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		t += 1
		out := w("Turn: %d\n", t)
		for j := 0; j < HEIGHT; j++ {
			for i := 0; i < WIDTH; i++ {
				s := " "
				for _, bee := range bees {
					if bee.x == i && bee.y == j && bee.live {
						s = bee.base.color("*")
					}
				}
				for _, base := range bases {
					if base.x == i && base.y == j {
						s = b("@")
					}
				}
				out += s
			}
			out += "\n"
		}
		clear_console()
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

func Born(base *Base) (bee *Bee) {
	return &Bee{Point: base.Point, base: base, live: true}
}

func randomColor() func(format string, a ...interface{}) string {
	randomColor := colors[rand.Intn(len(colors))]
	return color.New(randomColor).SprintfFunc()
}

func MakeBase() (* Base) {

	return &Base{
		Point: Point{rand.Intn(WIDTH), rand.Intn(HEIGHT)},
		color: randomColor(),
	}
}

func clear_console() {
	fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
}
