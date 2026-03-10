# Foojank Protocol Schemas and Bindings

This repository contains the [Cap'n Proto](https://capnproto.org/) schema definitions and generated bindings for the `foojank` protocol.

**Project Structure**

- `go/`: Generated Go bindings.
- `java/`: Generated Java bindings.

## Generate Bindings

To maintain a consistent environment, this project uses **Devbox**.

1. [Install Devbox](https://www.jetify.com/docs/devbox/installing-devbox).
2. Install missing dependencies: `devbox install`
3. Generate the bindings for all supported languages: `devbox run build`.

## Using the Bindings

### Go

To use the Go bindings in your project, add the module to your `go.mod`:

```bash
go get github.com/foohq/foojank-proto/go/agent
```

Example usage:

```go
import (
    "fmt"
    "github.com/foohq/foojank-proto/go/agent"
    "capnproto.org/go/capnp/v3"
)

func main() {
    // Assuming 'data' contains the raw Cap'n Proto message
    msg, _ := capnp.Unmarshal(data)
    req, _ := agent.ReadRootStartWorkerRequest(msg)

    cmd, _ := req.Command()
    args, _ := req.Args()

    fmt.Printf("Command: %s\n", cmd)
    for i := 0; i < args.Len(); i++ {
        arg, _ := args.At(i)
        fmt.Printf("Arg: %s\n", arg)
    }

    // Example usage (encoding response):
    respMsg, respSeg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
    resp, _ := agent.NewRootStartWorkerResponse(respSeg)
    resp.SetError("Optional error message")

    // Assuming 'conn' is an io.Writer
    capnp.NewEncoder(conn).Encode(respMsg)
}
```

### Java

The Java bindings are generated in `io.github.foohq.foojank.agent`.

Example usage:

```java
import io.github.foohq.foojank.Agent;
import org.capnproto.MessageReader;
import org.capnproto.Serialize;
import org.capnproto.TextList;

public class Main {
    public static void main(String[] args) throws java.io.IOException {
        // Assuming 'input' is a ReadableByteChannel or similar
        MessageReader message = Serialize.read(input);
        Agent.StartWorkerRequest.Reader reader = message.getRoot(Agent.StartWorkerRequest.factory);

        System.out.println("Command: " + reader.getCommand());
        TextList.Reader reqArgs = reader.getArgs();
        for (int i = 0; i < reqArgs.size(); i++) {
            System.out.println("Arg: " + reqArgs.get(i));
        }

        // Example usage (encoding response):
        org.capnproto.MessageBuilder respMessage = new org.capnproto.MessageBuilder();
        Agent.StartWorkerResponse.Builder respBuilder = respMessage.initRoot(Agent.StartWorkerResponse.factory);
        respBuilder.setError("Optional error message");

        // Assuming 'output' is a WritableByteChannel
        org.capnproto.Serialize.write(output, respMessage);
    }
}
```

## Schema Overview

The `agent.capnp` schema defines the communication protocol for agent management, including:

- Starting/Stopping workers.
- Handling worker `stdin` and `stdout` streams.
- Updating worker status and client information.
