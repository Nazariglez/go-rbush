/**
 * Created by nazarigonzalez on 6/10/16.
 */

package rbush

import (
	"testing"
	"math"
	"math/rand"
)

var R RBush

func testCompareMinX(a, b *Box) float64 {
	return a.MinX - b.MinX
}

func TestQuickSelectDefault(t *testing.T) {
	a := []float64{65, 28, 59, 33, 21, 56, 22, 95, 50, 12, 90, 53, 28, 77, 39}
	arr2 := []Box{}
	for i:=0; i < len(a); i++ {
		arr2 = append(arr2, Box{MinX:a[i]})
	}

	QuickSelectDefault(arr2, 8, testCompareMinX)

	valid := true
	res := []float64{39, 28, 28, 33, 21, 12, 22, 50, 53, 56, 59, 65, 90, 77, 95}
	for i := 0; i < len(arr2); i++ {
		if arr2[i].MinX != res[i] {
			valid = false
			break
		}
	}

	if !valid {
		t.Error("Invalid int sequence")
	}
}

func TestQuickSelect(t *testing.T) {
	a := []float64{65, 28, 59, 33, 21, 56, 22, 95, 50, 12, 90, 53, 28, 77, 39}
	arr2 := []Box{}
	for i:=0; i < len(a); i++ {
		arr2 = append(arr2, Box{MinX:a[i]})
	}

	QuickSelect(arr2, 8, 0, len(arr2)-1, testCompareMinX)

	valid := true
	res := []float64{39, 28, 28, 33, 21, 12, 22, 50, 53, 56, 59, 65, 90, 77, 95}
	for i := 0; i < len(arr2); i++ {
		if arr2[i].MinX != res[i] {
			valid = false
			break
		}
	}

	if !valid {
		t.Error("Invalid int sequence")
	}
}

func TestNewRBush(t *testing.T) {
	R = NewRBush(9)

	if len(R.data.children) != 0 {
		t.Error("Expected len children = 0")
	}else if !R.data.leaf {
		t.Error("Expected leaf = false")
	}else if R.data.height != 1 {
		t.Error("Expected height = 1")
	}
}

func TestRbush_InsertBox(t *testing.T) {
 bbox := Box{
	 MinX: 10,
	 MinY: 20,
	 MaxX: 30,
	 MaxY: 40,
 }

	R.InsertBox(&bbox)

	if len(R.data.children) != 1 {
		t.Error("Expected len children = 1")
	}else if R.data.height != 1 {
		t.Error("Expected height = 1")
	}else if !R.data.leaf {
		t.Error("Expected leaf = false")
	}else if (R.data.MinX != bbox.MinX || R.data.MinY != bbox.MinY || R.data.MaxX != bbox.MaxX|| R.data.MaxY != bbox.MaxY) {
		t.Error("Expected Same bbox in node and box")
	}
}

func TestRbush_Clear(t *testing.T) {
	if len(R.data.children) != 1 {
		t.Error("Expected len children 1")
	}

	R.Clear()

	if len(R.data.children) != 0 {
		t.Error("Expected len children 0")
	}
}

func TestRbush_InsertBox2(t *testing.T) {
	for i := 0.0; i < 20; i++{
		R.InsertBox(&Box{
			MinX: i*10 - (i+1),
			MinY: i*10 - (i+1),
			MaxX: i*15 - (i+1),
			MaxY: i*15 - (i+1),
			Data: i,
		})
	}

	finalValues := Box{
		MinX:-1,
		MinY:-1,
		MaxX:265,
		MaxY:265,
	}

	if len(R.data.children) != 4 {
		t.Error("Expected len children 4")
	}else if R.data.height != 2 {
		t.Error("Expected height 2")
	}else if R.data.leaf {
		t.Error("Expected leaf false")
	}else if (R.data.MinX != finalValues.MinX || R.data.MinY != finalValues.MinY || R.data.MaxX != finalValues.MaxX|| R.data.MaxY != finalValues.MaxY) {
		t.Error("Invalid final Box values")
	}

	box1 := Box{MaxX:41,MaxY:41,MinX:-1,MinY:-1}
	box2 := Box{MaxX:97,MaxY:97,MinX:35,MinY:35}
	box3 := Box{MaxX:153,MaxY:153,MinX:71,MinY:71}
	box4 := Box{MaxX:265,MaxY:265,MinX:107,MinY:107}

	var bbox *Box

	for i := 0; i < len(R.data.children); i++ {
		children := R.data.children[i]

		if (i<3 && len(children.children) != 4){
			t.Error("Expected 4 childrens")
			break
		}else if (i==3 && len(children.children) != 8) {
			t.Error("Expected 8 childrens")
			break
		}else if children.height != 1 {
			t.Error("Expected children height 1")
			break;
		}else if !children.leaf {
			t.Error("Expected children leaf true")
			break;
		}

		switch i {
		case 0:
			bbox = &box1
		case 1:
			bbox = &box2
		case 2:
			bbox = &box3
		case 3:
			bbox = &box4
		}

		if (children.MinX != bbox.MinX || children.MinY != bbox.MinY || children.MaxX != bbox.MaxX|| children.MaxY != bbox.MaxY) {
			t.Error("Expected valid box values")
			break
		}
	}
}

