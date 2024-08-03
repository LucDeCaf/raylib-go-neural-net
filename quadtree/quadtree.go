package quadtree

import (
	"cool-ai/particle"

	"github.com/gen2brain/raylib-go/raylib"
)

type Quad struct {
	BBox rl.Rectangle

	Particles []particle.Particle
	Capacity  int

	Topleft     *Quad
	Topright    *Quad
	Bottomleft  *Quad
	Bottomright *Quad
}

func NewQuad(bbox rl.Rectangle, cap int) Quad {
	return Quad{
		BBox: bbox,

		Particles: make([]particle.Particle, 0, cap),
		Capacity:  cap,

		Topleft:     nil,
		Topright:    nil,
		Bottomleft:  nil,
		Bottomright: nil,
	}
}

func (q *Quad) Quadrants() [4]*Quad {
	return [4]*Quad{
		q.Topleft,
		q.Topright,
		q.Bottomleft,
		q.Bottomright,
	}
}

func (q *Quad) AddParticles(p []particle.Particle) {
	for _, p := range p {
		q.AddParticle(p)
	}
}

func (q *Quad) AddParticle(p particle.Particle) {
	// If subdivided, add to correct subdivision
	if q.Capacity == 0 {
		q.QuadrantOf(p).AddParticle(p)
		return
	}

	// If not full, add particle
	if len(q.Particles) < q.Capacity {
		q.Particles = append(q.Particles, p)
		return
	}

	// Capacity is reached, subdivide
	q.SubdivideAndAdd(p)
}

func (q *Quad) QuadrantOf(p particle.Particle) *Quad {
	w := q.BBox.Width / 2
	h := q.BBox.Height / 2

	left := p.Position.X < q.BBox.X+w
	top := p.Position.Y < q.BBox.Y+h

	if top {
		if left {
			return q.Topleft
		} else {
			return q.Topright
		}
	}
	if left {
		return q.Bottomleft
	}
	return q.Bottomright
}

func (q *Quad) SubdivideWithoutMoving() {
	w := q.BBox.Width / 2
	h := q.BBox.Height / 2

	t := q.BBox.Y
	l := q.BBox.X
	b := t + h
	r := l + w

	tl := NewQuad(rl.NewRectangle(l, t, w, h), q.Capacity)
	tr := NewQuad(rl.NewRectangle(r, t, w, h), q.Capacity)
	bl := NewQuad(rl.NewRectangle(l, b, w, h), q.Capacity)
	br := NewQuad(rl.NewRectangle(r, b, w, h), q.Capacity)

	q.Topleft = &tl
	q.Topright = &tr
	q.Bottomleft = &bl
	q.Bottomright = &br
}

func (q *Quad) Subdivide() {
	q.SubdivideWithoutMoving()

	for _, p := range q.Particles {
		quadrant := q.QuadrantOf(p)
		quadrant.AddParticle(p)
	}

	q.Particles = nil
	q.Capacity = 0
}

func (q *Quad) SubdivideAndAdd(p particle.Particle) {
	q.Subdivide()
	q.QuadrantOf(p).AddParticle(p)
}
