# Fix the macos TouchBar

Many Macbook Pro models from 2016-2020 have a "TouchBar". A lot of these start to
malfunction after 2 or 3 years. Sometimes the hardware locks up, and when the
macbook wakes up from sleep and fails to re-initialize the touchbar the
system simply panics and reboots.

This utility hooks into the "sleep" and "wake" notifications that macos sends,
and on wakeup it simply kills the `TouchBarServer` process. That's often
enough to fix hanging TouchBar hardware - unfortunately not always.

The `fix-touchbar` utility should be run with setuid-root priviliges. Unfortunately,
macos doesn't allow setuid executabled anymore, so that doesn't work :(

For now, it has to be started manually after login with

```
sudo ./fix-touchbar -bg
```

## TODO

Run a seperate daemon as root which is responsible for killing the TouchBarServer
process. Maybe an
[XPC service](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingXPCServices.html)

