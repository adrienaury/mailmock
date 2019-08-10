Mailmock
========

Mailmock is a lightweight SMTP server designed for testing. It exposes a REST API which will enable your CI/CD to check what transaction were made with the SMTP server and what was sent to who.

Mailmock is inspired by Mailhog, but doesn't expose a graphical user interface.

Features
--------

- SMTP Server implementing RFC5321
- HTTP REST API to list transactions and mails the SMTP server handles

Installation
------------

For now, just run it:

    mailmock

Next releases will come with Docker images and more configuration options.

Contribute
----------

- Issue Tracker: github.com/adrienaury/mailmock/issues
- Source Code: github.com/adrienaury/mailmock

Support
-------

If you are having issues, please let me know.
I'm Adrien and my mail is adrien.aury@gmail.com

License
-------

### Main license

The project is licensed under the [GNU GENERAL PUBLIC LICENSE v3](https://www.gnu.org/licenses/gpl-3.0.html).

---

### Exception notices

Some files contains a [GPL linking exception](https://en.wikipedia.org/wiki/GPL_linking_exception) to allow linking modules in any project (everything under pkg folder).

---

### External user librairies

This project uses the following Go librairies :
- github.com/go-chi/chi v4.0.2+incompatible
- github.com/go-chi/render v1.0.1
- github.com/stretchr/testify v1.3.0

Here are the required copyright notices and this permission notices :

    Copyright (c) 2015-present Peter Kieltyka (https://github.com/pkieltyka), Google Inc.
    Copyright (c) 2016-Present https://github.com/go-chi authors
    Copyright (c) 2012-2018 Mat Ryer and Tyler Bunnell

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

---

### Go standard librairies

This project uses some standard librairies of the Go language which is licensed by the following copyright notice, list of conditions and disclaimer :

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

