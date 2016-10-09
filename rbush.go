/**
 * Created by nazarigonzalez on 21/9/16.
 */

package rbush

import (
	"math"
)

type RBush struct {
	maxEntries float64
	minEntries float64
	data       *node
}

func NewRBush(maxEntries float64) RBush {
	if maxEntries == 0 {
		maxEntries = 9
	}

	return RBush{
		maxEntries: math.Max(4.0, maxEntries),
		minEntries: math.Max(2.0, math.Ceil(maxEntries*0.4)),
		data:       createNode([]*node{}),
	}
}

func (rbush *RBush) Load(data []Box) {
	l := len(data)
	if l < int(rbush.minEntries) {
		for i := 0; i < l; i++ {
			rbush.InsertBox(&data[i])
		}
		return
	}

	right := l - 1
	_node := rbush.build(copySliceBox(data), 0, 0.0, right, float64(right), 0)
	if len(rbush.data.children) == 0 {
		rbush.data = _node
	} else if rbush.data.height == _node.height {
		rbush.splitRoot(rbush.data, _node)
	} else {
		if rbush.data.height < _node.height {
			tmpNode := rbush.data
			rbush.data = _node
			_node = tmpNode
		}

		rbush.insert(_node, rbush.data.height-_node.height-1)
	}
}

func (rbush *RBush) build(items []Box, left int, lf float64, right int, rf float64, height int) *node {
	var _node *node

	N := rf - lf + 1
	M := rbush.maxEntries

	if N <= M {
		_node = createNode(boxToNodes(items[left : right+1]))
		_node.calcBBox()
		return _node
	}

	if height == 0 {
		heightFloat := math.Ceil(math.Log(N) / math.Log(M))
		height = int(heightFloat)

		M = math.Ceil(N / math.Pow(M, heightFloat-1))
	}

	_node = createNode([]*node{})
	_node.leaf = false
	_node.height = height

	N2 := math.Ceil(N / M)
	N1 := N2 * math.Ceil(math.Sqrt(M))

	var right2, right3 float64

	multiselect(items, lf, rf, N1, compareNodeMinX)

	for i := lf; i <= rf; i += N1 {
		right2 = math.Min(i+N1-1, rf)
		multiselect(items, i, right2, N2, compareNodeMinY)

		for j := i; j <= right2; j += N2 {
			right3 = math.Min(j+N2-1, right2)

			_node.children = append(_node.children, rbush.build(items, int(j), j, int(right3), right3, height-1))
		}
	}

	_node.calcBBox()
	return _node
}

func (rbush *RBush) splitRoot(_node *node, newNode *node) {
	rbush.data = createNode([]*node{_node, newNode})
	rbush.data.height = _node.height + 1
	rbush.data.leaf = false
	rbush.data.calcBBox()
}

func (rbush *RBush) InsertBox(box *Box) {
	_node := createNode([]*node{})
	_node.copyBox(box)
	rbush.insert(_node, rbush.data.height-1)
}

func (rbush *RBush) insert(item *node, level int) {
	insertPath := []*node{}
	_node := rbush.chooseSubtree(item, rbush.data, level, &insertPath)

	_node.children = append(_node.children, item)
	_node.extend(item)

	maxEntriesInt := int(rbush.maxEntries)

	for level >= 0 {
		if len(insertPath[level].children) > maxEntriesInt {
			rbush.split(insertPath, &level)
			level--
		} else {
			break
		}
	}

	rbush.adjustParentBBoxes(item, &insertPath, &level)
}

