package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type people struct {
	Name string
	Age  int
	Yo   string
}

func (p *people) SayHello(t time.Time) (text string) {
	afternoon := t.Hour() <= 12
	if afternoon {
		text = fmt.Sprintf("大家上午好，我是%s", p.Name)
	} else {
		text = fmt.Sprintf("大家下午好，我是%s", p.Name)
	}
	return
}

func YoungOld(peopleage int) (yo string) {
	if peopleage > 40 {
		yo = "old"
	} else {
		yo = "young"
	}
	return
}

func main() {
	var p1 people
	p1 = people{
		Name: "赵",
		Age:  22,
	}
	p1.Yo = YoungOld(p1.Age)

	var p2 = people{
		Name: "钱",
		Age:  55,
	}
	p2.Yo = YoungOld(p2.Age)

	p3 := people{
		Name: "孙",
		Age:  40,
	}
	p3.Yo = YoungOld(p3.Age)

	p4 := people{
		Name: "李",
		Age:  33,
	}
	p4.Yo = YoungOld(p4.Age)
	t1 := time.Now()
	text := p4.SayHello(t1)
	fmt.Println("打招呼", text)

	var p5 []people
	p5 = append(p5, p1, p2, p3, p4)
	p5l := len(p5)
	j := rand.Intn(9)
	if j >= 0 && j < p5l {
		pj := p5[j]
		pj.Name = "第" + strconv.Itoa(j+1) + "个人"
		p5[j] = pj
		fmt.Println("slice", p5[j])
	}

	var p6 map[int]people
	p6 = make(map[int]people, 0)

	for i, v := range p5 {
		v.SayHello(t1)
		v = p5[i]
		p5[i] = v
		fmt.Println("slice", p5[i])

		p6[i+1] = v
	}

	fmt.Println("map", p6)

	k := rand.Intn(9)
	pk, ok := p6[k]
	if ok {
		pk.Yo = "old" //
		p6[k] = pk
		fmt.Println("map", p6[k])
	}

	for k, v := range p6 {
		v.SayHello(t1)
		v = p6[k]
		p6[k] = v
		fmt.Println("map", p6[k])
	}
}
