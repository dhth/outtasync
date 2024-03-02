# outtasync

✨ Overview
---

`outtasync` helps users quickly identify the CloudFormation stacks that have
gone out of sync from the state represented by their counterpart stack files.
This can occur when someone updates a stack but fails to commit the latest stack
file to the codebase. Alternatively, it may happen when a stack is updated on
one deployment environment but not on others. 🤷

<p align="center">
  <img src="./outtasync.gif?raw=true" alt="Usage" />
</p>

🛠️ Pre-requisites
---

- `git` (used to show the diff for out-of-sync changes)
    - `outtasync` doesn't change or override git's pager, so the diff will
        follow your `.gitconfig` settings (if present).

⚡️ Usage
---

1. Create a configuration file that looks like the following.

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
        local: ~/projects/brd-dll-service/cloudformation/service.yml
        region: eu-central-1
        refreshCommand: aws sso login --profile rgb-prod
      - name: galactus-service-prod
        local: ~/projects/galactus-service/cloudformation/service.yml
        region: eu-central-1
    ```

2. Place this file at `~/.config/outtasync.yml` *(optional)*

3. Install `outtasync` by running `go install github.com/dhth/outtasync@latest`

4. Run the TUI as follows:

```bash
outtasync
# or
outtasync --config-file /path/to/config.yml
# or
outtasync -c /path/to/config.yml
```

5. Press `?` to view keyboard shortcuts to use the TUI.

Acknowledgements
---

`outtasync` is built using the awesome TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
