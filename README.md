# dp-library-consortium
Project for the Master Course in Distributed Programming for Web, IoT and Mobile Systems @unifi MS Software: Science and Technology

## Description

Distributed web application written in Go for managing a consortium of libraries.

## Technologies

- Go
- GORM
- SQLite
- Docker
- HTML Templates

## Architecture

```text
                       Browser
                           │
                     Central Server
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
  Library Service    Library Service    Library Service
      Torino             Milano              Roma
        │                  │                  │
      SQLite            SQLite             SQLite
```

## How to run

docker compose up --build

To access the web application: in your browser navigate to: http://locahost

## Services

- Central Server
- Library Service

## Authors

Samuele Mancini
