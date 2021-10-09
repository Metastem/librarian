<img src="https://codeberg.org/imabritishcow/librarian/raw/branch/main/templates/static/img/librarian.svg" width="96" height="96" />
# librarian
An alternative frontend for LBRY/Odysee. Inspired by [Invidious](https://github.com/iv-org/invidious).

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
  <img alt="License: AGPLv3+" src="https://shields.io/badge/License-AGPL%20v3+-blue.svg">
</a>
<a href="https://matrix.to/#/#librarian:bcow.xyz">
  <img alt="Matrix" src="https://img.shields.io/matrix/librarian:bcow.xyz?label=Matrix&color=blue&server_fqdn=m.bcow.xyz">
</a>

## Features

### User features
- Lightweight
- No ads
- No tracking
- No crypto garbage

### Technical features
- Copylefted libre software under the AGPL
- No Code of Conduct
- No Contributor License Agreement or Developer Certificate of Origin

## Demo

[Video](https://librarian.bcow.xyz/@MusicARetro:e/Rick+Astley+Never+Gonna+Give+You+Up:4)<br>
[Channel](https://librarian.bcow.xyz/@DistroTube:2)

## Instances

Open an issue to have your instance listed here!

| Website                                                     | Country             | Cloudflare |
| ----------------------------------------------------------- | ------------------- | ---------- |
| [librarian.bcow.xyz](https://librarian.bcow.xyz) (official) | 🇨🇦 CA |            |
| [lbry.itzzen.net](https://lbry.itzzen.net) | 🇺🇸 US |            |
| [odysee.076.ne.jp](https://odysee.076.ne.jp) ([edited source code](https://git.076.ne.jp/TechnicalSuwako/Librarian-mod)) | 🇯🇵 JP |            |
| [ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion](http://ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion/) |  | |
| [vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion](http://vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion/) |  | |

## Install
Librarian can run on any platform Go compiles on, memory usage varies on instance usage due to caching.

### Requirements
- Go v1.15 or later
- libvips

### Build from source
Clone the repository and `cd` into it.
```
git clone https://codeberg.org/imabritishcow/librarian
cd librarian
```

Build Librarian, make sure `libvips` is installed.
```
go build .
```

Edit the config file using your preferred editor.
```
cp config.example.yml config.yml
nvim config.yml
```

You can now run Librarian.
```
./librarian
```

### `go install`
You can install Librarian using Go.
```
go install codeberg.org/imabritishcow/librarian@latest
```

Edit the config file using your preferred editor.
```
cp config.example.yml config.yml
nvim config.yml
```

You can now run Librarian.
```
librarian # If GOBIN is in your PATH
$HOME/go/bin/librarian # If GOBIN is not in PATH
```

### Docker
Install Docker and docker-compose, then clone this repository.
```
git clone https://codeberg.org/imabritishcow/librarian
cd librarian
```

Edit the config file using your preferred editor.
```
cp config.example.yml config.yml
nvim config.yml
```

You can now run Librarian.
```
sudo docker-compose up -d
```
