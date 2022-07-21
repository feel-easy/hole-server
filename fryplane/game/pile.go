package game

type Pile struct {
	planes []int
}

func NewPile() *Pile {
	return &Pile{planes: make([]int, 0, 144)}
}

func (p *Pile) Add(plane int) {
	p.planes = append(p.planes, plane)
}

func (p *Pile) Tiles() []int {
	planes := make([]int, len(p.planes))
	copy(planes, p.planes)
	return planes
}

func (p *Pile) ReplaceTop(plane int) {
	p.planes[len(p.planes)-1] = plane
}

func (p *Pile) Top() int {
	pileSize := len(p.planes)
	if pileSize == 0 {
		return 0
	}
	return p.planes[pileSize-1]
}

func (d *Pile) BottomDrawOne() int {
	plane := d.planes[len(d.planes)-1]
	d.planes = d.planes[0 : len(d.planes)-1]
	return plane
}
