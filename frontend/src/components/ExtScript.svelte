<script>
  import { createEventDispatcher, onMount, tick } from "svelte";
  import { theme } from "../stores/theme.js";

  export let config;

  let script;
  let envs;

  const dispatch = createEventDispatcher();

  onMount(async () => {
    //
    // Get the list of external scripts from the server.
    //
    // extScripts {
    //    name     - User given name for the script
    //    script   - File name of the script
    //    path     - directory of the script
    //    env      - name of the environment
    // }
    //
    if (config.script !== "new") {
      script = await getExtScript(config.script);
    } else {
      script = {
        name: "new",
        script: "",
        path: "",
        env: "",
      };
    }
    envs = await getEnvNames();
  });

  async function getExtScript(name) {
    //
    // Get the named external script and return it.
    //
    let resp = await fetch(`http://localhost:9978/api/scripts/ext/${name}`, {
      headers: {
        "Content-type": "application/json",
      },
    });
    let scpt = await resp.json();
    return scpt;
  }

  async function getEnvNames() {
    //
    // Get the list of environments from the server.
    //
    let resp = await fetch(`http://localhost:9978/api/scripts/env/list`, {
      headers: {
        "Content-type": "application/json",
      },
    });
    let scpt = await resp.json();

    return scpt;
  }

  async function changeScript() {
    if (script.name !== "" && script.name !== null && script.name !== "new") {
      //
      // TODO: Add/Update the script
      //
      await fetch(`http://localhost:9978/api/scripts/ext/${script.name}`, {
        method: "PUT",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify(script),
      });
    }
  }

  async function deleteScript() {
    //
    // Remove the script
    //
    await fetch(`http://localhost:9978/api/scripts/ext/${script.name}`, {
      method: "DELETE",
      headers: {
        "Content-type": "application/json",
      },
    });
    tick();
    goback();
  }

  function goback() {
    dispatch("changeView", {
      view: "lists",
      config: {},
    });
  }
</script>

<div
  id="script"
  style="color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
>
  {#if typeof script !== "undefined" && typeof envs !== "undefined"}
    <label id="scriptName" for="scriptName"> Name of Script </label>
    <input
      id="scriptName"
      name="scriptName"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:blur={changeScript}
      bind:value={script.name}
    />
    <label id="scriptScript" for="scriptScript">
      What is the name of the Script file?
    </label>
    <input
      id="scriptScript"
      name="scriptScript"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:blur={changeScript}
      bind:value={script.script}
    />
    <label id="scriptPath" for="scriptPath">
      What is the directory for the script?
    </label>
    <input
      id="scriptPath"
      name="scriptPath"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:blur={changeScript}
      bind:value={script.path}
    />
    <label id="scriptEnv" for="scriptEnv">
      What is the environment for the script?
    </label>
    <select
      id="scriptEnv"
      name="scriptEnv"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:blur={changeScript}
      bind:value={script.env}
    >
      {#each envs as env}
        <option value={env}>{env}</option>
      {/each}
    </select>
  {/if}
  <div id="buttonRow">
    <button
      class="buttonStyle"
      type="button"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:click={() => {
        changeScript();
        goback();
      }}
    >
      Return
    </button>
    <button
      class="buttonStyle"
      type="button"
      style="border-radius: 5px; border-color: ${$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:click={deleteScript}
    >
      Delete
    </button>
  </div>
</div>

<style>
  #script {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px;
    height: 100%;
    width: 100%;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    width: 100%;
    padding: 0px;
    margin: 20px auto 0px auto;
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

  label,
  select {
    margin-top: 10px;
  }
</style>
