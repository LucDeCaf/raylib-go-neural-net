package quadtree_test

import (
	"testing"

	"cool-ai/particle"
	"cool-ai/quadtree"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestSubdivideWithoutMoving(t *testing.T) {
	q := quadtree.NewQuad(rl.NewRectangle(0, 0, 800, 600), 1)

	q.SubdivideWithoutMoving()

	expected := [...]rl.Rectangle{
		rl.NewRectangle(0, 0, 400, 300),     // tl
		rl.NewRectangle(400, 0, 400, 300),   // tr
		rl.NewRectangle(0, 300, 400, 300),   // bl
		rl.NewRectangle(400, 300, 400, 300), // br
	}
	result := [...]rl.Rectangle{
		q.Topleft.BBox,
		q.Topright.BBox,
		q.Bottomleft.BBox,
		q.Bottomright.BBox,
	}

	for i := 0; i < 4; i++ {
		if expected[i] != result[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], result[i])
		}
	}
}

func TestQuadrantOf(t *testing.T) {
	q := quadtree.NewQuad(rl.NewRectangle(0, 0, 800, 600), 99)
	q.SubdivideWithoutMoving()

	tl := particle.NewParticle(100, 100, 5, rl.Black)
	tr := particle.NewParticle(700, 100, 5, rl.Black)
	bl := particle.NewParticle(100, 500, 5, rl.Black)
	br := particle.NewParticle(700, 500, 5, rl.Black)

	if v := q.QuadrantOf(tl); v != q.Topleft {
		t.Errorf("expected '%v' but got '%v'", &q.Topleft, &v)
	}
	if v := q.QuadrantOf(tr); v != q.Topright {
		t.Errorf("expected '%v' but got '%v'", &q.Topright, &v)
	}
	if v := q.QuadrantOf(bl); v != q.Bottomleft {
		t.Errorf("expected '%v' but got '%v'", &q.Bottomleft, &v)
	}
	if v := q.QuadrantOf(br); v != q.Bottomright {
		t.Errorf("expected '%v' but got '%v'", &q.Bottomright, &v)
	}
}

func TestSubdivide(t *testing.T) {
	q := quadtree.NewQuad(rl.NewRectangle(0, 0, 800, 600), 3)

	tl := particle.NewParticle(100, 100, 5, rl.Black)
	tr := particle.NewParticle(700, 100, 5, rl.Black)
	bl := particle.NewParticle(100, 500, 5, rl.Black)
	br := particle.NewParticle(700, 500, 5, rl.Black)

	q.Particles = []particle.Particle{tl, tr, bl, br}

	q.Subdivide()

	if q.Topleft.Particles[0] != tl {
		t.Errorf("expected '%v' but got '%v'", tl, q.Topleft.Particles[0])
	}
	if q.Topright.Particles[0] != tr {
		t.Errorf("expected '%v' but got '%v'", tr, q.Topright.Particles[0])
	}
	if q.Bottomleft.Particles[0] != bl {
		t.Errorf("expected '%v' but got '%v'", bl, q.Bottomleft.Particles[0])
	}
	if q.Bottomright.Particles[0] != br {
		t.Errorf("expected '%v' but got '%v'", br, q.Bottomright.Particles[0])
	}
}

func TestParticleOverflow(t *testing.T) {
	q := quadtree.NewQuad(rl.NewRectangle(0, 0, 800, 600), 3)

	tl := particle.NewParticle(100, 100, 5, rl.Black)
	tr := particle.NewParticle(700, 100, 5, rl.Black)
	bl := particle.NewParticle(100, 500, 5, rl.Black)
	br := particle.NewParticle(700, 500, 5, rl.Black)

	// Fill capacity
	q.AddParticle(tl)
	q.AddParticle(tr)
	q.AddParticle(bl)

	expected := []particle.Particle{tl, tr, bl}
	result := q.Particles
	for i := 0; i < 3; i++ {
		if expected[i] != result[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], result[i])
		}
	}

	// Cause subdivision
	q.AddParticle(br)

	if q.Capacity != 0 {
		t.Errorf("expected '%v' but got '%v'", 0, q.Capacity)
	}

	expected = []particle.Particle{tl, tr, bl, br}
	result = []particle.Particle{
		q.Topleft.Particles[0],
		q.Topright.Particles[0],
		q.Bottomleft.Particles[0],
		q.Bottomright.Particles[0],
	}

	for i := 0; i < 4; i++ {
		if expected[i] != result[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], result[i])
		}
	}
}
