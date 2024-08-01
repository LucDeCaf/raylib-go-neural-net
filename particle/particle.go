package particle

import rl "github.com/gen2brain/raylib-go/raylib"

type Particle struct {
	Position rl.Vector2
	Radius   float32
	Color    rl.Color
}

func OverlapSize(p1, p2 Particle) float32 {
	distance := rl.Vector2Distance(p1.Position, p2.Position)
	radiusSum := p1.Radius + p2.Radius

	return radiusSum - distance
}

func NewParticleV(v rl.Vector2, r float32, c rl.Color) Particle {
	return Particle{
		Position: v,
		Radius:   r,
		Color:    c,
	}
}

func NewParticle(x, y, r float32, c rl.Color) Particle {
	v := rl.Vector2{X: x, Y: y}
	return NewParticleV(v, r, c)
}

func SolveCollisions(particles *[]Particle) {
	for i, p1 := range *particles {
		for j, p2 := range *particles {
			if i == j {
				continue
			}

			if overlap := OverlapSize(p1, p2); overlap > 0 {
				direction := rl.Vector2Subtract(p1.Position, p2.Position)
				directionNorm := rl.Vector2Normalize(direction)

				moveAmount := overlap / 2

				// p1
				(*particles)[i].Position = rl.Vector2Add(p1.Position, rl.Vector2Scale(directionNorm, moveAmount))
				// p2
				(*particles)[j].Position = rl.Vector2Subtract(p2.Position, rl.Vector2Scale(directionNorm, moveAmount))
			}
		}
	}
}

func SolveCollisionsSubsteps(particles *[]Particle, substeps int) {
	for i := 0; i < substeps; i++ {
		SolveCollisions(particles)
	}
}
