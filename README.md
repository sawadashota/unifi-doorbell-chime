UniFi Doorbell Chime
===

Notify to PC when Doorbell is rung.

Getting Started
---

```
$ unifi-doorbell-chime init
```

Then, please edit config file.

Then exec command.

```
$ unifi-doorbell-chime start
```

Daemonize
---

### Mac

```
$ brew services start sawadashota/tap/unifi-doorbell-chime
```

### Linux

Configure for systemctl

```
$ sudo vi /etc/systemd/system/unifi-doorbell-chime.service
``` 

```
[Unit]
Description=unifi-doorbell-chime
After=network.service

[Service]
Environment=DISPLAY=:0
ExecStart=/path/to/unifi-doorbell-chime start
Restart=always
Type=simple
User=<USER>
Group=<GROUP>


[Install]
WantedBy=multi-user.target
```

```
$ sudo systemctl enable --now unifi-doorbell-chime.service
``` 


Installation
---

```
$ brew install sawadashota/tap/unifi-doorbell-chime
```
