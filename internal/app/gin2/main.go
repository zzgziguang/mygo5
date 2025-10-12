package main

import "fmt"

type people struct {
	Name string
	Age  int
}

func main() {
	var p1 people
	p1 = people{
		Name: "赵",
		Age:  11,
	}
	var p2 = people{
		Name: "钱",
		Age:  22,
	}
	p3 := people{
		Name: "孙",
		Age:  33,
	}
	var p4 []people
	p4 = append(p4, p1, p2, p3)
	p4len := len(p4)
	for i, v := range p4 {
		if i == p4len-2 {
			p4[i-1] = v
			v.Age = 5
			fmt.Println(p4[i-1],v)
		}
	}
	var p5 map[string]people
	p5 = make(map[string]people, 0)
	p5["a"]=p1
	p5["b"]=p2
	p6 := map[string]people{
		"c":people{
			Age: 44,
			Name: "李",
		}
	}
	for k,v:=range p5{
		pm1,ok:=p5[]
	}
	
}
