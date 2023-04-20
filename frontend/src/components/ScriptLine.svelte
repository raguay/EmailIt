<script>
  import { onMount, afterUpdate } from "svelte";
  import { termscripts } from "../stores/termscripts.js";
  import { runscript } from "../stores/runscript.js";
  import { aliases } from "../stores/aliases.js";
  import { state } from "../stores/state.js";
  import { config } from "../stores/config.js";
  import { theme } from "../stores/theme.js";
  import * as App from "../../wailsjs/go/main/App.js";

  let commands = [];
  let showHtmlDiv = false;
  let showOutputDiv = false;
  let showError = false;
  let errorOutput = "";
  let textOutput = "";
  let htmlOutput = "";
  let inputdiv = null;
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

  onMount(async () => {
    homeDir = await App.GetHomeDir();
    wd = homeDir;

    //
    // Load the aliases.
    //
    loadAliases();
  });

  afterUpdate(() => {
    if (!showOutputDiv && !showHtmlDiv) inputdiv.focus();
    else if (showOutputDiv) {
      outputDivInput.focus();
    } else {
      htmlDivInput.focus();
    }
  });

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
    // Show the user information that is in text format. Convert it
    // to HTML and then show in the information area.
    //
    showOutputDiv = true;
    if (clear) textOutput = "";
    textOutput += data.toString();
  }

  function showOutputHTML(data, clear) {
    //
    // Show HTML formated data to the user.
    //
    showHtmlDiv = true;
    if (clear) htmlOutput = "";
    htmlOutput += data;
  }

  function lineInput(e) {
    //
    // Get the information from the line input and see if we can expand it to a
    // valid command.
    // TODO: Not Implemented
    //
  }

  async function ProcessLine(text) {
    //
    // clear out the output items.
    //
    showHtmlDiv = false;
    showOutputDiv = false;
    showError = false;
    htmlOutput = "";
    textOutput = "";

    //
    // Get the words of the line.
    //
    var words = text.trim().split(" ");
    lastData.line = text;
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
        showOutputText(`<p style="margin: 3px;">${item.Name}</p>`, false);
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
        showOutputText(`<p>${text}</p>`, true);
      } else {
        showErrorHTML("Not a valid Path.");
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
      text = new String(text.trim());

      //
      // Process the text based on what it is.
      //
      let newText = await $runscript(scriptName, text);
      showOutputText(`<p>${newText}</p>`, true);
    }
    lastData.valid = false;
  }

  async function editCommand(text) {
    //
    // fix the file name.
    //
    text = new String(text.trim());
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text[0] !== "/") {
      text = await App.AppendPath(wd, text);
    }
    console.log("edit: ", text);

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
      showOutputText("", true);
      for (const item of $aliases) {
        //
        // Show each alias
        //
        showOutputText(`<p>${item.name}</p>`, false);
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
      showOutputText(`<p>${commands[i]}</p>`, false);
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
        showOutputText(`<p>${item.Name}</p>`, false);
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
      showErrorHTML(`The path "${path}" already exists.`);
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
      showErrorHTML(`The file "${path}" already exists!`);
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
    var args = ["/bin/zsh", "-c", line];
    var cmd = "/bin/zsh";

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
        ProcessLine(e.target.value);
        break;

      case "Tab":
        //
        // Expand the suggestion.
        //
        e.stopPropagation();
        break;

      case "Escape":
        //
        // Reset everything.
        //
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        e.stopPropagation();
        break;
    }
  }

  function outputKeyProcess(e) {
    switch (e.key) {
      case "UpArrow":
      case "j":
        break;

      case "DownArrow":
      case "k":
        break;

      case "Enter":
        //
        // The user hit enter, so process the line.
        //
        e.stopPropagation();
        break;

      case "Escape":
        //
        // Reset everything.
        //
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        e.stopPropagation();
        break;
    }
  }

  function htmlKeyProcess(e) {
    switch (e.key) {
      case "UpArrow":
      case "j":
        break;

      case "DownArrow":
      case "k":
        break;

      case "Enter":
        //
        // The user hit enter, so process the line.
        //
        e.stopPropagation();
        break;

      case "Escape":
        //
        // Reset everything.
        //
        inputdiv.value = "";
        showOutputDiv = false;
        showHtmlDiv = false;
        showError = false;
        e.stopPropagation();
        break;
    }
  }

  async function quit() {
    App.Quit();
  }
</script>

<div id="container">
  <div
    id="ScriptTermDiv"
    style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
  >
    <input
      id="CommandInput"
      on:keydown={processKey}
      on:input={lineInput}
      bind:this={inputdiv}
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; font-size: {$theme.fontSize}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    />
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
      <div
        id="scriptOutput"
        style="width: {showHtmlDiv
          ? '40%'
          : '100%'}; background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        <input bind:this={outputDivInput} on:keydown={outputKeyProcess} />
        {@html textOutput}
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
    width: 1000px;
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
    height: 440px;
    margin: 10px;
    padding: 10px;
    overflow: auto;
    border-radius: 10px;
  }

  #scriptOutput input {
    display: none;
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

  #htmlOutput input {
    display: none;
  }

  #ScriptTermDiv {
    display: flex;
    flex-direction: column;
    padding: 0px 10px 10px 10px;
    margin: 0px;
    width: 1000px;
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

  #scriptOutput p {
    margin: 0px;
  }
</style>
