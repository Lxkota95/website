A study of binary encoding: equity order messages with FIX semantics, in
[Simple Binary Encoding](https://github.com/aeron-io/simple-binary-encoding) (SBE), in Go.

- An SBE schema compiled to Go codecs, with every byte of the wire format checked by tests
- Measured against JSON: 2–3× faster to encode and decode, but *bigger*, because fixed-size
  fields trade bytes for predictability
- Found that the Go generator's range check rejects SBE's own NUL padding in text fields
- Concluded; a TCP client and server stay ideas
