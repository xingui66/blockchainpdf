package main

type People1 struct {
	Name string
	Age  int
}

//不逃逸
func getInfo1(p *People1) *People1 {
	p.Name = "Jim"
	return p
}

func getInfo2(p People1) *People1 {
	return &p
}

func main() {
	p1 := People1{
		Name: "Lily",
		Age:  20,
	}
	_ = getInfo1(&p1) //没逃逸

	//_ = getInfo2(p1) //逃逸
}
