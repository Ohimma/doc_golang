## 2.2 面向对象: 类型断言(type)

     
1. 类型断言 类型断言就是将接口类型的值x,装换成类型T。x.(T) or v:=x.(T) or v, ok := x.(T)
2. 类型断言的必要条件就是x是接口类型,非接口类型的x不能做类型断言

3. 类型切换(type switch)，用于比较类型而不是值，x.(type) 用于检测x的类型


**a) 之前多态案例：电脑usb接口，实现phone除了start和stop，外加一个call方法。camera不变**   
```
package main
import "fmt"

type Usb interface {
    Start()
    // Call()    这样不行，因为camera不用实现该方法  
    Stop()
}

type Phone struct {
    Name string
}

type Camera struct {
    name string
}
// Phone 是实现了接口的两个方法
func (p Phone) Start() {
    fmt.Println("手机开始工作.......")
}
func (p Phone) Call() {
    fmt.Println("手机开始打电话.......")
}
func (p Phone) Stop() {
    fmt.Println("手机停止工作.......")
}
// Camera 是实现了接口的两个方法
func (c Camera) Start() {
    fmt.Println("相机开始工作.......")
}
func (c Camera) Stop() {
    fmt.Println("相机停止工作.......")
}

type Computer struct {

}

// 编写一个方法实现usb接口类型，只要是实现usb接口，就是指usb接口声明的所有方法
// 而以上的 Phone 和 Camera 都实现了 Usb 的接口(start和stop方法)
func (c Computer) Working(usb Usb) {
    usb.Start()
    phone, ok := usb.(Phone)       // 类型断言
    if ok {
        phone.Call()
    }
    usb.Stop()
}

func main() {
    // 1. usb 多态参数特性
    computer := Computer{}

    // 2. usb 数组特性
    var usbArr [3]Usb
    usbArr[0] = Phone{"vivo"}
    usbArr[1] = Phone{"小米"}
    usbArr[2] = Camera{"尼康"}

    for _, v := range usbArr {
        computer.Working(v)
        fmt.Println()
    }

    
}

>>>
手机开始工作.......
手机开始打电话.......
手机停止工作.......

手机开始工作.......
手机开始打电话.......
手机停止工作.......

相机开始工作.......
相机停止工作.......
```

**b) 写一个函数，循环判断传入参数的类型**    

```
package main
import "fmt"

type Student struct {

}

func TypeJudge(items ...interface{}) {
    for index, v := range items {
        switch v.(type) {                  // v.(type) 是类型断言
            case bool:
                fmt.Printf("第%v个参数是bool类型，值是%v \n", index, v)
            case float32:
                fmt.Printf("第%v个参数是float32类型，值是%v \n", index, v)
            case float64:
                fmt.Printf("第%v个参数是float64类型，值是%v \n", index, v)   
            case int, int32, int64:
                fmt.Printf("第%v个参数是int类型，值是%v \n", index, v)
            case string:
                fmt.Printf("第%v个参数是string类型，值是%v \n", index, v)
            case Student:
                fmt.Printf("第%v个参数是Student类型，值是%v \n", index, v)
            case *Student:
                fmt.Printf("第%v个参数是*Student类型，值是%v \n", index, v)
            default:
                fmt.Printf("第%v个参数类型不确定，值是%v \n", index, v)
        }
    }
}


func main() {
    var n1  float32 = 1.1
    var n2  float64 = 2.2
    var n3  int32 = 32
    var n4  string = "tom"
    n5 := 500
    n6 := "北京"

    stu1 := Student{}
    stu2 := &Student{}

    TypeJudge(n1, n2, n3, n4, n5, n6, stu1, stu2)
}

>>>
第0个参数是float32类型，值是1.1
第1个参数是float64类型，值是2.2
第2个参数是int类型，值是32
第3个参数是string类型，值是tom
第4个参数是int类型，值是500
第5个参数是string类型，值是北京
第6个参数是Student类型，值是{}
第7个参数是*Student类型，值是&{}
```
