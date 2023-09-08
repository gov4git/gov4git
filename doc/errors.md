# Handling errors

The command-line client reports errors in a standardized way. Errors are strings which can be interpreted and acted upon by the caller.

In the event of an error:
- If the verbose flag `-v` is _not set_, the client will print the error string to stderr and exit with a non-zero exit code.
- If the verbose flag `-v` is _set_, the client will print a stack trace to stderr, followed by the error string, and exit with a non-zero exit code.

## Interpreting and handling errors

### Remote is ahead on push

*Manifestation:* The error string contains the substring `"non-fast-forward update"` 

*Interpretation:* Your local client's cache is stale, and your client has attempted to update a remote which is ahead.

*Incidence:* This error is expected to occur in normal operation at a very low frequency. Known causes:
- Two clients, authenticated as the organizer, race to update the governance repositories. This can happen occasionally when different GitHub integration actions execute at the same time, which can happen due to normal GitHub service latency variations.

*Resolution:* To resolve this error:
1. Refresh the client's cache by running `gov4git cache update`
2. Retry the operation
