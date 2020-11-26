## 1.3 派生类型: 接口(interface)

golang 的核心就是接口

#### 1. 接口定义
面向对象领域中: 接口是指定义一个对象的行为。至于如何实现这个行为，则由对象本身去确定

* 接口是方法的集合，接口指定了一个类型应该具有的方法，然后由该类型决定的如何实现
* 一个自定义类型需要把接口的方法都实现，才能说该类型实现了该接口
* 接口可以定义一组方法，但可以不用全实现，到某个自定义类型(比如结构体)用到的时候在实现就行
* 接口是引用类型

1. 接口本身不能创建实例，都需要指向一个实现了该接口方法的自定义类型(struct)
2. 只要是自定义类型都可以实现接口，不仅仅是结构体
3. 一个自定义类型可以实现多个接口
4. golang 接口中不能有任何变量
5. 一个接口可以继承多个接口，如果实现该接口，需要把继承的接口全都实现
6. 空接口没有任何方法，所以所有类型都实现了空接口，空接口也可以接收任何变量

```
package main
import "fmt"

/* 定义接口 (必须实现所有方法才算实现该接口) */
// type A interface {
//     method_name1() [return_type]
//     method_namen() [return_type]
//  }

// func (struct_name_variable struct_name) method_name1() [return_type] {
//     .....
// }
 
 
 type A interface {
    Say()
 }
 
 type B interface {
    Hello()
 }

 type C interface {
    A
    B
    test()
 }
 
 type D interface {
 }

 // 1. 定义结构体 实现接口
 type struct_name struct {
 }
 
 func (s struct_name) Say() {
     fmt.Println("struct_name say()")
 }
 
 func (s struct_name) Hello() {
     fmt.Println("struct_name Hello()")
 }
 
 func (s struct_name) test() {
    fmt.Println("struct_name test()")
}

 // 2. 定义整型 实现接口
 type integer int
 
 func (i integer) Say() {
     fmt.Println("integer say i=", i)
 }
 
 
 func main() {
     // var i A   报错，本身不能创建实例
 
     var s struct_name
     var a1 A = s        // s 结构体调用 A 接口Say方法
     var b  B = s        // s 结构体调用 B 接口Hello方法
     a1.Say()
     b.Hello()

     var c C = s          // s 结构体调用 C 接口的 A 接口的Say()
     c.Say()              // == c.A.Say()
     
     var i integer = 10
     var a2 A = i       // i 整型 类型调用 A 接口Say()方法
     a2.Say()

     var d D = s        // s 结构体实现了 D 空接口 (== var d interface{} = s)  
     fmt.Println("main() d is ", d)
     var num1 float64 = 8.8
     d = num1
     fmt.Println("main() d is ", d)
 }

>>>
struct_name say()
struct_name Hello()
struct_name say()
integer say i= 10
main() d is  {}
main() d is  8.8
```



#### 2. 案例

**a) 实现对 hero 结构体切片的排序： sort.Sort(date interface)**
```
package main
import (
    "fmt"
    "sort"
    "math/rand"
)

// 1. 定义结构体
type Hero struct {
    Name string
    Age int
}
// 2. 定义结构体切片
type Heroslice []Hero 

// 3. 实现接口--方法 获取长度
func (hs Heroslice) Len() int {
    return len(hs)
}

// 4. 实现接口--less方法  使用什么标准进行排序(从小到大)
func (hs Heroslice) Less(i int, j int) bool {
    return hs[i].Age > hs[j].Age
}

func (hs Heroslice) Swap(i int, j int) {
    temp := hs[i]
    hs[i] = hs[j]
    hs[j] = temp
}

func main() {
    // 1. 对一个数组/切片 进行排序
    var slice = []int{0, -1, 10, 20, 50}

    sort.Ints(slice)
    fmt.Printf("slice is %v\n", slice)

    // 2. 对一个结构体切片进行排序
    // 冒泡排序 + 
    var heros Heroslice
    for i := 0; i < 5; i++ {
        hero := Hero {
            Name: fmt.Sprintf("英雄~~%d", rand.Intn(100)),
            Age: rand.Intn(100),
        } 
        heros = append(heros, hero)
    }

    fmt.Println("排序前.......", heros)
    for _ , v := range heros {
        fmt.Println(v)
    }

    sort.Sort(heros)

    fmt.Println("排序后.......", heros)
    for _ , v := range heros {
        fmt.Println(v)
    }
}

>>>
slice is [-1 0 10 20 50]
排序前....... [{英雄~~81 87} {英雄~~47 59} {英雄~~81 18} {英雄~~25 40} {英雄~~56 0}]
{英雄~~81 87}
{英雄~~47 59}
{英雄~~81 18}
{英雄~~25 40}
{英雄~~56 0}
排序后....... [{英雄~~81 87} {英雄~~47 59} {英雄~~25 40} {英雄~~81 18} {英雄~~56 0}]
{英雄~~81 87}
{英雄~~47 59}
{英雄~~25 40}
{英雄~~81 18}
{英雄~~56 0}
```

