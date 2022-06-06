<script>
  import { onMount, afterUpdate } from "svelte";
  import { templates } from "../stores/templates.js";
  import { showTemplates } from "../stores/showTemplates.js";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";
  import { noteEditor } from "../stores/noteEditor.js";
  import { emailEditor } from "../stores/emailEditor.js";
  import { templateEditor } from "../stores/templateEditor.js";
  import { scriptEditor } from "../stores/scriptEditor.js";

  let search = "";
  let list = [];
  let searchInput;
  let first = true;
  let cursor = 0;

  $: list = searchTempaltes(search);
  $: setDefaults($showTemplates);

  onMount(() => {});

  function setDefaults(flag) {
    if (flag) {
      cursor = 0;
      search = "";
      searchInput.focus();
    }
  }

  afterUpdate(() => {
    if ($showTemplates && first) {
      list = searchTempaltes(search);
      first = false;
      searchInput.focus();
    }
  });

  function searchTempaltes(text) {
    var tmp = [];
    if (text === "" || text === undefined) {
      tmp = $templates;
    } else {
      text = text.toLowerCase();
      tmp = $templates.filter((item) => {
        if (item !== undefined && item !== null) {
          return item.toLowerCase().includes(text);
        }
        return false;
      });
    }
    tmp = tmp.sort((a, b) => a.toLowerCase() > b.toLowerCase());
    return tmp;
  }

  function runTemplate(template) {
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
    fetch("http://localhost:9978/api/template/run", {
      method: "PUT",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        template: template,
        text: text,
      }),
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        console.log(data);
        if ($state === "emailit") {
          //
          // Paste the template in the body of the email.
          //
          if (selection) {
            $emailEditor.replaceSelection(data.text);
          } else {
            $emailEditor.setValue(data.text);
          }
          $emailEditor.focus();
        } else if ($state === "notes") {
          //
          // Paste the template in the current note at the location.
          //
          if (selection) {
            $noteEditor.replaceSelection(data.text);
          } else {
            $noteEditor.setValue(data.text);
          }
          $noteEditor.focus();
        } else if ($state === "scripts") {
          if (selection) {
            $scriptEditor.replaceSelection(data.text);
          } else {
            $scriptEditor.setValue(data.text);
          }
          $scriptEditor.focus();
        } else if ($state === "templates") {
          if (selection) {
            $templateEditor.replaceSelection(data.text);
          } else {
            $templateEditor.setValue(data.text);
          }
          $templateEditor.focus();
        }
        $showTemplates = false;
        search = "";
      });
  }

  function keyDownProcessor(e) {
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
        runTemplate(list[cursor]);
        break;

      case "Escape":
        $showTemplates = false;
        break;
    }
  }
</script>

<div
  id="templatemenu"
  style="left: {$showTemplates
    ? '10px'
    : '-900px'}; background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
>
  <input
    type="text"
    bind:value={search}
    bind:this={searchInput}
    on:keydown={keyDownProcessor}
    style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
  />
  <div id="list">
    <ul>
      {#if typeof $templates === "object"}
        {#each list as template, key}
          <li
            on:click={() => {
              runTemplate(template);
            }}
            style="background-color: {cursor === key
              ? $theme.Purple
              : 'transparent'};"
            data-key={key}
          >
            {template}
          </li>
        {/each}
      {/if}
    </ul>
  </div>
</div>

<style>
  #templatemenu {
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
