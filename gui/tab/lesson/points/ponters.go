package points

import (
	"fmt"
)

type Point struct {
	X, Y int
}

func (p *Point) Move(dx int, dy int) {
	p.X = p.X + dx
	p.Y = p.Y + dy
}

func RunPointers() {
	pt1 := Point{X: 10, Y: 20}
	pt1.Move(1, 1)
	fmt.Println(pt1)
}
