## 4.2 核心编程: 管道(channel)

#### 1. 需求背景

需求：使用 协程 方式对20内的数值求幂，并输出         

```
package main
import (
    "fmt"
    "time"
)

var myMap = make(map[int]int, 10)

func test(n int) {
    res := 1
    for i := 1; i <= n; i++ {
        res *= i
    }
    myMap[n] = res     // 报错1：并发写没有加锁导致资源争夺  error: concurrent map writes
}

func main() {
    for i := 1; i <= 20; i++ {
        go test(i)
    }

    time.Sleep(time.Second * 10)   // 报错2： 防止写成还没执行完，主线程执行退出了

    for i, v := range myMap {
        fmt.Printf("map[%d] = %v\n", i, v)
    }
}
```

**a) 解决一：全局变量互斥锁**
<https://golang.org/pkg/sync/>
```
package main
import (
    "fmt"
    "time"
    "sync"
)

var (
    myMap = make(map[int]int, 10)
    // 声明一个全局互斥锁， sync 包内的 Mutex 互斥
    lock sync.Mutex
)
func test(n int) {
    res := 1
    for i := 1; i <= n; i++ {
        res *= i
    }
    lock.Lock()  
    myMap[n] = res     // 报错1：并发写没有加锁导致资源争夺  error: concurrent map writes
    lock.Unlock()
}

func main() {
    for i := 1; i <= 20; i++ {
        go test(i)
    }

    // 休眠10 秒，是为了让协程完成任务，否则主线程很快退出了，myMap 中还没有结果
    time.Sleep(time.Second * 10)   // 报错2：

    lock.Lock()
    for i, v := range myMap {
        fmt.Printf("map[%d] = %v\n", i, v)
    }
    lock.Unlock()
}
```


**b) 解决二：管道channel**
1. 全局变量加锁 可以解决 grouting 通信问题，但不完美，比如主线程等待协程的时间不能精确
2. 全局变量加锁 不利于多个协程对全局变量的读写



#### 2. channel 介绍与使用

1. channel 本质是一个数据结构--队列
2. 数据是先进先出的 (FIFO: first in first out)
3. channel 本身就是线程安全，多 gorouting 访问不需要加锁
4. channel 是有类型的，一个 string 只能存放 string 类型数据

5. channel 是引用类型
6. channel 必须初始化后才能写入数据

```
package main
import (
    "fmt"
)

func main() {
    // 1. 初始化一个管道  int 类型，make 有缓冲得 chan 并且容量为 3
    var intChan chan int 
    intChan = make(chan int, 3)

    fmt.Printf("intChan value is %v, addres is %p\n", intChan, &intChan)

    // 2. 写入数据
    intChan<- 10
    n1 := 20
    intChan<- n1
    intChan<- 30

    fmt.Printf("intChan len=%v  cap=%v \n", len(intChan), cap(intChan))

    // 3. 写入数据不能超过 cap，否则会报错 fatal error: all goroutines are asleep - deadlock!
    // intChan<- 40

    // 4. 取数据  推出数据
    n2 := <-intChan
    <-intChan
    n4 := <-intChan
    fmt.Println("n2=", n2, "n4=", n4)

    // 5. 在没有使用协程情况下，如果管道数据全部取出，再取也会报错  fatal error: all goroutines are asleep - deadlock!
    // n5 := <-intChan
    // fmt.Println("n5=", n5)
}

>>>
intChan value is 0xc000084080, addres is 0xc000006028
intChan len=3  cap=3
n2= 10  n4= 30
```


#### 3. channel 的关闭与遍历

1. 使用 close 函数关闭 channel
2. channel关闭后，不能再写入数据，但仍可以读数据
3. 使用 for-range 进行遍历，不能用 for-len (因为长度会变化)
4. 遍历时，channel没关闭会出现deadlock错误；channel关闭后，会正常遍历完数据并退出

