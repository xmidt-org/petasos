# petasos

[![Build Status](https://travis-ci.org/xmidt-org/petasos.svg?branch=master)](https://travis-ci.org/xmidt-org/petasos)
[![codecov.io](http://codecov.io/github/xmidt-org/petasos/coverage.svg?branch=master)](http://codecov.io/github/xmidt-org/petasos?branch=master)
[![Code Climate](https://codeclimate.com/github/xmidt-org/petasos/badges/gpa.svg)](https://codeclimate.com/github/xmidt-org/petasos)
[![Issue Count](https://codeclimate.com/github/xmidt-org/petasos/badges/issue_count.svg)](https://codeclimate.com/github/xmidt-org/petasos)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/petasos)](https://goreportcard.com/report/github.com/xmidt-org/petasos)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/petasos/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/xmidt-org/petasos.svg)](CHANGELOG.md)

The HTTP Redirector Component

"Thanks for calling. I will connect you to the next available handler."

The main package for this application is petasos.

# How to Install

## Centos 6

1. Import the public GPG key (replace `0.0.1-65` with the release you want)

```
rpm --import https://github.com/xmidt-org/petasos/releases/download/0.0.1-65/RPM-GPG-KEY-comcast-xmidt
```

2. Install the rpm with yum (so it installs any/all dependencies for you)

```
yum install https://github.com/xmidt-org/petasos/releases/download/0.0.1-65/petasos-0.0.1-65.el6.x86_64.rpm
```
