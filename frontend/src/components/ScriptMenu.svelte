<div 
  id='scriptmenu'
  style="left: { $showScripts ? '10px' : '-900px'}; background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
>
  <input
    type="text"
    bind:value={searchIn}
    bind:this={searchInput}
    on:keydown={keyDownProcessor}
    style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border: solid 3px {$theme.borderColor};"
  />
  <div 
    id='list'
  >
    <ul>
      {#if typeof $scripts === 'object'}
        {#each list as script, key}
          <li
            on:click={() => { runScript(script); }}
            style="background-color: {cursor === key ? $theme.Purple : 'transparent'};"
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
    diplay: flex;
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

<script>
  import { onMount, afterUpdate } from 'svelte';
  import { scripts } from '../stores/scripts.js';
  import { showScripts } from '../stores/showScripts.js';
  import { theme } from '../stores/theme.js';
  import { state } from '../stores/state.js';
  import { noteEditor } from '../stores/noteEditor.js';
  import { emailEditor } from '../stores/emailEditor.js';
  import { templateEditor } from '../stores/templateEditor.js';
  import { scriptEditor } from '../stores/scriptEditor.js';

  let searchIn = "";
  let list = [];
  let searchInput;
  let first = true;
  let cursor = 0;

  $: list = searchScripts(searchIn);
  $: setDefaults($showScripts);

  onMount(() => {
  });

  function setDefaults(flag) {
    if(flag) {
      cursor = 0;
      searchIn = '';
      searchInput.focus();
    }
  }

  afterUpdate(() => {
    if($showScripts&&first) {
      list = searchScripts(searchIn);
      first = false;
      searchInput.focus();
    }
  });

  function searchScripts(text) {
    var tmp = [];
    if((text === '')||(text === undefined)) {
      tmp = $scripts;
    } else {
      text = text.toLowerCase();
      tmp = $scripts.filter(item => {
        if((item !== undefined)&&(item !== null)) {
          return item.name.toLowerCase().includes(text);
        }
        return false;
      });
    }
    tmp = tmp.sort((a,b) => a.name.toLowerCase() > b.name.toLowerCase());
    return(tmp);
  }

  function runScript(script) {
    var text = '';
    var selection = false;
    if($state === 'emailit') {
      //
      // Get the text in the email body.
      //
      if($emailEditor.somethingSelected()) {
        selection = true;
        text = $emailEditor.getSelection();
      } else {
        text = $emailEditor.getValue();
      }
    } else if($state === 'notes') {
      //
      // Get the text from the current note.
      //
      if($noteEditor.somethingSelected()) {
        selection = true;
        text = $noteEditor.getSelection();
      } else {
        text = $noteEditor.getValue();
      }
    } else if($state === 'scripts') {
      if($scriptEditor.somethingSelected()) {
        selection = true;
        text = $scriptEditor.getSelection();
      } else {
        text = $scriptEditor.getValue();
      }
    } else if($state === 'templates') {
      if($templateEditor.somethingSelected()) {
        selection = true;
        text = $templateEditor.getSelection();
      } else {
        text = $templateEditor.getValue();
      }
    }
    fetch('http://localhost:9978/api/script/run', {
        method: 'PUT',
        headers: {
          'Content-type': 'application/json'
        },
        body: Body.json({
          script: script.name,
          text: text
        })
      }).then(resp => {
        return resp.data;
      })
      .then(data => {
        if($state === 'emailit') {
          //
          // Paste the script in the body of the email.
          //
          if(selection) {
            $emailEditor.replaceSelection(data.text);
          } else {
            if(script.insert) {
              $emailEditor.insertAtCursor(data.text);
            } else {
              $emailEditor.setValue(data.text);
            }
          }
          $emailEditor.focus();
        } else if($state === 'notes') {
          //
          // Paste the script in the current note at the location.
          //
          if(selection) {
            $noteEditor.replaceSelection(data.text);
          } else {
            if(script.insert) {
              $noteEditor.insertAtCursor(data.text);
            } else {
              $noteEditor.setValue(data.text);
            }
          }
          $noteEditor.focus();
        } else if($state === 'scripts') {
          if(selection) {
            $scriptEditor.replaceSelection(data.text);
          } else {
            if(script.insert) {
              $scriptEditor.insertAtCursor(data.text);
            } else {
              $scriptEditor.setValue(data.text);
            }
          }
          $scriptEditor.focus();
        } else if($state === 'templates') {
          if(selection) {
            $templateEditor.replaceSelection(data.text);
          } else {
            if(script.insert) {
              $templateEditor.insertAtCursor(data.text);
            } else {
              $templateEditor.setValue(data.text);
            }
          }
          $templateEditor.focus();
        }
        $showScripts = false;
        searchIn = '';
      });
  }

  function keyDownProcessor(e) {
    switch(e.key) {
      case 'ArrowDown':
        e.preventDefault();
        cursor += 1;
        if(cursor >= list.length) cursor = list.length - 1;
      break;

      case 'ArrowUp':
        e.preventDefault();
        cursor -= 1;
        if(cursor < 0) cursor = 0;
      break;

      case 'Enter':
        e.preventDefault();
        runScript(list[cursor]);
      break;

      case 'Escape':
        $showScripts = false;
      break;
    }
  }
</script>
