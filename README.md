![EmailIt](https://socialify.git.ci/raguay/EmailIt/image?description=1&descriptionEditable=A%20email%20and%20text%20utility%20program%20with%20many%20unique%20features.&font=Source%20Code%20Pro&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark)

[![Richard's GitHub stats](https://github-readme-stats.vercel.app/api?username=raguay)](https://github.com/anuraghazra/github-readme-stats)

# EmailIt

A simple program for sending markdown formatted emails quickly. It can expand macros defined using [Handlebars](https:///handlebarsjs.com/). It also has notes, scripts to act on the text, and Handlebar templates. You can set default data to be used in the expansions. I'm still actively developing this program, but the basics are functional. I've only built it as a macOS universal build. I'll be adding other builds in the future. This program is built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev/).

![EmailIt Main Screen](/images/mainscreen.png)

It also has a easy to use TUI for sending emails from the command line. By using flags to fill in the address for sending, subject line, and body, you can quickly send an email. It will then open in a TUI or you can have it send it directly. This is what the TUI looks like:

![EmailIt TUI](/images/emailTUIdemo.gif)

The email that was sent by this recording is shown here in my Gmail account:

![EmailIt TUI email result](/images/EmailSentByTUI.png)

Along with all of this, it has a command prompt that is easily set to follow your directory changes in the terminal for allowing you to run scripts on files or simply launch files with your editor or whatever program is set on your system. I call it ScriptLine since it's main purpose is for running scripts.

![ScriptLine](/images/scriptline.png)

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

The maskfile supplied will compile the universal build on macOS. You can invoke it with:

```sh
mask build
```

You can look at the maskfile.md file to see what the build command does if you don't want to use mask itself.

I develop on a MacBook Pro M1 system. Therefore, all features do work on a similar system. I'm trying to test and use it on Linux and Windows, but I don't currently have a Windows system to work with.

## Documentation

EmailIt came about when my favorite email sending program went vaporware on me (The Let.ter application for macOS). I still needed a way to send markdown based emails to people quickly. I then merged in the functionality of [ScriptPad](https://github.com/raguay/SvelteScriptPad) program I created to have multiple notes, text altering scripts, and templates. I often need to reference or copy something from a note to an email. I also have many email templates that I use in the Template interface. 

I use this program everyday and is very helpful to my workflow. I hope you enjoy it as well.

### Table of Contents

- [Setup as Mailto Protocal Handler](#setup-as-mailto-protocal-handler)
- [How To Use](#how-to-use)
	- [Address Book](#address-book)
	- [Notes](#notes)
	- [Scripts](#scripts)
	- [Templates](#templates)
	- [ScriptLine](#scriptline)
  - [TUI and the command line](#tui-and-the-command-line)
	- [Preferences](#preferences)
	  - [General](#general)
	  - [Theme](#theme)
	  - [External Scripts](#external-scripts)
	  - [Environments](#environments)
	  - [GitHub](#github)
	- [Change Log](#change-log)
	
## Setup as Mailto Protocal Handler

I have yet to figure out how to make EmailIt a normal mailto protocal handler, but you can do it with the help of [LinCastor](https://onflapp.github.io/blog/pages/LinCastor.html). When you download the program and create a mailto protocal handler, add this script:

```sh
#!/bin/zsh
# the script will execute with following environmental variables defined:
#   URL           => my-http://myhost.domain.com:8080/mysite/a.html?search=blah#myanchor
#   URL_SCHEME    => my-http
#   URL_HOST      => myhost.domain.com
#   URL_PORT      => 8080
#   URL_QUERY     => ?search=blah
#   URL_PATH      => /mysite/a.html
#   URL_PATH_EXT  => html
#   URL_PATH_NAME => a.html
#   URL_FRAGMENT  => #myanchor
#   URL_VALUE     => everything that comes after the 'scheme:'
#   URL_B64VALUE  => the same as URL_VALUE but decoded using based64
#
# NOTE: The system env variable (e.g. PATH) may have different values
# then what you would get if you would run the script in the terminal.
#
# if you are mapping a file extension, use URL_PATH to pass the file path
#
# ---
#
# If you intended to trigger this handler from a bookmarket,
# you can 'print out' valid javascript code to be execute in the browser when this scripts returns.
# Unless you execute the command asynchronosly in the Terminal.
#
# ---
#
# you should exit with 0, this means the handler has finished succesfully
# non-zero value indicates an error

/Applications/EmailIt.app/Contents/macOS/EmailIt $URL
```

And then click the 'Save and Activate' button in the lower right corner.

![Lincastor Mailto Config](/images/lincastor.png)

Once done, anytime you click on a mailto link, it will open in EmailIt for your additions before sending.

## How to Use

When you launch the program, the EmailIt main screen is shown:

![EmailIt Main Screen](/images/mainscreen.png)

Here, you can easily type in an email address, a subject line, and the text body. The text body is assumed to be markdown format. It will be sent to the Handlebars macro expander and then converted to HTML and plain text email. The buttons along the bottom show the different functions that can be accessed at this screen. Any email address used is saved into the [Address Book](#address-book). When you enter the `To` field, a list of addresses from the Address Book will be displayed. As you type, the list will contract to the ones that match what you are typing.

![EmailIt Addresses Display](/images/addressdialog.png)

Before you can send an email, the email SMPT servers need to be setup. When you first run the program without any server, the button that says `CustomCT` will say `Create a New Account`. When you click it, you will be shown the accounts dialog.

:![EmailIt Accounts Dialog](/images/account.png)

I have several accounts setup, but for your first time, it will only have the buttons on the bottom. When you click on of the account buttons, that account is set to the `Current Account`. If you click `Cancel`, the current account will go back to the one when you first entered this dialog. If you click `Delete`, the `Current Account` will be deleted and changed to the topmost account in the list. If you click `Save`, you will go to the main screen with the account button showing the name of the account you selected. If you click `New`, the New Account dialog will be shown.

![EmailIt New Account Top](/images/newaccount1.png)

The `Default` checkbox will make the new account the default account. All the other fields are fairly self explanitory. They are the different information needed to connect to most SMPT servers.

![EmailIt New Account Bottom](/images/newaccount2.png)

When you scroll down the new accounts dialog, the bottom fields are seen. The `Signiture` field is added to the end of all your emails. This should be in HTML format. The `Header HTML` and `Footer HTML` will be placed at the top and bottom of your email. This is where you can customise the look of your email with font size, color changes, etc. You would put an opening `<div>` with your custom styling in the header HTML and the footer HTML would close it out.

When you are done, simply click `Save` or `Cancel` if you don't want to save it. Either one will take you back to the main screen.

Once your message is typed, addressed, and has a subject, you can click the `Preview` button to see what it will look like:

![EmailIt Preview](/images/preview.png)

If you like what you see, you can click `Send`. Since my configuration doesn't setup any background or foreground colors, you see the default for the EmailIt program (or if you change the theme, whatever the new theme uses). The person that receives it will see the colors their email program uses (usually white background with black type). If you click `Edit`, you will go back to the main screen with the information you typed in still there. If you press `Clear`, you are sent to the main screen without any information. Back to the beginning to start over.

## Address Book

When you click `Address Book` on the main window, you will get the Address Book dialog.

![Address Book](/images/addressbook.png)

Here, you can edit an address by clicking the pencil, delete an address by clicking the `X`, and create a new address by clicking `New` button. These addresses are collected when you write an email or when you add them in this dialog.

## Notes

When you click `Notes` button on the main screen, you will see the notes you are keeping in the program.

![Notes](/images/notes.png)

The ten color circles to the right select the different notes. I use one as a scratch pad for running text altering scripts. The buttons on the bottom go to different parts of the program.

## Scripts

In the EmailIt editor or the Notes editor, you can click the `Scripts` button or `<ctrl>-m` to toggle the scripts menu. 

![Scripts Menu](/images/scriptmenu.png)

The input at the top allows you to search the list by removing all items that don't match the letters you type. When you select a script and have text selected, it will process that text with the script and replace it with the results. If you don't select text, it will process the entire note.

You can create or edit your scripts when you press the `Edit Scripts` button.

![Scripts Editor](/images/scripteditor.png)

The script editor is where you can create or edit your different scripts. System scripts or External scripts can't be edited here. When you type the `Name` field, it will show scripts similar to what you type. If you select one of them, you can edit that script. If you give a unique name, you can create that script.

The `Insert?` checkbox is for telling the editor to insert the output of the script. 

The `Terminal Script` tells the program that the script is for the Script Terminal and the output is a JSON structure for displaying information that can be acted upon. This is described more in the [Script Terminal](#script-terminal) section.

The `Description` field is for a brief description of what the script does. For a Terminal Script, it is displayed in the `help` command.

The text editor area is for writing your JavaScript routine. The `SP` object is how you get information and run different functions. 

The variable `SP.text` contains either all the text in the note if there wasn't a selection, or just the selected text. You read this variable to get the text to process and write into the variable what you want to replace either the selection or all the text in the notepad.

There are some predefined libraries in variables for your scripts to use. You can make use of the following:

|  |  |
| -- | ---- |
| SP.mathjs | [Math.js library](http://mathjs.org/) |
| SP.mathParser | The Math.js parser used to process text |
| SP.DateTime | [luxon library](https://moment.github.io/luxon/#/) |
| SP.Handlebars | [Handlebars library](https://handlebarsjs.com/) |

The are some predefined function available as well:

|  |  |
| -- | ---- |
| SP.insertCharacters(`<num>`,`<char>`) | This function makes a string of `<num>` `<char>`. |
| SP.returnNote(`<id>`) | This function returns the note with the `<id>`. |
| SP.runScript(`<scrpt>`,`<text>`) | This function runs the `<scrpt>` as a string on the `<text>` text. |
| SP.ProcessMathSelection(`<text>`) | Runs the given text through the Math.js parser. |
| SP.ProcessMathNote(`<id>`) | Run the given note id through the Math.js parser. |

These functions will not be usable for external scripts. They are for internal JavaScript scripts and internal Scriptline commands.

When creating terminal commands for Scriptline, the `SP` structure will have these added data points:

| | |
| -- | ---- |
| SP.data.previousLines | This will contain the previously ran command. |
  SP.data.registers | This is an array of text registers that can hold information from one chained command to the next. |
  SP.data.line | This is the command line for the currently executing command. |
  SP.data.result | This is the previous commands result. |

This structure can change over time. I'm still experimenting with many ideas for it. I have plans for creating more functions to be used with this scheme. Stay tuned!

You can create script in one note and use a different note for the input.  For example, in a note, place the following code:

```js

try {
  var lines = SP.text.split('\n');
  SP.text = '';
  for(var i = 0; i < lines.length ; i++) {
      var match = lines[i].trim().match(/^\d+\. (.*)$/);
    if (match != null)
        SP.text += match[1] + '\n';
  }
} catch (e) {
   SP.text += "\n\n" + e.toString()
}

```

Then go to a different note and place several lines of text. Run the
script `Bullet lines with Numbers`. Every line will have the proper number at the front of it. Now, run the script `Evaluate Note # as Script` with `#` the number of the note where you put the script. The numbers at the beginning will now be removed! You can use the script editor to save this script and use it from the script menu.

There are scripts for processing math: the 'Basic Math' and 'Evaluate Page for Math' scripts. The 'Basic Math' script is for processing arbitrary pieces of math in a selection. The 'Evaluate Page for Math' script is for processing the entire note with a nice running result along the right. The 'Basic Math' script doesn't reset the state of the math library (ie: variable definitions and functions), but the 'Evaluate Page for Math' does each time invoked so as to not create multiple copies of function and variables.

Copy the following note to a notepad:

```md

# Using the ‘Evaluate Page for Math’ script

Your notes can have any text you need. But when a line starts with a ‘>’, t
hat whole line is processed for math. The line is processed and the answer pushed
to the right with a ‘|’ symbol.

> 6 * 95
> x=6*8-10
> x

Text in the middle doesn’t clear out the variable or function assignments before it.

> f(x)=x^2-5*x+12
> f(60)
> f(x)

The length of the note isn’t a concern either.

> f(100)

> bank=34675
> check=5067
> balance = bank - check

> sin(45)

The math package used doesn’t do conversions or symbols inside of the mat
h expressions. The math library used is [mathjs 4.0](http://mathjs.org/).

```

Then press `<ctrl>-m` and select the 'Evaluate Page for Math' script. Each line with the '>' as the front character now has the results to the right.  When you change the text lines and re-run the script, the math lines are all updated. All other lines are not effected by the script. You can change any equation or variable and it's effects will trickle down the page.

## Templates

In the EmailIt editor or the Notes editor, you can click the `Templates` button or `<ctrl>-t` to toggle the template menu. 

![Templates Menu](/images/templatemenu.png)

The input at the top will filter the list according to the letters typed. When you select a template, it is placed at the cursor point in the note.

You can create/delete/edit templates in the Template Editor. You access the Template Editor by the `Edit Template` button at the bottom of the Note Editor.

![Templates Editor](/images/templateeditor.png)

When you type in the `Name` field, a list of templates that match the text is given. Selecting one of the templates allow you to edit that template. By typing a new name, you can create a new template.

The `Description` field is for giving a short description.

To save the new or edited template, you press the `Save` button at the bottom. To delete the current template in the editor, press the `Delete` button at the bottom of the editor.

The text editor is where you create your template. Templates are processed using [Handlebars](https://handlebarsjs.com/) parser. Handlebars allows you to have a text file with anything in it along with some macros. Some macros need auguments and some do not. The macros are expanded and placed into the text. Macros with arguments are called helpers. Along with the standard Handlebar helpers, several others have been added as well. The following is an explanation of the additional helpers:

| | |
|--|----|
| `{{date <format>}}` | This will format the current date and time as per the format string given. See the help document that is loaded upon initialization. |
| `{{cdate <date/time> <format>}}` | This takes the date/time string and formats it according to the format given. See the help document that is loaded upon initialization. |
| `{{last <weeks> <dow> <format>}}` | This will print the date `<weeks>` ago in the `<format>` format for the `<dow>` day of week. |
| `{{next <weeks> <dow> <format>}}` | This will print the date `<weeks>` in the future using the `<format>` format for the `<dow>` day of week. |

The following data expansions are defined as well:

| | |
|--|----|
| `{{cDateMDY}}` | gives the current date in Month Day, 4-digit year format |
| `{{cDateDMY}}` | gives the current date in Day Month 4-digit Year format |
| `{{cDateDOWDMY}}` | gives the current date in Day of Week, Day Month 4-digit year format |
| `{{cDateDOWMDY}}` | gives the current date in Day of Week Month Day, 4-digit year format |
| `{{cDay}}` | gives the current date in Day format |
| `{{cMonth}}` | gives the current date in Month format |
| `{{cYear}}` | gives the current date in 4-digit year format |
| `{{cMonthShort}}` | gives the current date in Short Month name format |
| `{{cYearShort}}` | gives the current date in 2-digit year format |
| `{{cDOW}}` | gives the current date in Day of Week format |
| `{{cMDthYShort}}` | gives the current date in Month day 2-digit year format |
| `{{cMDthY}}` | gives the current date in Month Day 4-digit year format |
| `{{cHMSampm}}` | gives the current date in h:mm:ss a format |
| `{{cHMampm}}` | gives the current date in h:mm a format |
| `{{cHMS24}}` | gives the current date in H:mm:ss 24 hour format |
| `{{cHM24}}` | gives the current date in H:mm 24 hour format |

I'm working on a way to have user defined Handlebar helpers and macros. 

## Scriptline

The Scriptline is a terminal that runs scripts in the EmailIt program on text given or on files given. You access the Scriptline from any other screen by pressing `<ctrl>-s`. It can also be opened from a command sent to the internal server.

![Scriptline](/images/scriptline.png)

By typing help, you will see all script commands and scripts. The current list of builtin commands are:

|  |  |
|--|----|
| help |   Show the help command in the terminal. |
| ls   |   List the items in the specified directory. |
| cd   |   A terminal script for changing directories. |
| open |   Open the file given. |
| runscript |   Run a script on a file or text string. |
| edit  |  Edit the given file. |
| alias |   The alais command allows you to substitute a single command for a list of commands. |
| hist  |  The hist command lists the previous command lines that worked. You can activate one to rerun the command. |
| rm    |  The rm command will delete the given file or directory. If nothing is given, it will list the current directory and an item can be delete using the command mode and r. |
| mkdir |   The mkdir command makes the given directory if it doesn't already exist. |
| mkfile |  The mkfile command makes the given file if it doesn't already exist. |
| rmalias | This command removes the given alias. |
| notes | This command will go to the notes. |
| scripteditor | This command will go to the Script Editor. |
| emailit | This command will go to the EmailIt. |
| quit | This command will quit the EmailIt program. |

This is in no way a full terminal emulator. It is a simple line based command runner. No pipes, redirections, flow control, etc using glyphs, but Scriptline commands can use some flow control by taking advantage of the data given. I mostly use this to run scripts on text files. You can run several commands at once by separating them with a `;` or by using functional programming style of function chaining. When you have a series of commands, each command can access what the previous command was and it's output in the `SP` data structure. Even external commands can get this information. Only the last command ran can be used for showing results. The command chaining can be used in the commands for the Command state or `tcommand` fields in the terminal script JSON output.

All scripts for use in the Scriptline have to take in a JSON structure called `SP` (same structure described in the scripts section but with some added fields) and output a JSON structure that tells the script terminal what to do. The JSON output structure is:

```json
{
  tcommand: <terminal command to run>
  lines: [{
    text: <text to display>,
    color: <color to show>,
    command: <command line string to run when selected>
  }, {
     <next line structure>
  }, ...]
}
```

The `tcommand` is a builtin command that is ran directly when the command is ran.

The `lines` structure is an array of objects containing a `text` field with the text to be printed to the terminal, `color` field is the color to print the output in, and a `command` field to be executed when the `r` key is pressed in Command state.

The valid colors are: red, black, green, orange, blue, magenta, cyan, gray, and default. These colors are controlled by the current theme for EmailIt. The default color is the text color.

## TUI and the command line

You can enter act with EmailIt from the command line. It has some TUI programs that make using the overall program from the command line easy. The command line for EmailIt has flags and commands. In order to use the program from the command line, you need to make an alias to the programs executable file for the macOS:

```bash
alias em=/Applications/EmailIt.app/Contents/MacOS/EmailIt
```
With the above line in your `.bashrc` or `.zshrc` file (you will have to make a similar line for other shells), you can run the EmailIt program from any directory. The different command line options are shown with the basic `-h` or `--help` flag or the `help` or `h` command:

```bash
> em help
NAME:
   EmailIt - A program for sending emails, working with text, and scripts.

USAGE:
   EmailIt [global options] command [command options] [arguments...]

VERSION:
   v2.1.0

AUTHOR:
   Richard Guay <raguay@customct.com>

COMMANDS:
   mkemail, me     Create an email using a TUI
   notes, n        Open the notes.
   emailit, e      Open the EmailIt email sending application.
   scriptline, sl  Open the ScriptLine application.
   sendemail, se   Send the email directly. No GUI or TUI.
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -a value       Address to send an email
   -s value       Subject for the email
   -b value       Body of the email
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   (c) 2022 Richard Guay
 ```


## Preferences

The preferences can be reached by pressing `<ctrl>-p` or `<cmd>-,` anywhere in the program. There are four sections currently: General, Theme, Node-Red, External Scripts, Environments, and GitHub.

### General

![General Preferences](/images/generalpref.png)

The general section currently is empty. I'm adding functionality to install workflows for Alfred, Launchpad, Keyboard Meastro, etc. 

### Theme

![Theme Preferences](/images/themepref.png)

The Theme preference is where you can create and change your theme. By selecting `new` in the theme selector, you can create a new theme based on the current theme. The circle colors are for the Notes editor note selector on the right. The `Update Theme` button will save your changes. The `Delete Theme` button will delete the currently selected theme and set the default theme. You can change the color by clicking on the color circle.

![Theme Preferences - Color Picker](/images/colorpicker.png)

You can adjust the color by selecting an area in the color wheel. When you get the color you like, press the `Select` button. Pressing the `Cancel` button will go back to the original color.

### External Scripts

You can use external scripts to process text and for use in the Scripts Terminal. To setup a external script for using in EmailIt, simply go to the External Scripts Preferences.

![External Scripts Preferences](/images/externalscripts.png)

The main screen shows the currently created external scripts. Clicking on a script name will open it in the External Script Editor. Clicking the `New` button will open the External Script Editor in a blank script.

![Adding an External Script](/images/addingexternalscript.png)

Here you setup the different setting for the external script.

The `Name` field is the name used in the Scripts Menu and for the `SPRunScript` Node-Red node. It doesn't have to anything to do with the actual name of the script file.

The `Description` field is used to describe the action the script will perform. For Terminal Scripts, it will be shown in the general help builtin command.

The `Help` field is displayed for terminal commands when you add the script name to the `help` builtin command.

The `What's the name of the Script File?` is where you put the real name of the script. This should be the main file for running the script. It should be executable and have the proper `#!` (shebang) as the first line of the script.

The `What's the directory of the script?` field is the full path to the directory that contains the script.


![Adding an External Script](/images/addingexternalscript2.png)

The `What is the environment for the script?` is a drop-down selector for selecting the environment to use for running the script.

The `Terminal Script?` checkbox is for specifying if this is a terminal script.

### Environments

All external scripts are ran using an environment. The environment describes the environment variables that are set before running the script. This is very useful for setting up a path for a particular version of a program. Scripts are always ran in the directory of the script.

![Environments Preferences](/images/environments.png)

The main screen shows the different environments defined already. If there isn't a `Default` environment, then a button is shown to create it. Clicking on an environment name opens it in the environment editor.

![Adding an Environment](/images/addingenvironment.png)

The environment editor has the name at the top just like in the Theme Editor. Changing the name allows you to create a new environment based on the currently selected one. You can delete rows by clinking the `X` or edit the row by clicking the pencil. The `Delete` button will delete the currently displayed environment. The `Return` button will save all changes and go back to the environment list.

### GitHub

![GitHub Importer Screen](/images/github.png)

The GitHub importing screen allows you go import external scripts and themes that are stored on GitHub. To install the script or theme, click the `install` button. Themes then have an extra button to `load` the theme downloaded.

In order for this importer to find the script and themes, your script or theme has to have the `emailit` topic. A script will need the `script` topic as well, while a theme will need the `theme` topic included.

When a script or theme is downloaded, a `delete` button will be shown. You can remove the script or theme by pressing the `delete` button.

In order to load the script properly, the `package.json` file in the script's repository has to have these added fields:

```json
    "script": {
        "name": <name of the script>,
        "description": <a short description of what the script does>,
        "help": <help text for the script>,
        "termscript": <boolean flag for a Script Terminal script (true if it is a script terminal script. Otherwise, false.)>,
        "script": <the scripts full name in the repository including extension>,
    }
```

In order to load the theme properly, the `package.json` file for the theme's repository has to have these added fields:

```json
    "theme": {
        "name": <name of the theme>,
        "description": <description of the theme>,
        "main": <the json file for the theme>,
    }
```

### Web API

EmailIt has a web based API for interfacing with other applications, 
command line tools, or whatever else would help. The base address for the APIs is `http://localhost:9978/api`. Every endpoint in this table builds on this base. These endpoints only allow access from requests on the same machine.

| Endpoint | Description |
| --- | ------ |
| /note/<number>/(a or w) | A PUT request will assert the note value while the GET request will return the note. Both use a JSON structure with the element `note`. The `a` on the end will append to the note while a `w` will over write. |
| /script/list | A GET request will return a list of scripts that can be ran on EmailIt. |
| /script/run | A PUT request requires a JSON body with a `script` element and a `text` element. The `script` script will be ran with the `text` and returned in a JSON structure with a `text` element. |
| /template/list | A GET request will return the name of all the templates in EmailIt. |
| /template/run | A PUT request requires a JSON body with a 'template' element and a 'text' element. The `template` will be ran with the `text` as an input. The results are return in a JSON structure with the result in the 'text' element. |

All the endpoints are used to make the plugins for Alfred, Keyboard Maestro, Dropzone, PopClip, and Launchpad. Also, the ScriptBar program uses these endpoints as well. I'm planning to add serving pages on the user's computer and other functionality as well. I have lots of ideas.

## Change Log

### Features to add/fix - not in order

- Finish and update existing Help pages.
- Running a script from inside of a script isn't working.
- Add the ability to get a web page from the server to launch scripts. ScriptServer in the works.
- Add scripts from GitHub that are marked for EmailIt.
- Add themes from GitHub that are marked for EmailIt.
- Get it working on Linux and Windows.
- Add in the regular expressions editor and runner.
- Launching the adding of different workflows/extensions to other programs
- Test running on a new system.
- User created handlebar helpers and macros.
- Undo history working. 

### Change Log:

#### October 27, 2022 - Version 2.0.0 Release Candidate 1
- Currently, if external program changes a note, EmailIt doesn't know about it. -- Fixed
- BullitenBoard: A program to display messages and dialogs to the user. Currently a NW.js program. Move it to Wails. -- It's a separate project and is a full Wails project.
- Removed the EmailIt Server. Everything is done in Go and Svelte.
- Alfred workflow has been updated to work with the new version.


#### August 18, 2022

- Removed auto launch of ScriptBar. Latest update to macOS will not allow me to without code signing.
- Moved documentation of ScriptBar to the ScriptBar repository.
- Added GitHub importing of External Scripts and Themes.

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
    - Moved the NP variable to SP for `EmailIt`
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


