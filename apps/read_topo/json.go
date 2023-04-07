package main
import (
    "fmt"
)

type Host struct {
    IP_Addr string
    Home string
    Name string
    Number int
}

type Network struct {
    Level string
    Name string
    Number int
}

type Intrfc struct {
    Bandwidth int
    Faces string
}

type Router struct {
    Name string
    Number int
    Intrfcs []Intrfc
}



func main() {
    type IArray []Intrfc
    IA := IArray {Intrfc {Bandwidth:1000,Faces:"network1"}, Intrfc{Bandwidth:5000,Faces:"network2"}}
 
    rtr := Router{Name:"router1", Number:66, Intrfcs:IA}


    fmt.Println(rtr.Name, rtr.Intrfcs[0].Faces)

}

