qq机器人

支持指令 
command包下写处理指令的逻辑，一个文件代表一条指令

现有指令
- /算数
 使用说明：@机器人 选择对应指令，后面输入算数表达式（算数表达式要是严格的数学表达式，只支持加减乘除，可以有
  大括号、中括号、小括号，注意使用英文输入），使用前替换代码里的Token为自己的机器人的token,结果统一保留五位小数
  
  ##### 正确的表达式：
   - 1+3+9/3
   - 1+(4*9)+{6/2}
  
  ##### 错误的表达式
    
    - +1-9
    - 6*9
   

扩展说明

```go
switch messType {
    case "MESSAGE_CREATE":
        command.Cluctue(mess)
}
```
在socket对应的事件类型下加入事件处理函数，传递mess消息，由指令函数自己处理