func (rbush *RBush) Search(box *Box) []*Box {
	_node := rbush.data
	result := []*Box{}

	if !intersects(box, _node) {
		return result
	}

	nodesToSearch := []*node{}

	for {
		ll := len(_node.children)
		for i := 0; i < ll; i++ {
			child := _node.children[i]

			if intersects(box, child) {
				if _node.leaf {
					result = append(result, child.box)
				} else if contains(box, child) {
					rbush.all(child, &result)
				} else {
					nodesToSearch = append(nodesToSearch, child)
				}
			}
		}

		if len(nodesToSearch) != 0 {
			_node, nodesToSearch = nodesToSearch[len(nodesToSearch)-1], nodesToSearch[:len(nodesToSearch)-1]
		} else {
			break
		}
	}

	return result
}

func (rbush *RBush) all(item *node, result *[]*Box) *[]*Box {
	nodesToSearch := []*node{}

	for {

		if item.leaf {
			for i := 0; i < len(item.children); i++ {
				*result = append(*result, item.children[i].box)
			}
		} else {
			nodesToSearch = append(nodesToSearch, item.children...)
		}

		if len(nodesToSearch) != 0 {
			item, nodesToSearch = nodesToSearch[len(nodesToSearch)-1], nodesToSearch[:len(nodesToSearch)-1]
		} else {
			break
		}
	}

	return result
}

func (rbush *RBush) adjustParentBBoxes(item *node, path *[]*node, level *int) {
	for i := *level; i >= 0; i-- {
		(*path)[i].extend(item)
	}
}

func (rbush *RBush) split(path []*node, level *int) {
	_node := path[*level]
	M := len(_node.children)
	m := int(rbush.minEntries)

	rbush.chooseSplitAxis(_node, m, M)
	splitIndex := chooseSplitIndex(_node, m, M)

	nodeCopy := make([]*node, len(_node.children)-splitIndex)
	copy(nodeCopy, _node.children[splitIndex:])
	newNode := createNode(nodeCopy)
	newNode.height = _node.height
	newNode.leaf = _node.leaf

  llen := len(_node.children)
  for i:=splitIndex; i < llen; i++ {
    _node.children[i] = nil
  }
	_node.children = _node.children[:splitIndex]

	_node.calcBBox()
	newNode.calcBBox()

	if *level > 0 {
		child := &path[*level-1].children
		*child = append(*child, newNode)
	} else {
		rbush.splitRoot(_node, newNode)
	}
}

func (rbush *RBush) chooseSplitAxis(_node *node, m, M int) {
	xMargin := allDistMargin(_node, m, M, "x")
	yMargin := allDistMargin(_node, m, M, "y")

	if xMargin < yMargin {
		sortByMinX(&_node.children)
	}
}

func (rbush *RBush) chooseSubtree(item *node, _node *node, level int, path *[]*node) *node {
	var (
		targetNode *node
	)

	definedNode := false

	for {
		*path = append(*path, _node)

		if _node.leaf || len(*path)-1 == level {
			break
		}

		minArea := math.Inf(1)
		minEnlargement := math.Inf(1)

		ll := len(_node.children)
		for i := 0; i < ll; i++ {

			child := _node.children[i]
			area := bboxArea(child)
			enlargement := enlargedArea(item, child) - area

			if enlargement < minEnlargement {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
				definedNode = true
			} else if enlargement == minEnlargement {
				if area < minArea {
					minArea = area
					targetNode = child
					definedNode = true
				}
			}

		}

		if definedNode {
			_node = targetNode
		} else {
			_node = _node.children[0]
		}

	}

	return _node
}

func (rbush *RBush) Clear() {
	rbush.data = createNode([]*node{})
}

func multiselect(arr []Box, left float64, right float64, n float64, compare callback) {
	stack := []float64{left, right}

	for len(stack) > 0 {
		right, left = stack[len(stack)-1], stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if right-left <= n {
			continue
		}

		mid := left + math.Ceil((right-left)/n/2)*n
		QuickSelect(arr, int(mid), int(left), int(right), compare)
		stack = append(stack, left, mid, mid, right)
	}
}

func copySliceBox(arr []Box) []Box {
	arr2 := make([]Box, len(arr))
	copy(arr2, arr)
	return arr2
}
