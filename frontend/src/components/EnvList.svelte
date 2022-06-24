<script>
  import { createEventDispatcher, onMount } from "svelte";
  import { theme } from "../stores/theme.js";

  let envs = [];
  let createDefault = false;

  const dispatch = createEventDispatcher();

  onMount(async () => {
    //
    // Get the list of environments from the server.
    //
    await getEnvList();

    //
    // See if the default has been created or not.
    //
    var def = envs.filter((item) => {
      return item == "Default";
    });
    if (Array.isArray(def) && def.length === 0) {
      createDefault = true;
    }
  });

  async function getEnvList() {
    //
    // Get the list of environment names from the server.
    //
    var resp = await fetch("http://localhost:9978/api/scripts/env/list");
    envs = await resp.json();
  }

  async function addEnv(env) {
    //
    // Add the new environment.
    //
    await fetch(`http://localhost:9978/api/scripts/env/${env.name}`, {
      method: "PUT",
      body: JSON.stringify(env),
    });
  }

  async function addNew() {
    let newEnv = {
      view: "env",
      config: {
        env: "new",
        envVar: [],
      },
    };
    dispatch("changeView", newEnv);
  }

  function openEnv(nm) {
    dispatch("changeView", {
      view: "env",
      config: {
        env: nm,
      },
    });
  }

  async function createDefaultEnv() {
    //
    // Get the default environment
    //

    await fetch("http://localhost:9978/api/scripts/env/Default", {
      method: "PUT",
    });
    //
    // Switch to the editing of the default environment.
    //
    dispatch("changeView", {
      view: "env",
      config: {
        env: "Default",
      },
    });
  }
</script>

<div id="envlist">
  <h2>Environments</h2>
  <ol>
    {#each envs as env}
      <li>
        <span
          class="envName"
          on:click={() => {
            openEnv(env);
          }}>{env}</span
        >
      </li>
    {/each}
  </ol>
  <div id="buttonRow">
    <button
      id="new"
      class="buttonStyle"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      type="button"
      on:click={addNew}
    >
      New Environment
    </button>
    {#if createDefault}
      <button
        id="default"
        class="buttonStyle"
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
        type="button"
        on:click={createDefaultEnv}
      >
        Create Default
      </button>
    {/if}
  </div>
</div>

<style>
  #envlist {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px;
    user-select: none;
  }

  #envlist ol {
    margin: 0px;
    padding: 0px;
    overflow-y: auto;
    height: 260px;
  }

  #envlist ol il {
    margin: 5px 0px;
    padding: 0px;
  }

  #envlist h2 {
    text-align: center;
    user-select: none;
    margin: 5px;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    margin: 10px 0px 0px 0px;
    padding: 0px;
  }

  .envName {
    user-select: none;
    cursor: pointer;
  }

  .buttonStyle {
    border-radius: 5px 5px 5px 5px;
    border-color: black;
    font-size: 15px;
    height: 40px;
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
