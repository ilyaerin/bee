package main

import (
	"time"
	"math/rand"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go" // TODO maybe change fatih/color
	"strings"
)

const COUNT = 250
const BASES_COUNT = 5
const WIDTH = 150
const HEIGHT = 35
const TURN_TIME = 50
const HIT = 80
const HIT_BACK = 30
const BASE_HEALTH = 500
const BEE_HEALTH = 200
const BASE_MOVING_PERCENT = 10
const TIME_TO_REVIVAL = 10

type Point struct {
	x int
	y int
}

type Bee struct{
	Point
	base *Base
	live bool
	health int
}

type Base struct {
	Point
	color func(format string, a ...interface{}) string
	count int
	live bool
	health int
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
var bees [COUNT*BASES_COUNT]*Bee

func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < BASES_COUNT; i++ {
		bases[i] = MakeBase(i)
	}

	for i := 0; i < BASES_COUNT * COUNT; i++ {
		bees[i] = Born(bases[i % BASES_COUNT])
		go bees[i].Living()
	}

	for _, base := range bases {
		go base.BaseLiving()
	}

	go Monitoring()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		loop(scanner.Text())
	}
}

func loop(input string) {
	switch input {
	case "q":
		os.Exit(0)
	default:
		id, _ := strconv.Atoi(input)
		bases[id - 1].BaseKill()
	}
}

func Monitoring()  {
	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		out := fmt.Sprintf(Stat())
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
						s = base.color("@")
					}
				}
				out += s
			}
			out += "\n"
		}
		clearConsole()
		fmt.Print(out)

		if checkGameOver() {
			fmt.Println("Game over")
			os.Exit(0)
		}
	}
}

func Stat() string {
	beesCount := 0
	for _, bee := range bees {
		if bee.live { beesCount += 1 }
	}

	outBases := []string{}
	for _, base := range bases {
		outBases = append(outBases, base.color("[Bees: %d Health: %d]", base.count, base.health))
	}
	return fmt.Sprintf("All bees: %d %s\n", beesCount, strings.Join(outBases, " "))
}

func checkGameOver() bool {
	liveBases := 0
	for _, base := range bases {
		if base.live { liveBases += 1 }
	}

	return liveBases == 1
}


func (bee *Bee) Move(x int, y int) {
	newX := bee.x + x
	newY := bee.y + y
	if newX >= 0 && newX <= WIDTH { bee.x = newX }
	if newY >= 0 && newY <= HEIGHT { bee.y = newY }

	for _, base := range bases {
		if base.x == newX && base.y == newY && bee.base != base && base.live {
			bee.health -= HIT_BACK
			base.health -= HIT
		}
	}

	for _, beeLoop := range bees {
		if beeLoop.x == newX && beeLoop.y == newY && bee.base != beeLoop.base && beeLoop.live {
			bee.health -= HIT_BACK
			beeLoop.health -= HIT
		}
	}
}

func (bee *Bee) Living() {
	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		if bee.health < 0 {
			bee.Kill()
			//return
		}

		if bee.live {
			bee.Move(rand.Intn(3)-1, rand.Intn(3)-1)
		} else {
			bee.Revival(TIME_TO_REVIVAL + rand.Intn(TIME_TO_REVIVAL))
		}
	}
}

func (base *Base) BaseLiving() {
	ticker := time.NewTicker(time.Millisecond * TURN_TIME)
	for range ticker.C {
		if base.live && rand.Intn(100) <= BASE_MOVING_PERCENT {
			base.BaseMove(rand.Intn(3)-1, rand.Intn(3)-1)
		}

		if base.health < 0 || base.count == 0 {
			base.BaseKill()
			return
		}
	}
}

func (base *Base) BaseMove(x int, y int) {
	newX := base.x + x
	newY := base.y + y
	if newX >= 0 && newX <= WIDTH { base.x = newX }
	if newY >= 0 && newY <= HEIGHT { base.y = newY }

	for _, baseLoop := range bases {
		if baseLoop.x == newX && baseLoop.y == newY && baseLoop != base && baseLoop.live {
			base.health -= HIT_BACK * 2
			baseLoop.health -= HIT * 2
		}
	}

	for _, bee := range bees {
		if bee.x == newX && bee.y == newY && bee.base != base && bee.live {
			bee.health -= HIT_BACK * 2
			//base.health -= HIT
		}
	}
}


func (base *Base) BaseKill() {
	for _, bee := range bees {
		if bee.live && bee.base == base {
			bee.Kill()
		}
	}

	base.live = false
	base.color = color.New(color.Faint).SprintfFunc()
}

func (bee *Bee) Kill() {
	bee.live = false
	bee.health = BEE_HEALTH
	bee.x = bee.base.x
	bee.y = bee.base.y
	bee.base.count -= 1
}

func (bee *Bee) Revival(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
	if bee.base.live {
		bee.base.count += 1
		bee.live = true
	}
}

func Born(base *Base) (bee *Bee) {
	base.count += 1
	return &Bee{Point: base.Point, base: base, live: true}
}

//func randomColor() func(format string, a ...interface{}) string {
//	randomColor := colors[rand.Intn(len(colors))]
//	return color.New(randomColor).SprintfFunc()
//}

func MakeBase(i int) (* Base) {
	return &Base{
		Point: Point{rand.Intn(WIDTH), rand.Intn(HEIGHT)},
		color: color.New(colors[i]).SprintfFunc(),
		//color: randomColor(),
		count: 0,
		live: true,
		health: BASE_HEALTH,
	}
}

func clearConsole() {
	fmt.Print("n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()
}
