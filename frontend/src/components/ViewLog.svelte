
<div 
  id='log' 
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  {#if (typeof log) === 'undefined'}
    <p>No Logs</p>
  {:else}
    <div 
      id='logItems' 
      bind:this={logDiv}
    >
      {@html log }
    </div>
  {/if}
  <div 
    id='buttonRow' 
  >
    <button
      on:click={viewEmailIt}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      EmailIt
    </button>
    <button
      on:click={viewNotes}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Notes
    </button>
    <button
      on:click={viewNodeRed}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Node-Red
    </button>
    <button
      on:click={viewNodeRedDashboard}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Node-Red Dashboard
    </button>
  </div>
</div>

<style> 
  #log {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
  }

  :global(#log p) {
    margin: 0px;
    padding: 0px;
  }

  #logItems {
    display: flex;
    flex-direction: column;
    overflow: scroll;
    flex-grow: 5;
    margin: 20px 20px 0px 20px;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    margin: 10px auto;
  }
  
  #buttonRow button {
    border-radius: 10px;
    padding: 5px 20px 5px 20px;
    margin: 0px 20px;
    width: 100%;
    max-height: 40px;
    height: 40px;
    width: auto;
    cursor: pointer;
  }

</style> 

<script> 
  import { onMount, afterUpdate } from 'svelte';
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";

  let log;
  let logDiv;

  onMount(() => {
    getLog();
  });

  afterUpdate(() => {
    // 
    // Keep the log scrolled to the bottom.
    //
    if(typeof logDiv !== 'undefined') {
      logDiv.scrollTop = logDiv.scrollHeight;
    }
  })

  function getLog(callback) {
    // 
    // Get the emails from the server.
    // 
    fetch('http://localhost:9978/api/messages', {
      method: 'GET',
      headers: {
        'Content-type': 'application/json'
      }
    }).then(resp => {
      return resp.data;
    })
    .then(dataArray => {
      log = dataToHtml(dataArray.data);
      if(typeof callback !== 'undefined') callback();
    });
  }

  function dataToHtml(data) {
    var result = '';
    data.map(item => result += '<p>' + item.replaceAll(' ','&nbsp;').replaceAll('\n','<br />').replaceAll('\t','&nbsp;&nbsp;') + '</p>')
    return(result);
  }

  function viewEmailIt() {
    state.set('emailit');
  }

  function viewNotes() {
    state.set('notes');
  }

  async function viewNodeRed() {
    // open('http://localhost:9978/red/admin');
  }
  
  async function viewNodeRedDashboard() {
    // open('http://localhost:9978/red/api/ui');
  }
</script> 

