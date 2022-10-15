<script>
  import { onMount, afterUpdate } from "svelte";
  import EmailIt from "./components/EmailIt.svelte";
  import ViewLog from "./components/ViewLog.svelte";
  import Notes from "./components/Notes.svelte";
  import ScriptMenu from "./components/ScriptMenu.svelte";
  import TemplateMenu from "./components/TemplateMenu.svelte";
  import ScriptsEditor from "./components/ScriptsEditor.svelte";
  import TemplatesEditor from "./components/TemplatesEditor.svelte";
  import Preferences from "./components/Preferences.svelte";
  import ScriptTerminal from "./components/ScriptTerminal.svelte";
  import { state } from "./stores/state.js";
  import { scripts } from "./stores/scripts.js";
  import { termscripts } from "./stores/termscripts.js";
  import { showScripts } from "./stores/showScripts.js";
  import { templates } from "./stores/templates.js";
  import { showTemplates } from "./stores/showTemplates.js";
  import { theme } from "./stores/theme.js";
  import { config } from "./stores/config.js";
  import { extScripts } from "./stores/extScripts.js";
  import { systemScripts } from "./stores/systemScripts.js";
  import { userScripts } from "./stores/userScripts.js";
  import { userTemplates } from "./stores/userTemplates.js";
  import { systemTemplates } from "./stores/systemTemplates.js";
  import { runscript } from "./stores/runscript.js";
  import { notes } from "./stores/notes.js";
  import * as App from '../wailsjs/go/main/App.js';
  import { DateTime } from "luxon";
  import { create, all } from 'mathjs';
  import Handlebars from "handlebars";

  let starting = true;
  let notestruct = {
    notes: ["","","","","","","","","",""],
    loadNotes: async function() {
      //
      // Load the notes.json file.
      //
      let notejsonloc = await App.AppendPath($config.configDir, "notes.json");
      if(await App.FileExists(notejsonloc)) {
        notestruct.notes = await App.ReadFile(notejsonloc);
        notestruct.notes = JSON.parse(this.notes);
      } else {
        notestruct.saveNotes();
      }
    },
    saveNotes: async function() {
      //
      // Save the notes.json file.
      //
      let notejsonloc = await App.AppendPath($config.configDir, "notes.json");
      await App.WriteFile(notejsonloc, JSON.stringify(notestruct.notes));
    },
    getNote: function(noteid) {
      return notestruct.notes[noteid];
    },
    putNote: async function(noteid, note) {
      notestruct.notes[noteid] = note;
      await notestruct.saveNotes();
    }
  }
  const mathconfig = {
    epsilon: 1e-12,
    matrix: 'Matrix',
    number: 'number',
    precision: 64,
    predictable: false,
    randomSeed: null
  };
  const mathjs = create(all, mathconfig);
  var SP = {
      text: '',
      luxon: DateTime,
      mathjs: mathjs,
      mathParser: mathjs.parser(),
      that: {},
      Handlebars: Handlebars,
      ProcessMathSelection: function(txt) {
        var result = this.mathParser.evaluate(txt)
        return result
      },
      ProcessMathNotes: function(Note) {
        var result = ''
        this.mathParser.clear()
        let lines = Note.match(/^.*((\\r\\n|\\n|\\r)|$)/gm)
        let numLines = lines.length
        for (let i = 0; i < lines.length; i++) {
          var line = lines[i];
          var lineResult = ''
          line = line.trim()
          if (line.indexOf('>') === 0) {
            line = line.substr(2)
            var inx = line.indexOf('|')
            if (inx !== -1) {
              line = line.substr(0, inx - 1)
              line = line.trim()
            }
            lineResult = this.mathParser.evaluate(line)
            if ((typeof lineResult) === 'function') {
              lineResult = 'Definition'
            }
            var whiteSP = 32 - (line.length + 3);
            if (whiteSP < 1) {
              whiteSP = 1;
            }
            result += '> ' + line + this.insertCharacters(whiteSP, ' ') + ' | ' + lineResult
          } else {
            result += line
          }
          if (--numLines !== 0) result += '\n'
        }
        return result
      },
      runScript: function(script, text) {
        return that.runJavaScriptScripts(that.getJavaScriptScript(script), text)
      },
      returnNote: function(id) {
        var result = '';
        if ((id >= 0) && (id <= 9)) result = that.NOTES[id];
        return result;
      },
      insertCharacters: function(num, ichar) {
        var result = ''
        if (num < 0) return result
        for (var i = 0; i < num; i++) {
          result += ichar
        }
        return result
      }
    };


  onMount(async () => {
    //
    // Wait to give the server time to start.
    //
    $state = "notready";

    //
    // Get stuff from the server.
    //
    await getConfig();
    await getScriptsList();
    getTermScriptsList();
    await getTemplatesList();
    await getTheme();
    await getNotes();

    $runscript = runScript;

    //
    // Set the state to emailit.
    //
    $state = "emailit";
  });

  function wait(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(ms);
      }, ms);
    });
  }

  afterUpdate(() => {
    if (starting) {
      starting = false;
    }
  });

  async function getNotes() {
    $notes = notestruct;
    await $notes.loadNotes();
  }

  async function getConfig() {
    //
    // Get the config and store in in the config store.
    //
    let hmdir = await App.GetHomeDir();
    let configdir = await App.AppendPath(hmdir,".config");
    configdir = await App.AppendPath(configdir, "EmailIt");
    if(await App.DirExists(configdir)) {
      let configloc = await App.AppendPath(configdir,"config.json");
      if(await App.FileExists(configloc)) {
        let configJSON = await App.ReadFile(configloc);
        $config = JSON.parse(configJSON);
      } else {
        //
        // Create the default config file.
        //
        await createDefaultConfig(hmdir, configdir);
      }
    } else {
      //
      // The config directory needs to be created and populated.
      //
      await App.MakeDir(configdir);
      await createDefaultConfig(hmdir, configdir);
    }
  }

  async function createDefaultConfig(homedir, configdir) {
    //
    // Create the default Configuration file.
    //
    let themedir = await App.AppendPath(configdir,"themes");
    let scriptsdir = await App.AppendPath(configdir,"scripts");
    let configloc = await App.AppendPath(configdir,"config.json");
    let systemdir = await App.GetExecutable();
    $config = {
      homeDir: homedir,
      systemDir: systemdir,
      themeDir: themedir,
      theme: "Default",
      configDir: configdir,
      scriptsDir: scriptsdir
    }
    await App.WriteFile(configloc, JSON.stringify($config));

    //
    // Create the default Theme.
    //
    let defaultThemeDir = await App.AppendPath($config.themeDir,"Default");
    if(!await App.DirExists(defaultThemeDir)) {
      await App.MakeDir(defaultThemeDir)
    }
    let defaultThemefile = await App.AppendPath(defaultThemeDir,"Default.json");
    await App.WriteFile(defaultThemefile, JSON.stringify({
      name: "Default",
      font: "Fira Code, Menlo",
      fontSize: "16pt",
      textAreaColor: '#454158',
      backgroundColor: '#22212C',
      textColor: '#80ffea',
      borderColor: '#1B1A23',
      Cyan: "#80FFEA",
      Green: "#8AFF80",
      Orange: "#FFCA80",
      Pink: "#FF80BF",
      Purple: "#9580FF",
      Red: "#FF9580",
      Yellow: "#FFFF80",
      functionColor: "#9580FF",
      stringColor: "#8AFF80",
      constantColor: "#FFCA80",
      keywordColor: "#FFFF80",
      highlightBackgroundColor: "#4f4f5f",
      selectionColor: "#22212C",
      buttons: [
        {
          color: '#80FFEA',
          id: 0
        },
        {
          color: '#8AFF80',
          id: 1
       },
        {
          color: '#FFCA80',
          id: 2
        },
        {
          color: '#FF80BF',
          id: 3
        },
        {
          color: '#9580FF',
          id: 4
        },
        {
          color: '#FF9580',
          id: 5
        },
        {
          color: 'blue',
          id: 5
        },
        {
          color: 'green',
          id: 7
        },
        {
          color: 'red',
          id: 8
        },
        {
          color: 'purple',
          id: 9
        },
      ]
    }));
    let themecnfg = await App.AppendPath(defaultThemeDir,"package.json");
    await App.WriteFile(themecnfg,JSON.stringify(
      {
        "name": "Default",
        "version": "1.0.0",
        "description": "The EmailIt program default theme.",
        "keywords": [
          "emailit", "theme"
        ],
        "author": "Richard Guay",
        "license": "MIT",
        "theme": {
          "name": "Default",
          "description": "The EmailIt program default theme.",
          "type": 0,
          "github": "",
          "main": "Default.json"
        }
    }));

    //
    // Create the default Notes.
    //
    notesloc = await App.AppendPath(configdir,"notes.json");
    await App.WriteFile(notesloc, JSON.stringify(["","","","","","","","","",""]));
  }

  async function getTheme() {
    let thm = await App.AppendPath($config.configDir,"theme.json");
    let thmjson = await App.ReadFile(thm);
    $theme = JSON.parse(thmjson);
  }

  async function getScriptsList() {
    let userScriptsLoc = await App.AppendPath($config.configDir,"scripts.json");
    let extScriptsLoc = await App.AppendPath($config.configDir,"extscripts.json");
    let systemScriptsLoc = await App.AppendPath($config.systemDir,"../Resources/systemscripts.json")
    let userScriptFile = {};
    if(await App.FileExists(userScriptsLoc)) {
      //
      // Load the user scripts file.
      //
      userScriptFile = await App.ReadFile(userScriptsLoc);
      $userScripts = JSON.parse(userScriptFile);
    } else {
      //
      // Save a default user scripts file.
      //
      $userScripts = [];
    }
    if(await App.FileExists(extScriptsLoc)) {
      //
      // Load the external scripts file.
      //
      let extScriptFile = await App.ReadFile(extScriptsLoc);
      $extScripts = JSON.parse(extScriptFile);
    } else {
      //
      // Save a default external scripts file.
      //
      $extScripts = [];
    }
    $scripts = $userScripts.filter(value => value.termscript === false).map(value => {
        return { name: value.name, insert: value.insert }
      }).concat($systemScripts.filter(value => value.termscript === false).map(value => {
        return { name: value.name, insert: value.insert }
      })).concat($extScripts.filter(value => value.termscript === false).map(value => {
        return { name: value.name, insert: false }
      }));
  }

  function getTermScriptsList() {
    $termscripts = $userScripts.filter(value => value.termscript === true)
      .concat($systemScripts.filter(value => value.termscript === true))
      .concat($extScripts.filter(value => value.termscript === true));
  }

  async function getTemplatesList() {
    let templateloc = await App.AppendPath($config.configDir,"templates.json");
    if(await App.FileExists(templateloc)) {
      let templatefile = await App.ReadFile(templateloc);
      $userTemplates = JSON.parse(templatefile);
      $templates = $systemTemplates.map(item => item.name).concat($userTemplates.map(item => item.name));
    } else {
      $userTemplates = [];
      $templates = $systemTemplates.map(item => item.name);
    }
  }

  //
  // Function:         runJavaScriptScriptsFile
  //
  // Description:      This will run the JavaScript function on the contents of a file.
  //
  // Inputs:
  //                   script          The script.
  //                   text            The text to process.
  //
  async function runJavaScriptFile(script, file) {
    let result = "Invalid File";
    if (await App.FileExists(file)) {
      result = await App.ReadFile(file);
      result = runJavaScript(script, result);
      let parts = await App.SplitFile(file);
      let nfile = await App.AppendPath(parts.Dir, `${parts.Name}-processed${parts.Extension}`);
      await App.WriteFile(nfile, result);
      result = `${parts.Name}-processed${parts.Extension} was created!`;
    }
    return (result);
  }

  //
  // Function:         runJavaScript
  //
  // Description:      This will run some given text with a script.
  //
  // Inputs:
  //                   script          The script.
  //                   text            The text to process.
  //
  function runJavaScript(script, text) {
    var originalText = text;
    SP.that = SP;
    SP.text = originalText;

    //
    // Try to evaluate the expression.
    //
    try {
      var scriptFunction = new Function('SP', `${script} ; return SP;`);
      SP = scriptFunction(SP);
    } catch (error) {
      console.error(error);
      SP.text = originalText;
    }
    //
    // Make sure we have a string and not an array or object.
    //
    if (typeof SP.text != 'string') {
      SP.text = SP.text.toString();
    }
    return (SP.text);
  }

  async function runExtScript(extScrpt, text) {
    //
    // extScript.name     - File name of the script
    // extScript.script   - User name for the script
    // extScript.path     - directory of the script
    // extScript.env      - name of the environment
    // extScript.termscript - Boolean true if a script terminal script.
    // extScript.description - Description of what the script does.
    // extScript.help     - A help message for the script
    //
    var result = '';
    var env = {};

    //
    // Get the environment.
    //
    if (extScrpt.env !== '') {
      let envlistloc = await App.AppendPath($config.configDir, "environments.json");
      let envlist = await App.ReadFile(envlistloc);
      envlist = JSON.parse(envlist);
      env = envlist.find(item => item.name === extScrpt.env);
      if (env !== 'undefined') {
        env = env.envVar;
      } else {
        env = {};
      }
    }
    try {
      let args = [];
      args.push(text);
      result = await App.RunCommandLine('./' + extScrpt.script, args, env, extScrpt.path); 
      if (result !== null && typeof result === 'object') result = result.toString();
    } catch (e) {
      result = "Error: " + e.message;
    }
    return (result);
  }

  async function runScript(script, text) {
    var result = "Script not found!";
    var scriptIndex = $userScripts.find((ele) => { return ele.name === script })
    if (typeof scriptIndex !== 'undefined') {
      script = scriptIndex.script;
      let isfile = await App.FileExists(text);
      if (isfile) {
        result = await runJavaScriptFile(script, text);
      } else {
        result = runJavaScript(script, text);
      }
    } else {
      scriptIndex = $systemScripts.find(ele => { return ele.name === script});
      if(typeof scriptIndex !== 'undefined') {
        script = scriptIndex.script;
        let isfile = await App.FileExists(text);
        if (isfile && text !== '..' && text !== '.') {
          result = await runJavaScriptFile(script, text);
        } else {
          result = runJavaScript(script, text);
        }
      } else {
        scriptIndex = $extScripts.find(ele => { return ele.name === script});
        if (typeof scriptIndex !== 'undefined') {
          //
          // It's an external script.
          //
          result = await runExtScript(scriptIndex, text);
        }
      } 
    }
    return(result);
  }
  
  function keyDownProcessor(e) {
    if (e.metaKey && e.key === ",") {
      e.preventDefault();
      $state = "preferences";
    } else if (e.ctrlKey) {
      switch (e.key) {
        case "e":
          $state = "emailit";
          e.preventDefault();
          break;

        case "v":
          $state = "viewlog";
          e.preventDefault();
          break;

        case "n":
          $state = "notes";
          e.preventDefault();
          break;

        case "m":
          $showScripts = !$showScripts;
          e.preventDefault();
          break;

        case "t":
          $showTemplates = !$showTemplates;
          e.preventDefault();
          break;

        case "l":
          $state = "scriptterm";
          e.preventDefault();
          break;

        case "p":
          $state = "preferences";
          e.preventDefault();
      }
    }
  }
</script>

<svelte:window on:keydown={keyDownProcessor} />
<div id="dragbar" />

{#if $state === "emailit"}
  <EmailIt />
{:else if $state === "viewlog"}
  <ViewLog />
{:else if $state === "notes"}
  <Notes />
{:else if $state === "scripts"}
  <ScriptsEditor />
{:else if $state === "templates"}
  <TemplatesEditor />
{:else if $state === "preferences"}
  <Preferences />
{:else if $state === "scriptterm"}
  <ScriptTerminal />
{:else}
  <div id="error">
    <h1>Something went wrong!</h1>
  </div>
{/if}

<ScriptMenu />
<TemplateMenu />

<style>
  :global(body) {
    background-color: rgb(34, 33, 44);
  }

  #dragbar {
    height: 30px;
    width: 100%;
    --wails-draggable: drag;
  }
</style>
