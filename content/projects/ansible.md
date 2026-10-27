A Rust library for working with [Ansible](https://docs.ansible.com/) inventories: load any inventory
Ansible understands and get typed access to its hosts and their variables.

- **0.3.0** reworked it into idiomatic Rust: a typed error enum instead of panics, a borrowing API, and
  runnable documentation examples
- Released from CI with crates.io *trusted publishing*: merging a version bump publishes the crate,
  with no long-lived API token anywhere
