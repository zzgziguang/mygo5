package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Student struct {
	Id    int
	Name  string
	Class int
}

func (s *Student) SayHello(t time.Time) (text string) {
	afternoon := t.Hour() <= 12
	if afternoon {
		text = fmt.Sprintf("大家上午好，我是 %s", s.Name)
	} else {
		text = fmt.Sprintf("大家下午好，我是 %s", s.Name)
	}
	return
}

func AssignClass(studentid int) (class int) {
	switch studentid % 3 {
	case 1:
		class = 1
	case 2:
		class = 2
	default:
		class = 3
	}
	return
}

func main() {
	var s1 Student
	s1 = Student{
		Id:   1,
		Name: "aaa",
	}
	s1.Class = AssignClass(s1.Id)

	var s2 = Student{
		Id:   2,
		Name: "bbb",
	}
	s2.Class = AssignClass(s2.Id)

	s3 := Student{
		Id:   3,
		Name: "ccc",
	}
	s3.Class = AssignClass(s3.Id)

	s4 := Student{
		Id:   4,
		Name: "ddd",
	}
	s4.Class = AssignClass(s4.Id)
	t1 := time.Now()
	text := s4.SayHello(t1)
	fmt.Println(text)

	var ss []Student
	ss = append(ss, s1, s2, s3, s4)

	ssLen := len(ss)

	i := rand.Intn(9)
	if i >= 0 && i < ssLen {
		studenti := ss[i]
		studenti.Name = "学生" + strconv.Itoa(i+1)
		ss[i] = studenti
	}

	var sm map[int]Student
	sm = make(map[int]Student, 0)
	studentmap := make(map[int]Student, 0)
	fmt.Println(studentmap)

	for i, v := range ss {
		v.SayHello(t1)
		v = ss[i]
		ss[i] = v
		sm[v.Id] = v
	}

	k := rand.Intn(9)
	studentk, ok := sm[k]
	if ok {
		studentk.Class = 1
		sm[k] = studentk
	}
	for k, v := range sm {
		v.SayHello(t1)
		v = sm[k]
		sm[k] = v
	}
	// s2 = s[1]
	// k["a"]
}
