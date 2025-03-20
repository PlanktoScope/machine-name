# machine-name
A tool for mapping serial numbers to language-localized names

## Introduction

Devices often assign a (pseudo-)unique serial number to each unit. However, each device may also need a (pseudo-) unique human-readable name for easy access, e.g. to network services. This package provides a one-to-one mapping between serial numbers (within a finite range) and Heroku-style names (i.e. two simple words followed by a short number, such as `friendly-cat-4827`). The names are constructed from word lists within a specific locale, so that the language used for the name can be localized. Additionally, the maximum allowed length of a generated name is 24 characters (strictly speaking, 24 bytes); this allows a short 8-character (strictly speaking, 8-byte) prefix to be prepended to the machine name (e.g. `pkscope-friendly-cat-4827`) for use as a Wi-Fi SSID, without violating the maximum allowed length of 32 bytes for a Wi-Fi SSID.

## Usage

### Deployment

First, you will need to download machine-name, which is available as a single self-contained executable file. You should visit this repository's [releases page](https://github.com/PlanktoScope/machine-name/releases/latest) and download an archive file for your platform and CPU architecture; for example, on a Raspberry Pi 4, you should download the archive named `machine-name_{version number}_linux_arm.tar.gz` (where the version number should be substituted). You can extract the machine-name binary from the archive using a command like:
```
tar -xzf machine-name_{version number}_{os}_{cpu architecture}.tar.gz machine-name
```

Then you may need to move the machine-name binary into a directory in your system path, or you can just run the machine-name binary in your current directory (in which case you should replace `machine-name` with `./machine-name` in the commands listed below).

Once you have machine-name, you should run it as follows:
```
machine-name name --format=hex --sn=0xd6b82659
```
but replacing `0xd6b82659` with a 32-bit hex string representing a serial number. The program will then print the machine name corresponding to that serial number. For example, the following serial numbers will result in the following machine names:

| Serial Number | Machine Name        |
|---------------|---------------------|
| `0xdeadc0de`  | metal-slope-23501   |
| `0xd6b82659`  | chain-list-27764    |
| `0x0`         | able-account-0      |
| `0x1`         | small-ball-26954    |
| `0x2`         | safe-minute-6738    |
| `0x3`         | linen-opinion-33692 |
| `0x4`         | cool-pocket-1684    |
| `0x8123`      | clear-field-33719   |

## Licensing

Except where otherwise indicated, source code provided here is covered by the following information:

Copyright Ethan Li and PlanktoScope project contributors

SPDX-License-Identifier: `Apache-2.0 OR BlueOak-1.0.0`

You can use the source code provided here either under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0) or under the [Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0); you get to decide. We are making the software available under the Apache license because it's [OSI-approved](https://writing.kemitchell.com/2019/05/05/Rely-on-OSI.html), but we like the Blue Oak Model License more because it's easier to read and understand.
