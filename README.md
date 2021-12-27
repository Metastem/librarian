<img src="https://codeberg.org/librarian/librarian/raw/branch/main/templates/static/img/librarian.svg" width="96" height="96" />
# librarian
An alternative frontend for LBRY/Odysee. Inspired by [Invidious](https://github.com/iv-org/invidious).

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
  <img alt="License: AGPLv3+" src="https://shields.io/badge/License-AGPL%20v3+-blue.svg">
</a>
<a href="https://matrix.to/#/#librarian:nitro.chat">
  <img alt="Matrix" src="https://img.shields.io/badge/chat-matrix-blue">
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
| [librarian.bcow.xyz](https://librarian.bcow.xyz) (official) | 🇨🇦 CA |           |
| [lbry.itzzen.net](https://lbry.itzzen.net) | 🇺🇸 US |            |
| [odysee.076.ne.jp](https://odysee.076.ne.jp) ([edited source code](https://git.076.ne.jp/TechnicalSuwako/Librarian-mod)) | 🇯🇵 JP |            |
| [librarian.davidovski.xyz](https://librarian.davidovski.xyz/) | 🇬🇧 UK | |
| [lbry.ix.tc](https://lbry.ix.tc/) | 🇬🇧 UK | |
| [ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion](http://ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion/) |  | |
| [vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion](http://vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion/) |  | |

## Install
Librarian can run on any platform Go compiles on, memory usage varies on instance usage due to caching.

> Librarian is still in beta and changes frequently, building the latest version from source is recommended.

### Requirements
- Go v1.15 or later

### Build from source
> For more detailed instructions, follow the [guide](https://codeberg.org/librarian/librarian/wiki/Setup-guide-%28manual%29).

Clone the repository and `cd` into it.
```
git clone https://codeberg.org/librarian/librarian
cd librarian
```

Build Librarian.
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
go install codeberg.org/librarian/librarian@latest
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
git clone https://codeberg.org/librarian/librarian
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
