<script>
  import { onMount } from "svelte";
  import { Terminal } from "xterm";
  import "xterm/css/xterm.css";
  import { WebLinksAddon } from "xterm-addon-web-links";
  import { FitAddon } from "xterm-addon-fit";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";
  import { termscripts } from "../stores/termscripts.js";

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
    homeDir = await window.go.main.App.GetHomeDir();
    wd = homeDir;
    term = new Terminal({
      rendererType: "canvas",
      convertEol: true,
      cursorBlink: true,
      cursorStyle: "bar",
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
            term.write(`${lastData.data[lcommandRow].tcommand}\n\r`);

            //
            // Get the command to run and run it.
            //
            RunTerminalCommand(lastData.data[lcommandRow].tcommand);
            break;
          default:
            break;
        }
      }
    });
    //
    // Make sure the terminal is focused.
    //
    term.focus();
  });

  function splitAt(index, xs) {
    return [xs.slice(0, index), xs.slice(index)];
  }

  function ProcessLine(text) {
    //
    // Get the words of the line.
    //
    var words = text.split(" ");
    lastData.line = text;
    lastData.valid = false;
    curr_line = "";

    //
    // See if the first word is a valid command.
    //
    let scrpt = $termscripts.filter((item) => item.name === words[0]);
    if (scrpt.length === 0) {
      //
      // Command not found. Print the error and give a new prompt.
      //
      term.write(
        `\r\n\r\n    ${termAtb.red}<Error>${termAtb.default} The command "${words[0]}" wasn't found!\r\n\r\n`
      );
      term.prompt();
    } else {
      //
      // Command found. Run it!
      //
      fetch("http://localhost:9978/api/script/run", {
        method: "PUT",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify({
          script: scrpt[0].name,
          text: words.slice(1).join(" "),
          envVar: { SCRIPTTERMCWD: wd },
        }),
      })
        .then((resp) => {
          return resp.json();
        })
        .then((data) => {
          ProcessScriptReturn(data.text);
        });
    }
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
    //      lcommand: <command line string to run>,
    //      tcommand: <terminal Command to run>
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
    // Get the command and it's arguments separated.
    //
    let words = text.split(" ");
    if (words.length > 0) {
      await tCommands[words[0]].command(words.slice(1).join(" "));
    }
    term.prompt();
  }

  async function cdCommand(text) {
    //
    // Remove inclosing quotes
    //
    if (text[0] === '"' || text[0] === "'") {
      text = text.slice(1, text.length - 1);
    }
    if (text.length > 0) {
      let path = new String(text);
      if (text[0] !== "/") {
        let ndir = new String(text);
        let nwd = new String(wd);
        path = await window.go.main.App.AppendPath(nwd, ndir);
      }
      let exists = await window.go.main.App.DirExists(path);
      if (exists) {
        wd = path;
      } else {
        term.write(`    <Error> The path "${path} doesn't exist!\n\n`);
      }
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
        var description = [item.description];
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
          description = description.join("\n\r          ");
        }
        term.write(`    ${item.name}   ${description}\n\r`);
      });
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
        term.write(`    ${spt.name}  -  ${spt.help}\n\r`);
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
        path = await window.go.main.App.AppendPath(path, text);
      }
    }
    var dirReal = await window.go.main.App.DirExists(path);
    if (dirReal) {
      var result = await window.go.main.App.ReadDir(path);
      var lines = [];
      for (let i = 0; i < result.length; i++) {
        //
        // Rewrite lastData.lines to have a tcommand for each entry printed.
        //
        let item = result[i];
        let npath = await window.go.main.App.AppendPath(item.Dir, item.Name);
        dirReal = await window.go.main.App.DirExists(npath);
        if (dirReal) {
          lines.push({
            name: item.Name,
            tcommand: `cd '${npath}'`,
          });
        } else {
          lines.push({
            name: item.Name,
            tcommand: `open '${npath}'`,
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
      let fileReal = await window.go.main.App.FileExists(path);
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
      text = await window.go.main.App.AppendPath(wd, text);
    }
    text = new String(text.trim());

    //
    // Run the open command on the file.
    //
    await window.go.main.App.RunCommandLine(
      "/usr/bin/open",
      ["-t", text],
      [],
      ""
    );
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
      var isText = false;
      if (text[0] === '"' || text[0] === "'") {
        text = text.slice(1, text.length - 1);
        isText = true;
      }
      if (text[0] !== "/" && !isText) {
        text = await window.go.main.App.AppendPath(wd, text);
        isText = false;
      }
      text = new String(text.trim());

      //
      // Process the text based on what it is.
      //
      if (isText) {
        //
        // Send the text to run the script on.
        //
        await fetch("http://localhost:9978/api/script/run", {
          method: "PUT",
          headers: {
            "Content-type": "application/json",
          },
          body: JSON.stringify({
            script: scriptName,
            text: text,
          }),
        })
          .then((resp) => {
            return resp.json();
          })
          .then((data) => {
            term.write(`\n\r     ${data.text}\n\r`);
          });
      } else {
        //
        // It is a file to run the script on.
        //
        await fetch("http://localhost:9978/api/script/run", {
          method: "PUT",
          headers: {
            "Content-type": "application/json",
          },
          body: JSON.stringify({
            script: scriptName,
            text: "",
            file: text,
          }),
        })
          .then((resp) => {
            return resp.json();
          })
          .then((data) => {
            term.write(`\n\r     ${data.text}\n\r`);
          });
      }
    }
  }

  function viewEmailIt() {
    $state = "emailit";
  }

  function viewNotes() {
    $state = "notes";
  }

  function viewLog() {
    $state = "viewlog";
  }

  function viewScriptEditor() {
    $state = "scripts";
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
      on:click={viewLog}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Log
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
    height: 100%;
    width: 100%;
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
    margin: 10px auto;
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
