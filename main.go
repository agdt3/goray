package main

import (
    "os"
    "fmt"
    "image"
    "math"
    "image/color"
    "image/jpeg"
    "image/draw"
)

// May not be a reason to distinguish these
type Point3 struct {
    X float64
    Y float64
    Z float64
}

type Vector3 struct {
    X float64
    Y float64
    Z float64
    Magnitude float64
    magnitudeInverse float64 // 1/Mag precomputed
}

func MakeVector3(x, y, z float64) *Vector3 {
    vec3 := new(Vector3)
    vec3.X = x
    vec3.Y = y
    vec3.Z = z
    vec3.Magnitude, vec3.magnitudeInverse = vec3.CalculateMagnitude()

    return vec3
}

func (v *Vector3) Normalize() {
    v.X = v.X * v.magnitudeInverse
    v.Y = v.Y * v.magnitudeInverse
    v.Z = v.Z * v.magnitudeInverse
    v.Magnitude = 1
    v.magnitudeInverse = 1
}

func (v Vector3) CalculateMagnitude() (float64, float64) {
    // Recalculates magnitude
    mag := math.Sqrt(
        (v.X * v.X) +
        (v.Y * v.Y) +
        (v.Z * v.Z))

    return mag, 1/mag
}

type CameraPerspective struct {
    Origin Point3
    Dir Vector3
}


type CameraTHEOTHERONE struct {
    Origin Point3
    Dir Vector3
}

type RayTraceConfig struct {
    UseLight bool
    UseShadows bool
    MaxReflections uint
    ImageWidth int
    ImageHeight int
}

type World struct {
    Cam CameraPerspective
    Img draw.Image // use the draw interface
//    Obj []Objects
    Config RayTraceConfig
}

func NewWorld () *World {
    // World with sane defaults
    world := new(World)
    vec3 := MakeVector3(0, 0, -1)
    world.Cam = CameraPerspective{Point3{0, 0, 0}, *vec3}
    world.Config = RayTraceConfig{true, true, 1, 640, 480}
    world.Img = image.NewRGBA(image.Rect(0, 0, world.Config.ImageWidth, world.Config.ImageHeight))
//    world.Obj = nil // TODO: Make objects

    return world
}

func LoadOrCreateFile() {

}

func (w World) traceRay(origin Point3, dir Vector3) color.RGBA {
    return color.RGBA{0, 0, 255, 0}
}

func (w World) generateRay(x, y, z float64) *Vector3 {
    // generates a vector3 that is assumed
    // have an origin at (0,0,0)
    vec3 := MakeVector3(x, y, z)
    return vec3
}

func (w *World) Trace() {
    b := w.Img.Bounds()
    fmt.Println(b)
    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            vec3 := MakeVector3(0, 0, -1)
            pixelColor := w.traceRay(Point3{0, 0, 0}, *vec3)
            w.Img.Set(x, y, pixelColor)
//            fmt.Println(w.Img.At(x,y))
        }
    }

    // Get rid of this later
    f, err := os.Create("./test.jpg")
    if err != nil {
        fmt.Println(err)
        return
    }
    jpeg.Encode(f, w.Img, &jpeg.Options{100})
}


func main() {
    world := NewWorld()
    world.Trace()
}
