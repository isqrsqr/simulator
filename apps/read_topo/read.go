package main
import (
    "fmt"
)

type Inside struct {
    Name string
    Active bool
    Value float64
}

type Outside struct {
    Name string
    Id int
    Interfaces []Inside 
}


func main() {
    type IArray []Inside
    IA := IArray {Inside{"if1", true, float64(1.0),}, Inside{"if2", false, float64(2.0),}}

    example := Outside{Name:"example", Id:5,Interfaces:IA}

    fmt.Println(example.Name)

}

