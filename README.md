# About
A web app to help me learn Go & HTMX

## Progress Tracker
- [X] Convert Markdown content to HTML
- [X] Render HTML pages from templates
- [X] Add basic content that describes my engineering projects
- [X] Add code samples for each project
- [X] Create Dockerfile and image
- [X] Set up CI/CD deployment
- [ ] Look into syntax highlighting for code samples
- [ ] Add drop-down sections for each project using HTMX
- [ ] Look into Go's `embed` module to include static files in the binary

## Development
- The `air` framework is used for local development & testing
- Build and execute with `air -c configs/.air.toml`

## Docker
- Configuration file: `build/Dockerfile`
- Build: `sudo docker build -t lxkota-test -f build/Dockerfile .`
- Run: `sudo docker run lxkota-test`
- Enter: `sudo docker exec -it <CONTAINER ID> /bin/bash`
- Stop: `sudo docket stop <CONTAINER ID>`

## Deployment
- Fly.io is used to deploy and host this site
- Configuration file: `configs/fly.toml`
- Fly's command line utility `flyctl` is used to orchestrate deployments:
  - Deploy with `flyctl deploy --config configs/fly.toml`
  - Check status with `flyctl status --config configs/fly.toml`
  - Check logs with `flyctl logs --config configs/fly.toml`
- For SSH:
  - Create token: `flyctl tokens create ssh -a lxkota > lxkota.token.ssh`
  - SSH: `FLY_API_TOKEN=$(cat lxkota.token.ssh) flyctl ssh console -a lxkota`