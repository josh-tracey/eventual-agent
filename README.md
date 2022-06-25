### Eventual Agent

Golang Pubsub

The eventual agent allows cluster local apps to subscribe and publish over network.

#### Goals

- Provide a WebSocket server for handling Pubsub for cloudevents
- Simple with small footprint

#### SDK

Typescript SDK [reactive-node/eventual-sdk](https://gitlab.com/adriftdev1/reactive-node/-/tree/master/packages/eventual-sdk) Coming soon!
```ts
export declare type dispose = () => void;
export declare type ObjectLiteral<T = any> = {
    [key: string]: T;
};
/**
 * @param specversion '1.0.2'
 * @param type 'com.example.some-event'
 * @param source '/mycontext'
 * @param subject
 * @param id uuid
 * @param time
 * @param datacontenttype 'application/json'
 * @param data ObjectLiteral
 */
export declare type ICloudEvent<T = any> = {
    specversion: '1.0.2';
    type: string;
    source?: string;
    subject?: string;
    id: string;
    time: string;
    datacontenttype: 'application/json';
    data: T;
    meta?: ObjectLiteral;
};
declare type CompleteEvent = {
    reason?: string;
    code?: number;
    wasClean?: boolean;
    target?: WebSocket;
};
export declare type Observable<T> = {
    next: (message: T) => Promise<void>;
    complete: (event: CompleteEvent) => Promise<void>;
    error: (error: Error) => void;
};
export declare type IEventualClientOptions = {
    waitForOpenConnection?: {
        maxNumberOfAttempts?: number;
        attemptInterval?: number;
    };
};
export declare const createEventualClient: (url: string, config?: IEventualClientOptions | undefined) => {
    subscribe: <T>(channels: string | string[], obs: Observable<ICloudEvent<T>>) => dispose;
    publish: <T_1>(channels: string | string[], data: ICloudEvent<T_1>) => Promise<void>;
};
export declare const createCloudEvent: <T = any>(
  type: string, 
  source: string, 
  data: T, 
  subject?: string | undefined, 
  meta?: ObjectLiteral<any> | undefined
) => ICloudEvent<T>;


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
