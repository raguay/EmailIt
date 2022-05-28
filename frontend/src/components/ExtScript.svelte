<div id='script' style="color: {styles.textcolor}; background-color: {styles.appBackground};">
  {#if (typeof script !== 'undefined') && (typeof envs !== 'undefined')}
    <label id='scriptName'
           for='scriptName'>
      Name of Script
    </label>
    <input id='scriptName'
           name='scriptName'
           on:blur={changeScript}
           bind:value={script.name}>
    <label id='scriptScript'
           for='scriptScript'>
      What is the name of the Script file?
    </label>
    <input id='scriptScript'
           name='scriptScript'
           on:blur={changeScript}
           bind:value={script.script}>
    <label id='scriptPath'
           for='scriptPath'>
      What is the directory for the script?
    </label>
    <input id='scriptPath'
           name='scriptPath'
           on:blur={changeScript}
           bind:value={script.path}>
    <label id='scriptEnv'
           for='scriptEnv'>
      What is the environment for the script?
    </label>
    <select id='scriptEnv'
           name='scriptEnv'
           on:blur={changeScript}
           bind:value={script.env}>
      {#each envs as env}
        <option value="{env}">{env}</option>
      {/each}
    </select>
  {/if}
  <div id='buttonRow'>
    <button 
      class='buttonStyle'
      type='button'
      on:click={() => {
        changeScript();
        goback();
      }}
      style="background-color: {styles.editorBackground}; color: {styles.textcolor};"
    >
      Return
    </button>
    <button 
      class='buttonStyle'
      type='button'
      style='background-color: {styles.editorBackground}; color: {styles.textcolor};'
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

  .buttonRow {
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
  import { createEventDispatcher, onMount, tick  } from 'svelte';

  export let config;
  export let ScriptPad;
  export let styles;

  let script;
  let envs;
  
  const dispatch = createEventDispatcher();

  onMount(() => {
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
    if(config.script !== 'new') {
      script = ScriptPad.getExtScript(config.script);
    } else {
      script = {
        name: 'new',
        script: '',
        path: '',
        env: ''
      }
    }
    envs = ScriptPad.getEnvNames();
  });

  function changeScript() {
    if((script.name !== '')&&(script.name !== null)&&(script.name !== 'new')) {
      ScriptPad.addExtScript(script);
    }
  }
  
  function deleteScript() {
    ScriptPad.removeExtScript(script.name);
    tick();
    goback();
  }

  function goback() {
    dispatch('changeView',{
      view: 'lists',
      config: {}
    })
  }
</script>
