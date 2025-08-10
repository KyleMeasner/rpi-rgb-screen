# rpi-rgb-screen

## Startup Script

1. Follow the [Github documentation](https://docs.github.com/en/authentication/connecting-to-github-with-ssh) to setup SSH keys on the Raspberry Pi and create an SSH key on your Github account.
2. Clone the repo somewhere on the Raspberry Pi.
3. Update `LOCAL_REPO_PATH` in `startup.sh` to point to the location you cloned the repo.
4. Add a crontab entry to run the script on startup.

```
crontab -e
```

With logging:
```
@reboot sleep 10 && rm -f /path/to/store/logs/logs.txt && /path/to/repo/rpi-rgb-screen/startup.sh >> /path/to/store/logs/logs.txt 2>&1
```

Without logging:
```
@reboot sleep 10 && /path/to/repo/rpi-rgb-screen/startup.sh
```

## Hard-coded Dependencies

1. The dependency `github.com/KyleMeasner/go-rpi-rgb-led-matrix` assumes that https://github.com/hzeller/rpi-rgb-led-matrix has been pulled down and built at `/home/kyle/rpi-rgb-led-matrix`.
