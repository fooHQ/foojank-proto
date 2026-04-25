@0xdcccaa5d36aa8b70;

using Go = import "/go.capnp";
$Go.package("agent");
$Go.import("github.com/foohq/foojank-proto/go/agent");

using Java = import "/capnp/java.capnp";
$Java.package("io.github.foohq.foojank");
$Java.outerClassname("Agent");

# Worker exit codes.
# Exit code is sent by a worker in UpdateWorkerStatus.
const exitSuccess         :Int64 = 0;
const exitFailure         :Int64 = 1;
const exitCommandNotFound :Int64 = 127;
const exitInterrupted     :Int64 = 130;

# cmdStartWorkerT is the NATS topic used to send start worker commands to an agent.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const cmdStartWorkerT  :Text = "FJ.AGENT.%s.CMD.WORKER.%s.START";
# cmdStopWorkerT is the NATS topic used to send stop worker commands to an agent.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const cmdStopWorkerT   :Text = "FJ.AGENT.%s.CMD.WORKER.%s.STOP";
# cmdWriteStdinT is the NATS topic used to send stdin data to a worker process.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const cmdWriteStdinT   :Text = "FJ.AGENT.%s.CMD.WORKER.%s.STDIN";

# evtStartWorkerT is the NATS topic where worker start events are published.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const evtStartWorkerT  :Text = "FJ.AGENT.%s.EVT.WORKER.%s.START";
# evtStopWorkerT is the NATS topic where worker stop events are published.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const evtStopWorkerT   :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STOP";
# evtWorkerStatusT is the NATS topic where worker status updates are published.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const evtWorkerStatusT :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STATUS";
# evtWorkerStdoutT is the NATS topic where worker stdout data is published.
# The "%s" placeholders must be replaced with the agent ID and worker ID, respectively.
const evtWorkerStdoutT :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STDOUT";
# evtAgentInfoT is the NATS topic where agent identification info is published.
# The "%s" placeholder must be replaced with the agent ID.
const evtAgentInfoT    :Text = "FJ.AGENT.%s.EVT.INFO";

# StartWorkerRequest is sent by a client to an agent to start a new worker process.
struct StartWorkerRequest {
    # command is the path to the executable to run.
    command @0 :Text;
    # args is the list of command-line arguments.
    args @1 :List(Text);
    # env is the list of environment variables in the form "KEY=VALUE".
    env @2 :List(Text);
}

# StartWorkerResponse is sent by an agent in response to a StartWorkerRequest.
struct StartWorkerResponse {
    # error contains a description of the error if the worker failed to start, or an empty string on success.
    error @0 :Text;
}

# StopWorkerRequest is sent by a client to an agent to stop a running worker process.
struct StopWorkerRequest {}

# StopWorkerResponse is sent by an agent in response to a StopWorkerRequest.
struct StopWorkerResponse {
    # error contains a description of the error if the worker failed to stop, or an empty string on success.
    error @0 :Text;
}

# UpdateWorkerStatus is sent by an agent to notify the client of a change in worker status.
struct UpdateWorkerStatus {
    # status is the exit code of the worker process, or other status indicators.
    status @0 :Int64;
}

# UpdateWorkerStdio is sent by a client when writing to a worker's stdin and by an agent when streaming stdout.
struct UpdateWorkerStdio {
    # data is the raw data.
    data @0 :Data;
}

# UpdateClientInfo is sent by an agent to identify itself to a client.
struct UpdateClientInfo {
    # username is the operating system username of the agent.
    username @0 :Text;
    # hostname is the name of the machine where the agent is running.
    hostname @1 :Text;
    # system is the operating system name (e.g., linux, darwin, windows).
    system @2 :Text;
    # address is the public address of the agent.
    address @3 :Text;
}

# Message is a top-level container for all protocol messages.
struct Message {
    content :union {
        startWorkerRequest @0 :StartWorkerRequest;
        startWorkerResponse @1 :StartWorkerResponse;
        stopWorkerRequest @2 :StopWorkerRequest;
        stopWorkerResponse @3 :StopWorkerResponse;
        updateWorkerStatus @4 :UpdateWorkerStatus;
        updateWorkerStdio @5 :UpdateWorkerStdio;
        updateClientInfo @6 :UpdateClientInfo;
    }
}
