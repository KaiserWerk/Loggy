# Loggy
A tiny log collector that writes to disk. Accepts log entries via TCP and UDP.
This is mainly a Dev / DevOps tool.

Are you tired of setting up a logging structure for every new project and log files
cluttering your project folders?
Just send everything to *Loggy* instead! Every message will be prefixed with a nano-precision
timestamp and has a maximum size of 1024 bytes (without timestamp).

### Setup

Download a release for your operating system and place it somewhere on your system
doesn't really matter where.
Make sure the folder supposed to hold your log files is writeable by the user
running *Loggy* and that two ports are open to bind to.

### Usage

By default, log files are places beside the executable.

The default ports are:
- TCP: 7441
- UDP: 7442

Just run the executable, possibly in the background or as a service. If you need a 
different log path or want to bind to different ports, run it with the required parameters:
```cmd
# On Windows:
> .\loggy-vX.X.X-win64.exe --logpath="C:\\Logs" --udp=7442 --tcp=7441
```

```cmd
# On Linux:
> ./loggy-vX.X.X-linux64 --logpath="/var/logs/loggy" --udp=7442 --tcp=7441
```

### Produced Log Files

Log files will be written with the permissions **0644**.
If the current log file reaches the maximum file size of 10 MB, it will be rotated.
Also, no rotated files will currently be removed.
