@0xd330108b08fd90db;

using Go = import "/go.capnp";
$Go.package("gateway");
$Go.import("github.com/foohq/foojank-proto/go/gateway");

using Java = import "/capnp/java.capnp";
$Java.package("io.github.foohq.foojank");
$Java.outerClassname("Gateway");

# cmdRegisterT is a message subject used to send agent registration commands to a gateway.
# The "%s" placeholders must be replaced with the gateway ID and agent ID, respectively.
const cmdRegisterT   :Text = "FJ.GATEWAY.%s.CMD.AGENT.%s.REGISTER";
# cmdUnregisterT is a message subject used to send agent unregistration commands to a gateway.
# The "%s" placeholders must be replaced with the gateway ID and agent ID, respectively.
const cmdUnregisterT :Text = "FJ.GATEWAY.%s.CMD.AGENT.%s.UNREGISTER";

# evtRegisterT is a message subject where agent registration events are published.
# The "%s" placeholders must be replaced with the gateway ID and agent ID, respectively.
const evtRegisterT   :Text = "FJ.GATEWAY.%s.EVT.AGENT.%s.REGISTER";
# evtUnregisterT is a message subject where agent unregistration events are published.
# The "%s" placeholders must be replaced with the gateway ID and agent ID, respectively.
const evtUnregisterT :Text = "FJ.GATEWAY.%s.EVT.AGENT.%s.UNREGISTER";

# Property is a key-value pair used to describe agent registration metadata.
struct Property {
    # key is the property name.
    key   @0 :Text;
    # value is the property value.
    value @1 :Text;
}

# RegisterAgentRequest is sent by a client to register an agent with a gateway.
struct RegisterAgentRequest {
    # properties is the list of agent properties submitted during registration.
    properties @0 :List(Property);
}

# RegisterAgentResponse is sent by a gateway in response to a RegisterAgentRequest.
struct RegisterAgentResponse {
    # properties is the list of gateway-assigned properties returned to the agent.
    properties @0 :List(Property);
    # error contains details when registration fails, or a success code otherwise.
    error      @1 :Error;
}

# UnregisterAgentRequest is sent by a client to unregister an agent from a gateway.
struct UnregisterAgentRequest {
    # properties is the list of agent properties submitted during unregistration.
    properties @0 :List(Property);
}

# UnregisterAgentResponse is sent by a gateway in response to an UnregisterAgentRequest.
struct UnregisterAgentResponse {
    # properties is the list of gateway properties returned to the agent.
    properties @0 :List(Property);
    # error contains details when unregistration fails, or a success code otherwise.
    error      @1 :Error;
}

# Error describes the result of a gateway operation.
struct Error {
    # code is the operation result code.
    code    @0 :Int32;
    # message is the operation result message.
    message @1 :Text;
}

# Message is a container for all protocol messages.
struct Message {
    content :union {
        registerAgentRequest @0 :RegisterAgentRequest;
        registerAgentResponse @1 :RegisterAgentResponse;
        unregisterAgentRequest @2 :UnregisterAgentRequest;
        unregisterAgentResponse @3 :UnregisterAgentResponse;
    }
}

# Envelope is a top-level container for all protocol messages.
struct Envelope {
    subject @0 :Text;
    payload @1 :Message;
}
