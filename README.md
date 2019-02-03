# CHIPotherm server

## API
State format: `[enabled | disabled] [temperature]`

### poll () -> state_format
Blocks until the state on the server changes, in which case responds with either:
    Formatted state (see "State format" above) or
    "refresh": Expects a subsequent call to `updateState(state_format)` notifying the server of the current state


### updateState (state_format) -> ()
Called by the CHIP: Updates the server's cached state.

### refreshState () -> ()
Called by the server: notifies the CHIP that the server requests an updateState.

### getCachedState () -> state_format
Immediately returns the server's cached state as state_format

### setState (state_format) -> ()
Called by clients to the server: updates the server's cached state *and* notifies observers (via `poll`) about the new state.

