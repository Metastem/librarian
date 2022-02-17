<img src="https://codeberg.org/librarian/librarian/raw/branch/main/static/img/librarian.svg" width="96" height="96" />
# librarian
An alternative frontend for LBRY/Odysee. Inspired by [Invidious](https://github.com/iv-org/invidious) and [Libreddit](https://github.com/spikecodes/libreddit).

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
  <img alt="License: AGPLv3+" src="https://shields.io/badge/License-AGPL%20v3-blue.svg">
</a>
<a href="https://matrix.to/#/#librarian:nitro.chat">
  <img alt="Matrix" src="https://img.shields.io/badge/chat-matrix-blue">
</a>
<a href="https://gitlab.com/overtime-zone-wildfowl/librarian">
  <img alt="CI" src="https://gitlab.com/overtime-zone-wildfowl/librarian/badges/main/pipeline.svg">
</a>

## Features

### User features
- Lightweight
- JavaScript not required*
- No ads
- No tracking
- No crypto garbage

### Technical features
- Copylefted libre software under the AGPL
- No Code of Conduct
- No Contributor License Agreement or Developer Certificate of Origin

\* JavaScript is required to play livestreams except on Apple devices.

## Demo

[Video](https://lbry.bcow.xyz/@RetroMusic:d/1987-Rick-Astley-Never-Gonna-Give-You-Up-1920x1080:f)<br>
[Channel](https://lbry.bcow.xyz/@DistroTube:2)

## Comparison
Comparing Librarian to Odysee. 

### Speed
Tested using [Google PageSpeed Insights](https://pagespeed.web.dev/).

|             | [Librarian](https://pagespeed.web.dev/report?url=https%3A%2F%2Flbry.bcow.xyz%2F) | [Odysee](https://pagespeed.web.dev/report?url=https%3A%2F%2Fodysee.com%2F) |
| ----------- | --------- | ------ |
| Performance | 99 | 27 |
| Request count | 17 | 470 |
| Resource Size | 702 KiB | 2,457 KiB |
| Time to Interactive | 0.9s | 18.4s |

### Privacy

#### Odysee
<a href="https://tosdr.org/en/service/2391">
  <img alt="Odysee Privacy Grade" src="https://shields.tosdr.org/en_2391.svg">
</a>

Odysee has admitted to using browser fingerprinting for ads and loads multiple ads, trackers, and an annoying cookie banner.

> We and our partners process data to provide:
Use precise geolocation data. Actively scan device characteristics for identification. Store and/or access information on a device. Personalised ads and content, ad and content measurement, audience insights and product development.

They also use your data for these purposes and you cannot opt-out of it.
- Ensure security, prevent fraud, and debug
- Technically deliver ads or content
- Match and combine offline data sources
- Link different devices
- Receive and use automatically-sent device characteristics for identification

**Ads/trackers:** (as of Feb 1, 2022)
- Google
- Vidcrunch
- and many more listed on the list of partners in the cookie banner.

And they have previously used:
- Traffic Junky (P***Hub)
- Unruly Media

#### Librarian
Privacy varies by instance. You can view a "privacy nutrition label" by clicking on the "Privacy" link at the bottom. The official [lbry.bcow.xyz](https://lbry.bcow.xyz/privacy) instance does not collect any data.

## Instances

Open an issue to have your instance listed here!

| Website                                                     | Country             | Cloudflare |
| ----------------------------------------------------------- | ------------------- | ---------- |
| [lbry.bcow.xyz](https://lbry.bcow.xyz) (official) | 🇨🇦 CA |           |
| [lbry.itzzen.net](https://lbry.itzzen.net) | 🇺🇸 US |            |
| [odysee.076.ne.jp](https://odysee.076.ne.jp) ([edited source code](https://git.076.ne.jp/TechnicalSuwako/Librarian-mod)) | 🇯🇵 JP |            |
| [librarian.davidovski.xyz](https://librarian.davidovski.xyz/) | 🇬🇧 UK | |
| [lbry.ix.tc](https://lbry.ix.tc/) | 🇬🇧 UK | |
| [librarian.pussthecat.org](https://librarian.pussthecat.org/) | 🇩🇪 DE | |
| [lbry.mutahar.rocks](https://lbry.mutahar.rocks/) | 🇫🇷 FR | |
| [librarian.esmailelbob.xyz](https://librarian.esmailelbob.xyz/) | 🇨🇦 CA | |

| [ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion](http://ecc5mi5ncdw6mxhjz6re6g2uevtpbzxjvxgrxia2gyvrlnil3srbnhyd.onion/) |  | |
| [vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion](http://vrmbc4brkgkaysmi3fenbzkayobxjh24slmhtocambn3ewe62iuqt3yd.onion/) |  | |

## Install
Librarian can run on any platform Go compiles on, memory usage varies on instance usage due to caching.

### Docker (recommeded)
Install Docker and docker-compose, then clone this repository.
```
git clone https://codeberg.org/librarian/librarian
cd librarian
```

Edit the config file using your preferred editor.
```
mkdir data
cp config.example.yml data/config.yml
nvim data/config.yml
```
You can also edit `docker-compose.yml` if you want to change ports or use the image instead of building it.

You can now run Librarian. 🎉
```
sudo docker-compose up -d
```

### Build from source
> For more detailed instructions, follow the [guide](https://codeberg.org/librarian/librarian/wiki/Setup-guide-%28manual%29).

#### Requirements
- Go v1.16 or later

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

You can now run Librarian. 🎉
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

You can now run Librarian. 🎉
```
librarian # If GOBIN is in your PATH
$HOME/go/bin/librarian # If GOBIN is not in PATH
```
