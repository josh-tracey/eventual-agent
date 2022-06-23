### Eventual Agent

Golang Pubsub

The eventual agent allows cluster local apps to subscribe and publish over network.

#### Goals

- Provide a WebSocket server for handling Pubsub for cloudevents
- Simple with small footprint

#### Interface

Typescript SDK [reactive-node/eventual-sdk](https://gitlab.com/adriftdev1/reactive-node/-/tree/master/packages/eventual-sdk) Coming soon!
```ts

type Observable<T> = {
  next: (message: T) => Promise<void>
  complete: (event: CompleteEvent) => Promise<void>
  error: (error: Error) => void
}

createEventualClient = (
  url: string,
  config?: IEventualClientOptions
):  {
    publish: async <T>(channels: string | string[], data: ICloudEvent<T>), 
    subscribe: <T>( channels: string | string[], obs: Observable<ICloudEvent<T>> ): dispose 
}

createCloudEvent = <T = any>(
  type: string,
  source: string,
  data: T,
  subject?: string,
  meta?: ObjectLiteral
): ICloudEvent<T> 

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
