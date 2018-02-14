# go-home

[![Build Status](https://travis-ci.org/galactic-filament/go-home.svg?branch=master)](https://travis-ci.org/galactic-filament/go-home)
[![Go Report Card](https://goreportcard.com/badge/github.com/galactic-filament/go-home)](https://goreportcard.com/report/github.com/galactic-filament/go-home)
[![Coverage Status](https://coveralls.io/repos/github/galactic-filament/go-home/badge.svg?branch=)](https://coveralls.io/github/galactic-filament/go-home?branch=)

## Libraries

Kind | Name
--- | ---
Web Framework | [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux)
SQL ORM | [SQLX](http://jmoiron.github.io/sqlx/)
Logging | [Logrus](https://github.com/Sirupsen/logrus)
Test Framework | stdlib + [Testify](https://github.com/stretchr/testify)
Test Coverage | [Goveralls](https://github.com/mattn/goveralls)

## Features Implemented

- [x] Hello world routes
- [x] CRUD routes for persisting posts
- [x] Database access
- [x] Request logging to /srv/app/log/app.log
- [x] Unit tests
- [ ] Unit test coverage reporting
- [x] Automated testing using TravisCI
- [ ] Automated coverage reporting using Coveralls
- [ ] CRUD routes for user management
- [ ] Password encryption using bcrypt
- [ ] Routes protected via HTTP authentication
- [ ] Routes protected via ACLs
- [x] Validates environment (env vars, database host and port are accessible)
