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
  import * as App from '../wailsjs/go/main/App.js';

  let starting = true;

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
    let defaultThemeDir = await App.AppendPath(themedir,"Default");
    if(!await App.DirExists(defaultThemeDir)) {
      await App.MakeDir(defaultThemeDir)
    }
    defaultThemefile = await App.AppendPath(defaultThemeDir,"Default.json");
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
    themecnfg = await App.AppendPath(defaultThemedir,"package.json");
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
    $termscripts = $userScripts.filter(value => value.termscript === true).map(value => {
        return { name: value.name, insert: value.insert }
      }).concat($systemScripts.filter(value => value.termscript === true).map(value => {
        return { name: value.name, insert: value.insert }
      })).concat($extScripts.filter(value => value.termscript === true).map(value => {
        return { name: value.name, insert: false }
      }));
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
