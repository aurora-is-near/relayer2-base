package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	tp "relayer2-base/tinypack"
)

type Person struct {
	Name         tp.VarData
	Age          uint64
	LivesOnEarth bool
	WeightKg     float64
}

func (p *Person) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&p.Name,
		&p.Age,
		&p.LivesOnEarth,
		&p.WeightKg,
	}, nil
}

type len3 struct{}

func (len3) GetTinyPackLength() int {
	return 3
}

type TeamOfThree struct {
	Members         tp.List[len3, tp.Pointer[Person]]
	FriendshipScore int64
}

func (t *TeamOfThree) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&t.Members,
		&t.FriendshipScore,
	}, nil
}

func main() {
	t := &TeamOfThree{
		Members: tp.CreateList[len3](
			tp.CreatePointer(&Person{
				Name:         tp.CreateVarData([]byte("Alex One")...),
				Age:          32,
				LivesOnEarth: true,
				WeightKg:     70.1,
			}),
			tp.CreatePointer(&Person{
				Name:         tp.CreateVarData([]byte("Alex Two")...),
				Age:          25,
				LivesOnEarth: false,
				WeightKg:     65.7,
			}),
			tp.CreatePointer(&Person{
				Name:         tp.CreateVarData([]byte("Alex Three")...),
				Age:          35,
				LivesOnEarth: true,
				WeightKg:     75.3,
			}),
		),
		FriendshipScore: 999999999999999,
	}

	data, err := tp.DefaultEncoder().Marshal(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(data))
	fmt.Println(data)

	t2 := new(TeamOfThree)
	err = tp.DefaultDecoder().Unmarshal(data, t2)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(t2)
}
