UniFi Doorbell Chime
===

Notify to Mac when Doorbell is rung.

Getting Started
---

Put config file `$HOME/.unifi-doorbell-chime.yaml` and edit all property of following.

```yaml
unifi:
  ip: "192.168.1.1"
  username: "username"
  password: "password"

message:
  templates:
    - "I'm on my way"
    - "I'm busy now"
```

Then exec command.

```
$ unifi-doorbell-chime
```

Installation
---

### Build From Source

```
$ git pull https://github.com/sawadashota/unifi-doorbell-chime.git
$ make install
$ make build
```

### Via Brew

```
$ brew install sawadashota/unifi-doorbell-chime/unifi-doorbell-chime
```
