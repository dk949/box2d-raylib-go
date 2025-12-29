package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	b2 "github.com/oliverbestmann/box2d-go"
)

// This is a line-for-line port of https://github.com/erincatto/box2d-raylib/blob/main/main.c

// This shows how to use Box2D v3 with raylib.
// It also show how to use Box2D with pixel units.
type Entity struct {
	body    b2.Body
	extent  b2.Vec2
	texture rl.Texture2D
}

func DrawEntity(entity *Entity) {
	// The boxes were created centered on the bodies, but raylib draws textures starting at the top left corner.
	// b2Body_GetWorldPoint gets the top left corner of the box accounting for rotation.
	p := entity.body.GetWorldPoint(b2.Vec2{X: -entity.extent.X, Y: -entity.extent.Y})
	rotation := entity.body.GetRotation()
	radians := rotation.Angle()

	ps := rl.Vector2{X: p.X, Y: p.Y}
	rl.DrawTextureEx(entity.texture, ps, (180/math.Pi)*radians, 1.0, rl.White)
}

const (
	GROUND_COUNT = 14
	BOX_COUNT    = 10
)

func main() {
	const (
		width  = 1920
		height = 1080
	)
	rl.InitWindow(width, height, "testing")

	rl.SetTargetFPS(60)
	// 128 pixels per meter is a appropriate for this scene. The boxes are 128 pixels wide.
	var lengthUnitsPerMeter float32 = 128.0
	b2.SetLengthUnitsPerMeter(lengthUnitsPerMeter)

	worldDef := b2.DefaultWorldDef()

	// Realistic gravity is achieved by multiplying gravity by the length unit.
	worldDef.Gravity.Y = 9.8 * lengthUnitsPerMeter
	world := b2.CreateWorld(worldDef)

	groundTexture := rl.LoadTexture("ground.png")
	boxTexture := rl.LoadTexture("box.png")

	groundExtent := b2.Vec2{X: 0.5 * float32(groundTexture.Width), Y: 0.5 * float32(groundTexture.Height)}
	boxExtent := b2.Vec2{X: 0.5 * float32(boxTexture.Width), Y: 0.5 * float32(boxTexture.Height)}

	// These polygons are centered on the origin and when they are added to a body they
	// will be centered on the body position.
	groundPolygon := b2.MakeBox(groundExtent.X, groundExtent.Y)
	boxPolygon := b2.MakeBox(boxExtent.X, boxExtent.Y)

	groundEntities := [GROUND_COUNT]Entity{}
	for i := range GROUND_COUNT {
		entity := &groundEntities[i]
		bodyDef := b2.DefaultBodyDef()
		bodyDef.Position = b2.Vec2{X: (2.0*float32(i) + 2.0) * groundExtent.X, Y: height - groundExtent.Y - 100.0}

		entity.body = world.CreateBody(bodyDef)
		entity.extent = groundExtent
		entity.texture = groundTexture
		shapeDef := b2.DefaultShapeDef()
		entity.body.CreatePolygonShape(shapeDef, groundPolygon)
	}

	boxEntities := [BOX_COUNT]Entity{}
	var boxIndex int
	for i := range 4 {
		y := height - groundExtent.Y - 100.0 - (2.5*float32(i)+2.0)*boxExtent.Y - 20.0

		for j := i; j < 4; j++ {
			x := 0.5*width + (3.0*float32(j)-float32(i)-3.0)*boxExtent.X
			if boxIndex >= BOX_COUNT {
				panic(fmt.Sprintf("boxIndex(%v) >= BOX_COUNT(%v)", boxIndex, BOX_COUNT))
			}

			entity := &boxEntities[boxIndex]
			bodyDef := b2.DefaultBodyDef()
			bodyDef.Type1 = b2.DynamicBody
			bodyDef.Position = b2.Vec2{X: x, Y: y}
			entity.body = world.CreateBody(bodyDef)
			entity.texture = boxTexture
			entity.extent = boxExtent
			shapeDef := b2.DefaultShapeDef()
			entity.body.CreatePolygonShape(shapeDef, boxPolygon)

			boxIndex += 1
		}
	}

	pause := false

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyP) {
			pause = !pause
		}

		if pause == false {
			deltaTime := rl.GetFrameTime()
			world.Step(deltaTime, 4)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		message := "Hello Box2D!"
		const fontSize = 36
		textWidth := rl.MeasureText("Hello Box2D!", fontSize)
		rl.DrawText(message, (width-textWidth)/2, 50, fontSize, rl.LightGray)

		for i := range GROUND_COUNT {
			DrawEntity(&groundEntities[i])
		}

		for i := range BOX_COUNT {
			DrawEntity(&boxEntities[i])
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(groundTexture)
	rl.UnloadTexture(boxTexture)

	rl.CloseWindow()

}
