# Notes

- Provide examples on fan-in channels ( multiple channels streaming to one ), as this is in a loop can't use the select..case structure 

```
		for _, plugin := range plugins {
			waitgroup.Add(1)
			go func(p *Plugin, wg *sync.WaitGroup){
				defer wg.Done() 
				pluginChannel := (*p).run(ctx)
				for {
				  request, more := <-pluginChannel
				  if !more {
					fmt.Printf("[%s] fan in loop terminating, channel closed\n", (*p).getName())  
					return
				  }
				  channel<-request
				} 
			}(&plugin, &waitgroup)		
		}
```

The code above also has example of how waitgroup can be used