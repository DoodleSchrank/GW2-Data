package entities
import("time")

type Item struct {
	Id int
	Name string
	Buyprice float64
	Sellprice float64
	Avgbuyprice float64
	Avgsellprice float64
	Lastupdate time.Time
}
