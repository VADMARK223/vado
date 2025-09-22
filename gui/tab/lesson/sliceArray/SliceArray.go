package sliceArray

import (
	"fmt"
)

func RunSliceArray() {
	fmt.Println("RunSliceArray")
	var a = [5]int{11, 21, 31, 41, 51}
	var b = [5]int{1, 2, 3, 4, 5}
	a = b
	fmt.Println(a, b)

	arr := [5]int{11, 21, 31, 41, 51}
	sl := arr[1:5]
	sl = append(sl, 99)
	fmt.Println("sl", sl, len(sl), cap(sl))

	var sl1 = make([]int, 0, 40)
	sl1 = append(sl1, 13)
	fmt.Println(sl1, len(sl1), cap(sl1))

	var arr1 = [5]int{11, 21, 31, 41, 51}
	changeArray1(&arr1)
	fmt.Println(arr1)
	var sl2 = []int{1, 2, 3, 4, 5}
	changeSlice2(sl2)
	fmt.Println("sl2", sl2)

	var sl3 = make([]int, 3, 40)
	fmt.Println(sl3, len(sl3), cap(sl3))
}

func changeSlice2(sl2 []int) {
	sl2[0] = 999
}

func changeArray1(v *[5]int) {
	v[0] = 99
}
