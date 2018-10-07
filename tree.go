package hambidgetree

// Serialize this like so
//
// ratios: (tree.ratios.ratios) // should marshal themselves with some kind of type flag, if strings, then expressions, otherwise floats
// xyRatioIndex:
// zyRatioIndex:
// root: id
// node: [id: id, split: hvd, index, left: id, right: id]


type TreeRegions interface {
	Offset() *Vector
	Scale() float64
	Region(id NodeID) Region
}
