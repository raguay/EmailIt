<script>
  import { onMount } from "svelte";
  import EmailIt from "./components/EmailIt.svelte";
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
  import { runhandlebars } from "./stores/runhandlebars.js";
  import { runtemplate } from "./stores/runtemplate.js";
  import * as App from '../wailsjs/go/main/App.js';
  import { DateTime } from "luxon";
  import { create, all } from 'mathjs';
  import Handlebars from "handlebars";

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
      DateTime: DateTime,
      mathjs: mathjs,
      mathParser: mathjs.parser(),
      that: SP,
      Handlebars: $runhandlebars,
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
        return runJavaScript(script, text)
      },
      returnNote: function(id) {
        var result = '';
        if ((id >= 0) && (id <= 9)) result = $notes.getNote(id);
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
    setupHandlebarHelpers();

    $runscript = runScript;
    $runhandlebars = Handlebars;
    SP.Handlebars = Handlebars;
    $runtemplate = runTemplate;

    //
    // Run the StartUp script if one exists.
    //
    let startup = $scripts.find(item => item.name === "StartUp");
    if(typeof startup !== 'undefined') {
      //
      // A StartUp script is defined. Run it!
      //
      $runscript(startup,"");
    }

    //
    // Set the state to emailit.
    //
    $state = "emailit";
  });

  function setupHandlebarHelpers() {
    //
    // Create the helpers functions for Handlebars.
    //
    Handlebars.registerHelper('save', function(name, text) {
      Handlebars.registerHelper(name, function() {
        return text;
      });
      return text;
    });

    Handlebars.registerHelper('clipboard', async function() {
      return await App.GetClip();
    });

    Handlebars.registerHelper('date', function(dFormat) {
      return DateTime.now().toFormat(dFormat);
    });

    Handlebars.registerHelper('cdate', function(cTime, dFormat) {
      return DateTime.fromISO(cTime).toFormat(dFormat);
    });

    Handlebars.registerHelper('last', function(weeks, fmat) {
      return DateTime.now().minus({ weeks: weeks}).toFormat(fmat);
    });

    Handlebars.registerHelper('next', function(weeks, fmat) {
      return DateTime.now().plus({weeks: weeks}).toFormat(fmat);
    });
  }

  function runTemplate(template, text) {
    //
    // process the template.
    //
    let result = '';
    try {
      //
      // Create various default targets for the templater. These have
      // to be created since they are various types of time/date stamps.
      //
      var data = [];
      data["cDateMDY"] = DateTime.now().toFormat("LLLL dd, yyyy");
      data["cDateDMY"] = DateTime.now().toFormat("dd LLLL yyyy");
      data["cDateDOWDMY"] = DateTime.now().toFormat("cccc, dd LLLL yyyy");
      data["cDateDOWMDY"] = DateTime.now().toFormat("cccc LLLL dd, yyyy");
      data["cDay"] = DateTime.now().toFormat("dd");
      data["cMonth"] = DateTime.now().toFormat("LLLL");
      data["cYear"] = DateTime.now().toFormat("yyyy");
      data["cMonthShort"] = DateTime.now().toFormat("LLL");
      data["cYearShort"] = DateTime.now().toFormat("yy");
      data["cDOW"] = DateTime.now().toFormat("cccc");
      data["cMDthYShort"] = DateTime.now().toFormat("LLL d yy");
      data["cMDthY"] = DateTime.now().toFormat("LLLL d yyyy");
      data["cHMSampm"] = DateTime.now().toFormat("h:mm:ss a");
      data["cHMampm"] = DateTime.now().toFormat("h:mm a");
      data["cHMS24"] = DateTime.now().toFormat("H:mm:ss");
      data["cHM24"] = DateTime.now().toFormat("H:mm");
      data["Templatename"] = template.name;
      data["text"] = text;

      //
      // Get the User's default data.
      //
      var defaultData = $templates.find(item => item.name === "Defaults");
      template = $templates.find(item => item.name === template);
      if(typeof template !== 'undefined') {
        template = template.template;
      }
      if (defaultData !== undefined) {
        data = MergeRecursive(data, JSON.parse(defaultData.template));
      }
      //
      // Parse the Template
      //
      var tpTemplate = $runhandlebars.compile(template);

      //
      // Run the template with the data.
      //
      result = tpTemplate(data);

      //
      // Make sure we have a string and not an array or object.
      //
      if (typeof result != 'string') {
        result = result.toString();
      }
    } catch (error) {
      console.error("Handlebars Error: " + error);
      result = "There was an error.";
    }
    return(result);
  }

  //
  //  Function:        MergeRecursive
  //
  //  Description:     Recursively merge properties of two objects
  //
  //  Inputs:
  //                   obj1    The first object to merge
  //                   obj2    The second object to merge
  //
  function MergeRecursive(obj1, obj2) {
    for (var p in obj2) {
      try {
        // Property in destination object set; update its value.
        if (obj2[p].constructor == Object) {
          obj1[p] = MergeRecursive(obj1[p], obj2[p]);
        } else {
          obj1[p] = obj2[p];
        }
      } catch (e) {
        // Property in destination object not set; create it and set its value.
        obj1[p] = obj2[p];
      }
    }
    return obj1;
  }

  function wait(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(ms);
      }, ms);
    });
  }

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
      $templates = $systemTemplates.concat($userTemplates);
    } else {
      $userTemplates = [];
      $templates = $systemTemplates;
    }
  }

  //
  // Function:         runJavaScriptScriptFile
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
    SP.text = text;
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
    if(typeof script === 'object') {
      if(typeof script.name !== 'undefined') {
        script = script.name;
      }
    }
    var scriptIndex = $userScripts.find((ele) => { return ele.name === script })
    if (typeof scriptIndex !== 'undefined') {
      script = scriptIndex.script;
      let tmptext = text;
      let isfile = await App.FileExists(tmptext);
      let isDir = await App.DirExists(tmptext);
      if (isfile && !isDir) {
        result = await runJavaScriptFile(script, text);
      } else {
        result = runJavaScript(script, text);
      }
    } else {
      scriptIndex = $systemScripts.find(ele => { return ele.name === script});
      if(typeof scriptIndex !== 'undefined') {
        script = scriptIndex.script;
        let tmptext = text;
        let isfile = await App.FileExists(tmptext);
        let isDir = await App.DirExists(tmptext);
        if (isfile && !isDir) {
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
