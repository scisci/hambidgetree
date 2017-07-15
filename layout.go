package hambidgetree

type Rectangle struct {
}

type Layout interface {
	Rectangles() []Rectangle
}
