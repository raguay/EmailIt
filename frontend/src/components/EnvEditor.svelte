<script>
  import EnvList from "./EnvList.svelte";
  import Env from "./Env.svelte";
  import { theme } from "../stores/theme.js";

  let view = "lists";
  let confg = {
    env: "Default",
  };

  function changeView(newView) {
    view = newView.detail.view;
    confg = newView.detail.config;
  }
</script>

<div
  id="EnvEditor"
  style="color: {$theme.textColor}; background-color: {$theme.backgroundColor};"
>
  {#if view === "lists"}
    <div id="listviews">
      <EnvList on:changeView={changeView} />
    </div>
  {:else if view === "env"}
    <Env {confg} on:changeView={changeView} />
  {/if}
</div>

<style>
  #EnvEditor {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 5px;
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
