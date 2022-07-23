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

  let starting = true;

  onMount(() => {
    getScriptsList();
    getTermScriptsList();
    getTemplatesList();
    getTheme();
  });

  afterUpdate(() => {
    if (starting) {
      starting = false;
    }
  });

  function getTheme() {
    fetch("http://localhost:9978/api/theme", {
      method: "GET",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        $theme = data.theme;
        if (typeof callback !== "undefined") callback();
      });
  }

  function getScriptsList() {
    fetch("http://localhost:9978/api/scripts/list", {
      method: "GET",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        $scripts = data.data;
        if (typeof callback !== "undefined") callback();
      });
  }

  function getTermScriptsList() {
    fetch("http://localhost:9978/api/scripts/term/list", {
      method: "GET",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        $termscripts = data.data;
        console.log($termscripts);
        if (typeof callback !== "undefined") callback();
      });
  }

  function getTemplatesList() {
    fetch("http://localhost:9978/api/template/list", {
      method: "GET",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        $templates = data.templates;
        if (typeof callback !== "undefined") callback();
      });
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
<div id="dragbar" data-wails-drag />

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
  <div>
    <h1>Something went wront!</h1>
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
  }
</style>
