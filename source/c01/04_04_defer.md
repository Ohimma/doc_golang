## 4.04 流程控制：延迟语句(defer)

- defer 的延迟顺序与最终的执行顺序是反向的 (先入后出)
- defer 按代码块执行顺序放入栈中，相关的值也会进入栈中，即使后面值有更改，栈输出是也是原值
- defer 是所在函数结束时进行，函数发生报错宕机也会执行

#### 1. 多个 defer 执行顺序

```
package main
import (
    "fmt"
)
func main() {
    fmt.Println("defer begin")

    // 将defer放入延迟调用栈
    defer fmt.Println(1)
    defer fmt.Println(2)
    // 最后一个放入, 位于栈顶, 最先调用
    defer fmt.Println(3)

    fmt.Println("defer end")
}
>>>
defer begin
defer end
3
2
1
```

#### 2. defer 与 return 孰先孰后

- defer 是 return 后才调用的

```
import "fmt"

var name string = "go"

func myfunc() string {
    defer func() {
        name = "python"
    }()

    fmt.Printf("myfunc 函数里的name：%s\n", name)
    return name
}

func main() {
    myname := myfunc()
    fmt.Printf("main 函数里的name: %s\n", name)
    fmt.Println("main 函数里的myname: ", myname)
}
>>>
myfunc 函数里的name：go
main 函数里的name: python
main 函数里的myname:  go
```

#### 3. 常用于在函数退出时释放资源

处理业务或逻辑中涉及成对的操作是一件比较烦琐的事情，比如打开和关闭文件、接收请求和回复请求、加锁和解锁等。在这些操作中，最容易忽略的就是在每个函数退出处正确地释放和关闭资源。 (Python 中没有 defer ，却用了 with 上下文管理器)

defer 语句正好是在函数退出时执行的语句，所以使用 defer 能非常方便地处理资源释放问题。

**1) 使用延迟并发解锁**

```
var (
    // 一个演示用的映射
    valueByKey      = make(map[string]int)
    // 保证使用映射时的并发安全的互斥锁
    valueByKeyGuard sync.Mutex
)
// 根据键读取值
func readValue(key string) int {
    // 对共享资源加锁
    valueByKeyGuard.Lock()

   // 取值
    v := valueByKey[key]
    // 对共享资源解锁
    valueByKeyGuard.Unlock()

    // 返回值
    return v
}

------------------更改后----------------------
func readValue(key string) int {
    valueByKeyGuard.Lock()

    // defer后面的语句不会马上调用, 而是延迟到函数结束时调用
    defer valueByKeyGuard.Unlock()

    return valueByKey[key]
}

```

**2) 使用延迟释放资源(文件句柄等)**

```
func f() {
    r := getResource()  //0，获取资源
    ......
    if ... {
        r.release()  //1，释放资源
        return
    }
    if ... {
        r.release()  //2，释放资源
        return
    }
    if ... {
        r.release()  //3，释放资源
        return
    }
    r.release()     //4，释放资源
    return
}
----------------------更改后--------------------------
func f() {
    r := getResource()  //0，获取资源

    defer r.release()  //1，释放资源
    ......
    if ... {
        return
    }
    if ... {
        return
    }
    if ... {
        return
    }
    return
}
```
