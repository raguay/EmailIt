<div 
  id='scripts' 
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  <div 
    id='header'
  >
    <div
      class='headerRow'
    >
      <label
        for='scriptName'
      >
        Name:
      </label>
      <SimpleAutoComplete 
        inputId="scriptName"
        items={list} 
        bind:selectedItem={scriptSel}
        inputClassName='scriptInput'
        className='scriptDiv'
        create=true
        theme={$theme}
        onChange={changeName}
        onCreate={(name) => {
          if((name !== undefined)&&(name !== '')) {
            scriptName = name;
            description = '';
            script = '';
            $scriptEditor.setValue(script);
          }
        }}
      />
      <label
        id='insertChkLab'
        for='insertChk'
      >
        Insert?
      </label>
      <input
        id='insertChk'
        name='insertChk'
        type=checkbox
        bind:checked={insert}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
    <div
      class='headerRow'
    >
      <label
        for='description'
      >
        Description:
      </label>
      <input
        id='description'
        bind:value={description}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
  </div>
  <CodeMirror height='400px' 
              width='970px' 
              config={editorConfig}
              initFinished={initFinished}
              on:textChange='{(event) => {textChanged(event.detail.data)}}' 
              on:editorChange='{(event) => {editorChange(event.detail.data); }}'
  />
  <div 
    id='buttonRow' 
  >
    <button
      on:click={saveScript}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Save Script
    </button>
    <button
      on:click={deleteScript}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Delete Script
    </button>
    <button
      on:click={viewNotes}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Notes
    </button>
    <button
      on:click={viewScriptsMenu}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Scripts Menu
    </button>
    <button
      on:click={viewTemplateMenu}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Templates Menu
    </button>
  </div>
</div>

<style> 
  #scripts {
    display: flex;
    flex-direction: column;
    padding: 10px;
    height: 100%;
    width: 100%;
  }
  
  #header {
    display: flex;
    flex-direction: column;
    margin: 0px 0px 0px 0px;
    width: 100%;
  }
  
  #description {
    width: 800px;
    border: solid 1px transparent;
    border-radius: 10px;
    padding: 5px 11px;
    height: 37px;
  }

  #insertChkLab {
    width: 90px;
    margin: auto 5px auto 10px;
  }

  #insertChk {
    margin: auto 5px;
  }

  :global(.scriptDiv) {
    width: 680px;
  }

  :global(.scriptInput) {
    border-radius: 10px;
  }


  .headerRow {
    display: flex;
    flex-direction: row;
    max-height: 40px;
    margin: 0px 0px 20px 0px;
    width: 100%;
  }

  .headerRow label {
    display: grid;
    justify-content: right;
    margin: auto 20px auto 0px;
    width: 150px;
  }
  
  #buttonRow {
    display: flex;
    flex-direction: row;
    margin: 10px auto;
    position: absolute;
    bottom: 0px;
    width: 100%;
    height: 40px;
  }
 
  #buttonRow button {
    border-radius: 10px;
    padding: 5px 10px 5px 10px;
    margin: 0px 5px;
    max-height: 40px;
    height: 40px;
    cursor: pointer;
  }
</style> 

<script> 
  import { onMount, tick } from 'svelte';
  import CodeMirror from '../components/CodeMirror.svelte';
  import SimpleAutoComplete from '../components/SimpleAutoComplete.svelte';
  import { state } from '../stores/state.js';
  import { theme } from '../stores/theme.js';
  import { showScripts } from '../stores/showScripts.js';
  import { showTemplates } from '../stores/showTemplates.js';
  import { scripts } from '../stores/scripts.js';
  import { scriptEditor } from '../stores/scriptEditor.js';
  
  let editorConfig = {
    language: 'javascript',
    lineNumbers: true,
    lineWrapping: true,
    lineHighlight: true
  };
  let initFinished = false;
  let script;
  let scriptSel = '';
  let scriptName = '';
  let description = '';
  let insert = false;
  let list;

  onMount(() => {
    // 
    // Load everything for working with the notes:
    //
    initFinished = true;
    getUserScripts(() => {
    });
  });

  function getUserScripts() {
    fetch('http://localhost:9978/api/scripts/user', {
        method: 'GET',
        headers: {
          'Content-type': 'application/json'
        }
      }).then(resp => {
        return resp.data;
      })
      .then(data => {
        list = data.scripts.sort();
        if(typeof callback !== 'undefined') callback();
      });
  }

  function getScript(name, callback) {
    if((name !== undefined)&&(name !== '')) {
      fetch(`http://localhost:9978/api/scripts/${name}`, {
          method: 'GET',
          headers: {
            'Content-type': 'application/json'
          }
        }).then(resp => {
          return resp.data;
        })
        .then(data => {
          scriptName = data.script.name;
          script = data.script.script;
          description = data.script.description;
          insert = data.script.insert;
          $scriptEditor.setValue(script);
          if(typeof callback !== 'undefined') callback();
        });
    }
  }

  function saveScript() {
    if((scriptName !== undefined)&&(scriptName !== '')) {
      fetch(`http://localhost:9978/api/scripts/${scriptName}`, {
          method: 'PUT',
          headers: {
            'Content-type': 'application/json'
          },
          body: Body.json({
            script: {
              name: scriptName,
              insert: insert,
              description: description,
              script: script
            }
          })
        }).then(() => {
          scriptSel = '';
          scriptName = "";
          insert = "";
          description = "";
          script = '';
          $scriptEditor.setValue(script);
          getScriptsList();
          getUserScripts();
        });
    }
  }
  
  function getScriptsList() {
    fetch('http://localhost:9978/api/scripts/list', {
        method: 'GET',
        headers: {
          'Content-type': 'application/json'
        }
      }).then(resp => {
        return resp.data;
      })
      .then(data => {
        $scripts = data.data;
        if(typeof callback !== 'undefined') callback();
      });
  }

  function deleteScript() {
    if((scriptName !== undefined)&&(scriptName !== '')) {
      fetch(`http://localhost:9978/api/scripts/${scriptName}`, {
          method: 'DELETE',
          headers: {
            'Content-type': 'application/json'
          }
        }).then(() => {
          scriptName = "";
          insert = "";
          description = "";
          script = '';
          scriptSel = '';
          $scriptEditor.setValue(script);
          getScriptsList();
          getUserScripts();
        });
    }
  }

  function editorChange(e) {
    $scriptEditor = e;
  }

  function textChanged(textCursor) {
    script = textCursor.value;
  }

  function focus() {
    if($scriptEditor !== null) {
      $scriptEditor.focus();
    }
  }

  function viewNotes() {
    $state = 'notes';
  }

  function viewScriptsMenu() {
    $showScripts = ! $showScripts;
  }

  function viewTemplateMenu() {
    $showTemplates = ! $showTemplates;
  }

  function changeName(newName) {
    if((newName !== undefined)&&(newName !== '')) {
      scriptName = newName;
      getScript(scriptName);
    }
  }
</script> 

