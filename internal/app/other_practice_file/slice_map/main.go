package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name int
	Age  int
}

func (p Person) PersonName() (name int) {
	name = p.Name
	return
}

func main() {
	var slice1 []int = []int{1, 2, 3}
	var slice2 []int
	slice2 = append(slice2, 4)
	slice2 = append(slice2, slice1...)
	fmt.Println(slice2)
	slice2[0] = 5
	fmt.Println(slice2)
	slice3 := make([]int, 6)
	slice3[0] = 6
	slice3 = append(slice3, slice2...)
	fmt.Println(slice3)

	var map1 map[string]int
	var map2 map[string]int = map[string]int{
		"c": 3,
	}
	map3 := make(map[string]int)
	map2 = map[string]int{
		"a": 1,
	}
	// map1["b"] = 2
	map3["b"] = 2
	fmt.Println(map1, map2, map3)

	var p1 Person
	p1.Age = 1
	p1.Name = 1
	p2 := Person{
		Name: 2,
		Age:  2,
	}
	p3 := new(Person)
	p4 := &Person{}
	fmt.Println(p1.Age, p2, p3, *p4)

	return
	int1 := make(chan int, 1)
	int1 <- 10
	for i := 1; i <= 10; i++ {

		select {
		case i2 := <-int1:
			fmt.Println("i2", i2)
		case i1 := <-int1:
			fmt.Println(i1)
		case <-time.After(3 * time.Second):
			fmt.Println("timeafter")

		default:
			fmt.Println("default")
		}

	}
}
