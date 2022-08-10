![EmailIt](https://socialify.git.ci/raguay/EmailIt/image?description=1&descriptionEditable=A%20email%20and%20text%20utility%20program%20with%20many%20unique%20features.&font=Source%20Code%20Pro&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark)

[![Richard's GitHub stats](https://github-readme-stats.vercel.app/api?username=raguay)](https://github.com/anuraghazra/github-readme-stats)

# EmailIt

A simple program for sending emails quickly through the EmailItServer. It also has notes, scripts to act on the text, and templates. I'm still actively developing this program, but the basics are functional. I've only built it as a macOS universal build. I'll be adding other builds in the future. This program is built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev/).

If you have questions, you can ask them in the [Discussions](https://github.com/raguay/EmailIt/discussions). If you have a bug, please create an [Issue](https://github.com/raguay/EmailIt/issues) for it. Thanks.

## How to build

In order to build the program, you have to have [node and npm](https://nodejs.org/en/), [go](https://go.dev/) and [Wails 2](https://wails.io/) installed on your system.

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

This requires the [EmailItServer](https://github.com/raguay/EmailItServer.git) for the back-end logic if you are making it yourself or running the dev version.

EmailIt automatically launches the EmailItServer and [ScriptBar](https://github.com/raguay/ScriptBarApp) if you are using the latest download.

## Documentation

EmailIt came about when my favorite emails sending program went vaporware on me (The Let.ter application for macOS). I still needed a way to send markdown based emails to people quickly. I then merged in the [ScriptPad](https://github.com/raguay/SvelteScriptPad) program I created to have multiple notes, text altering scripts, and templates. I've also integrated the [ScriptBar](https://github.com/raguay/ScriptBarApp) program since it needs the EmailIt server to run.

I use this program everyday and is very helpful to my workflow. I hope you enjoy it as well.

### Table of Contents

- [How To Use](#how-to-use)
- [Scripts](#scripts)
- [Templates](#templates)
- [Script Terminal](#script-terminal)
- [Preferences](#preferences)
  - [General](#general)
  - [Theme](#theme)
  - [Node-Red](#node-red)
  - [External Scripts](#external-scripts)
  - [Environments](#environments)
- [Change Log](#change-log)

## How to Use

EmailIt is a program to email markdown based emails quickly, nine notes to organize thoughts, scripts to change the text, templates to create basic formats for notes or emails, a Node-Red server, and more. The scripts, templates, and notes can be used from other programs as well. I've created workflows for Alfred, Launchpad, Keyboard Maestro, and others.

NOTE: Work in progress.

## Scripts

NOTE: Work in progress.

## Templates

NOTE: Work in progress.

## Script Terminal

NOTE: Work in progress.

## Preferences

NOTE: Work in progress.

## General

NOTE: Work in progress.

## Theme

NOTE: Work in progress.

## Node-Red

NOTE: Work in progress.

## External Scripts

NOTE: Work in progress.

## Environments

NOTE: Work in progress.

## Change Log

### Feactures to add/fix - not in order

- Finish and update existing Help pages.
- Currently, if external program changes a note, EmailIt doesn't know about it.
- Running a script from inside of a script isn't working.
- Add the ability to get a web page from the server to launch scripts. ScriptServer in the works.
- Add scripts from GitHub that are marked for EmailIt.
- Add themes from GitHub that are marked for EmailIt.
- Get it working on Linux and Windows.
- BullitenBoard: A program to display messages and dialogs to the user. Currently a NW.js program. Move it to Wails.
- Add in the regular expressions editor and runner.
- Make the automatic launching of ScriptBar programmable/optional.
- Get ScriptBar able to use xBar scripts.
- Launching the adding of different workflows/extensions to other programs (Alfred, LaunchBar, Keyboard Maestro, etc.)
- Test running on a new system.
- ScriptBar: Have a message to add items when there isn't a configuration.


### Change Log:

#### August 10, 2022

- EmailIt is a full Wails application that has the EmailIt server and ScriptBar embedded and launched from the program.
- Both are updated to the Wails Beta 2 version 43 and updated to latest javascript libraries.
- Script Terminal now works in EmailIt. Many of the basic termainal commands have been implimented as well. Now you can run your scripts on any file in your system easily.
- You can create terminal scripts in JavaScript or as an external script. Just click the `Terminal Script` checkbox.

#### September 21, 2021

- ScriptBar
  - Fully functional as a Wails project except for menubar icon and proper hiding.

### September 13, 2021

#### ScriptBar

- Working on making it a Wails project.
- Update to using Stores instead of events for passing values.
- Changed how resizing is done (much simpler!)

#### September 02, 2021

- Working on Preferences: Themes.
- Working on saving and creating new themes.

#### August 25, 2021

- Changed the input for emails to allow multiple email addresses. Using my own made input instead of SimpleAutoComplete.
- Making second release with this change.

#### August 21, 2021

- Created first release of EmailIt and EmailItServer.
- Fixed a few minor bugs as well.
- Started creating/importing a Preference Editor for changing/creating themes, environments, and external scripts.

#### August 20, 2021

- Adding/Editing/Deleting Scripts and Templates now work.
- Fixed bug in the Script and Template menu.
- Sorted items in the Script Menu, Template menu, Script editor, and Template editor.

#### August 17, 2021

- Added the ability to run templates
- Added theming the editor dynamically
- Fixed several bugs with script and template running
- Email body editor is now CodeMirror 6

#### August 14, 2021

- The email body editor is not CodeMirror 6
- Fixed the saving of the email body and restoring it when going to the notes and back to the email.

#### August 11, 2021

Moving the whole project to Tauri and CodeMirror 6 along with an email utility. I'm renaming it all to EmailIt. Notes, running scripts, displaying the log, and the EmailIt parts are working.

##### Things to work on:

- Integrate Templates again
- Setup and expand the script server
- Update and expand the help pages
- EmailIt to handle names in the email address
- EmailIt to allow multiple email to addresses
- Edit/Add to the address book
- Change the Email Alfred workflow to use the builtin addresses

#### July 02, 2019

- The main browser and bulletin board are now compiled. No more loading plain JavaScript.
- The server is still running JavaScript directly. I can't get it to load compiled code for it. Most likely need to move to an external program.

#### July 01, 2019
- CodeMirror Component
	- Fixed the setting of a height and width
- Script Editor
	- Fixed annoying shift in the editor location
    - Added a Cancel button to gracefully exit
    - Tidied up the looks
    - Added the `Insert` label
    - Moved the NP variable to SP for `ScriptPad`
- RegularExp Editor
	- Added a Cancel button and changed the "Change All" button to just "Change".
    - Fixed some looks and feel stuff
- Todos
	- Added the ability to exit out of the file watchers when exiting the program.

#### June 21, 2019
- Todos
	- Fix the server code to not die when NotePlan hasn't made the new file for today yet.
	- Added two way syncing of todo lists. External program changes a todo, it will get loaded into the Todo list and not loose current location.
    - Fix the server code to not die when NotePlan hasn't made the new file for today yet.

#### June 20, 2019
- Todo Preferences
	- Fixed file selector to use system file selector.
    - Editing a Taskpaper file shows the path in a input box.
    - It now properly propagates to the Todo panel.
#### June 17, 2019
- Todo Preferences
	- Need a better file selector.
	- If preferences is open, disable hiding on blur. Needed for the file requester to work.

#### June 14, 2019
- Fixed regression where handle to editor wasn't obtainable in the old way. Cursor placement and history fixed.
- Fixed endless loop in Regular Expressions editor.

#### June 10, 2019
- Preferences
	- Need to style the selector for keymaps.
- Documentation
	- lots have been done, but not finished.
#### May 01, 2019
- Finished the Preferences page - first pass
	- General Preferences - first pass
    - Todos is done and seems to work fine.
    - Change the input box to a selector for Keymaps.

#### April 30, 2019
- Rest of the Theme preferences
- Todo preferences
- The todo list doesn't update until you exit and reload the program. Need to reload todo stuff when a change happens with them.

#### April 29, 2019
- Todo list for Preferences are working! Adding, deleting, and editing TaskPaper todos works!

#### April 25, 2019
- Finished the theme section of Preferences.
- Upgraded to the lastest NW.js. Much faster!
- Upgraded to the full Svelte v3 and the latest CodeMirror.

#### April 23, 2019

- ColorPicker is now working great and positioned okay.
- All of the theme color circles are there, but not connected to the color picker yet.
- Setting new editor themes now works!

#### April 22, 2019
- Saving a style saves the current button as selected also. When saving, set all buttons to false and set the proper one to true.
    - Finish the color picker component.
- Fixed the Preference button to look sunken when active.

#### April 18, 2019
- Working on the theme preferences. Showing the different button colors and allows the changing of the color. Need to make a better color picker, but it's getting there. Not saving styles yet.

#### April 16, 2019

- Preferences and styles are coming from configuration files now and not hard coded.
- Some work done towards creating the preferences dialog interactions.
- If todos are all turned off, then the clipboard icon isn't displayed.

#### April 09, 2019
##### Fixed:
- Get TeaCode integration working again.
- No history sharing between notes. But, no preserving history when changing notes either.
- Sets editor config from main config. This gives different keymaps, line number, etc for each editor. It is now set from the preferences in the main component.
- Added bracket completions, search/replace, vim keymap, sublime keymap, emacs keymap, jump-to-line

#### April 08, 2019
##### Fixed:
- Typing in the editor puts the current line to the bottom of the viewport if it has been scrolled down some.

#### March 25, 2019
##### fixed
- The replace in the regular expression editor 
- should properly insert the \n and \r characters.

#### March 21, 2019

- Todo list looks better 
- Done items are automatically placed at the bottom, but still in order in the list.
- The regular expression editor no longer removes all spaces. Captures are colorized while non-captured items aren't.

##### Fixed:
- Get the RegExp Editor to show text inbetween inner captures.
- Get the RegExp Editor to show captures inline, not in block form.


#### March 20, 2019

- The Svelte version is fully feature compatable with the Mint version.
- Also, single regular expression fixes done.
- But, lost my old notes here with testing.


