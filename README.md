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

The project is licensed under the GNU GENERAL PUBLIC LICENSE v3.

Some files contains a [GPL linking exception](https://en.wikipedia.org/wiki/GPL_linking_exception) to allow linking modules in any project (everything under pkg folder).

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

