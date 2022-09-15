<img src="https://codeberg.org/librarian/librarian/raw/branch/main/static/img/librarian.svg" width="96" height="96" />

# librarian
An alternative frontend for LBRY/Odysee. Inspired by [Invidious](https://github.com/iv-org/invidious) and [Libreddit](https://github.com/spikecodes/libreddit).

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
  <img alt="License: AGPLv3+" src="https://shields.io/badge/License-AGPL%20v3-blue.svg">
</a>
<a href="https://matrix.to/#/#librarian:nitro.chat">
  <img alt="Matrix" src="https://img.shields.io/badge/chat-matrix-blue">
</a>

## Features

- Lightweight
- JavaScript not required*
- No ads
- No tracking
- No crypto garbage

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

See what trackers and cookies they use: https://themarkup.org/blacklight.?url=odysee.com

#### Librarian
Librarian itself does not collect any data but instance operators may collect data. You can view a "privacy nutrition label" by clicking on the "Privacy" link at the bottom.

## Instances

Open an issue to have your instance listed here!

### Clearnet

| URL                                                             | Country      | Provider         | Privacy               | Livestreams | Notes |
| :-------------------------------------------------------------- | :----------- | :--------------- | :-------------------- | :---------- | :---- |
| [lbry.bcow.xyz](https://lbry.bcow.xyz) (official)               | 🇨🇦️ CA, 🇳🇱 NL | Fly.io           | ⚠️ Data collected     | ✅️ |  |
| [odysee.076.ne.jp](https://odysee.076.ne.jp)                    | 🇯🇵 JP        | GMOグローバルサイン | ✅ Data not collected | ✅️ | [Edited source code](https://git.076.ne.jp/TechnicalSuwako/Librarian-mod) |
| [librarian.pussthecat.org](https://librarian.pussthecat.org/)   | 🇩🇪 DE        | Hetzner          | ⚠️ Data collected     | ✅️ |  |
| [lbry.projectsegfau.lt](https://lbry.projectsegfau.lt/)         | 🇳🇱 NL        | BuyVM            | ✅ Data not collected | ✅️ |  |
| [librarian.esmailelbob.xyz](https://librarian.esmailelbob.xyz/) | 🇨🇦 CA        | OVH              | ⚠️ Data collected         | ❌️ |  |
| [lbry.vern.cc](https://lbry.vern.cc/)                           | 🇨🇦 CA        | OVHCloud         | ✅ Data not collected | ❌️ | [Edited theme](https://git.vern.cc/root/modifications/src/branch/master/librarian) |
| [lbry.slipfox.xyz](https://lbry.slipfox.xyz)                    | 🇺🇸 US        | Hetzner          | ✅ Data not collected | ❌️ |  |
| [lbry.sytes.net](https://lbry.sytes.net)                        | 🇲🇩 MD        | AlexHost         | ✅ Data not collected | ❌️ |  |
| [lbry.webhop.me](https://lbry.webhop.me)                        | 🇨🇿 CZ        | Smartape         | ✅ Data not collected | ✅️ |  |
| [lbry.mywire.org](https://lbry.mywire.org)                      | 🇷🇺 RU        | justhost.ru      | ✅ Data not collected | ❌️ |  |
| [lbry.ooguy.com](https://lbry.ooguy.com)                        | 🇸🇰 SK        | STARK INDUSTRIES | ✅ Data not collected | ❌️ |  |
 
### Tor

| URL | Privacy               | Live streams | Notes |
| :-- | :-------------------- | :----------- | :---- |
| [librarian.esmail5pdn24shtvieloeedh7ehz3nrwcdivnfhfcedl7gf4kwddhkqd.onion](http://librarian.esmail5pdn24shtvieloeedh7ehz3nrwcdivnfhfcedl7gf4kwddhkqd.onion/) | ⚠️ Data collected | ❌️ | Onion of librarian.esmailelbob.xyz |
| [lbry.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion](http://lbry.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.vern.cc. [Edited theme](https://git.vern.cc/root/modifications/src/branch/master/librarian) |
| [5znbzx2xcymhddzekfjib3isgqq4ilcyxa2bsq6vqmnvbtgu4f776lqd.onion](http://5znbzx2xcymhddzekfjib3isgqq4ilcyxa2bsq6vqmnvbtgu4f776lqd.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.slipfox.xyz |
| [zjmls4urccc5wxc64ilbv4oqarjab6pbanvmhd2qqt6p7fgqrjjnryqd.onion](http://zjmls4urccc5wxc64ilbv4oqarjab6pbanvmhd2qqt6p7fgqrjjnryqd.onion/) | ✅ Data not collected | ✅️ | Onion of lbry.webhop.me |
| [bxewpsswttslepw27w2hhxhlizwm7l7y54x3jw5cfrb64hb6lgc557ad.onion](http://bxewpsswttslepw27w2hhxhlizwm7l7y54x3jw5cfrb64hb6lgc557ad.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.ooguy.com |

### Automatically redirect links

#### LibRedirect
Use [LibRedirect](https://github.com/libredirect/libredirect) to automatically redirect Odysee links to Librarian! This needs to be enabled in settings.
- [Firefox](https://addons.mozilla.org/firefox/addon/libredirect/)
- [Chromium-based browsers (Brave, Google Chrome)](https://github.com/libredirect/libredirect#install-in-chromium-brave-and-chrome)
- [Edge](https://microsoftedge.microsoft.com/addons/detail/libredirect/aodffkeankebfonljgbcfbbaljopcpdb)

#### GreaseMonkey script
There is a script to redirect Odysee links to Librarian.
[https://codeberg.org/zortazert/GreaseMonkey-Redirect/src/branch/main/odysee-to-librarian.user.js](https://codeberg.org/zortazert/GreaseMonkey-Redirect/src/branch/main/odysee-to-librarian.user.js)

## Install
Librarian can run on any platform Go compiles on, memory usage varies on instance usage due to caching.

> Make sure to join our [Matrix chat](https://matrix.to/#/#librarian:nitro.chat) to get notified on updates for Odysee API changes.

### Docker (recommended)
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
You can also edit `docker-compose.yml` if you want to change ports or build from source.

You can now run Librarian. 🎉
```
sudo docker-compose up -d
```

#### Automatic updates
[Watchtower](https://containrrr.dev/watchtower/) can automatically update your Docker containers.

Create a new `docker-compose.yml` file or add the watchtower section to your existing `docker-compose.yml` file.
```yml
version: "3"
services:
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
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

To include version information use:
```
go build -ldflags "-X codeberg.org/librarian/librarian/pages.VersionInfo=$(date '+%Y-%m-%d')-$(git rev-list --abbrev-commit -1 HEAD)"
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
