# Tracker Server

A simple Go backend for [Tracker](https://github.com/mkellnhofer/tracker-android).

## Downloads

Downloads can be found at [releases](https://github.com/mkellnhofer/tracker-server/releases).
(Binaries are only provided for Linux. For other systems you have to build them yourself. See
[this blog post](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) for how to cross
compile Go code.)

## Installation 

1. Extract the archive with `unzip tracker-1.0.0.zip -d tracker`
2. Go to the directory that was just created
3. Execute `./tracker`

(You don't have to provide a database. At the first start a SQLite database is created which stores
all data.)

## Configuration

The configuration can be changed in file `/config/config.ini`. By default port 8080 and no password
is used. 

Besides setting a password, I would recommend to us a reverse proxy e.g. Nginx which does TLS
offloading. (See
[Nginx documentation](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/) for how to
configure Nginx as a reverse proxy and
[this tutorial](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-18-04)
for how to secure Nginx with a Let's Encrypt certificate.)

## API Documentation

The server provides a REST API which is available under path `/api/v1`. A detailed documentation can
be found [here](API.md).

## Copyright and License

Copyright Matthias Kellnhofer. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in
compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is
distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing permissions and limitations under the
License.