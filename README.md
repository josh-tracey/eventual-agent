### Eventual Agent

Golang Pubsub

The eventual agent allows cluster local apps to subscribe and publish over network.

#### Goals

- Provide a WebSocket server for handling Pubsub for cloudevents
- Simple with small footprint

#### SDK

Typescript SDK [reactive-node/eventual-sdk](https://gitlab.com/adriftdev1/reactive-node/-/tree/master/packages/eventual-sdk) Source Coming soon!

Install with
```sh
yarn add @adriftdev/eventual-sdk
```

SDK Example Usage

```ts
import { createCloudEvent, createEventualClient } from '@adriftdev/eventual-sdk'

const run = async () => {
  const URL = 'ws://localhost:8080'

  const client = createEventualClient(URL)
  const dispose = client.subscribe<{ temp: number }>('tempUpdates', {
    next: async (event) => {
      if (event.subject === 'CpuTemps') {
        console.log(`Cpu Temp: ${event.data.temp}`)
      } else {
        console.log(`Fan Temp: ${event.data.temp}`)
      }
    },
    error: async (error) => {
      console.log(error)
    },
    complete: async (event) => {
      console.log(event.code)
    },
  })
  const iterations = 4

  let iteration = [...Array(iterations).keys()]
  for await (const i of iteration) {
    await client.publish(
      'tempUpdates', // channel / channels
      createCloudEvent(
        'tempUpdates', //type
        'com.adriftdev.events', // source,
        { temp: i + 50 }, //data
        'CpuTemps'
      )
    )
    if (i == iterations - 1) {
      dispose()
    }
  }
}
run()


/* 
#####---OUTPUT---######

Cpu Temp: 50
Cpu Temp: 51
Cpu Temp: 52
Cpu Temp: 53

#######################
*/
```


#### Implementation
Written in golang for simplicity, minimal footprint and faster processing

```go
type PublishEvent struct {
	Type     string        `json:"type"`
	Channels []interface{} `json:"channels"`
	Event    CloudEvent `json:"event"`
}

type SubscribeMessage struct {
	Type     string        `json:"type"`
	Channels []interface{} `json:"channels"`
}

type CloudEvent struct {
	ID              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	Subject         string                 `json:"subject"`
	Data            map[string]interface{} `json:"data"`
	DataContentType string                 `json:"datacontenttype"`
	Time            string                 `json:"time"`
	SpecVersion     string                 `json:"specversion"`
	Meta            map[string]interface{} `json:"meta"`
}



```


CloudEvents.io
https://github.com/cloudevents/spec/blob/v1.0.1/spec.md
