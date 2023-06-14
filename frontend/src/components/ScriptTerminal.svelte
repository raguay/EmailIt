<script>
  import { onMount } from "svelte";
  import { Terminal } from "xterm";
  import "xterm/css/xterm.css";
  import { WebLinksAddon } from "xterm-addon-web-links";
  import { FitAddon } from "xterm-addon-fit";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";
  import { termscripts } from "../stores/termscripts.js";
  import { runscript } from "../stores/runscript.js";
  import { aliases } from "../stores/aliases.js";
  import { config } from "../stores/config.js";
  import * as App from "../../wailsjs/go/main/App.js";

  let term = null;
  let commands = [];
  let curr_line = "";
  let tCommands = {
    cd: {
      command: cdCommand,
    },
    help: {
      command: helpCommand,
    },
    ls: {
      command: lsCommand,
    },
    open: {
      command: openCommand,
    },
    runscript: {
      command: runscriptCommand,
    },
    edit: {
      command: editCommand,
    },
    alias: {
      command: aliasCommand,
    },
    hist: {
      command: histCommand,
    },
    rm: {
      command: rmCommand,
    },
    mkfile: {
      command: mkfileCommand,
    },
    mkdir: {
      command: mkdirCommand,
    },
    rmalias: {
      command: removeAliasCommand,
    },
    term: {
      command: viewScriptTerminal,
    },
    notes: {
      command: viewNotes,
    },
    emailit: {
      command: viewEmailIt,
    },
    scriptline: {
      command: viewScriptLine,
    },
  };
  let mode = "insert";
  let wd = "~";
  let esc = String.fromCharCode(0x1b);
  let lastData = {
    line: "",
    data: [],
    valid: false,
  };
  let lcommandRow = 0;
  let lcommandCol = 0;
  let termAtb = {
    red: `${esc}[31m`,
    black: `${esc}[30m`,
    green: `${esc}[32m`,
    orange: `${esc}[33m`,
    blue: `${esc}[34m`,
    magenta: `${esc}[35m`,
    cyan: `${esc}[36m`,
    gray: `${esc}[37m`,
    default: `${esc}[39m`,
    up: `${esc}[A`,
    down: `${esc}[B`,
    left: `${esc}[D`,
    right: `${esc}[C`,
  };
  let homeDir = "";

  onMount(async () => {
    homeDir = await App.GetHomeDir();
    wd = homeDir;
    term = new Terminal({
      rendererType: "canvas",
      convertEol: true,
      cursorBlink: true,
      cursorStyle: "bar",
      allowProposedApi: true,
      theme: {
        background: $theme.textAreaColor,
        black: $theme.backgroundColor,
        blue: $theme.Cyan,
        brightBlack: $theme.backgroundColor,
        brightBlue: $theme.Cyan,
        brightCyan: $theme.Cyan,
        brightGreen: $theme.Green,
        brightMagenta: $theme.Pink,
        brightRed: $theme.Red,
        brightWhite: $theme.textColor,
        brightYellow: $theme.Yellow,
        cursor: $theme.textColor,
        cursorAccent: $theme.highlightBackgroundColor,
        cyan: $theme.Cyan,
        foreground: $theme.textColor,
        green: $theme.Green,
        magenta: $theme.Pink,
        red: $theme.Red,
        selection: $theme.highlightBackgroundColor,
        white: $theme.textColor,
        yellow: $theme.Yellow,
      },
    });
    const webLinksAddon = new WebLinksAddon();
    const fitAddon = new FitAddon();
    term.loadAddon(webLinksAddon);
    term.loadAddon(fitAddon);

    term.open(document.getElementById("terminal"));
    fitAddon.fit();

    term.write(" Welcome to Script Terminal!\n");
    term.prompt = () => {
      term.write(" $ ");
    };
    term.prompt();

    term.onKey(({ key, domEvent }) => {
      if (mode === "insert") {
        //
        // This is the insert mode keymappings.
        //
        if (domEvent.key === "Enter") {
          //
          // If it is an enter key, process the line
          //
          if (curr_line) {
            ProcessLine(curr_line);
          }
        } else if (domEvent.key === "Backspace") {
          //
          // If it is a backspace, delete the last character form the line.
          //
          if (curr_line) {
            if (curr_line.length >= 1) {
              curr_line = curr_line.slice(0, curr_line.length - 1);
              term.write("\b \b");
            }
          }
        } else if (
          domEvent.key === "ArrowUp" ||
          domEvent.key === "ArrowDown" ||
          domEvent.key === "ArrowLeft" ||
          domEvent.key === "ArrowRight"
        ) {
          //
          // Only do something when in command mode.
          //
        } else if (domEvent.key === "Tab") {
          //
          // If it is a tab, add two spaces to the line and terminal.
          //
          term.write("  ");
          curr_line += "  ";
        } else if (domEvent.key === "Escape") {
          //
          // Set toggle the mode.
          //
          if (mode === "insert") {
            if (lastData.valid === true) {
              mode = "command";
              lcommandRow = lastData.data.length - 1;
              lcommandCol = 0;
              term.write(termAtb.up);
              term.write(termAtb.right);
            }
          } else {
            mode = "insert";
          }
        } else if (domEvent.key === "l" && domEvent.ctrlKey) {
          term.clear();
        } else {
          //
          // Every other key, add to the line and terminal.
          //
          curr_line += key;
          term.write(key);
        }
      } else {
        //
        // This is the command mode keymappings.
        //
        let depth = 0;
        switch (key) {
          case "i":
            depth = lastData.data.length - lcommandRow;
            for (var i = 0; i < depth; i++) {
              term.write(termAtb.down);
            }
            for (i = 0; i < lcommandCol; i++) {
              term.write(termAtb.left);
            }
            term.write(termAtb.left);
            mode = "insert";
            break;
          case "k":
            if (lcommandRow !== 0) {
              lcommandRow -= 1;
              term.write(termAtb.up);
            }
            break;
          case "j":
            if (lcommandRow < lastData.data.length - 1) {
              term.write(termAtb.down);
              lcommandRow += 1;
            }
            break;
          case "l":
            lcommandCol += 1;
            term.write(termAtb.right);
            break;
          case "h":
            if (lcommandCol !== 0) {
              term.write(termAtb.left);
              lcommandCol -= 1;
            }
            break;
          case "r":
            //
            // first, go back to the command line.
            //
            depth = lastData.data.length - lcommandRow;
            for (var i = 0; i < depth; i++) {
              term.write(termAtb.down);
            }
            for (i = 0; i < lcommandCol; i++) {
              term.write(termAtb.left);
            }
            term.write(termAtb.left);
            mode = "insert";

            //
            // Add the command to the command line.
            //
            term.write(`${lastData.data[lcommandRow].command}\n\r`);

            //
            // Get the command to run and run it.
            //
            ProcessLine(lastData.data[lcommandRow].command);
            break;
          default:
            break;
        }
      }
    });

    //
    // Load the aliases.
    //
    loadAliases();

    //
    // Make sure the terminal is focused.
    //
    term.focus();
  });

  function splitAt(index, xs) {
    return [xs.slice(0, index), xs.slice(index)];
  }

  async function ProcessLine(text) {
    //
    // Get the words of the line.
    //
    var words = text.trim().split(" ");
    lastData.line = text;
    lastData.valid = false;
    curr_line = "";

    //
    // See if the first word is a valid command.
    //
    let scrpt = $termscripts.filter((item) => item.name === words[0]);
    if (scrpt.length === 0) {
      let alias = $aliases.filter((item) => item.name === words[0]);
      if (alias.length === 0) {
        //
        // Command not found. Print the error and give a new prompt.
        //
        term.write(
          `\r\n\r\n    ${termAtb.red}<Error>${termAtb.default} The command "${words[0]}" wasn't found!\r\n\r\n`
        );
        term.prompt();
      } else {
        //
        // run the aliases commands.
        //
        var parts = alias[0].line.split(";");
        for (const item of parts) {
          await ProcessLine(item);
        }
      }
    } else {
      //
      // Command found. Run it!
      //
      await runCommandLineScripts(scrpt[0], words.slice(1).join(" "));
    }
  }

  async function runCommandLineScripts(scrpt, args) {
    //
    // Command found. Run it! This is the script command running.
    //
    let result = await $runscript(scrpt.name, args);
    ProcessScriptReturn(result);
  }

  function ProcessScriptReturn(data) {
    //
    // Give the carrage return before the data.
    //
    term.write("\r\n");

    //
    // Process the JSON structure given to perform the command. The JSON structure
    // should be:
    //
    //  {
    //    tcommand: <terminal command to run>
    //    lines: [{
    //      text: <text to display>,
    //      color: <color to show>,
    //      command: <command line string to run>
    //    }, {
    //      <next line structure>
    //    }, ...]
    //  }
    //
    data = JSON.parse(data);
    lastData.data = data.lines;
    lastData.valid = true;

    //
    // See if there is a terminal command to run.
    //
    if (typeof data.tcommand !== "undefined" && data.tcommand.length > 0) {
      //
      // We have a terminal command to run.
      //
      RunTerminalCommand(data.tcommand);
    }

    //
    // Process each line returned.
    //
    data.lines.forEach((line) => {
      //
      // Write the line text with the color change.
      //
      term.write(
        `    ${termAtb[line.color]}${line.text}${termAtb.default}\r\n`
      );
    });
  }

  async function RunTerminalCommand(text) {
    //
    // Keep a list of valid commands.
    //
    commands.push(lastData.line);

    //
    // Get the command and it's arguments separated. There can by any number of commands separated by a semicolon.
    //
    let com = text.split(" ");
    if (com[0].trim() !== "alias") {
      let parts = text.split(";");
      for (let i = 0; i < parts.length; i++) {
        let words = parts[i].split(" ");
        if (words.length > 0) {
          await tCommands[words[0]].command(words.slice(1).join(" "));
        }
      }
    } else {
      //
      // this is the alias command. Run it without decomposing the line.
      //
      await tCommands[com[0]].command(com.slice(1).join(" "));
    }
    term.prompt();
  }

  async function cdCommand(text) {
    //
    // Remove inclosing quotes
    //
    while (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text.length > 0) {
      let path = new String(text);
      if (text[0] !== "/") {
        let ndir = new String(text);
        let nwd = new String(wd);
        path = await App.AppendPath(nwd, ndir);
      }
      let exists = await App.DirExists(path);
      if (exists) {
        wd = path;
      } else {
        term.write(`    <Error> The path "${path} doesn't exist!\n\n`);
      }
    } else {
      wd = homeDir;
    }
    lastData.valid = false;
  }

  async function helpCommand(text) {
    text = text.trim().split(" ")[0];
    if (text.length === 0) {
      //
      // show the commands available.
      //
      $termscripts.forEach((item) => {
        //
        // Make sure the description lines are not too long.
        //
        let description = truncateLines(item.description);
        term.write(`    ${item.name}    ${description}\n\r`);
      });
      //
      // Show the aliases if any.
      //
      if ($aliases.length > 0) {
        //
        // We have aliases. Show them.
        //
        term.write(`\n\r    Aliases:\n\r\n\r`);
        $aliases.forEach((item) => {
          term.write(`    ${item.name}="${item.line}"\n\r`);
        });
      }
    } else {
      //
      // show the help string for the command.
      //
      let spt = $termscripts.find((item) => item.name === text);
      if (spt === "undefined") {
        term.write(
          `\n\r    ${termAtb.red}<Error>${termAtb.default} ${text} is an invalid Command.\n\r`
        );
      } else {
        let help = truncateLines(spt.help);
        term.write(`    ${spt.name}  -  ${help}\n\r`);
      }
    }
    lastData.valid = false;
  }

  function truncateLines(text) {
    var description = [text];
    let index = 80;
    let subin = 0;
    while (description[subin].length > index) {
      while (description[subin][index] !== " ") index--;
      let nsub = splitAt(index, description[subin]);
      if (subin === 0) {
        description = nsub;
      } else {
        description = description.slice(0, subin).concat(nsub);
      }
      subin++;
      index = 80;
    }
    if (subin === 0) {
      description = description[0];
    } else {
      description = description.join("\n\r           ");
    }
    return description;
  }

  async function lsCommand(text) {
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = new String(text.trim());
    var path = new String(wd);
    if (text !== "") {
      if (text[0] === "/") {
        path = text;
      } else {
        path = await App.AppendPath(path, text);
      }
    }
    var dirReal = await App.DirExists(path);
    if (dirReal) {
      var result = await App.ReadDir(path);
      var lines = [];
      for (let i = 0; i < result.length; i++) {
        //
        // Rewrite lastData.lines to have a tcommand for each entry printed.
        //
        let item = result[i];
        let npath = await App.AppendPath(item.Dir, item.Name);
        dirReal = await App.DirExists(npath);
        if (dirReal) {
          lines.push({
            name: item.Name,
            command: `cd '${npath}'`,
          });
        } else {
          lines.push({
            name: item.Name,
            command: `open '${npath}'`,
          });
        }

        //
        // Print the item name.
        //
        term.write(`    ${item.Name}\n\r`);
      }
      lastData.data = lines;
      lastData.valid = true;
    } else {
      let fileReal = await App.FileExists(path);
      if (fileReal) {
        term.write(`    ${text}\n\r`);
        lastData.data = [
          {
            name: text,
            tcommand: `ls '${path}'`,
          },
        ];
      } else {
        term.write(
          `\n\r    ${termAtb.red}<Error>${termAtb.default} ${path} is an invalid Directory.\n\r`
        );
      }
    }
  }

  async function openCommand(text) {
    //
    // Fix up the file path given to the open command.
    //
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text[0] !== "/") {
      text = await App.AppendPath(wd, text);
    }
    text = new String(text.trim());

    //
    // Run the open command on the file.
    //
    await runCommandLine(`/usr/bin/open -t ${text}`, [], () => {}, wd);
    lastData.valid = false;
  }

  async function runscriptCommand(text) {
    //
    // Run a script on a file or text. Make sure enough arguments were given. It requires two arguments separated by a comma:
    // the script name and the file name or text to be processed.
    //
    lastData.valid = false;
    var scriptName = text.split(",");
    if (scriptName.length < 2) {
      term.write(
        `\n\r    ${termAtb.red}<Error>${termAtb.default} runscript wasn't given enough arguments.\n\r`
      );
    } else {
      //
      // Get the script name and file name or text separated.
      //
      text = scriptName.splice(1).join(",").trim();
      scriptName = scriptName[0].trim();

      //
      // Fix up the file path given to the open command.
      //
      if (text[0] === '"' || text[0] === "'") {
        text = text.slice(1, text.length - 1);
      }

      let textfile = await App.AppendPath(wd, text);
      if (await App.FileExists(textfile)) {
        text = textfile;
      }
      text = new String(text.trim());

      //
      // Process the text based on what it is.
      //
      let newText = await $runscript(scriptName, text);
      term.write(`\n\r     ${newText}\n\r`);
    }
    lastData.valid = false;
  }

  async function editCommand(text) {
    //
    // fix the file name.
    //
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text[0] !== "/") {
      text = await App.AppendPath(wd, text);
    }
    text = new String(text.trim());

    //
    // Setup the user editor data file.
    //
    let userEditor = await App.AppendPath(homeDir, ".myeditorchoice");
    if (!(await App.FileExists(userEditor))) {
      //
      // They don't have this file setup. Open in the system's default editor.  TODO: Not usable on non-macOS systems.
      //
      await runCommandLine(`/usr/bin/open '${text}'`, [], () => {}, wd);
    } else {
      //
      // Get the user editor.
      //
      var editor = await App.ReadFile(userEditor);
      editor = editor.toString().trim();
      if (editor.endsWith(".app")) {
        //
        // Open the file with a program. TODO: Not usable on non-macOS systems.
        //
        await runCommandLine(
          `/usr/bin/open -a ${editor} '${text}'`,
          [],
          () => {},
          wd
        );
      } else {
        //
        // It is a command line editor. Open specially.
        //
        if (editor === "emacs") {
          //
          // Open emacs.
          //
          await runCommandLine(
            'emacsclient -n -q "' + text + '"',
            [],
            (err, result) => {},
            wd
          );
        }
      }
    }
    lastData.valid = false;
  }

  async function aliasCommand(text) {
    //
    // Add to the Aliases.
    //
    let parts = text.trim().split("=");
    if (parts.length > 1) {
      //
      // Add to the aliases array.
      //
      let ln = parts[1];
      if (ln[0] === '"' || ln[0] === "'") {
        ln = text.slice(1, ln.length - 1);
      }
      let aliasname = parts[0].trim();
      $aliases = $aliases.filter((item) => item.name !== aliasname);
      $aliases.push({
        name: aliasname,
        line: ln,
      });

      //
      // Save the new alias.
      //
      await saveAliases();
    } else if (text.trim() === "") {
      //
      // List the aliases.
      //
      term.write("   Aliases:\n\r");
      for (const item of $aliases) {
        term.write(`    ${item.name} = "${item.line}"\n\r`);
      }
    } else {
      term.write(
        `\n\r    ${termAtb.red}<Error>${termAtb.default} Not enough parameters for an alias.\n\r`
      );
    }
    lastData.valid = false;
  }

  async function loadAliases() {
    let userAliases = await App.AppendPath(homeDir, ".myaliases");
    if (await App.FileExists(userAliases)) {
      let tmp = await App.ReadFile(userAliases);
      $aliases = JSON.parse(tmp);
    }
  }

  async function removeAliasCommand(name) {
    $aliases = $aliases.filter((item) => item.name !== name);
    await saveAliases();
  }

  async function saveAliases() {
    let userAliases = await App.AppendPath(homeDir, ".myaliases");
    await App.WriteFile(userAliases, JSON.stringify($aliases));
  }

  async function histCommand(text) {
    //
    // Get needed variables ready.
    //
    let lines = [];
    let depth = 5;

    //
    // See if we were given a depth to show.
    //
    text = parseInt(text.trim());
    if (Number.isInteger(text)) {
      depth = text;
    }

    //
    // Make sure we have that many commands to display.
    //
    if (depth > commands.length - 1) {
      depth = commands.length - 1;
    }

    //
    // Display the command and create the command for it. Don't show the last
    // command as it will be the history command.
    //
    for (let i = commands.length - (depth + 1); i < commands.length - 1; i++) {
      term.write(`    ${commands[i]}\n\r`);
      lines.push({
        name: commands[i],
        command: commands[i],
      });
    }
    lastData.data = lines;
    lastData.valid = true;
  }

  async function rmCommand(text) {
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd);
    if (text !== "") {
      textblank = false;
      text = new String(text);
      if (text[0] === "/") {
        path = text;
      } else {
        path = await App.AppendPath(path, text);
      }
    } else {
      textblank = true;
    }
    var dirReal = await App.DirExists(path);
    if (dirReal && textblank) {
      var result = await App.ReadDir(path);
      var lines = [];
      for (let i = 0; i < result.length; i++) {
        //
        // Rewrite lastData.lines to have a tcommand for each entry printed.
        //
        let item = result[i];
        let npath = await App.AppendPath(item.Dir, item.Name);
        lines.push({
          name: item.Name,
          command: `rm '${npath}'`,
        });

        //
        // Print the item name.
        //
        term.write(`    ${item.Name}\n\r`);
      }
      lastData.data = lines;
      lastData.valid = true;
    } else {
      let fileReal = await App.FileExists(path);
      if (fileReal || dirReal) {
        //
        // Remove the file or directory.
        //
        await App.DeleteEntries(path);
      } else {
        term.write(
          `\n\r    ${termAtb.red}<Error>${termAtb.default} ${path} is an invalid Directory.\n\r`
        );
      }
    }
  }

  async function mkdirCommand(text) {
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd);
    if (text !== "") {
      textblank = false;
      text = new String(text);
      if (text[0] === "/") {
        path = text;
      } else {
        path = await App.AppendPath(path, text);
      }
    } else {
      textblank = true;
    }
    var dirReal = await App.DirExists(path);
    if (!dirReal && !textblank) {
      //
      // Create the directory.
      //
      await App.MakeDir(path);
    } else {
      //
      // It already exists.
      //
      term.write(
        `\r\n\r\n    ${termAtb.red}<Error>${termAtb.default} The directory "${path}" already exists!\r\n\r\n`
      );
      term.prompt();
    }
  }

  async function mkfileCommand(text) {
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd);
    if (text !== "") {
      textblank = false;
      text = new String(text);
      if (text[0] === "/") {
        path = text;
      } else {
        path = await App.AppendPath(path, text);
      }
    } else {
      textblank = true;
    }
    var fileReal = await App.FileExists(path);
    if (!fileReal && !textblank) {
      //
      // Create the file.
      //
      await App.MakeFile(path);
    } else {
      //
      // It already exists.
      //
      term.write(
        `\r\n\r\n    ${termAtb.red}<Error>${termAtb.default} The file "${path}" already exists!\r\n\r\n`
      );
      term.prompt();
    }
  }

  function viewEmailIt() {
    $state = "emailit";
  }

  function viewNotes() {
    $state = "notes";
  }

  function viewScriptEditor() {
    $state = "scripts";
  }

  function viewScriptTerminal() {
    $state = "scriptterm";
  }

  function viewScriptLine() {
    $state = "scriptline";
  }

  async function runCommandLine(line, rEnv, callback, dir) {
    //
    // Get the environment to use.
    //
    let envlistloc = await App.AppendPath(
      $config.configDir,
      "environments.json"
    );
    let envlist = await App.ReadFile(envlistloc);
    envlist = JSON.parse(envlist);
    let nEnv = envlist.map((item) => "Default");
    if (typeof rEnv !== "undefined") {
      nEnv = { ...nEnv, ...rEnv };
    }

    //
    // Fix the environment from a map to an array of strings.
    //
    var penv = [];
    for (const key in nEnv) {
      penv.push(`${key}=${nEnv[key]}`);
    }

    //
    // Make sure dir has a value.
    //
    if (typeof dir === "undefined") dir = ".";

    //
    // Run the command line in a shell. #TODO: make the shell configurable.
    //
    var args = [$config.shell, "-c", line];
    var cmd = $config.shell;

    //
    // Run the command line.
    //
    var result = await App.RunCommandLine(cmd, args, penv, dir);
    var err = await App.GetError();
    //
    // If callback is given, call it with the results.
    //
    if (typeof callback !== "undefined" || callback !== null) {
      callback(err, result);
    }
  }
