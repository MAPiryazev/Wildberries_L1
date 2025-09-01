package main

import "fmt"

type Human struct {
	weight int
	height int
}

func (human *Human) showParams() {
	fmt.Printf("Hey im %d tall and my weight is %d\n", human.height, human.weight)
}

type Worker struct {
	Human
	proffession string
}

func (worker *Worker) showProfession() {
	fmt.Printf("Im a %s\n", worker.proffession)
}

func main() {
	wrkr := Worker{Human: Human{weight: 70,
		height: 180},
		proffession: "baker"}

	wrkr.showParams()
	wrkr.showProfession()
}
