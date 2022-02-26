# Formica Core

Base backend services for formica industrial analytic platform

## build
`docker build --build-arg APP=signal . -t formica-signal`


## Signal
Entry point of incoming data to the system

`docker run formica-signal -p 8080:8080 --name formica-signal --network formica`

## Report
Meaningful consumption of saved signals into the system