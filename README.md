# goqueue

## TaskQueue

一个可重试任务队列, 队列存储调度系统中的各种操作任务.

队列本身并不会在后台消费任务, 而是通过`Pop`接口将任务暴露给调用方.

Pop 接口在调用时, 支持传入处理函数, 那么认为任务的消费是投过 Pop 函数.

- TaskQueu
    - Add
    - Pop
    - Close

### Task

```go
type Task struct {
    
}
```

### Add

每个任务有一个`druation 参数`, 表示任务的延迟时间.
