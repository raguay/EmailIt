<script>
  import { createEventDispatcher, onMount, tick } from "svelte";
  import { theme } from "../stores/theme.js";
  import { config } from "../stores/config.js";
  import { extScripts } from "../stores/extScripts.js";
  import * as App from '../../wailsjs/go/main/App.js';

  export let cfg;

  let script = {
    name: "",
    script: "",
    path: "",
    env: "",
    description: "",
    help: "",
    termscript: false,
  };
  let envs = [];
  let envlist = null;

  const dispatch = createEventDispatcher();

  onMount(async () => {
    //
    // Get the list of external scripts from the server.
    //
    // extScripts {
    //    name         - User given name for the script
    //    script       - File name of the script
    //    path         - directory of the script
    //    env          - name of the environment
    //    description  - A description of what the script does.
    //    help         - A help message for the script
    //    termscript   - True if it's a termianl script. Otherwise, it's false.
    // }
    //
    if (cfg.script !== "new") {
      script = $extScripts.find(item => item.name === cfg.script);
    } else {
      script = {
        name: "new",
        script: "",
        path: "",
        env: "",
        termscript: false,
        description: "",
        help: "",
      };
    }

    //
    // Get the list of environments from the server.
    //
    let envlistloc = await App.AppendPath($config.configDir, "environments.json");
    envlist = await App.ReadFile(envlistloc);
    envlist = JSON.parse(envlist);
    envs = envlist.map(item => item.name);
  });

  async function changeScript() {
    if (script.name !== "" && script.name !== null && script.name !== "new") {
      //
      // Add/Update the script
      //
      $extScripts = $extScripts.filter(item => item.name !== script.name);
      $extScripts.push(script);
      let extScriptsLoc = await App.AppendPath($config.configDir,"extscripts.json");
      await App.WriteFile(extScriptsLoc, JSON.stringify($extScripts));
    }
  }

  async function deleteScript() {
    //
    // Remove the script
    //
    $extScripts = $extScripts.filter(item => item.name !== script.name);
    let extScriptsLoc = await App.AppendPath($config.configDir,"extscripts.json");
    await App.WriteFile(extScriptsLoc, JSON.stringify($extScripts));
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
    <label id="scrptNameLab" for="scrptName"> Name of Script </label>
    <input
      id="scrptName"
      name="scrptName"
      bind:value={script.name}
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:blur={changeScript}
    />
    <label id="scrptDesLab" for="scrptDes">Description</label>
    <input
      id="scrptDes"
      name="scrptDes"
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      bind:value={script.description}
      on:blur={changeScript}
    />
    <label id="helpLab" for="help">Help</label>
    <input
      id="help"
      name="help"
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      bind:value={script.help}
      on:blur={changeScript}
    />
    <label id="scriptScript" for="scriptScript">
      What is the name of the Script file?
    </label>
    <input
      id="scriptScript"
      name="scriptScript"
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      bind:value={script.script}
      on:blur={changeScript}
    />
    <label id="scriptPath" for="scriptPath">
      What is the directory for the script?
    </label>
    <input
      id="scriptPath"
      name="scriptPath"
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      bind:value={script.path}
      on:blur={changeScript}
    />
    <label id="scriptEnv" for="scriptEnv">
      What is the environment for the script?
    </label>
    <select
      id="scriptEnv"
      name="scriptEnv"
      style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      bind:value={script.env}
      on:blur={changeScript}
    >
      {#each envs as env}
        <option value={env}>{env}</option>
      {/each}
    </select>
    <label id="termScriptLab" for="termScriptChk"> Terminal Script? </label>
    <input
      id="termScriptChk"
      name="termScriptChk"
      type="checkbox"
      bind:checked={script.termscript}
      on:blur={changeScript}
    />
</div>
<div id="buttonRow">
  <button
    class="buttonStyle"
    type="button"
    style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
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
    style="border-radius: 5px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
    on:click={deleteScript}
  >
    Delete
  </button>
</div>

<style>
  #script {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px;
    max-height: 100%;
    width: 100%;
    overflow-y: auto;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    width: 100%;
    padding: 0px;
    margin: 20px auto 0px auto;
  }

  #scrptNameLab {
    padding-top: 0px;
    margin-top: 0px;
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
