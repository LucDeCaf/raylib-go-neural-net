package main

import (
	// Standard lib
	// "log"
	"math/rand"

	// Submodules
	"cool-ai/particle"
	"cool-ai/quadtree"

	// Third-party
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Raylib setup
	rl.InitWindow(800, 600, "Super cool AI Simulation")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	particles := []particle.Particle{}

	// Main loop
	for !rl.WindowShouldClose() {
		// Add particle on click
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			c := rl.NewColor(uint8(rand.Intn(256)), 0, 0, 255)
			p := particle.NewParticleV(rl.GetMousePosition(), 5, c)

			particles = append(particles, p)
		}

		// Create new quadtree
		qt := quadtree.NewQuad(rl.NewRectangle(0, 0, 800, 600), 1)
		qt.AddParticles(particles)

		// Add new particles
		// const particleCount = 5
		// for i := 0; i < particleCount; i++ {
		// 	x := rand.Float32() * float32(rl.GetScreenWidth())
		// 	y := rand.Float32() * float32(rl.GetScreenHeight())
		// 	r := 2 + (rand.Float32()*3)
		// 	c := rl.NewColor(uint8(rand.Intn(256)), 0, 0, 255)

		// 	particles = append(particles, particle.NewParticle(x, y, r, c))
		// }

		// Solve particle collisions
		// particle.SolveCollisionsSubsteps(&particles, 8)

		// Drawing logic
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		// Visualise quadtree
		debugQuadtree(&qt)

		// Draw particles
		drawParticles(particles)

		rl.EndDrawing()
	}
}

func debugQuadtree(q *quadtree.Quad) {
	// Draw rectangle
	rl.DrawRectangleLinesEx(q.BBox, 3, rl.Black)

	for _, quad := range q.Quadrants() {
		if quad != nil {
			debugQuadtree(quad)
		}
	}
}

func drawParticles(p []particle.Particle) {
	for _, p := range p {
		rl.DrawCircleV(p.Position, p.Radius, p.Color)
	}
}
