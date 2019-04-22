
An attempt to write my first go program to help cement understanding of the language ( understanding of pointers and pass by value is still a little sketchy ). 

Concurrency is a difficult consideration :
- noticed that sequential thinking some times gets in the way

        ```
        nextGoRoutineInputChannel <- output
		go nextGoRoutine()
		```

    Sequential thinking when chainning channels and goroutines conditioned me to code such that the input values for the next goroutine is pushed into the channel first before a call to the next goroutine, this is an issue as the program will just block without launching off the next goroutine. The right sequence to follow is to launch the next goroutine first and then only followed by sending the input values into the channel

        ```
        go nextGoRoutine() 
        nextGoRoutineInputChannel <- output
		```

- It's possible resources are closed out before your goroutine manage to start using it
    ```
        http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		    go middleware.Run(&w, request)
        })
        http.ListenAndServe(":8080", nil)
    ```

    In the case above `w` has been closed or disposed of after the end of the HandleFunc, i.e. right after launching the middleware.Run goroutine. In order to prevent this from happening, we need to stop the function from exiting before the middleware.Run has completed. One solution is
    to block by reading a done channel

    ```
        http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		    go middleware.Run(&w, request)
            <-middleware.GetDoneChannel()
        })
        http.ListenAndServe(":8080", nil)
    ```