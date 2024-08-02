package particle

import (
	"errors"
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

type Quad struct {
	// Size
	Bbox rl.Rectangle

	// Data
	Particles []Particle
	Capacity  int

	// Subsections
	Topleft     *Quad
	Topright    *Quad
	Bottomleft  *Quad
	Bottomright *Quad
}

func NewQuad(x, y, w, h float32, cap int) Quad {
	return Quad{
		Bbox:        rl.NewRectangle(x, y, w, h),
		Particles:   make([]Particle, 0, cap),
		Capacity:    cap,
		Topleft:     nil,
		Topright:    nil,
		Bottomleft:  nil,
		Bottomright: nil,
	}
}

func (q *Quad) Quadrants() [4]*Quad {
	return [...]*Quad{q.Topleft, q.Topright, q.Bottomleft, q.Bottomright}
}

func (q *Quad) quadrant(p Particle) (*Quad, error) {
	// Check if point in range
	if p.Position.X < q.Bbox.X || p.Position.X > q.Bbox.X+q.Bbox.Width {
		return nil, fmt.Errorf(
			"out of bounds: %v",
			q,
		)
	}

	left := p.Position.X < (q.Bbox.X + q.Bbox.Width)
	top := p.Position.Y < (q.Bbox.Y + q.Bbox.Height)

	if left {
		if top {
			return q.Topleft, nil
		} else {
			return q.Bottomleft, nil
		}
	} else {
		if top {
			return q.Topright, nil
		} else {
			return q.Bottomright, nil
		}
	}
}

func (q *Quad) AddParticle(p Particle) error {
	// Determine quadrant + bounds check
	quadrant, err := q.quadrant(p)
	if err != nil {
		return err
	}

	// If there is space, add particle to particle list
	if len(q.Particles) < q.Capacity {
		q.Particles = append(q.Particles, p)
		return nil
	}

	// If exceeding capacity, subdivide
	if quadrant != nil {
		// Quadrants already exist, allow quadrant to handle placing particle

		err = quadrant.AddParticle(p)
		if quadrant.AddParticle(p) != nil {
			return err
		}
	} else {
		// No sub-quadrants yet, create new ones
		err := q.Subdivide()
		if err != nil {
			return fmt.Errorf("error subdividing: %v", err)
		}
	}

	return nil
}

func (q *Quad) AddParticles(p []Particle) error {
	for _, p := range p {
		if err := q.AddParticle(p); err != nil {
			return err
		}
	}

	return nil
}

func (q *Quad) Subdivide() error {
	width := q.Bbox.Width / 2
	height := q.Bbox.Height / 2

	top := q.Bbox.Y
	bottom := top + width
	left := q.Bbox.X
	right := left + height

	// Don't modify data until confirming that no errors will occur
	if q.Topleft != nil {
		return errors.New("topleft already exists")
	}
	if q.Topright != nil {
		return errors.New("topright already exists")
	}
	if q.Bottomleft != nil {
		return errors.New("bottomleft already exists")
	}
	if q.Bottomright != nil {
		return errors.New("bottomright already exists")
	}

	// Have to create quadrants like this for some reason
	tl := NewQuad(top, left, width, height, q.Capacity)
	q.Topleft = &tl

	tr := NewQuad(top, right, width, height, q.Capacity)
	q.Topright = &tr

	bl := NewQuad(bottom, left, width, height, q.Capacity)
	q.Bottomleft = &bl

	br := NewQuad(bottom, right, width, height, q.Capacity)
	q.Bottomright = &br

	// Place particles
	for _, p := range q.Particles {
		quad, err := q.quadrant(p)
		if err != nil {
			return fmt.Errorf("error reassigning particles: %v", err)
		}

		err = quad.AddParticle(p)
		if err != nil {
			return fmt.Errorf("error reassigning particles: %v", err)
		}
	}

	// Remove old particle list
	(*q).Particles = []Particle{}

	return nil
}
