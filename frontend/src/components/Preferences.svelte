<script>
  import GenPrefs from "./GenPrefs.svelte";
  import ThemePrefs from "./ThemePrefs.svelte";
  import ExtScriptEditor from "./ExtScriptEditor.svelte";
  import EnvEditor from "./EnvEditor.svelte";
  import GitHub from "./GitHub.svelte";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";

  let currentPref = "General";
  let prefs = [
    {
      name: "General",
    },
    {
      name: "Theme",
    },
    {
      name: "External Scripts",
    },
    {
      name: "Environments",
    },
    {
      name: "GitHub",
    },
  ];

  function setNewCurrent(newCurrent) {
    currentPref = newCurrent;
  }

  function viewEmailit() {
    $state = "emailit";
  }

  function viewNotes() {
    $state = "notes";
  }
</script>

<div
  id="PreferencePanel"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  <ul id="tabs">
    {#each prefs as pref}
      {#if pref.name === currentPref}
        <li
          class="tabName"
          style="border-bottom: 3px {$theme.backgroundColor}; z-index: 100;"
        >
          {pref.name}
        </li>
      {:else}
        <li
          class="tabName"
          style="border-color: gray;"
          on:click={() => {
            setNewCurrent(pref.name);
          }}
        >
          {pref.name}
        </li>
      {/if}
    {/each}
  </ul>
  <div id="prefListWrap">
    {#if currentPref === "General"}
      <GenPrefs />
    {:else if currentPref === "Theme"}
      <ThemePrefs />
    {:else if currentPref === "External Scripts"}
      <ExtScriptEditor />
    {:else if currentPref === "Environments"}
      <EnvEditor />
    {:else if currentPref === "GitHub"}
      <GitHub />
    {/if}
  </div>
  <div id="buttonPanel">
    <button
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
      on:click={() => {
        viewEmailit();
      }}
    >
      EmailIt+
    </button>
    <button
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
      on:click={(event) => {
        viewNotes();
      }}
    >
      Notes
    </button>
  </div>
</div>

<style>
  #PreferencePanel {
    display: flex;
    flex-direction: column;
    height: 560px;
    width: 1000px;
    overflow: hidden;
    list-style: none;
    padding: 10px;
    margin: 0px;
    user-select: none;
    outline-style: none;
    z-index: 1;
    border-radius: 10px;
    border: 0px transparent;
  }

  #tabs {
    display: flex;
    flex-direction: row;
    margin: 5px 0px 5px 0px;
    padding: 0px;
    list-style-type: none;
    width: 460px;
    min-width: 460px;
    max-width: 460px;
    padding: 5px 0px 0px 10px;
  }

  #prefListWrap {
    overflow-y: scroll;
    overflow-x: hidden;
    height: 100%;
    padding: 20px;
    border: solid 3px;
    z-index: 0;
    border-radius: 10px;
  }

  #buttonPanel {
    flex-direction: row;
    align-content: center;
    margin: 10px auto;
  }

  #buttonPanel button {
    border-radius: 10px;
    border-color: transparent;
    outline: none;
    margin: 0px 10px;
    padding: 5px 10px 5px 10px;
    user-select: none;
    outline-style: none;
    cursor: pointer;
  }

  .tabName {
    list-style: none;
    margin: 0px 10px -8px 0px;
    padding: 10px 20px 5px 20px;
    cursor: pointer;
    border-top-left-radius: 30px;
    border-top-right-radius: 10px;
    border-top: solid 3px;
    border-right: solid 3px;
    border-left: solid 3px;
    border-bottom: solid 0px;
  }
</style>
