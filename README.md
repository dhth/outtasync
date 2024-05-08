# outtasync

✨ Overview
---

`outtasync` helps users quickly identify the CloudFormation stacks that have
gone out of sync with the state represented by their counterpart stack files.
This can occur when someone updates a stack but fails to commit the latest stack
file to the codebase. Alternatively, it may happen when a stack is updated on
one deployment environment but not on others. 🤷

<p align="center">
  <img src="./outtasync.gif?raw=true" alt="Usage" />
</p>


💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/outtasync
```

**go**:

```sh
go install github.com/dhth/outtasync@latest
```

🛠️ Pre-requisites
---

- `git` (used to show the diff for out-of-sync changes)
    - `outtasync` doesn't change or override git's pager, so the diff will
        follow your `.gitconfig` settings (if present).

🛠️ Configuration
---

Create a configuration file that looks like the following. By default,
`outtasync` will look for this file at `~/.config/outtasync.yml`.

```yaml
globalRefreshCommand: aws sso login --sso-session sessionname
profiles:
- name: qa
  stacks:
  - name: bingo-service-qa
    local: ~/projects/bingo-service/cloudformation/infrastructure.yml
    region: eu-central-1
    refreshCommand: aws sso login --profile qa1
  - name: papaya-service-qa
    local: ~/projects/papaya-service/cloudformation/service.yml
    region: eu-central-1
  - name: racoon-service-qa
    local: ~/projects/racoon-service/cloudformation/service.yml
    region: eu-central-1
- name: prod
  stacks:
  - name: brb-dll-prod
    local: ~/projects/brb-dll-service/cloudformation/service.yml
    region: eu-central-1
    refreshCommand: aws sso login --profile rgb-prod
  - name: galactus-service-prod
    local: ~/projects/galactus-service/cloudformation/service.yml
    region: eu-central-1
```

`refreshCommand` overrides `globalRefreshCommand` whereever set.

*Note: The `globalRefreshCommand` and `refreshCommand` settings are only needed
if you want to invoke the command that refreshes your AWS credentials via the
TUI directly.*

⚡️ Usage
---

`outtasync` can run in two modes: A TUI mode (ideal for running locally), and a
CLI mode (ideal for running in a CI pipeline). TUI mode is the default.

### TUI Mode

```bash
outtasync
outtasync -config-file /path/to/config.yml
outtasync -profiles qa,prod
outtasync -p '<regex-pattern-for-stack-names>'
outtasync -p '.*(qa|staging)$'
```

### CLI Mode

```bash
outtasync -mode=cli
```

This will print an output like the following to stdout.

```
2 stacks are outtasync:

qa:eu-central-1:bingo-service-qa
prod:eu-central-1:galactus-service-prod
```

### Downloading in a CI pipeline

`outtasync` can be downloaded from Github releases and used as follows:

```bash
curl -s -OL https://github.com/dhth/outtasync/releases/download/v0.3.0/outtasync_v0.3.0_linux_amd64.tar.gz
tar -xzvf outtasync_v0.3.0_linux_amd64.tar.gz
./outtasync -mode=cli
```

TODO
---

- [ ] Add a command to generate a sample config file
- [x] Add CLI mode

Acknowledgements
---

`outtasync` is built using the awesome TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
