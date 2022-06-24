<script>
  import ExtScriptList from "./ExtScriptList.svelte";
  import ExtScript from "./ExtScript.svelte";
  import { theme } from "../stores/theme.js";

  let view = "lists";
  let config = {};

  function changeView(newView) {
    view = newView.detail.view;
    config = newView.detail.config;
  }
</script>

<div
  id="ExtScriptEditor"
  style="color: {$theme.textcolor}; background-color: {$theme.appBackground};"
>
  {#if view === "lists"}
    <div id="listviews">
      <ExtScriptList on:changeView={changeView} />
    </div>
  {:else if view === "script"}
    <ExtScript {config} on:changeView={changeView} />
  {/if}
</div>

<style>
  #ExtScriptEditor {
    display: flex;
    flex-direction: column;
    padding: 10px;
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
