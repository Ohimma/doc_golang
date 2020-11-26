## 5.2 函数：错误处理

#### 1) defer recover，捕获异常不终止程序
```
package main
import "fmt"
func test() {
    n1 := 100
    n2 := 0
    res := n1 / n2
    fmt.Printf("print test() res=%v\n", res) 
}

func main() {
    test()
    fmt.Println("print main()......")
}
>>>
panic: runtime error: integer divide by zero
xxxxxxxxx
---------------------------------------
package main
import "fmt"


func test() {
    //defer 匿名函数
    defer func() {
        err := recover()    // recover内置函数，捕获异常
        if err != nil {
            fmt.Printf("err=%v\n", err)
            fmt.Println("发送邮件给xxxxx@qq.com")
        }
    }()
    n1 := 100
    n2 := 0
    res := n1 / n2
    fmt.Printf("print test() res=%v\n", res) 
}

func main() {
    test()
    fmt.Println("print main()......")
}

>>>
err=runtime error: integer divide by zero
发送邮件给xxxxx@qq.com
print main()......
```



#### 2) 自定义错误, panic 终止程序


```
package main
import (
    "fmt"
    "errors"
)

func readConf(name string) (ere error) {
    if name == "config.ini" {
        return nil
    } else {
        return errors.New("读取文件错误")    // 自定义一个错误并返回
    }
}

func test() {
    err := readConf("config")
    if err != nil {
        fmt.Printf("print test() err %v\n", err)
        panic(err)     // 读取输出错误，并终止程序
    }
    fmt.Println("print test()......")
}

func main() {
    test()
    fmt.Println("print main()......")
}

>>>
print test() err 读取文件错误
panic: 读取文件错误

goroutine 1 [running]:
main.test()
        E:/me/code/go/test.go:19 +0xd5
main.main()
        E:/me/code/go/test.go:25 +0x29
exit status 2
```


