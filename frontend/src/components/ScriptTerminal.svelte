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
    cd: cdCommand,
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

  onMount(() => {
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
          case "j":
            if (lcommandRow !== 0) {
              lcommandRow -= 1;
              term.write(termAtb.up);
            }
            break;
          case "k":
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
            // Get the command to run and run it.
            //
            console.log("run the line");
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
      if (typeof line.tcommand !== "undefined" && line.command.length > 0) {
        //
        // We have a terminal command to perform.
        //
        RunTerminalCommand(line.tcommand);
      }
    });

    //
    // Keep a list of valid commands.
    //
    commands.push(lastData.line);

    //
    // Go to next line and give a prompt.
    //
    term.prompt();
  }

  function RunTerminalCommand(text) {
    //
    // Get the command and it's arguments separated.
    //
    let words = text.split(" ");
    if (words.length > 0) {
      tCommands[words[0]](words.slice(1).join(" "));
    }
  }

  function cdCommand(text) {
    //
    // Remove inclosing quotes
    //
    if (text[0] === '"' || text[0] === '"') {
      text = text.slice(1, text.length - 1);
    }
    if (text.length > 0) {
      wd = text;
    }
    lastData.valid = false;
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
