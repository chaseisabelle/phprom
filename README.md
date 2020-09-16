# phprom

_a little tool i build for php apps to work with prometheus metrics_

---

## usage

coming soon - link the php composer package

---

## the api

the api uses a modified resp protocol:
- https://github.com/chaseisabelle/goresp
- https://github.com/chaseisabelle/resphp

each command is delimited by a null byte

### examples

all examples assume phprom is running on `localhost:3333`

- register a counter
    ```
    printf '+M\r\n\0' | nc localhost 3333
    ```  
- fetch the metrics
    ```
    printf '+M\r\n\0' | nc localhost 3333
    ```