A Discord bot written to learn async Rust, on [tokio](https://tokio.rs) with poise and serenity.

- `/r6rank` gives members the server role for their Rainbow Six Siege rank
- `/r6verify` looks the rank up instead, through a pluggable `RankSource` trait (the first provider is
  Arenyze), with caching and a cooldown to stay inside a free API quota
- Clean shutdown on signals, typed lookup errors, and a public repo with CI