</script>

<div
  id="ScriptTermDiv"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  <div
    id="terminal"
    style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
  />
  <div
    id="statusline"
    style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
  >
    <span
      id="modeIndicator"
      style="background-color : {mode === 'insert'
        ? $theme.Cyan
        : $theme.Purple}; color: {$theme.backgroundColor}">{mode}</span
    >
    <span id="workingdir">
      WD: {wd}
    </span>
  </div>
  <div id="buttonRow">
    <button
      on:click={viewEmailIt}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      EmailIt
    </button>
    <button
      on:click={viewNotes}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Notes
    </button>
    <button
      on:click={viewScriptEditor}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Edit Scripts
    </button>
  </div>
</div>

<style>
  #ScriptTermDiv {
    display: flex;
    flex-direction: column;
    padding: 0px;
    margin: 0px;
    border: 0px transparent;
    border-radius: 0px 0px 10px 10px;
  }

  #terminal {
    display: flex;
    flex-direction: column;
    margin: 0px 10px 0px 10px;
    width: 1000px;
    height: 470px;
  }

  #statusline {
    display: flex;
    flex-direction: row;
    height: 30px;
    width: 1000px;
    margin: 10px;
  }

  #modeIndicator {
    padding: 3px 10px;
    border-radius: 5px;
    width: 94px;
  }

  #workingdir {
    padding: 3px 20px;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    margin: 13px auto;
  }

  #buttonRow button {
    border-radius: 10px;
    padding: 5px 20px 5px 20px;
    margin: 0px 20px;
    width: 100%;
    max-height: 40px;
    height: 40px;
    width: auto;
    cursor: pointer;
  }
</style>
