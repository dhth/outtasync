# outtasync

✨ Overview
---

`outtasync` helps users quickly identify the CloudFormation stacks that have
gone out of sync with the state represented by their counterpart stack files.
This can occur when someone updates a stack but fails to commit the latest stack
file to the codebase. Alternatively, it may happen when a stack is updated on
one deployment environment but not on others. 🤷

[![Demo Video](https://tools.dhruvs.space/images/outtasync/outtasync-video-1.png)](https://www.youtube.com/watch?v=BjJcBquIyk8)

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

Or get the binaries directly from a [release][2]. Read more about verifying the
authenticity of released artifacts [here](#-verifying-release-artifacts).

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
    tags:
    - qa
    - auth
  - name: racoon-service-qa
    local: ~/projects/racoon-service/cloudformation/service.yml
    region: eu-central-1
    tags:
    - qa
    - payments
- name: prod
  stacks:
  - name: brb-dll-prod
    local: ~/projects/brb-dll-service/cloudformation/service.yml
    region: eu-central-1
    refreshCommand: aws sso login --profile rgb-prod
    tags:
    - prod
    - integrations
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
outtasync -t 'comma,separated,list,of,tags,to,filter,for'
outtasync -p '<regex-pattern-for-stack-names>'
outtasync -p '.*(qa|staging)$'
outtasync -c # to check status for all stacks on startup
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

⌨️ Keymaps
---

```
↑/k                                  up
↓/j                                  down
→/l/pgdn                             next page
←/h/pgup                             prev page
g/home                               go to start
G/end                                go to end
ctrl+f/enter                         check status
a                                    check status for all
r                                    refresh aws credentials
ctrl+d/v                             show diff
o                                    filter outtasync stacks
i                                    filter in-sync stacks
e                                    filter stacks with errors
q                                    return to previous page/quit
/                                    filter
?                                    show/close help
```

🖥️ Screenshots
---

![Usage-1](https://tools.dhruvs.space/images/outtasync/outtasync-1.png)

![Usage-2](https://tools.dhruvs.space/images/outtasync/outtasync-2.png)

![Usage-3](https://tools.dhruvs.space/images/outtasync/outtasync-3.png)

🔐 Verifying release artifacts
---

In case you get the `outtasync` binary directly from a [release][2], you may
want to verify its authenticity. Checksums are applied to all released
artifacts, and the resulting checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace `x.y.z` in the commands listed below with the version
you want):

1. Download the following files from the release:

   - outtasync_x.y.z_checksums.txt
   - outtasync_x.y.z_checksums.txt.pem
   - outtasync_x.y.z_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob outtasync_x.y.z_checksums.txt \
       --certificate outtasync_x.y.z_checksums.txt.pem \
       --signature outtasync_x.y.z_checksums.txt.sig \
       --certificate-identity-regexp 'https://github\.com/dhth/outtasync/\.github/workflows/.+' \
       --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/outtasync/releases/download/vx.y.z/outtasync_x.y.z_linux_amd64.tar.gz
   sha256sum --ignore-missing -c outtasync_x.y.z_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf outtasync_x.y.z_linux_amd64.tar.gz
   ./outtasync
   # profit!
   ```

TODO
---

- [ ] Add a command to generate a sample config file
- [x] Add CLI mode

Acknowledgements
---

`outtasync` is built using the awesome TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/dhth/outtasync/releases
