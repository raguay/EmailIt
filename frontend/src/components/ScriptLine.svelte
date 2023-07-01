<script>
  import { onMount, afterUpdate, tick } from "svelte";
  import { termscripts } from "../stores/termscripts.js";
  import { runscript } from "../stores/runscript.js";
  import { aliases } from "../stores/aliases.js";
  import { state } from "../stores/state.js";
  import { commands } from "../stores/commands.js";
  import { config } from "../stores/config.js";
  import { theme } from "../stores/theme.js";
  import { sp } from "../stores/sp.js";
  import * as App from "../../wailsjs/go/main/App.js"; // My application API
  import * as rt from "../../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  let showHtmlDiv = false;
  let showOutputDiv = false;
  let showHint = false;
  let hints = [];
  let hintCursor = 0;
  let scrollDiv = null;
  let showError = false;
  let errorOutput = "";
  let textOutput = [];
  let htmlOutput = [];
  let inputdiv = null;
  let cursorDiv = null;
  let containerDiv = null;
  let AltAdj = "";
  let InterActive = true;
  let inputState = 0;
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
    editscript: {
      command: viewScriptEditor,
    },
    quit: {
      command: quit,
    },
  };
  let wd = "~";
  let lastData = {
    line: "",
    data: [],
    valid: false,
  };
  let homeDir = "";
  let outputDivInput;
  let htmlDivInput;
  let currentLine = 0;
  let inputTimer = null;
  let oldheight = 0;

  onMount(async () => {
    homeDir = await App.GetHomeDir();
    wd = homeDir;

    //
    // Load the aliases.
    //
    loadAliases();
  });

  afterUpdate(async () => {
    await tick();
    FocusInput();
    await AdjustHeight();
  });

  async function AdjustHeight() {
    if (containerDiv !== null) {
      let height = containerDiv.clientHeight + 20; // Have to add the height of the dragbar.
      let width = containerDiv.clientWidth;
      if (height !== oldheight) {
        await rt.WindowSetSize(width, height);
        oldheight = height;
      }
    }
  }

  function FocusInput() {
    if (inputdiv !== null) {
      if (inputTimer !== null) clearTimeout(inputTimer);
      if (!showOutputDiv && !showHtmlDiv) inputdiv.focus();
      else if (showOutputDiv) {
        outputDivInput.focus();
      } else {
        htmlDivInput.focus();
      }
      inputTimer = setTimeout(FocusInput, 100);
    }
  }

  function splitAt(index, xs) {
    return [xs.slice(0, index), xs.slice(index)];
  }

  function showErrorHTML(msg) {
    //
    // Display an error message to the user.
    //
    showError = true;
    errorOutput = `<span style="color: ${$theme.Red};">${msg}</span>`;
  }

  function showOutputText(data, clear) {
    //
    // Show the information from running a command. This goes to the main output
    // flow.
    //
    showOutputDiv = true;
    if (clear) textOutput = [];
    if (data !== "") {
      textOutput.push(data.toString());
      textOutput = textOutput;
    }
  }

  function showOutputHTML(data, clear) {
    //
    // Add content to the secodary output area. This will be HTML formated already.
    //
    showHtmlDiv = true;
    if (clear) htmlOutput = [];
    htmlOutput.push(data);
  }

  function lineInput() {
    //
    // Get the information from the line input and see if we can expand it to a
    // valid command.
    //
    let curin = inputdiv.value.trim();
    switch (inputState) {
      case 0:
        //
        // This state is for finding the right command.
        //
        if (curin.length === 0) {
          showHint = false;
        } else {
          hints = $aliases
            .map((ele) => ele.name)
            .filter((ele) => ele.includes(curin))
            .concat(
              $termscripts
                .map((ele) => ele.name)
                .filter((ele) => ele.includes(curin))
            )
            .sort();
          if (hints.length > 0) {
            showHint = true;
            hintCursor = 0;
            if (hints.length === 1) {
              if (typeof hints[hintCursor] !== "undefined") {
                inputdiv.value = hints[hintCursor];
                showHint = false;
                inputState = 1;
              }
            }
          }
        }
        break;
      case 1:
        //
        // This state is for sub command strings. TODO: Figure out how to implement.
        //
        break;
    }
  }

  async function ProcessLine(text) {
    //
    // clear out the output items.
    //
    showHtmlDiv = false;
    showOutputDiv = false;
    showError = false;
    htmlOutput = [];
    textOutput = [];
    InterActive = true;
    $sp.data.previousLines = [];
    $sp.data.registers = [];
    $sp.data.line = "";
    $sp.data.result = "";
    inputState = 0; // An input was actioned. Make sure the state is back to 0.

    text = text.trim();
    let chains = text.split(").");
    if (chains.length === 1 && !text.includes("(")) {
      var commandparts = text.split(";");
      for (var i = 0; i < commandparts.length; i++) {
        //
        // It is a command line style. Get the words of the line.
        //
        var words = commandparts[i].split(" ");
        lastData.line = commandparts[i];
        lastData.valid = false;

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
            showErrorHTML(`Command '${words[0]}' doesn't exist!`);
            break;
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
          // Command found. Run it! This is the script command running.
          //
          $sp.data.previousLines.push($sp.data.line);
          $sp.data.line = commandparts[i];

          //
          // Keep a list of valid commands.
          //
          $commands.push($sp.data.line);

          //
          // Run the command.
          //
          $sp.data.result = await $runscript(
            scrpt[0].name,
            words.slice(1).join(" ")
          );

          //
          // Process the results.
          //
          ProcessScriptReturn($sp.data.result);
        }
      }
    } else {
      //
      // It is a chaining command style.
      //
      for (let i = 0; i < chains.length; i++) {
        if (chains[i][chains[i].length - 1] === ")") {
          //
          // Remove the closing parenthesis.
          //
          chains[i] = chains[i].slice(0, -1);
        }
        let parts = chains[i].split("(");
        let scrpt = $termscripts.filter((item) => item.name === parts[0]);
        if (scrpt.length > 0) {
          //
          // It's a valid command. Let run it.
          //
          let args = parts[1].replaceAll(",", " ");
          $sp.data.previousLines.push($sp.data.line);
          $sp.data.line = chains[i] + ")";
          //
          // Keep a list of valid commands.
          //
          $commands.push($sp.data.line);

          //
          // Run the command.
          //
          $sp.data.result = await $runscript(scrpt[0].name, args);

          //
          // Process the results.
          //
          ProcessScriptReturn($sp.data.result);
        } else {
          //
          // An invalid command throws the whole chain out.
          //
          showHtmlDiv = false;
          showOutputDiv = false;
          showError = false;
          htmlOutput = [];
          textOutput = [];
          InterActive = false;
          $sp.data.previousLines = [];
          $sp.data.registers = [];
          $sp.data.line = "";
          $sp.data.result = "";
          showErrorHTML(`Command '${parts[0]}' doesn't exist!`);
          break;
        }
      }
    }
  }

  function ProcessScriptReturn(data) {
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
    currentLine = 0;
    data = JSON.parse(data);
    lastData.data = data.lines;
    lastData.valid = true;
    InterActive = true;

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
      showOutputText(
        `<p style="margin: 5px;"><span style="text: ${line.color};">${line.text}</span></p>`,
        false
      );
    });
  }

  async function RunTerminalCommand(text) {
    //
    // Get the command and it's arguments separated. There can by any number of commands separated by a semicolon.
    //
    let com = text.split(" ");
    let parts = text.split(";");
    for (let i = 0; i < parts.length; i++) {
      let words = parts[i].split(" ");
      if (words.length > 0) {
        if (typeof tCommands[words[0]] !== "undefined") {
          //
          // It's a valid command. Run it.
          //
          await tCommands[words[0]].command(words.slice(1).join(" "));
        } else {
          //
          // It's not a command. Tell the user.
          //
          showErrorHTML(`Terminal Command ${words[0]} doesn't exist!`);
        }
      }
    }
  }

  async function cdCommand(text) {
    InterActive = false;
    //
    // Remove inclosing quotes
    //
    while (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text.length > 0) {
      let path = new String(text).toString();
      if (text[0] !== "/") {
        let ndir = new String(text).toString();
        let nwd = new String(wd).toString();
        path = await App.AppendPath(nwd, ndir);
      }
      let exists = await App.DirExists(path);
      if (exists) {
        wd = path;
      } else {
        //
        // Error message
        //
        showErrorHTML(`The Path doesn't exist!`);
      }
    } else {
      wd = homeDir;
    }
    lastData.valid = false;
  }

  async function helpCommand(text) {
    InterActive = false;
    text = text.trim().split(" ")[0];
    if (text.length === 0) {
      //
      // show the commands available.
      //
      showOutputText("<h2>Available Commands</h2><table>", true);
      $termscripts.forEach((item) => {
        //
        // Make sure the description lines are not too long.
        //
        showOutputText(
          `<tr><td>${item.name}</td><td>${item.description}</td></tr>`,
          false
        );
      });
      showOutputText("</table>", false);
      //
      // Show the aliases if any.
      //
      if ($aliases.length > 0) {
        //
        // We have aliases. Show them.
        //
        showOutputText("<h2>Available Aliases</h2><table>", false);
        $aliases.forEach((item) => {
          //
          // Show the aliases
          //
          showOutputText(
            `<tr><td>${item.name}</td><td>${item.line}</td></tr>`,
            false
          );
        });
        showOutputText("</table>", false);
      }
    } else {
      //
      // show the help string for the command.
      //
      let spt = $termscripts.find((item) => item.name === text);
      if (spt === "undefined") {
        //
        // Invalid command.aliases
        //
        showErrorHTML("Command not found!");
      } else {
        showOutputText(
          `<table><tr><td>${spt.name}</td><td>${spt.help}</td></tr></table>`,
          true
        );
      }
    }
    lastData.valid = false;
  }

  async function lsCommand(text) {
    InterActive = true;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = new String(text.trim()).toString();
    var path = new String(wd).toString();
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
      showOutputText("", true);
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
        showOutputText(`${item.Name}`, false);
      }
      lastData.data = lines;
      lastData.valid = true;
    } else {
      let fileReal = await App.FileExists(path);
      if (fileReal) {
        lastData.data = [
          {
            name: text,
            tcommand: `ls '${path}'`,
          },
        ];
        showOutputText(`${text}`, true);
      } else {
        showErrorHTML("Not a valid Path.");
      }
    }
  }

  async function openCommand(text) {
    InterActive = false;
    //
    // Fix up the file path given to the open command.
    //
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text[0] !== "/") {
      text = await App.AppendPath(wd, text);
    }
    text = new String(text.trim()).toString();

    //
    // Run the open command on the file.
    //
    await runCommandLine(`/usr/bin/open -t ${text}`, [], () => {}, wd);
    lastData.valid = false;
  }

  async function runscriptCommand(text) {
    InterActive = false;
    //
    // Run a script on a file or text. Make sure enough arguments were given. It requires two arguments separated by a comma:
    // the script name and the file name or text to be processed.
    //
    lastData.valid = false;
    var scriptName = text.split(",");
    if (scriptName.length < 2) {
      showErrorHTML("runscript wasn't given enough parameters.");
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
      text = new String(text.trim()).toString();

      //
      // Process the text based on what it is.
      //
      let newText = await $runscript(scriptName, text);
      showOutputText(`<p style="margin: 5px;">${newText}</p>`, true);
    }
    lastData.valid = false;
  }

  async function editCommand(text) {
    InterActive = false;
    //
    // fix the file name.
    //
    text = new String(text.trim()).toString();
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text[0] !== "/") {
      text = await App.AppendPath(wd, text);
    }

    //
    // Setup the user editor data file.
    //
    let userEditor = await App.AppendPath(homeDir, ".myeditorchoice");
    let editorExists = await App.FileExists(userEditor);
    if (!editorExists) {
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
    InterActive = false;
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
      showOutputText("<h2>Available Aliases</h2>", true);
      for (const item of $aliases) {
        //
        // Show each alias
        //
        showOutputText(`<p style="margin: 5px;">${item.name}</p>`, false);
      }
    } else {
      //
      // Error: not enough given
      //
      showErrorHTML("Not enough information is given!");
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
    InterActive = true;
    //
    // Get needed variables ready.
    //
    let lines = [];
    let depth = 999; // Most likely will not get that deep.

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
    if (depth > $commands.length - 1) {
      depth = $commands.length - 1;
    }

    //
    // Display the command and create the command for it. Don't show the last
    // command as it will be the history command.
    //
    for (
      let i = $commands.length - 1;
      i > $commands.length - (depth + 1);
      i--
    ) {
      showOutputText(`${$commands[i]}`, false);
      lines.push({
        name: $commands[i],
        command: $commands[i],
      });
    }
    lastData.data = lines;
    lastData.valid = true;
  }

  async function rmCommand(text) {
    InterActive = false;
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd).toString();
    if (text !== "") {
      textblank = false;
      text = new String(text).toString();
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
      InterActive = true;
      var result = await App.ReadDir(path);
      var lines = [];
      showOutputText("", true);
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
        showOutputText(`${item.Name}`, false);
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
        showErrorHTML(`The path "${path}" doesn't exist!`);
      }
    }
  }

  async function mkdirCommand(text) {
    InterActive = false;
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd).toString();
    if (text !== "") {
      textblank = false;
      text = new String(text).toString();
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
      showErrorHTML(`The path "${path}" already exists.`);
    }
  }

  async function mkfileCommand(text) {
    InterActive = false;
    let textblank = false;
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    text = text.trim();
    var path = new String(wd).toString();
    if (text !== "") {
      textblank = false;
      text = new String(text).toString();
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
      showErrorHTML(`The file "${path}" already exists!`);
    }
  }

  async function viewEmailIt() {
    //
    // Reset the window size to normal.
    //
    await rt.WindowSetSize($config.width, $config.height);
    $state = "emailit";
  }

  async function viewNotes() {
    //
    // Reset the window size to normal.
    //
    await rt.WindowSetSize($config.width, $config.height);
    $state = "notes";
  }

  async function viewScriptEditor() {
    //
    // Reset the window size to normal.
    //
    await rt.WindowSetSize($config.width, $config.height);
    $state = "scripts";
  }

  async function viewScriptTerminal() {
    //
    // Reset the window size to normal.
    //
    await rt.WindowSetSize($config.width, $config.height);
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
    // Run the command line in a shell.
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

  function processKey(e) {
    switch (e.key) {
      case "Enter":
        //
        // The user hit enter, so process the line.
        //
        e.stopPropagation();
        e.preventDefault();
        showHint = false;
        inputState = 0;
        ProcessLine(e.target.value);
        break;

      case "Tab":
        //
        // Expand the suggestion.
        //
        if (typeof hints[hintCursor] !== "undefined")
          inputdiv.value = hints[hintCursor];
        showHint = false;
        hintCursor = 0;
        inputState = 1;
        e.stopPropagation();
        e.preventDefault();
        break;

      case "Escape":
        //
        // Reset everything.
        //
        showHint = false;
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        inputState = 0;
        e.stopPropagation();
        e.preventDefault();
        break;

      case "ArrowUp":
        if (showHint) {
          hintCursor = hintCursor - 1;
          if (hintCursor < 0) hintCursor = 0;
        }
        break;

      case "ArrowDown":
        if (showHint) {
          hintCursor = hintCursor + 1;
          if (hintCursor >= hints.length) hintCursor = hints.length - 1;
        }
        break;
    }
  }

  function updateScroll(scrollDiv, amount) {
    if (scrollDiv !== null) {
      scrollDiv.scrollTop += amount;
      if (scrollDiv.scrollTop < 0) scrollDiv.scrollTop = 0;
    }
  }

  function checkInView(container, element) {
    //Get container properties
    let Offset = 0;
    let cTop = container.scrollTop + 140;
    let cBottom = cTop + container.clientHeight;
    cTop += 5; // take into account the padding on the top. Not needed for the bottom calculations.

    //Get element properties
    let eTop = element.offsetTop;
    let eBottom = eTop + element.clientHeight;

    //Check if in view
    let isTotal = eTop >= cTop && eBottom <= cBottom;
    let isPartial =
      (eTop < cTop && eBottom > cTop) || (eBottom > cBottom && eTop < cBottom);

    if (eBottom > cBottom) {
      Offset = eBottom - cBottom;
    } else if (eTop < cTop) {
      Offset = eTop - cTop;
    }

    //Return outcome
    return {
      cTop: cTop,
      cBottom: cBottom,
      eTop: eTop,
      eBottom: eBottom,
      istotal: isTotal,
      ispartial: isPartial,
      offset: Offset,
    };
  }

  async function checkPosition() {
    await tick();
    await tick();
    let pos = checkInView(scrollDiv, cursorDiv);
    if (!pos.istotal) updateScroll(scrollDiv, pos.offset);
  }

  async function outputKeyProcess(e) {
    switch (e.key) {
      case "0":
      case "1":
      case "2":
      case "3":
      case "4":
      case "5":
      case "6":
      case "7":
      case "8":
      case "9":
        AltAdj = `${AltAdj}${e.key}`;
        break;

      case "ArrowUp":
      case "k":
        e.stopPropagation();
        if (AltAdj !== "") {
          currentLine = currentLine - parseInt(AltAdj);
        } else {
          currentLine = currentLine - 1;
        }
        if (currentLine < 0) currentLine = 0;
        await checkPosition();
        AltAdj = "";
        break;

      case "ArrowDown":
      case "j":
        e.stopPropagation();
        if (AltAdj !== "") {
          currentLine = currentLine + parseInt(AltAdj);
        } else {
          currentLine = currentLine + 1;
        }
        if (currentLine >= lastData.data.length) {
          currentLine = lastData.data.length - 1;
        }
        await checkPosition();
        AltAdj = "";
        break;

      case "G":
        //
        // This implements a go to the top.
        //
        currentLine = 0;
        await checkPosition();
        AltAdj = "";
        break;

      case "g":
        //
        // This implements a go to the bottom.
        //
        currentLine = lastData.data.length - 1;
        await checkPosition();
        AltAdj = "";
        break;

      case "Enter":
        //
        // The user hit enter, so process the line.
        //
        showHint = false;
        inputState = 0;
        e.stopPropagation();
        ExecuteLine(currentLine);
        break;

      case "Escape":
        //
        // Reset everything.
        //
        showHint = false;
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        inputState = 0;
        AltAdj = "";
        e.stopPropagation();
        break;
    }
  }

  function htmlKeyProcess(e) {
    switch (e.key) {
      case "DownArrow":
      case "k":
        currentLine = currentLine - 1;
        if (currentLine < 0) currentLine = 0;
        e.stopPropagation();
        break;

      case "UpArrow":
      case "j":
        currentLine = currentLine + 1;
        if (currentLine >= lastData.data.length) {
          currentLine = lastData.data.length - 1;
        }
        e.stopPropagation();
        break;

      case "Enter":
        //
        // The user hit enter, so process the line.
        //
        showHint = false;
        inputState = 0;
        e.stopPropagation();
        ExecuteLine(currentLine);
        break;

      case "Escape":
        //
        // Reset everything.
        //
        showHint = false;
        inputState = 0;
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        e.stopPropagation();
        break;
    }
  }

  function ExecuteLine(index) {
    ProcessLine(lastData.data[index].command);
  }

  async function quit() {
    App.Quit();
  }
</script>

<div id="container" bind:this={containerDiv}>
  <div
    id="ScriptTermDiv"
    style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
  >
    <input
      id="CommandInput"
      autocomplete="off"
      spellcheck="false"
      autocorrect="off"
      on:keydown={processKey}
      on:input={lineInput}
      bind:this={inputdiv}
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; font-size: {$theme.fontSize}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    />
    {#if showHint}
      <div id="hintDiv" style="border-color: {$theme.Cyan};">
        {#each hints as hint, hindex}
          {#if hindex === hintCursor}
            <p style="background-color: {$theme.highlightBackgroundColor};">
              {hint}
            </p>
          {:else}
            <p>{hint}</p>
          {/if}
        {/each}
      </div>
    {/if}
    <div
      id="statusline"
      style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      WD: {wd}
    </div>
    {#if showError}
      <div id="errorContainer">
        {@html errorOutput}
      </div>
    {/if}
  </div>
  <div id="outputContainer">
    {#if showOutputDiv}
      <input
        style="height: 0px; width: 0px; margin: 0px; border: 0px; padding: 0px;"
        bind:this={outputDivInput}
        on:keydown|preventDefault={outputKeyProcess}
      />
      <div
        bind:this={scrollDiv}
        id="scriptOutput"
        style="width: {showHtmlDiv
          ? '40%'
          : '100%'}; background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        {#if InterActive}
          {#each textOutput as tout, index}
            {#if currentLine == index}
              <a
                href="/"
                data-index={index}
                bind:this={cursorDiv}
                class="outputLine"
                style="background-color: {$theme.highlightBackgroundColor}; color: {$theme.selectionColor};"
                on:click|preventDefault={() => {
                  ExecuteLine(index);
                }}
              >
                {@html tout}
              </a>
            {:else}
              <a
                href="/"
                data-index={index}
                class="outputLine"
                style="background-color: {$theme.backgroundColor}; color: {$theme.textColor};"
                on:click|preventDefault={() => {
                  ExecuteLine(index);
                }}
              >
                {@html tout}
              </a>
            {/if}
          {/each}
        {:else}
          {@html textOutput.join("\n")}
        {/if}
      </div>
    {/if}
    {#if showHtmlDiv}
      <div
        id="htmlOutput"
        style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        <input bind:this={htmlDivInput} on:keydown={htmlKeyProcess} />
        {@html htmlOutput}
      </div>
    {/if}
  </div>
</div>

<style>
  #container {
    display: flex;
    flex-direction: column;
    padding: 0px;
    margin: 0px;
    background-color: transparent;
    width: 1022px;
    height: 70%;
    border: 0px;
    border-color: transparent;
  }

  #errorContainer {
    display: flex;
    flex-direction: column;
  }

  #outputContainer {
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 100%;
    border: 0px transparent;
  }

  #scriptOutput {
    display: flex;
    flex-direction: column;
    height: 455px;
    margin: 5px 0px 0px 0px;
    padding: 5px;
    overflow: auto;
    border-radius: 10px;
  }

  #htmlOutput {
    display: flex;
    flex-direction: column;
    height: 440px;
    margin: 10px;
    padding: 10px;
    overflow: auto;
    border-radius: 10px;
  }

  #ScriptTermDiv {
    display: flex;
    flex-direction: column;
    padding: 0px 10px 10px 10px;
    margin: 0px;
    border-radius: 0px 0px 10px 10px;
  }

  #CommandInput {
    margin: 10px;
    padding: 5px;
    border-radius: 10px;
    border: 3px;
  }

  #statusline {
    display: flex;
    flex-direction: row;
    height: 30px;
    width: 900px;
    margin: 10px;
  }

  #hintDiv {
    display: flex;
    flex-direction: column;
    padding: 5px;
    margin: 5px 0px 5px 20px;
    border-radius: 10px;
    border: 3px solid;
  }

  #hintDiv p {
    margin: 5px 0px 0px 5px;
    padding: 0px 0px 0px 5px;
    border-radius: 5px;
  }

  .outputLine {
    text-decoration: none;
    margin: 0px;
    padding: 5px 0px 0px 0px;
    border-radius: 5px;
  }
</style>
