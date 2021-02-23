package main
import (
    "fmt"
    "./model"
)

func main() {
    p := model.NewPerson("tom")
    fmt.Println(p)
    p.SetAge(18)
    p.SetSal(5000)
    fmt.Println(p)
    fmt.Printf(p.Name, "age=", p.GetAge() , "sal=", p.GetSal())
}