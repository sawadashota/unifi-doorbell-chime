UniFi Doorbell Chime
===

Notify to Mac when Doorbell is rung.

Getting Started
---

Put config file `$HOME/.unifi-doorbell-chime/config.yaml` and edit all property of following.

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

```
$ brew install sawadashota/tap/unifi-doorbell-chime
```
