# MongoTUI - MongoDB TUI client

[![Build Status](https://dev.azure.com/janrenken/MongoTui/_apis/build/status/renkman.mongotui?branchName=main)](https://dev.azure.com/janrenken/MongoTui/_build/latest?definitionId=3&branchName=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/renkman/mongotui)](https://goreportcard.com/report/github.com/renkman/mongotui)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://github.com/lachsfilet/Renkbench/blob/master/LICENSE)

![Screenshot](mongotui.png)

[MongoDB](https://www.mongodb.com/ "MongoDB") TUI client written in [Go](https://golang.org/ "Go"), using the [tview](https://github.com/rivo/tview/ "tview") library for UI, [Keyring](https://github.com/99designs/keyring) for storing connections and the [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver/ "MongoDB Go Driver").

The unit tests are using [Testify](https://github.com/stretchr/testify "Testify").

## Features

- Tree view of the connected instance, its databases and collections
- Command execution
- Result tree view

