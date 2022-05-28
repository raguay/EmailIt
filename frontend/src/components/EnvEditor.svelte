<div 
  id='EnvEditor' 
  style="color: {$theme.textColor}; background-color: {$theme.backgroundColor};"
>
  {#if view === 'lists'}
    <div id='listviews'>
      <EnvList
        on:changeView={changeView}
      />
    </div>
  {:else if view === 'script'}
    <Env
      on:changeView={changeView}
    />
  {/if}
</div>

<style>
  #EnvEditor {
    display: flex;
    flex-direction: column;
    padding: 5px;
    margin: 0px;
  }

  #listviews {
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 0px;
    margin: 0px;
    overflow-y: scroll;
  }
</style>

<script>
  import { createEventDispatcher, onMount, tick  } from 'svelte';
  import EnvList from './EnvList.svelte';
  import Env from './Env.svelte';
  import { theme } from '../stores/theme.js';
  
  const dispatch = createEventDispatcher();
  
  let view = 'lists';
  let config = {};


  function changeView(newView) {
    view = newView.detail.view;
    config = newView.detail.config;
  }
</script>
