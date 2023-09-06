# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
All dates in this file are given in the [UTC time zone](https://en.wikipedia.org/wiki/Coordinated_Universal_Time).

## Unreleased

## 0.1.3 - 2023-09-06

### Changed

- Replaced long words in the en_US.UTF-8 wordlists with shorter words (8-character limit for the first word, 9-character limit for the second word), so that the maximum possible length of a generated name is 24 characters. This change should only affect machine names which used those words; other machine names should be unaffected.

## 0.1.2 - 2023-05-20

### Added

- Embedded wordlists and name-generation functionality are now provided in public packages.

## 0.1.1 - 2023-05-17

### Fixed

- Removed extraneous newline from the stdout output of the command-line tool.

## 0.1.0 - 2023-05-17

### Added

- Programmatic generation of US English word lists.
- Generation of machine names from 32-bit serial numbers, using the US English word lists as the default.
