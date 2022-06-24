<script>
  import { createEventDispatcher, onMount, tick } from "svelte";
  import EnvTableRow from "./EnvTableRow.svelte";
  import { theme } from "../stores/theme.js";

  export let config;

  let env;
  let addNew = false;
  let KVname = "";
  let KVvalue = "";

  const dispatch = createEventDispatcher();

  onMount(async () => {
    //
    // Get the list of environments from the server.
    //
    // scriptEnv {
    //    name    - Name of the environment
    //    envVar  - key,value array of environment variables
    //}
    //
    if (typeof config.env !== "undefined") {
      if (config.env !== "new" && config.env !== null) {
        if (config.env === "default") {
          config.env = "Default";
        }
        env = await getEnv(config.env);
      } else {
        //
        // Create a new one for the user to change.
        //
        env = {
          name: "new",
          envVar: {},
        };
      }
    }
  });

  async function getEnv(name) {
    //
    // get the environment from the server.
    //
    var resp = await fetch(`http://localhost:9978/api/scripts/env/${name}`);
    var list = await resp.json();
    return list;
  }

  async function changeEnv() {
    if (
      typeof env !== "undefined" &&
      env !== null &&
      env.name !== "" &&
      env.name !== null &&
      env.name !== "new"
    ) {
      await addEnv(env);
    }
  }

  async function addEnv(ev) {
    //
    // Add the new environment.
    //
    await fetch(`http://localhost:9978/api/scripts/env/${ev.name}`, {
      method: "PUT",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify(ev),
    });
  }

  async function deleteEnv() {
    if (
      typeof env !== "undefined" &&
      env.name !== null &&
      env.name !== "new" &&
      env.name !== ""
    ) {
      await removeEnv(env.name);
      tick();
      goback();
    }
  }

  async function removeEnv(envName) {
    //
    // Remove environment name.
    //
    await fetch(`http://localhost:9978/api/scripts/env/${envName}`, {
      method: "DELETE",
    });
  }

  function goback() {
    dispatch("changeView", {
      view: "lists",
      config: {},
    });
  }

  async function addKV() {
    env.envVar[KVname] = KVvalue;
    addNew = false;
    KVname = "";
    KVvalue = "";
    await changeEnv();
  }

  async function deleteRow(kv) {
    delete env.envVar[kv[0]];
    await changeEnv();
    env = await getEnv(config.env);
  }
</script>

<div
  id="env"
  style="color: {$theme.textColor}; background-color: {$theme.textareaColor};"
>
  {#if typeof env !== "undefined"}
    <div style="display: flex; flex-direction: row;">
      <label id="envName" for="envName"> Name of the Environment </label>
      <input
        id="envName"
        name="envName"
        style="padding-left: 10px; margin-top: -5px; margin-left: 20px; border-radius: 10px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
        on:blur={changeEnv}
        bind:value={env.name}
      />
    </div>
    {#if env.envVar !== "undefined"}
      <div id="EnvTable">
        <table>
          <thead>
            <tr>
              <th />
              <th />
              <th> Name </th>
              <th> Value </th>
            </tr>
          </thead>
          <tbody>
            {#if Object.entries(env.envVar).length > 0}
              {#each Object.entries(env.envVar) as kv}
                <EnvTableRow
                  name={kv[0]}
                  value={kv[1]}
                  on:deleteRow={() => {
                    deleteRow(kv);
                  }}
                  on:editRow={(item) => {
                    env.envVar[item.detail.name] = item.detail.value;
                    addKV();
                  }}
                />
              {/each}
            {/if}
            {#if addNew}
              <tr>
                <td>
                  <input
                    class="inputKV"
                    type="text"
                    style="background-color: {$theme.textAreaColor}; text= {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
                    bind:value={KVname}
                  />
                </td>
                <td>
                  <input
                    class="inputKV"
                    type="text"
                    style="background-color: {$theme.textAreaColor}; text= {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
                    bind:value={KVvalue}
                    on:blur={addKV}
                  />
                </td>
              </tr>
            {:else}
              <tr>
                <td span="2">
                  <span
                    class="addNewItem"
                    on:click={() => {
                      addNew = true;
                    }}
                  >
                    +
                  </span>
                </td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    {/if}
  {/if}
  <div id="buttonRow">
    <button
      id="goback"
      class="buttonStyle"
      type="button"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:click={() => {
        changeEnv();
        goback();
      }}
    >
      Return
    </button>
    <button
      class="buttonStyle"
      type="button"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-name: {$theme.font}; font-size: {$theme.fontSize};"
      on:click={deleteEnv}
    >
      Delete
    </button>
  </div>
</div>

<style>
  #env {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px;
    height: 100%;
    width: 100%;
  }

  th,
  td {
    min-width: 10px;
  }

  #EnvTable {
    display: flex;
    flex-direction: column;
    overflow: auto;
    width: 100%;
    height: 300px;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    width: 100%;
    padding: 0px;
    margin: 10px 0px auto 0px;
  }

  .buttonStyle {
    border-radius: 5px;
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

  td,
  th {
    min-width: 10px;
    text-align: left;
  }

  .addNewItem {
    color: red;
    cursor: pointer;
    font-size: 20px;
  }
</style>
