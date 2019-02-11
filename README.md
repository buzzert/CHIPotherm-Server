# CHIPotherm server

This is the server side component of the [Chipotherm](https://github.com/buzzert/Chipotherm) thermostat. 

It's purpose is to basically act as a command-and-control server for the Chipotherm program.

## Building / Running
Build via `go build`.
Runs on port 43001. 

## API
State format: `[enabled | disabled] [temperature]`

### poll () -> state_format
Blocks until the state on the server changes, in which case responds with either:
    Formatted state (see "State format" above) or
    "refresh": Expects a subsequent call to `updateState(state_format)` notifying the server of the current state

### updateState (state_format) -> ()
Called by the CHIP: Updates the server's cached state.

### refreshState () -> (state as json)
Called by the server: notifies the CHIP that the server requests an updateState. Blocks until we hear back via `updateState`, then returns the new state.

## Client APIs
These APIs are used by the "client", which in this case is the static webpage served by the CHIPotherm Server. 

### getCachedState () -> (state as json)
Immediately returns the server's cached state as JSON

### setState (state_format) -> ()
Called by clients to the server: updates the server's cached state *and* notifies observers (via `poll`) about the new state.

