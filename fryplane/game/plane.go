package game

import "strings"

type Direction int

const (
	_ Direction = iota
	TOP
	BOTTOM
	LEFT
	RIGHT
)

//    0  1  2  3  4  5  6  7  8  9
// 0  O  O  X  O  O  O  O  O  O  O
// 1  X  X  X  X  X  O  O  O  O  O
// 2  O  O  X  O  O  O  O  *  O  O
// 3  O  X  X  X  O  *  *  *  *  *
// 4  O  O  O  O  O  O  O  *  O  O
// 5  O  O  O  O  O  O  *  *  *  O
// 6  O  O  O  O  O  O  O  O  O  O
// 7  O  O  O  O  O  O  O  O  O  O
// 8  O  O  O  O  O  O  O  O  O  O
// 9  O  O  O  O  O  O  O  O  O  O

// 不同朝向的飞机所对应的机头的横纵坐标的范围
// const PlaneHeadRangeMap map[Direction][] = {
// 	TOP: {
// 		row: [0, 6],
// 		col: [2, 7]
// 	},
// 	LEFT: {
// 		row: [2, 7],
// 		col: [0, 6]
// 	},
// 	BOTTOM: {
// 		row: [3, 9],
// 		col: [2, 7]
// 	},
// 	RIGHT: {
// 		row: [2, 7],
// 		col: [3, 9]
// 	}
// }

type Point struct {
	row int
	col int
}

func (p Point) Number() int {
	return p.row*10 + p.col
}

func (p Point) PlaneList(direction Direction) []*Point {
	row, col := p.row, p.col
	switch direction {
	case TOP:
		return []*Point{
			{row, col},
			{row + 1, col - 2},
			{row + 1, col - 1},
			{row + 1, col},
			{row + 1, col + 1},
			{row + 1, col + 2},
			{row + 2, col},
			{row + 3, col - 1},
			{row + 3, col},
			{row + 3, col + 1},
		}
	case BOTTOM:
		return []*Point{
			{row, col},
			{row - 1, col - 2},
			{row - 1, col - 1},
			{row - 1, col},
			{row - 1, col + 1},
			{row - 1, col + 2},
			{row - 2, col},
			{row - 3, col - 1},
			{row - 3, col},
			{row - 3, col + 1},
		}
	case LEFT:
		return []*Point{
			{row, col},
			{row - 2, col + 1},
			{row - 1, col + 1},
			{row, col + 1},
			{row + 1, col + 1},
			{row + 2, col + 1},
			{row, col + 2},
			{row - 1, col + 3},
			{row, col + 3},
			{row + 1, col + 3},
		}
	case RIGHT:
		return []*Point{
			{row, col},
			{row - 2, col - 1},
			{row - 1, col - 1},
			{row, col - 1},
			{row + 1, col - 1},
			{row + 2, col - 1},
			{row, col - 2},
			{row - 1, col - 3},
			{row, col - 3},
			{row + 1, col - 3},
		}
	}
	return nil
}

func (p *Point) ValidHead(direction Direction) bool {
	switch direction {
	case TOP:
		return 0 <= p.row && p.row <= 6 && 2 <= p.col && p.col <= 7
	case BOTTOM:
		return 3 <= p.row && p.row <= 9 && 2 <= p.col && p.col <= 7
	case LEFT:
		return 2 <= p.row && p.row <= 7 && 0 <= p.col && p.col <= 6
	case RIGHT:
		return 0 <= p.row && p.row <= 6 && 3 <= p.col && p.col <= 9
	}
	return false
}

type Plane struct {
	id        int
	direction Direction
	head      Point
}

func (plan *Plane) CanGeneratePlan() bool {
	pointList := plan.head.PlaneList(plan.direction)
	for _, point := range pointList {
		if point.row < 0 || point.col < 0 {
			return false
		}
	}
	return true
}

func (plan *Plane) NumberList() []int {
	numberList := make([]int, 0, 10)
	points := plan.head.PlaneList(plan.direction)
	for _, p := range points {
		numberList = append(numberList, p.Number())
	}
	return numberList
}

func (plan *Plane) String() string {
	plist := make([][]string, 0, 10)
	for i := 0; i < 10; i++ {
		tem := make([]string, 0, 10)
		for j := 0; j < 10; j++ {
			tem = append(tem, "*")
		}
		plist = append(plist)
	}
	points := plan.head.PlaneList(plan.direction)
	for _, p := range points {
		plist[p.row][p.col] = "O"
	}
	plist[plan.head.row][plan.head.col] = "X"
	ret := ""
	for _, rows := range plist {
		ret += strings.Join(rows, "  ")
	}
	return ret
}
