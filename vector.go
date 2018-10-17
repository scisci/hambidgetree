package hambidgetree

type Vector struct {
	x float64
	y float64
	z float64
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{x, y, z}
}

func (v *Vector) Add(other *Vector) *Vector {
	return NewVector(v.x+other.x, v.y+other.y, v.z+other.z)
}
