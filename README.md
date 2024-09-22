# Cabinet
_a simple and easy to install file server over http_

## Installation

run the following

```
go build && \

# will ask for a passcode to create
# creates a cabinet user and group
# creates a directory at /usr/local/share/Cabinet (owned by cabinet:cabinet)
# copies the cabinet binary to /usr/bin/cabinet
# copies the service file in ./setup/ to /systemd/service/cabinet.service
# enables the cabinet service in systemd
sudo ./cabinet -s
```

### If you want to start it right away
```
sudo systemctl start cabinet
```

Then I suggest you use something like caddy, nginx, or apache to proxy to it. If you want help with that I'll be around.
You can find me at @vvesley.bsky.social

### If you want to watch the logs you can just use journalctl because we have logging set to journal to in the service file
```
sudo journalctl -fu cabinet
```

## How does it work 

It will serve the files in /usr/local/share/Cabinet as long as the client provides the secret. 
If you hit a directory it will serve the links to the media in that directory.

## Why use it? 

* an easy to setup file share server with the ability to upload files  
* It's homegrown and made in my free time and you like homegrown software
* You think the name is cute

Or Don't! I don't care :)

## Contributing

* Fork it
* Fix it
* Fuck it

## License 

see the LICENSE file in the root of this directory
