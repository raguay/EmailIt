<div id='env' style="color: {styles.textcolor}; background-color: {styles.appBackground};">
  {#if typeof env !== 'undefined'}
    <label id='envName'
           for='envName'>
      Name of the Environment
    </label>
    <input id='envName'
           name='envName'
           on:blur={changeEnv}
           bind:value={env.name}>
    {#if env.envVar !== 'undefined'}
      <div id='EnvTable'>
        <table>
          <thead>
            <tr>
              <th>
              </th>
              <th>
              </th>
              <th>
                Name
              </th>
              <th>
                Value
              </th>
            </tr>
          </thead>
          <tbody>
            {#each Object.entries(env.envVar) as kv}
              <EnvTableRow name={kv[0]} value={kv[1]} />
            {/each}
            {#if addNew}
              <tr><td>
                <input class="inputKV" type='text' bind:value={KVname} />
              </td><td>
                <input class="inputKV" type='text' bind:value={KVvalue} on:blur={addKV} />
              </td></tr>
            {:else}
              <tr><td span=2><span class="addNewItem" on:click={() => { addNew = true; }}>+</span></td></tr>
            {/if}
          </tbody>
        </table>
      </div>
    {/if}
  {/if}
  <div id='buttonRow'>
    <button id='goback'
      class='buttonStyle'
      type='button'
      style='background-color: {styles.editorBackground}; color: {styles.textcolor};'
      on:click={() => {
        changeEnv();
        goback();
      }}
    >
      Return
    </button>
    <button 
      class='buttonStyle'
      type='button'
      style='background-color: {styles.editorBackground}; color: {styles.textcolor};'
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

  #EnvTable {
    display: flex;
    flex-direction: column;
    overflow: auto;
    width: 100%;
    height: auto;
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

  td, th {
    min-width: 100px;
    text-align: left;
  }

  .addNewItem {
    color: red;
    cursor: pointer;
    font-size: 20px;
  }
</style>

<script>
  import { createEventDispatcher, onMount, tick  } from 'svelte';
  import EnvTableRow from './EnvTableRow.svelte';
  
  export let config;
  export let ScriptPad;
  export let styles;

  let env;
  let addNew = false;
  let KVname = '';
  let KVvalue = '';

  const dispatch = createEventDispatcher();

  onMount(() => {
    //
    // Get the list of external scripts from the server.
    //
    // scriptEnv {
    //    name    - Name of the environment
    //    envVar  - key,value array of environment variables
    //}
    //
    if(typeof config.env !== 'undefined') {
      if ((config.env !== 'new')&&(config.env !== null)) {
        if(config.env === 'default') {
          config.env = 'Default';
        }
        env = ScriptPad.getEnv(config.env);
      } else {
        //
        // Create a new one for the user to change.
        //
        env = {
          name: 'new',
          envVar: []
        };
      }
    }
  })

  function changeEnv() {
    if((typeof env !== 'undefined')&&(env.name !== '')&&(env.name !== null)&&(env.name !== 'new')) {
      ScriptPad.addEnv(env);
    }
  }
  
  function deleteEnv() {
    if((typeof env !== 'undefined')&&(env.name !== null)&&(env.name !== 'new')&&(env.name !== '')) {
      ScriptPad.removeEnv(env.name);
      tick();
      goback();
    }
  }

  function goback() {
    dispatch('changeView',{
      view: 'lists',
      config: {}
    })
  }

  function addKV() {
    env.envVar[KVname] = KVvalue;
    addNew = false;
    KVname = '';
    KVvalue = '';
    changeEnv();
  }
</script>
