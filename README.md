[![Richard's GitHub stats](https://github-readme-stats.vercel.app/api?username=raguay)](https://github.com/anuraghazra/github-readme-stats)

# EmailIt

A simple program for sending emails quickly through the EmailItServer. It also has notes, scripts to act on the text, and templates. I'm still actively developing this program, but the basics are functional. I've only built it on a Macbook Pro with an M1 chip. I'll be adding other builds in the future. This program is built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev/).

## How to build

For the first time, you need to download the libraries used by:

```sh
cd frontend
npm install
npm run build
```

To build the program and package it properly, you run:

```sh
wails build
```

or, if you want to build the macOS universal application, you run:

```sh
wails build --platform "darwin/universal"
```

If you want to develop the program, you can run developer mode with:

``` sh
wails dev
```

This requires the [EmailItServer](https://github.com/raguay/EmailItServer.git) for the back-end logic currently. I've not gotten far enough to have it all integrated. This project is slowly being created. My plan is for the EmailItServer and [ScriptBar](https://github.com/raguay/ScriptBarApp) to be launched and used from EmailIt.

## Documentation

### Table of Contents

- [How To Use](#how-to-use)
- [Scripts](#scripts)
- [Templates](#templates)
- [Preferences](#preferences)
  - [General](#general)
  - [Theme](#theme)
  - [Node-Red](#node-red)
  - [External Scripts](#external-scripts)
  - [Environments](#environments)

### How to Use

EmailIt is a program to email markdown based emails quickly, nine notes to organize thoughts, scripts to change the text, templates to create basic formats for notes or emails, a Node-Red server, and more. The scripts, templates, and notes can be used from other programs as well. I've created workflows for Alfred, Launchpad, Keyboard Maestro
### Scripts

### Templates

### Preferences

#### General

#### Theme

#### Node-Red

#### External Scripts

#### Environments

