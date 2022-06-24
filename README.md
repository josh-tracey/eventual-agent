### Eventual Agent

Just a PubSub Websocket Server

The eventual agent allows cluster local apps to subscribe and publish over network.

#### Goals

- Provide a WebSocket server for handling Pubsub for cloudevents
- Simple with small footprint

#### Interface

```ts
interface EventualMessage {
  type: "publish" | "subscribe" | "unsubscribe"
}

// Payload for both subscribing and unsubscribing to/from channels
interface SubscribeMessage extends EventualMessage {
  channels: string | string[]
}

// Format of a publish message payload
interface PublishMessage extends EventualMessage {
  channels: string | string[]
  event: ICloudEvent
}

interface ICloudEvent<T = any> {
  specversion: '1.0'
  type: string // 'com.example.someevent'
  source: string // '/mycontext'
  id: string // uuid
  time: string // '2018-04-05T17:31:00Z'
  datacontenttype: 'application/json'
  data: T
  meta?: ObjectLiteral
  [key: string]: string | number | ObjectLiteral | undefined
}
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
	Id              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	Data            interface{} `json:"data"`
	DataContentType string      `json:"datacontenttype"`
	Time            int         `json:"time"`
	SpecVersion     string      `json:"specversion"`
}


```


CloudEvents.io
https://github.com/cloudevents/spec/blob/v1.0.1/spec.md
