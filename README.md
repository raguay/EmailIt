![EmailIt-Tauri](https://socialify.git.ci/raguay/EmailIt-Tauri/image?description=1&descriptionEditable=A%20simple%20and%20quick%20email%20program%20to%20replace%20Let.ter%20application%20with%20notes%20and%20scripting%20added%20in%20as%20a%20bonus.&font=Bitter&forks=1&issues=1&language=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark)

[![Richard's GitHub stats](https://github-readme-stats.vercel.app/api?username=raguay)](https://github.com/anuraghazra/github-readme-stats)

# EmailIt

A simple program for sending emails quickly through the EmailItServer. It also has notes, scripts to act on the text, and templates. I'm still actively developing this program, but the basics are functional. I've only built it on a Macbook Pro with an M1 chip. I'll be adding other builds in the future.

For the first time, you need to download the libraries used by:

```sh
npm install
```

To build the rust program and package it properly, you run:

```sh
npm run build
cargo tauri build
```

The `npm` line builds the frontend Svelte project. The `cargo` line builds the actual rust application that runs the program. 

In order to get the EmailIt icon, you must copy it from `graphics/emailit.icns` to `EmailIt.app/Contents/Resources/icon.icns`. The normal build doesn't like my icns file for some reason. An an M1 chips system, the `cargo tauri` program doesn't have the `icon` subcommand to generate the right icns file.

This requires the [EmailItServer](https://github.com/raguay/EmailItServer.git) for the backend logic currently. I've not gotten far enough to have it all itegrated. This project is slowly being created.

# Note:  This version is build with the Tauri framework. I've moved to using the [Wails](https://wails.io) framework instead since I'm better with go language than Rust.