```
package main
import (
    "fmt"
)

func main() {
    // 初始化一个管道
    intChan := make(chan int, 3)
    intChan<- 10
    intChan<- 20
    close(intChan)

    // 1. 关闭管道后，在写入数据会报错，但可以读数据
    // intChan<- 30
    n1 := <-intChan
    fmt.Println("n1=", n1)

    // 2. 使用 for-len 遍历管道，会少值，
    intChan2 := make(chan int, 10)
    for i := 0; i < 10; i++ {
         intChan2<- i
    }

    close(intChan2)
    for i := 0; i < len(intChan2); i++ {
        fmt.Println("intChan2 = ", <-intChan2, "  len = ", len(intChan2))
    } 

    fmt.Println("okokokok")
    // 3. 使用 for-range 遍历管道，写入完需要关闭，要不然取得时候会 fatal error: all goroutines are asleep - deadlock!
    intChan3 := make(chan int, 10)
    for i := 0; i < 10; i++ {
         intChan3<- i
    }

    close(intChan3)
    for v := range intChan3 {
        fmt.Println("intChan3 = ", v)
    }

}

>>>
n1= 10
intChan2 =  0   len =  9
intChan2 =  1   len =  8
intChan2 =  2   len =  7
intChan2 =  3   len =  6
intChan2 =  4   len =  5
okokokok
intChan3 =  0
intChan3 =  1
intChan3 =  2
intChan3 =  3
intChan3 =  4
intChan3 =  5
intChan3 =  6
intChan3 =  7
intChan3 =  8
intChan3 =  9
```

#### 4. channel 使用细节

1. 默认管道双向，可读可写  var ch1 chan int
2. 声明为只读  var ch3 <-chan struct{}
3. 声明为只写  var ch2 chan<- float          
              ch2 = make(chan float, 3)

4. select 可以解决 管道数据阻塞问题 (之前是用 close 函数解决)
5. grouting 中使用 recover ，解决程序中出现一个或几个 panic，导致程序崩溃


```
package main
import (
    "fmt"
    "time"
)

func sayHello(){
    for i := 0; i < 5; i++ {
        fmt.Println("recove hello, world")
    }
}

func test() {
    defer func() {
        err := recover()
        if err != nil {
            fmt.Println("recover test() 协程出现了错误 ", err)
        }
    }()

    var myMap map[int]string
    myMap[0] = "golang"
}

func main() {
    // 4. grouting 中使用 recover ，解决程序中出现一个或几个 panic，导致程序崩溃
    go sayHello()

    go test()

    for i := 0; i < 3; i++ {
        fmt.Println("recover main() ...... ", i)
        time.Sleep(time.Millisecond * 10)
    } 

    // 3. 传统方式需要先 close + for-range 解决阻塞读取，现在用 select 解决阻塞问题
    intChan := make(chan int, 5)
    for i := 0; i < 5; i++ {
        intChan<- i
    }
    strChan := make(chan string, 5)
    for i := 0; i < 5; i++ {
        strChan<- "Hello " + fmt.Sprintf("%d \n", i)
    }

    for {
        select {
        case v := <-intChan :
            fmt.Println("select intChan = ", v)
            // time.Sleep(time.Second)
        case v := <-strChan :
            fmt.Println("select strChan = ", v)
            // time.Sleep(time.Second)
        default:
            fmt.Println("select 取完了，结束......")
            return 
        }
    }
}

>>>
recover main() ......  0
recove hello, world
recove hello, world
recove hello, world
recove hello, world
recove hello, world
recover test() 协程出现了错误  assignment to entry in nil map
recover main() ......  1
recover main() ......  2
select strChan =  Hello 0

select strChan =  Hello 1

select strChan =  Hello 2

select intChan =  0
select strChan =  Hello 3

select strChan =  Hello 4

select intChan =  1
select intChan =  2
select intChan =  3
select intChan =  4
select 取完了，结束......
```

