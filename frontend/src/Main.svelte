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
  import { Jumper } from "svelte-loading-spinners";

  let starting = true;
  let count = 0;
  let loadItemsCount = 0;
  const loadNumber = 4;

  onMount(async () => {
    //
    // Wait to give the server time to start.
    //
    $state = "notready";
    await wait(1000);

    //
    // Get stuff from the server.
    //
    getScriptsList();
    getTermScriptsList();
    getTemplatesList();
    getTheme();
  });

  function wait(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        console.log("Done waiting");
        resolve(ms);
      }, ms);
    });
  }

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
        loadItemsCount++;
        if (loadItemsCount === loadNumber) {
          $state = "emailit";
        }
        if (typeof callback !== "undefined") callback();
      })
      .catch(async () => {
        count++;
        await wait(100);
        getTheme();
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
        loadItemsCount++;
        if (loadItemsCount === loadNumber) {
          $state = "emailit";
        }
        if (typeof callback !== "undefined") callback();
      })
      .catch(async () => {
        count++;
        await wait(100);
        getScriptsList();
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
        loadItemsCount++;
        if (loadItemsCount === loadNumber) {
          $state = "emailit";
        }
        if (typeof callback !== "undefined") callback();
      })
      .catch(async () => {
        count++;
        await wait(100);
        getTermScriptsList();
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
        loadItemsCount++;
        if (loadItemsCount === loadNumber) {
          $state = "emailit";
        }
        if (typeof callback !== "undefined") callback();
      })
      .catch(async () => {
        count++;
        await wait(100);
        getTemplatesList();
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
  <div
    style="display: flex; flex-direction: column; width: 1000px; height: 500px;"
  >
    <div
      style="display: flex; flex-direction: row; color: white; margin: auto;"
    >
      <h1 style="margin: auto 60px auto 100px;">Waiting on the Server!</h1>
      <span style="margin: 10px; display: none;">{count}</span>
      <Jumper size="60" color="#80ffea" unit="px" duration="1s" />
    </div>
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
