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

# Notes on modules

Discovered what's wrong with my noobie thought process. The right thing to do is probably to create a *new* repository and simply import the `envoyproxy/go-control-plane` project.

Old notes ( there are some useful bits, about modules )
- Was trying to get `envoyproxy/go-control-plane` to run but it's using an older module management called `glide` ( says on it's repo to move to godep - yet another older module management mechanism ). Could be wrong but I doubt latest version of go works with `glide` ( something to do with suppoprting `vendor` lookup only if it's using go newest module system )

The issue with using go latest module system when the source repository does not support it is that it will report the following error : 

```
    build github.com/envoyproxy/go-control-plane: cannot load envoy/api/v2/core: cannot find module providing package envoy/api/v2/core
```

The `cannot load` phrase seem to indicate that it's reaching out to the repository remotely. So in order to stop it from doing that add the `replace` directive in the `go.mod` file like this 

    `replace github.com/envoyproxy/go-control-plane/envoy => ./envoy`

