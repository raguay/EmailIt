<script>
  import { createEventDispatcher, onMount, tick } from "svelte";
  import { theme } from "../stores/theme.js";
  import { config } from "../stores/config.js";
  import * as App from '../../wailsjs/go/main/App.js';

  let envs = [];
  let envlist = null;
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
    var def = envs.find((item) => {
      return item.name === "Default";
    });
    if (typeof def === 'undefined') {
      createDefault = true;
    }
    return(()=>{
      envlist = null;
    })
  });

  async function getEnvList() {
    //
    // Get the list of environment names from the server.
    //
    if(envlist === null) {
      let envlistloc = await App.AppendPath($config.configDir, "environments.json");
      envlist = await App.ReadFile(envlistloc);
      envlist = JSON.parse(envlist);
    }
    envs = envlist;
  }

  async function saveEnvironments() {
    let envlistloc = await App.AppendPath($config.configDir, "environments.json");
    await App.WriteFile(envlistloc, JSON.stringify(envs));
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
        env: nm.name,
      },
    });
  }

  async function createDefaultEnv() {
    //
    // Get the default environment
    //
    let defenv = await App.GetEnvironment();
    let newArray = {};
    for(let i=0; i<defenv.length; i++){
      let parts = defenv[i].split('=');
      newArray[parts[0]] = parts[1];
    }
    envs.push({
      name: "Default",
      envVar: newArray
    });
    saveEnvironments();
    envlist = null;
    await tick();
    openEnv({
      name: "Default",
      envVar: [],
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
          }}>{env.name}</span
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
