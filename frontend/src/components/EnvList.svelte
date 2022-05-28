<div id='envlist'>
  <h1>Environments</h1>
  <ol>
    {#each envs as env}
      <li>
        <span class='envName'
              on:click={() => { openEnv(env); }}
        >{env}</span>
      </li>
    {/each}
  </ol>
  <div id='buttonRow'>
    <button 
      id='new'
      class='buttonStyle'
      style='background-color: {styles.editorBackground}; color: {styles.textcolor};'
      type='button'
      on:click={addNew}
    >
      New Environment
    </button>
    {#if createDefault}
      <button
        id='default'
        class='buttonStyle'
        style='background-color: {styles.editorBackground}; color: {styles.textcolor};'
        type='button'
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
    padding: 20px;
    user-select: none;
  }
  
  #envlist ol {
    margin: 0px;
    padding: 0px;
  }

  #envlist ol il {
    margin: 5px 0px;
    padding: 0px;
  }

  #envlist h1 {
    text-align: center;
    user-select: none;
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
    -webkit-tap-highlight-color:transparent;
    outline-style:none;
    cursor: pointer;
  }
</style>

<script>
  import { createEventDispatcher, onMount  } from 'svelte';
  
  export let config;
  export let ScriptPad;
  export let styles;
  
  let envs = [];
  let createDefault = false;

  const dispatch = createEventDispatcher();
  
  onMount(() => {
    //
    // Get the list of environments from the server.
    //
    getEnvList();

    //
    // See if the default has been created or not.
    //
    var def = envs.filter(item => item === 'Default');
    if(Array.isArray(def) && (def.length === 0)) {
      createDefault = true;
    }
  })
  
  function getEnvList() {
    envs = ScriptPad.getEnvNames();
  }
  
  function addNew() {
    dispatch('changeView', {
      view: 'env',
      config: {
        env: 'new',
        envVar: []
      }
    })
  }
  
  function openEnv(nm) {
    dispatch('changeView', {
      view: 'env',
      config: {
        env: nm
      }
    })
  }

  function createDefaultEnv() {
    //
    // Get the default environment
    //
    var defEnv = ScriptPad.createDefaultEnv();

    //
    // Save the default environment.
    //
    ScriptPad.addEnv(defEnv);

    //
    // Switch to the editing of the default environment.
    //
    dispatch('changeView', {
      view: 'env',
      config: {
        env: 'Default'
      }
    });
  }
</script>
