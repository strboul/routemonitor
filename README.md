# routemonitor

Minimalist route monitoring

## Reasoning

Why does this project exist?

- If you've got lots of VPNs on your computer for different servers, and you
  turn them on and off a bunch, it's important to make sure the routes are set
  up right.

- Checking the routing table doesn't work well because you're always changing
  connections. It just creates too much unnecessary info.

- You don't want to install anything big, e.g.
  [Prometheus](https://github.com/prometheus) which does a lot, but you just
  want simple alerts for your routes without any extra weight.

- Because it's fun!

## Usage

Save a config in e.g. `~/.config/routemonitor/config.yml`.

<!-- keep it mirrored with example/config.yml -->
```yml
fail_fast: true
route:
  - name: Single address
    ip: 10.10.10.24
    expect:
      - when:
          device: eth0
          source: 172.17.0.2

  - name: Block addresses
    ip: 10.10.10.0/30
    expect:
      - when:
          device: xeth1
          source: 10.10.10.25
      - when:
          device: eth0 # fallback if first device doesn't exist
```

Then run `routemonitor`:

```sh
routemonitor -config=~/.config/routemonitor/config.yml -verbose -json
# {"time":"...","msg":"checking route","name":"Single address","ip":"10.10.10.24/32","num IPs":1,"expects":[{"When":{"Device":"eth0","Gateway":"","Source":"172.17.0.2"}}]}
# {"time":"...","msg":"checking route","name":"Block addresses","ip":"10.10.10.0/30","num IPs":4,"expects":[{"When":{"Device":"xeth1","Gateway":"","Source":"10.10.10.25"}},{"When":{"Device":"eth0","Gateway":"","Source":""}}]}
# All routes are as expected.
```

In case there's a mismatch:

```sh
# ... ERROR name="Block addresses" ip="10.10.10.0/30" interface not exist expect="xeth1"
# exit status 1
```

### Periodically run

Create a service job on a repeating schedule e.g. every 5 minutes runs with
[systemd.timer](https://www.freedesktop.org/software/systemd/man/latest/systemd.timer.html)
in Linux:

```sh
#!/bin/bash
set -o pipefail
routemonitor -config=~/.config/routemonitor/config.yml 2>&1 | tee /tmp/routemonitor.log  || \
  (notify-send 'routemonitor failed' \"$(cat /tmp/routemonitor.log)\" && exit 1)"
```

## Installation

```sh
go install github.com/strboul/routemonitor@latest
```

## Roadmap

- [ ] Support IPv6

- [ ] Add DNS resolver check

- [ ] Test/support other systems e.g. \*BSD, Darwin, Windows (`go tool dist list`)