func TestRbush_Search(t *testing.T) {
	boxes := R.Search(&Box{
		MinX: 20, MinY: 20,
		MaxX: 45, MaxY: 45,
	})

	if len(boxes) != 4 {
		t.Error("Expected 4 collision boxes")
	}

	collisionId := []float64{2,3,4,5}

	for i := range boxes {
		id := boxes[i].Data.(float64)
		valid := false

		for j := range collisionId {
			if collisionId[j] == id {
				valid = true
				break
			}
		}

		if !valid {
			t.Error("Expected a valid collision box ID")
			break
		}
	}
}

func TestRbush_Search2(t *testing.T) {
	boxes := R.Search(&Box{
		MinX: 20, MinY: 20,
		MaxX: 300, MaxY: 300,
	})

	if len(boxes) != 18 {
		t.Error("Expected 18 collision boxes")
	}
}

func TestRbush_Load(t *testing.T) {
	R.Clear()

	boxes := []Box{}
	for i := 0.0; i < 10; i++ {
		boxes = append(boxes, Box{
			MinX: i*10-(i*2),
			MinY: i*10-(i*2),
			MaxX: i*10-(i*2) + i*15-(i*2),
			MaxY: i*10-(i*2) + i*15-(i*2),
			Data: int(i),
		})
	}

	R.Load(boxes)

	if len(R.data.children) != 2 {
		t.Error("Expected childre = 2")
	}else if(len(R.data.children[0].children) != 5 || len(R.data.children[1].children) != 5){
		t.Error("Expected node with 5 childrens")
	}

	for i := range R.data.children[0].children {
		if R.data.children[0].children[i].box.Data.(int) != i {
			t.Error("Expected a box with the same id and index")
			break
		}
	}

	for i := range R.data.children[1].children {
		n := i+5
		if R.data.children[1].children[i].box.Data.(int) != n {
			t.Error("Expected a box with the same id and index")
			break
		}
	}
}

func BenchmarkQuickSelect(b *testing.B) {
	b.SetBytes(2)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		arr := []Box{}

		for i := 0; i < b.N; i++ {
			arr = append(arr, Box{
				MinX: rand.Float64(),
				MaxX: rand.Float64(),
			})
		}

		l := len(arr)
		k := int(math.Floor(float64(b.N/2)))

		b.StartTimer()
		QuickSelect(arr, k, 0, l-1, testCompareMinX)
	}
}

func BenchmarkQuickSelectDefault(b *testing.B) {
	b.SetBytes(2)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		arr := []Box{}

		for i := 0; i < b.N; i++ {
			arr = append(arr, Box{
				MinX: rand.Float64(),
				MaxX: rand.Float64(),
			})
		}

		k := int(math.Floor(float64(b.N/2)))
		b.StartTimer()
		QuickSelectDefault(arr, k, testCompareMinX)
	}
}

func BenchmarkRbush_InsertBox(b *testing.B) {
	b.StopTimer()
	R.Clear()

	boxes := []Box{}
	for i:=0; i < b.N; i++{
		boxes = append(boxes, Box{
			MinX: rand.Float64()*50,
			MinY: rand.Float64()*50,
			MaxX: 25 + rand.Float64()*200,
			MaxY: 25 + rand.Float64()*200,
		})
	}

	for i:=0; i < b.N; i++{
		b.StartTimer()
		R.InsertBox(&boxes[i])
		b.StopTimer()
	}
}

func BenchmarkRbush_Load(b *testing.B) {
	b.StopTimer()
	R.Clear()

	boxes := []Box{}
	for i:=0; i < b.N; i++{
		boxes = append(boxes, Box{
			MinX: rand.Float64()*50,
			MinY: rand.Float64()*50,
			MaxX: 25 + rand.Float64()*200,
			MaxY: 25 + rand.Float64()*200,
		})
	}

	b.StartTimer()
	R.Load(boxes)
}

func BenchmarkRbush_Load2(b *testing.B) {
	b.StopTimer()
	R.Clear()

	boxes := []Box{}
	for i:=0; i < b.N; i++{
		boxes = append(boxes, Box{
			MinX: rand.Float64()*50,
			MinY: rand.Float64()*50,
			MaxX: 25 + rand.Float64()*200,
			MaxY: 25 + rand.Float64()*200,
		})
	}

	for i:=0; i < b.N; i++ {
		R.Clear()
		b.StartTimer()
		R.Load(boxes)
		b.StopTimer()
	}
}

func BenchmarkRbush_Search(b *testing.B) {
	b.StopTimer()

	boxes := []Box{}
	for i:=0; i < b.N; i++{
		boxes = append(boxes, Box{
			MinX: rand.Float64()*10,
			MinY: rand.Float64()*10,
			MaxX: 10 + rand.Float64()*250,
			MaxY: 10 + rand.Float64()*250,
		})
	}

	for i:=0; i < b.N; i++{
		b.StartTimer()
		_ = R.Search(&boxes[i])
		b.StopTimer()
	}
}