# Foojank Protocol Schemas and Bindings

This repository contains the [Cap'n Proto](https://capnproto.org/) schema definitions and generated bindings for the `foojank` protocol.

**Project Structure**

- `go/`: Go bindings.
- `java/`: Java bindings.
- `rust/`: Rust bindings.
- `cpp/`: C++ bindings.

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

### Rust

To use the Rust bindings in your project, add the `foojank` crate as a dependency in your `Cargo.toml`:

```toml
[dependencies]
foojank = { git = "https://github.com/foohq/foojank-proto", path = "rust" }
capnp = "0.25"
```

Example usage:

```rust
use foojank::agent_capnp::start_worker_request;
use foojank::agent_capnp::start_worker_response;
use capnp::message::ReaderOptions;
use capnp::serialize_packed;

fn main() -> capnp::Result<()> {
    // Assuming 'input' is a Read implementation
    let message_reader = serialize_packed::read_message(&mut input, ReaderOptions::new())?;
    let req = message_reader.get_root::<start_worker_request::Reader>()?;

    println!("Command: {}", req.get_command()?.to_str()?);
    let args = req.get_args()?;
    for i in 0..args.len() {
        println!("Arg: {}", args.get(i)?.to_str()?);
    }

    // Example usage (encoding response):
    let mut message = capnp::message::Builder::new_default();
    {
        let mut resp = message.init_root::<start_worker_response::Builder>();
        resp.set_error("Optional error message");
    }

    // Assuming 'output' is a Write implementation
    serialize_packed::write_message(&mut output, &message)?;
    Ok(())
}
```

### C++

The C++ bindings are generated in `cpp/agent.capnp.h` and `cpp/agent.capnp.c++`.

Example usage:

```cpp
#include <iostream>
#include <capnp/message.h>
#include <capnp/serialize-packed.h>
#include "agent.capnp.h"

int main() {
    // Assuming 'fd' is an open file descriptor
    capnp::PackedFdMessageReader message(fd);
    auto req = message.getRoot<StartWorkerRequest>();

    std::cout << "Command: " << req.getCommand().cStr() << std::endl;
    auto args = req.getArgs();
    for (auto arg : args) {
        std::cout << "Arg: " << arg.cStr() << std::endl;
    }

    // Example usage (encoding response):
    capnp::MallocMessageBuilder respMessage;
    auto resp = respMessage.initRoot<StartWorkerResponse>();
    resp.setError("Optional error message");

    // Assuming 'output' is a file descriptor
    capnp::writePackedMessageToFd(output, respMessage);

    return 0;
}
```

## Schema Overview

The `agent.capnp` schema defines the communication protocol for agent management, including:

- Starting/Stopping workers.
- Handling worker `stdin` and `stdout` streams.
- Updating worker status and client information.
