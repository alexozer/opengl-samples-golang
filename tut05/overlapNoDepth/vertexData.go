package overlapNoDepth

const (
	rightExtent  = 0.8
	leftExtent   = -rightExtent
	topExtent    = 0.20
	middleExtent = 0.0
	bottomExtent = -topExtent
	frontExtent  = -1.25
	rearExtent   = -1.75
)

var vertexData = []float32{
	// Object 1 positions
	leftExtent, topExtent, rearExtent,
	leftExtent, middleExtent, frontExtent,
	rightExtent, middleExtent, frontExtent,
	rightExtent, topExtent, rearExtent,

	leftExtent, bottomExtent, rearExtent,
	leftExtent, middleExtent, frontExtent,
	rightExtent, middleExtent, frontExtent,
	rightExtent, bottomExtent, rearExtent,

	leftExtent, topExtent, rearExtent,
	leftExtent, middleExtent, frontExtent,
	leftExtent, bottomExtent, rearExtent,

	rightExtent, topExtent, rearExtent,
	rightExtent, middleExtent, frontExtent,
	rightExtent, bottomExtent, rearExtent,

	leftExtent, bottomExtent, rearExtent,
	leftExtent, topExtent, rearExtent,
	rightExtent, topExtent, rearExtent,
	rightExtent, bottomExtent, rearExtent,

	// Object 2 positions
	topExtent, rightExtent, rearExtent,
	middleExtent, rightExtent, frontExtent,
	middleExtent, leftExtent, frontExtent,
	topExtent, leftExtent, rearExtent,

	bottomExtent, rightExtent, rearExtent,
	middleExtent, rightExtent, frontExtent,
	middleExtent, leftExtent, frontExtent,
	bottomExtent, leftExtent, rearExtent,

	topExtent, rightExtent, rearExtent,
	middleExtent, rightExtent, frontExtent,
	bottomExtent, rightExtent, rearExtent,

	topExtent, leftExtent, rearExtent,
	middleExtent, leftExtent, frontExtent,
	bottomExtent, leftExtent, rearExtent,

	bottomExtent, rightExtent, rearExtent,
	topExtent, rightExtent, rearExtent,
	topExtent, leftExtent, rearExtent,
	bottomExtent, leftExtent, rearExtent,

	// Object 1 colors
	0.75, 0.75, 1, 1,
	0.75, 0.75, 1, 1,
	0.75, 0.75, 1, 1,
	0.75, 0.75, 1, 1,

	0, 0.5, 0, 1,
	0, 0.5, 0, 1,
	0, 0.5, 0, 1,
	0, 0.5, 0, 1,

	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,

	0.8, 0.8, 0.8, 1,
	0.8, 0.8, 0.8, 1,
	0.8, 0.8, 0.8, 1,

	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,

	// Object 2 colors
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,

	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,
	0.5, 0.5, 0, 1,

	0, 0.5, 0, 1,
	0, 0.5, 0, 1,
	0, 0.5, 0, 1,

	0.75, 0.75, 1, 1,
	0.75, 0.75, 1, 1,
	0.75, 0.75, 1, 1,

	0.8, 0.8, 0.8, 1,
	0.8, 0.8, 0.8, 1,
	0.8, 0.8, 0.8, 1,
	0.8, 0.8, 0.8, 1,
}

var indexData = []uint16{
	0, 2, 1,
	3, 2, 0,

	4, 5, 6,
	6, 7, 4,

	8, 9, 10,
	11, 13, 12,

	14, 16, 15,
	17, 16, 14,
}
