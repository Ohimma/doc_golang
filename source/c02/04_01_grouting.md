## 4.1 核心编程: 协程(grouting)

#### 1. 基本概念

从语言层面天生支持并发是 golang 最强大的特性之一。      
其他主流编程语言，如果要提升并发能力，就需要学习创建进程池/线程池，学习各种并发库等，分别应用到不同的场景；但是 golang 原生的协程，能自动帮你做这些事。

**a) 进程和线程**     
1. 进程是操作系统的执行过程，系统进行资源分配的最小单位
2. 线程是进程的执行实例，是程序执行的最小单位，比进程更小能独立运行的基本单位
3. 一个程序至少有一个进程，一个进程至少有一个线程，一个进程可并发多个线程


**b) 并发和并行**
1. 并发是指 多线程的程序在单核上运行
2. 并行是指 多线程的程序在多个核上运行


**c) go主线程 和 go协程**
1. 主线程可以理解为进程/线程，go 线程可以起多个协程，协程可以理解为轻量级线程
2. 主线程是物理线程，作用在 cpu 上，重量级的耗资源

1. 协程 有独立的栈空间
2. 协程 共享程序的堆空间
3. 协程 调度由用户控制
4. 协程 是轻量级的线程，是逻辑态的资源消耗小(编译器的优化)
5. 协程 可轻松开启上万协程，其他语言用线程实现并发，基于内核态(java/c)资源消耗大

**d) gorouting 的调度模型 MPG**  
1. M 是操作系统主线程 (物理线程)  
2. P 是协程执行需要的上下文
3. G 是协程

MPG 有一种情况是当 M0 的 G0 协程阻塞，系统会创建 M1 主线程 执行其余的等待 协程



#### 2. 协程使用

**a) 基本使用 go xxx**
```
package main
import (
    "fmt"
    "time"
    "strconv"
)

func test() {
    for i := 0; i < 5; i++ {
        fmt.Println("test() Hello,world", strconv.Itoa(i))
        time.Sleep(time.Second)
    }
}

func test2() {
    for i := 0; i < 5; i++ {
        fmt.Println("test2() Hello,baby", i)
        time.Sleep(time.Second)
    }
}

func main() {
    go test()    // 前边加 go 就开启了 协程
    test2()
}

$ go run .\test.go
test2() Hello,baby 0
test() Hello,world 0
test() Hello,world 1
test2() Hello,baby 1
test() Hello,world 2
test2() Hello,baby 2
test() Hello,world 3
test2() Hello,baby 3
test() Hello,world 4
test2() Hello,baby 4
```


**b) 设置运行 cpu 数**  
<https://golang.org/pkg/runtime/#pkg-index> 

1. go1.8 之后程序默认就运行在多核上，不需要配置
  
```
package main
import (
    "fmt"
    "runtime"
)
func main() {
    // 获取当前系统cpu数量
    num := runtime.NumCPU()
    
    // 设置 num-1 的cpu 运行程序
    runtime.GOMAXPROCS(num)
    fmt.Println("num=", num)
}
>>>
num= 8
```

