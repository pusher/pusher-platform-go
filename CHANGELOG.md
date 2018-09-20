# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/pusher/pusher-platform-go/compare/0.1.2...HEAD)

## [0.1.2]

- Fix token issuer. The secret was being used to construct the issuer instead.

## [0.1.1]

- Fix token expiry time. It was previously a string of the time object, but we require a timestamp.

## [0.1.0]

Initial version
