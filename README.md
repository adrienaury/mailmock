# Mailmock

[![Go Report Card](https://goreportcard.com/badge/github.com/adrienaury/mailmock)](https://goreportcard.com/report/github.com/adrienaury/mailmock)
[![Github Release Card](https://img.shields.io/github/release/adrienaury/mailmock)](https://github.com/adrienaury/mailmock/releases)
[![codecov](https://codecov.io/gh/adrienaury/mailmock/branch/master/graph/badge.svg)](https://codecov.io/gh/adrienaury/mailmock)

Mailmock is a lightweight SMTP server designed for testing. It exposes a REST API which will enable your CI/CD to check what transaction were made and what was sent to who.

Mailmock is inspired by Mailhog, but doesn't expose a graphical user interface.

## Features

- SMTP Server implementing RFC5321
- HTTP REST API to list transactions and mails the SMTP server handles

## Installation

### Using Docker

```bash
# this will start mailmock, binding the SMTP server on local port 1025 and the HTTP server to local port 1080
docker run -ti --rm -p 1080:80 -p 1025:25 adrienaury/mailmock
```

### Not using Docker

Download the latest version for your OS from the [release page](https://github.com/adrienaury/mailmock/releases).

Then run it:

```bash
mailmock
```

Next releases will come with more configuration options.

## Contribute

Contributions to this project are very welcome.

If you want to contribute, please check CONTRIBUTING.md

### Links

- Issue Tracker: github.com/adrienaury/mailmock/issues
- Source Code: github.com/adrienaury/mailmock

## Support

If you are having issues, please let me know.
I'm Adrien and my mail is adrien.aury@gmail.com

## License

### Main license

The project is licensed under the [GNU GENERAL PUBLIC LICENSE v3](https://www.gnu.org/licenses/gpl-3.0.html).

#### Exception notices

Some files contains a [GPL linking exception](https://en.wikipedia.org/wiki/GPL_linking_exception) to allow linking modules in any project (everything under pkg folder).

### External user librairies

This project uses the following Go librairies :

- github.com/go-chi/chi v4.0.2+incompatible
- github.com/go-chi/render v1.0.1
- github.com/stretchr/testify v1.3.0
- github.com/heptio/workgroup v0.8.0-beta.1
- github.com/sirupsen/logrus v1.4.2
- logur.dev/logur v0.15.0
- logur.dev/adapter/logrus v0.2.0

Here are the required copyright and permission notices :

```text
Copyright (c) 2015-present Peter Kieltyka (https://github.com/pkieltyka), Google Inc.
Copyright (c) 2016-Present https://github.com/go-chi authors
Copyright (c) 2012-2018 Mat Ryer and Tyler Bunnell
Copyright (c) 2019 Márk Sági-Kazár <mark.sagikazar@gmail.com>
Copyright (c) 2014 Simon Eskildsen

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

```text
Copyright © 2017 Heptio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

### Go standard librairies

This project uses some standard librairies of the Go language which is licensed by the following copyright notice, list of conditions and disclaimer :

```text
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
