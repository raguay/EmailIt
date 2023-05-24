<script>
  import { onMount, afterUpdate } from "svelte";
  import { scripts } from "../stores/scripts.js";
  import { showScripts } from "../stores/showScripts.js";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";
  import { noteEditor } from "../stores/noteEditor.js";
  import { emailEditor } from "../stores/emailEditor.js";
  import { templateEditor } from "../stores/templateEditor.js";
  import { scriptEditor } from "../stores/scriptEditor.js";
  import { runscript } from "../stores/runscript.js";

  let searchIn = "";
  let list = [];
  let searchInput;
  let first = true;
  let cursor = 0;

  $: list = searchScripts(searchIn);
  $: setDefaults($showScripts);

  onMount(() => {});

  function setDefaults(flag) {
    if (flag) {
      cursor = 0;
      searchIn = "";
      searchInput.focus();
    }
  }

  afterUpdate(() => {
    if ($showScripts && first) {
      list = searchScripts(searchIn);
      first = false;
      searchInput.focus();
    }
  });

  function searchScripts(text) {
    var tmp = [];
    if (text === "" || text === undefined) {
      tmp = $scripts;
    } else {
      text = text.toLowerCase();
      tmp = $scripts.filter((item) => {
        if (item !== undefined && item !== null) {
          return item.name.toLowerCase().includes(text);
        }
        return false;
      });
    }
    tmp = tmp.sort((a, b) => a.name.toLowerCase() > b.name.toLowerCase());
    return tmp;
  }

  async function runScript(script) {
    var text = "";
    var selection = false;
    if ($state === "emailit") {
      //
      // Get the text in the email body.
      //
      if ($emailEditor.somethingSelected()) {
        selection = true;
        text = $emailEditor.getSelection();
      } else {
        text = $emailEditor.getValue();
      }
    } else if ($state === "notes") {
      //
      // Get the text from the current note.
      //
      if ($noteEditor.somethingSelected()) {
        selection = true;
        text = $noteEditor.getSelection();
      } else {
        text = $noteEditor.getValue();
      }
    } else if ($state === "scripts") {
      if ($scriptEditor.somethingSelected()) {
        selection = true;
        text = $scriptEditor.getSelection();
      } else {
        text = $scriptEditor.getValue();
      }
    } else if ($state === "templates") {
      if ($templateEditor.somethingSelected()) {
        selection = true;
        text = $templateEditor.getSelection();
      } else {
        text = $templateEditor.getValue();
      }
    }
    let datatext = await $runscript(script,text);
    if ($state === "emailit") {
          //
          // Paste the script in the body of the email.
          //
          if (selection) {
            $emailEditor.replaceSelection(datatext);
          } else {
            if (script.insert) {
              $emailEditor.insertAtCursor(datatext);
            } else {
              $emailEditor.setValue(datatext);
            }
          }
          $emailEditor.focus();
        } else if ($state === "notes") {
          //
          // Paste the script in the current note at the location.
          //
          if (selection) {
            $noteEditor.replaceSelection(datatext);
          } else {
            if (script.insert) {
              $noteEditor.insertAtCursor(datatext);
            } else {
              $noteEditor.setValue(datatext);
            }
          }
          $noteEditor.focus();
        } else if ($state === "scripts") {
          if (selection) {
            $scriptEditor.replaceSelection(datatext);
          } else {
            if (script.insert) {
              $scriptEditor.insertAtCursor(datatext);
            } else {
              $scriptEditor.setValue(datatext);
            }
          }
          $scriptEditor.focus();
        } else if ($state === "templates") {
          if (selection) {
            $templateEditor.replaceSelection(datatext);
          } else {
            if (script.insert) {
              $templateEditor.insertAtCursor(datatext);
            } else {
              $templateEditor.setValue(datatext);
            }
          }
          $templateEditor.focus();
        }
        $showScripts = false;
        searchIn = "";
  }

  async function keyDownProcessor(e) {
    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        cursor += 1;
        if (cursor >= list.length) cursor = list.length - 1;
        break;

      case "ArrowUp":
        e.preventDefault();
        cursor -= 1;
        if (cursor < 0) cursor = 0;
        break;

      case "Enter":
        e.preventDefault();
        await runScript(list[cursor]);
        break;

      case "Escape":
        $showScripts = false;
        break;
    }
  }
</script>

<div
  id="scriptmenu"
  style="left: {$showScripts
    ? '10px'
    : '-900px'}; background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
>
  <input
    type="text"
    bind:value={searchIn}
    bind:this={searchInput}
    on:keydown={keyDownProcessor}
    style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
  />
  <div id="list">
    <ul>
      {#if typeof $scripts === "object"}
        {#each list as script, key}
          <li
            on:click={async () => {
              await runScript(script);
            }}
            style="background-color: {cursor === key
              ? $theme.Purple
              : 'transparent'};"
            data-key={key}
          >
            {script.name}
          </li>
        {/each}
      {/if}
    </ul>
  </div>
</div>

<style>
  #scriptmenu {
    display: flex;
    flex-direction: column;
    padding: 10px;
    margin: 0px;
    position: absolute;
    z-index: 500;
    bottom: 70px;
    height: 500px;
    border-radius: 10px;
  }

  #list {
    display: flex;
    flex-direction: column;
    overflow: scroll;
    height: 450px;
  }

  ul {
    padding: 0px;
    margin: 0px;
  }

  li {
    text-decoration: none;
    list-style: none;
    cursor: pointer;
    padding: 0px 5px;
    margin: 0px;
    border-radius: 5px;
  }

  input {
    border-radius: 10px;
    width: 100%;
  }
</style>
