/**
 * Created by nazarigonzalez on 21/9/16.
 */

package rbush

import (
	"math"
)

type Box struct {
	MinX, MaxX, MinY, MaxY float64
	Data interface{}
}

type callback func(a, b *Box) float64

func QuickSelect(arr []Box, k int, left int, right int, compare callback){
	kf := float64(k)
	lf := float64(left)
	rf := float64(right)
	quickSelect(arr, k, kf, left, lf, right, rf, compare)
}

func QuickSelectDefault(arr []Box, k int, compare callback){
	kf := float64(k)
	right := len(arr)-1
	rf := float64(right)
	quickSelect(arr, k, kf, 0, 0.0, right, rf, compare)
}

func quickSelect(arr []Box, k int, kf float64, left int, lf float64, right int, rf float64, compare callback){

	for right > left {
		if right-left > 600 {
			var ss float64

			n := rf - lf+1
			m := kf - lf+1
			z := math.Log(n)
			s := 0.5 * math.Exp(2*z/3)

			if m-n/2 < 0 {
				ss = -1
			}else{
				ss = 1
			}

			sd := 0.5 * math.Sqrt(z*s*(n-s)*ss)
			nl := math.Floor(kf - m*s/n + sd)
			nr := math.Floor(kf + (n-m) * s/n +sd)
			newLeft := math.Max(lf, nl)
			newRight := math.Min(rf, nr)

			quickSelect(arr, k, kf, int(newLeft), newLeft, int(newRight), newRight, compare)
		}

		t := arr[k]
		i := left
		j := right

		swap(arr, left, k)
		if compare(&arr[right], &t) > 0{
			swap(arr, left, right)
		}

		for i < j {
			swap(arr, i, j)
			i++
			j--

			for compare(&arr[i], &t) < 0 {
				i++
			}

			for compare(&arr[j], &t) > 0 {
				j--
			}
		}

		if compare(&arr[left], &t) == 0 {
			swap(arr, left, j)
		}else{
			j++
			swap(arr, j, right)
		}

		if j <= k {
			left = j+1
		}

		if k <= j {
			right = j-1
		}

	}
}

func swap(arr []Box, i, j int){
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}