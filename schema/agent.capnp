@0xdcccaa5d36aa8b70;

using Go = import "/go.capnp";
$Go.package("agent");
$Go.import("github.com/foohq/foojank-proto/go/agent");

using Java = import "/capnp/java.capnp";
$Java.package("io.github.foohq.foojank");
$Java.outerClassname("Agent");

const cmdStartWorkerT  :Text = "FJ.AGENT.%s.CMD.WORKER.%s.START";
const cmdStopWorkerT   :Text = "FJ.AGENT.%s.CMD.WORKER.%s.STOP";
const cmdWriteStdinT   :Text = "FJ.AGENT.%s.CMD.WORKER.%s.STDIN";

const evtStartWorkerT  :Text = "FJ.AGENT.%s.EVT.WORKER.%s.START";
const evtStopWorkerT   :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STOP";
const evtWorkerStatusT :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STATUS";
const evtWorkerStdoutT :Text = "FJ.AGENT.%s.EVT.WORKER.%s.STDOUT";
const evtAgentInfoT    :Text = "FJ.AGENT.%s.EVT.INFO";

# Deprecated: use cmdStartWorkerT
const startWorkerT :Text = "FJ.API.WORKER.START.%s.%s";
# Deprecated: use cmdStopWorkerT
const stopWorkerT :Text = "FJ.API.WORKER.STOP.%s.%s";
# Deprecated: use cmdWriteStdinT
const writeWorkerStdinT :Text = "FJ.API.WORKER.WRITE.STDIN.%s.%s";
# Deprecated: use evtWorkerStdoutT
const writeWorkerStdoutT :Text = "FJ.API.WORKER.WRITE.STDOUT.%s.%s";
# Deprecated: use evtWorkerStatusT
const updateWorkerStatusT :Text = "FJ.API.WORKER.UPDATE.STATUS.%s.%s";
# Deprecated: use evtAgentInfoT
const updateClientInfoT :Text = "FJ.API.CLIENT.UPDATE.INFO.%s";
# Deprecated: use evtStartWorkerT or evtStopWorkerT with cmdSeq correlation
const replyMessageT :Text = "FJ.API.MESSAGE.REPLY.%s.%s";

struct StartWorkerRequest {
    command @0 :Text;
    args @1 :List(Text);
    env @2 :List(Text);
}

struct StartWorkerResponse {
    error @0 :Text;
}

struct StopWorkerRequest {}

struct StopWorkerResponse {
    error @0 :Text;
}

struct UpdateWorkerStatus {
    status @0 :Int64;
}

struct UpdateWorkerStdio {
    data @0 :Data;
}

struct UpdateClientInfo {
    username @0 :Text;
    hostname @1 :Text;
    system @2 :Text;
    address @3 :Text;
}

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
