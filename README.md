# About
A web app to help me learn Go & HTMX

## Progress Tracker
- [X] Convert Markdown content to HTML
- [X] Render HTML pages from templates
- [X] Add basic content that describes my engineering projects
- [X] Add code samples for each project
- [ ] Create Dockerfile and image
- [ ] Look into syntax highlighting for code samples
- [ ] Add drop-down sections for each project using HTMX
- [ ] Set up CI/CD deployment
- [ ] Look into Go's `embed` module to include static files in the binary

## Deployment | Fly io
- Fly.io is used to deploy and host this site
- Configuration file `fly.toml` lives in the root of this repo 
- Fly's command line utility `flyctl` is used to orchestrate deployments
  - Deploy with `flyctl deploy`
  - Check status with `flyctl status`
  - Check logs with `flyctl logs`

## Docker
- Build: `sudo docker build -t lxkota-test -f build/Dockerfile .`
- Run: `sudo docker run lxkota-test`