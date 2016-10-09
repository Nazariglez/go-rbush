/**
 * Created by nazarigonzalez on 5/10/16.
 */

package rbush

import (
	"math"
	"sort"
)

type node struct {
	children               []*node
	height                 int
	leaf                   bool
	MinX, MaxX, MinY, MaxY float64
	box                    *Box
}

type byMinX []*node
type byMinY []*node

func (a byMinX) Len() int           { return len(a) }
func (a byMinY) Len() int           { return len(a) }
func (a byMinX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMinY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMinX) Less(i, j int) bool { return a[i].MinX < a[j].MinX }
func (a byMinY) Less(i, j int) bool { return a[i].MinY < a[j].MinY }

func sortByMinX(nodes *[]*node) {
	sort.Sort(byMinX(*nodes))
}

func sortByMinY(nodes *[]*node) {
	sort.Sort(byMinY(*nodes))
}

func (node *node) calcBBox() {
	node.distBBox(node, 0, len(node.children))
}

func (destNode *node) distBBox(node *node, k int, p int) {
	destNode.MinX = math.Inf(1)
	destNode.MinY = math.Inf(1)
	destNode.MaxX = math.Inf(-1)
	destNode.MaxY = math.Inf(-1)

	for i := k; i < p; i++ {
		destNode.extend(node.children[i])
	}
}

func (node *node) extend(b *node) {
	node.MinX = math.Min(node.MinX, b.MinX)
	node.MinY = math.Min(node.MinY, b.MinY)
	node.MaxX = math.Max(node.MaxX, b.MaxX)
	node.MaxY = math.Max(node.MaxY, b.MaxY)
}

func compareNodeMinX(a, b *Box) float64 {
	return a.MinX - b.MinX
}

func compareNodeMinY(a, b *Box) float64 {
	return a.MinY - b.MinY
}

func bboxArea(a *node) float64 {
	return (a.MaxX - a.MinX) * (a.MaxY - a.MinY)
}

func bboxMargin(a *node) float64 {
	return a.MaxX - a.MinX + (a.MaxY - a.MinY)
}

func enlargedArea(a, b *node) float64 {
	return (math.Max(b.MaxX, a.MaxX) - math.Min(b.MinX, a.MinX)) * (math.Max(b.MaxY, a.MaxY) - math.Min(b.MinY, a.MinY))
}

func contains(a *Box, b *node) bool {
	return a.MinX <= b.MinX &&
		a.MinY <= b.MinY &&
		b.MaxX <= a.MaxX &&
		b.MaxY <= a.MaxY
}

func intersects(a *Box, b *node) bool {
	return b.MinX <= a.MaxX &&
		b.MinY <= a.MaxY &&
		b.MaxX >= a.MinX &&
		b.MaxY >= a.MinY
}

func intersectionArea(a, b *node) float64 {
	minX := math.Max(a.MinX, b.MinX)
	minY := math.Max(a.MinY, b.MinY)
	maxX := math.Min(a.MaxX, b.MaxX)
	maxY := math.Min(a.MaxY, b.MaxY)

	return math.Max(0, maxX-minX) * math.Max(0, maxY-minY)
}

func createNode(children []*node) *node {
	return &node{
		children: children,
		height:   1,
		leaf:     true,
		MinX:     math.Inf(1),
		MinY:     math.Inf(1),
		MaxX:     math.Inf(-1),
		MaxY:     math.Inf(-1),
	}
}

func allDistMargin(_node *node, m, M int, prop string) float64 {
	if prop == "x" {
		sortByMinX(&_node.children)
	} else {
		sortByMinY(&_node.children)
	}

	leftBBox := createNode([]*node{})
	leftBBox.distBBox(_node, 0, m)
	rightBBox := createNode([]*node{})
	rightBBox.distBBox(_node, M-m, M)

	margin := bboxMargin(leftBBox) + bboxMargin(rightBBox)

	var child *node

	for i := m; i < M-m; i++ {
		child = _node.children[i]
		leftBBox.extend(child)
		margin += bboxMargin(leftBBox)
	}

	for i := M - m - 1; i >= m; i-- {
		child = _node.children[i]
		rightBBox.extend(child)
		margin += bboxMargin(rightBBox)
	}

	return margin
}

func chooseSplitIndex(_node *node, m, M int) int {
	var bbox1, bbox2 *node
	var overlap, area float64
	var i, index int

	minOverlap := math.Inf(1)
	minArea := math.Inf(1)

	for i = m; i <= M-m; i++ {
		bbox1 = createNode([]*node{})
		bbox1.distBBox(_node, 0, i)
		bbox2 = createNode([]*node{})
		bbox2.distBBox(_node, i, M)

		overlap = intersectionArea(bbox1, bbox2)
		area = bboxArea(bbox1) + bboxArea(bbox2)

		if overlap < minOverlap {
			minOverlap = overlap
			index = i

			if area < minArea {
				minArea = area
			}
		} else if overlap == minOverlap {
			if area < minArea {
				minArea = area
				index = i
			}
		}
	}
	return index
}

func (_node *node) copyBox(box *Box) {
	_node.MinX = box.MinX
	_node.MinY = box.MinY
	_node.MaxX = box.MaxX
	_node.MaxY = box.MaxY
	_node.box = box
}

func boxToNodes(arr []Box) []*node {
	var _n *node
	nodes := []*node{}
	for i := 0; i < len(arr); i++ {
		_n = createNode([]*node{})
		_n.copyBox(&arr[i])
		nodes = append(nodes, _n)
	}
	return nodes
}
