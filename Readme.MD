# Golang

## Reasons for it 
Much like react in the ui space, a systems language was needed to quickly build highly distributed systems. It is really easy to learn a language like Lua because of it's limited keywords ( 21 ) and concepts. Similarly Go has minimal keywords ( 25 ) compare to C# ( 100 ). 

https://github.com/leighmcculloch/keywords

Go has first class support for concurrency and parallelism, and it goes a step further in addressing the difficulty in writing correct concurrent programs. It does this by favoring slightly higher level primitives like channel and goroutines where message sent to guaranteed to be only accessible by only one goroutine ( not both producer and consumer ) at any one time ( `Communicating Sequential Processes`) over low level system locks ( mutex, semaphore, etc ) and threads. 



# First impression
- A lot like python, with wierder keyword sequence
- Simplicity of coding for concurrency
- OO mainly via extension methods (C#) like
- Oh my gawd - pointers!
