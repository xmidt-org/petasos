# petasos

[![Build Status](https://travis-ci.org/Comcast/petasos.svg?branch=master)](https://travis-ci.org/Comcast/petasos) 
[![codecov.io](http://codecov.io/github/Comcast/petasos/coverage.svg?branch=master)](http://codecov.io/github/Comcast/petasos?branch=master)
[![Code Climate](https://codeclimate.com/github/Comcast/petasos/badges/gpa.svg)](https://codeclimate.com/github/Comcast/petasos)
[![Issue Count](https://codeclimate.com/github/Comcast/petasos/badges/issue_count.svg)](https://codeclimate.com/github/Comcast/petasos)
[![Go Report Card](https://goreportcard.com/badge/github.com/Comcast/petasos)](https://goreportcard.com/report/github.com/Comcast/petasos)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/Comcast/petasos/blob/master/LICENSE)

The HTTP Redirector Component

"Thanks for calling. I will connect you to the next available handler."

The main package for this application is petasos.

# How to Install

## Centos 6

1. Import the public GPG key (replace `v0.0.1-65alpha` with the release you want)

```
rpm --import https://github.com/Comcast/petasos/releases/download/v0.0.1-65alpha/RPM-GPG-KEY-comcast-webpa
```

2. Install the rpm with yum (so it installs any/all dependencies for you)

```
yum install https://github.com/Comcast/petasos/releases/download/v0.0.1-65alpha/petasos-0.0.1-65.el6.x86_64.rpm
```
