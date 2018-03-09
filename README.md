# analyticsServer

## Summary

**Status:** In development

It is minimalistic analytics server written on Go and uses MongoDB as storage.

## How it works

**Send data:**

1. *POST* request to / with one event content
2. Content parsed and saved in MongoDB

**Get data:**

1. *GET* request to /
2. All events returned