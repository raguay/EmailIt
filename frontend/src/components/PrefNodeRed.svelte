<div>
  <div class="labelRadialHolder">
    <label class="styleLabel">Auto Launch Node-Red?</label>
    <input class="PrefCheck" 
           type="checkbox" 
           bind:checked={preferences.REDAuto} 
           on:change="{() => {changeFire();}}" 
    />
  </div>
  <div class="labelRadialHolder">
    <label class="styleLabel">Dashboard configured?</label>
    <input class="PrefCheck" 
           type="checkbox" 
           bind:checked={preferences.REDDashboard} 
           on:change="{() => {changeFire();}}" 
    />
  </div>
</div>

<style>
  .styleLabel {
    font-size: 18px;
    margin: 0px 0px 0px 0px;
    user-select: none;
  }

  .labelRadialHolder {
    display: flex;
    flex-direction: row;
  }

  .PrefCheck {
    margin: 6px 0px 0px 10px;
    font-size: 20px;
  }
</style>

<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  var preferences = {
    REDAuto: true,
    REDDashboard: true
  }

  function fire(name, data) {
    dispatch(name, {
      data: data
    });
  }

  function changeFire() {
    fire('changeTodos', {
      auto: preferences.REDAuto,
      dashboard: preferences.REDDashboard
    });
  }
</script>
