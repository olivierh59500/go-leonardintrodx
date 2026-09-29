package demo

import (
	"math"
	"testing"

	"github.com/olivierh59500/democonstructionkit/geometry"
)

func TestCubeTransformUsesOriginalCamera(t *testing.T) {
	pose := cubeTransform(1)
	camera := geometry.Camera{Center: geometry.Vec2{X: Width / 2, Y: Height / 2}, Focal: Height * math.Sqrt(3) / 2, Near: .1}
	center, _, visible := camera.Project(pose.Position)
	if !visible || math.Abs(center.X-421.00261483906723) > 1e-9 || math.Abs(center.Y-165.51707880847889) > 1e-9 {
		t.Fatalf("projected center at 1s: %+v, visible=%t", center, visible)
	}
	if math.Abs(pose.Rotation.X+math.Pi/2) > 1e-12 || math.Abs(pose.Rotation.Y-1.1*math.Pi/2) > 1e-12 {
		t.Fatalf("native cube rotation at 1s: %+v", pose.Rotation)
	}
}
