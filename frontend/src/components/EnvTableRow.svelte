<script>
  import { createEventDispatcher, afterUpdate } from "svelte";
  import { theme } from "../stores/theme.js";

  const dispatch = createEventDispatcher();

  export let name;
  export let value;

  var editValues = false;
  var nameDOM;

  afterUpdate(() => {
    if (editValues) {
      nameDOM.focus();
    }
  });

  function changeData() {
    editValues = false;
    dispatch("editRow", {
      name: name,
      value: value,
    });
  }

  function delData() {
    dispatch("deleteRow", {
      name: name,
      value: value,
    });
  }
</script>

<tr>
  {#if editValues}
    <td /><td />
    <td>
      <input
        bind:this={nameDOM}
        type="text"
        style="padding-left: 10px; border-radius: 10px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-family: '{$theme.font}'; font-size: {$theme.fontSize};"
        bind:value={name}
      />
    </td>
    <td>
      <input
        type="text"
        style="padding-left: 10px; border-radius: 10px; border-color: {$theme.borderColor}; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; font-family: {$theme.font}; font-size: {$theme.fontSize};"
        bind:value
        on:blur={changeData}
      />
    </td>
  {:else}
    <td class="editTD">
      <span
        on:click={() => {
          editValues = true;
        }}
      >
        ✏️
      </span>
    </td>
    <td class="deleteTD">
      <span style="width: 10px; cursor: pointer;" on:click={delData}> ❌ </span>
    </td>
    <td>{name}</td>
    <td>{value}</td>
  {/if}
</tr>

<style>
  .editTD {
    max-width: 10px;
    width: 10px;
    cursor: pointer;
    user-select: none;
  }

  .deleteTD {
    max-width: 10px;
    width: 10px;
    cursor: pointer;
    user-select: none;
  }
</style>
