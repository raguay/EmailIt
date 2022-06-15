<script>
  import { createEventDispatcher, onMount } from "svelte";
  import { theme } from "../stores/theme.js";

  let scripts = [];

  const dispatch = createEventDispatcher();

  onMount(() => {
    //
    // Get the list of external scripts from the server.
    //
    scripts = listExtScripts();
  });

  function listExtScripts() {
    //
    // TODO: get the list of external scripts from the server.
    //
    return [];
  }

  function openScript(name) {
    dispatch("changeView", {
      view: "script",
      config: {
        script: name,
      },
    });
  }

  function addNew() {
    dispatch("changeView", {
      view: "script",
      config: {
        script: "new",
      },
    });
  }
</script>

<div id="scriptlist">
  <h1>External Scripts</h1>
  <ol>
    {#each scripts as script}
      <li>
        <span
          class="scriptName"
          on:click={() => {
            openScript(script);
          }}>{script}</span
        >
      </li>
    {/each}
  </ol>
  <div id="buttonRow">
    <button
      id="new"
      class="buttonStyle"
      type="button"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor};"
      on:click={addNew}
    >
      New Script
    </button>
  </div>
</div>

<style>
  #scriptlist {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 10px 20px 20px 20px;
    user-select: none;
  }

  #scriptlist ol {
    margin: 0px;
    padding: 0px;
  }

  #scriptlist ol il {
    margin: 5px 0px;
    padding: 0px;
  }

  #scriptlist h1 {
    text-align: center;
    user-select: none;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    margin: 10px 0px 0px 0px;
    padding: 0px;
  }

  .scriptName {
    cursor: pointer;
    user-select: none;
  }

  .buttonStyle {
    border-radius: 5px;
    border-color: black;
    font-size: 15px;
    height: 30px;
    text-shadow: 2px 2px 2px black;
    box-shadow: 2px 2px 5px 2px black;
    outline: none;
    margin: 0px 10px;
    padding: 3px 6px 6px 6px;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    -o-user-select: none;
    user-select: none;
    -webkit-tap-highlight-color: transparent;
    outline-style: none;
    cursor: pointer;
  }
</style>
