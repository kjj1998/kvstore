# Redis Lite Key Value Store

## Introduction

This is an implementation of a Redis Lite Key Value Store that is done for my personal learning purposes.

## Features

### SET Command
``` bash
SET <key> <value>
SET <key> <value> EX <seconds>
```
- Stores the value under the given key
- Overwrites the existing value if the key already exists
- Value of the key can also be stored with a time-to-live in seconds

### GET command
``` bash
GET <key>
```
- Gets the value stored under the given key
- `EXPIRED` returned if the key has expired
- `NULL` returned if the key does not exist

### DEL command
``` bash
DEL <key>
```
- Deletes the value stored under the given key