<script>
  import { onMount } from "svelte";
  import CodeMirror from "../components/CodeMirror.svelte";
  import SimpleAutocomplete from "../components/SimpleAutocomplete.svelte";
  import { state } from "../stores/state.js";
  import { theme } from "../stores/theme.js";
  import { showScripts } from "../stores/showScripts.js";
  import { showTemplates } from "../stores/showTemplates.js";
  import { scriptEditor } from "../stores/scriptEditor.js";
  import { config } from "../stores/config.js";
  import { scripts } from "../stores/scripts.js";
  import { userScripts } from "../stores/userScripts.js";
  import { extScripts } from "../stores/extScripts.js";
  import { systemScripts } from "../stores/systemScripts.js";
  import { termscripts } from "../stores/termscripts.js";
  import * as App from "../../wailsjs/go/main/App.js";

  let editorConfig = {
    language: "javascript",
    lineNumbers: true,
    lineWrapping: true,
    lineHighlight: true,
  };
  let initFinished = false;
  let script;
  let scriptSel = "";
  let scriptName = "";
  let description = "";
  let language = "JavaScript";
  let insert = false;
  let termscript = false;
  let list;

  onMount(() => {
    //
    // Load everything for working with the scripts
    //
    getUserScripts();
    initFinished = true;
  });

  function getUserScripts(callback) {
    list = $userScripts.map((item) => item.name).sort();
    if (typeof callback !== "undefined") callback();
  }

  async function getScript(name, callback) {
    if (name !== undefined && name !== "") {
      let data = $userScripts.find((item) => item.name === name);
      if (typeof data !== "undefined") {
        scriptName = data.name;
        script = data.script;
        description = data.description;
        language = data.language;
        insert = data.insert == "true" ? true : false;
        if (typeof data.termscript === "undefined") {
          termscript = false;
        } else {
          termscript = data.termscript;
        }
        if (language === "") {
          language = "JavaScript";
          await saveScript();
        }
        if (language !== "Lips") {
          editorConfig.language = language;
        } else {
          editorConfig.language = "Scheme";
        }
        $scriptEditor.setValue(script);
      }
      if (typeof callback !== "undefined") callback();
    }
  }

  async function saveScript() {
    if (scriptName !== undefined && scriptName !== "") {
      let scriptstruct = {
        name: scriptName,
        insert: insert,
        description: description,
        script: script,
        language: language,
        termscript: termscript,
        help: description,
      };
      $userScripts = $userScripts.filter((item) => item.name !== scriptName);
      $userScripts.push(scriptstruct);
      await saveUserScripts();
      getUserScripts();
    }
  }

  async function deleteScript() {
    $userScripts = $userScripts.filter((item) => item.name !== scriptName);
    scriptName = "";
    scriptSel = "";
    insert = true;
    description = "";
    language = "JavaScript";
    script = "";
    scriptSel = "";
    $scriptEditor.setValue(script);
    await saveUserScripts();
    getUserScripts();
  }

  async function saveUserScripts() {
    let userScriptsLoc = await App.AppendPath(
      $config.configDir,
      "scripts.json",
    );
    await App.WriteFile(userScriptsLoc, JSON.stringify($userScripts));
    $scripts = $userScripts
      .filter((value) => value.termscript === false)
      .map((value) => {
        return { name: value.name, insert: value.insert };
      })
      .concat(
        $systemScripts
          .filter((value) => value.termscript === false)
          .map((value) => {
            return { name: value.name, insert: value.insert };
          }),
      )
      .concat(
        $extScripts
          .filter((value) => value.termscript === false)
          .map((value) => {
            return { name: value.name, insert: false };
          }),
      );
    getUserScripts();
    getTermScriptsList();
  }

  function getTermScriptsList() {
    $termscripts = $userScripts
      .filter((value) => value.termscript === true)
      .concat($systemScripts.filter((value) => value.termscript === true))
      .concat($extScripts.filter((value) => value.termscript === true));
  }

  function editorChange(e) {
    $scriptEditor = e;
  }

  function textChanged(textCursor) {
    script = textCursor.value;
  }

  function focus() {
    if ($scriptEditor !== null) {
      $scriptEditor.focus();
    }
  }

  function viewNotes() {
    $state = "notes";
  }

  function viewScriptsMenu() {
    $showScripts = !$showScripts;
  }

  function viewTemplateMenu() {
    $showTemplates = !$showTemplates;
  }

  function viewScriptLine() {
    $state = "scriptline";
  }

  function changeName(newName) {
    if (newName !== undefined && newName !== "") {
      scriptName = newName;
      getScript(scriptName);
    }
  }
</script>

<div
  id="scripts"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  <div id="header">
    <div class="headerRow">
      <label for="scriptName"> Name: </label>
      <SimpleAutocomplete
        inputId="scriptName"
        items={list}
        bind:selectedItem={scriptSel}
        inputClassName="scriptInput"
        className="scriptDiv"
        create="true"
        theme={$theme}
        onChange={changeName}
        onCreate={(name) => {
          if (name !== undefined && name !== "") {
            scriptName = name;
            description = "";
            script = "";
            language = "JavaScript";
            $scriptEditor.setValue(script);
          }
        }}
      />
      <label id="insertChkLab" for="insertChk"> Insert? </label>
      <input
        id="insertChk"
        name="insertChk"
        type="checkbox"
        bind:checked={insert}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
    <div class="headerRow">
      <label for="description"> Description: </label>
      <input
        id="description"
        autocomplete="off"
        spellcheck="false"
        autocorrect="off"
        bind:value={description}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize}"
      />
      <label id="termScriptLab" for="termScriptChk"> Terminal Script? </label>
      <input
        id="termScriptChk"
        autocomplete="off"
        spellcheck="false"
        autocorrect="off"
        name="termScriptChk"
        type="checkbox"
        bind:checked={termscript}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize}"
      />
    </div>
    <div class="headerRow">
      <label id="scriptlangLab" for="scriptlang"> Language: </label>
      <select
        id="scriptlang"
        name="scriptlangSel"
        bind:value={language}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize}"
      >
        {#if language === "JavaScript"}
          <option name="JavaScript" value="JavaScript" selected>
            JavaScript
          </option>
          <option name="Prolog" value="Prolog">Prolog</option>
          <option name="Lips" value="Lips">Lips</option>
        {:else if language === "Prolog"}
          <option name="JavaScript" value="JavaScript">JavaScript</option>
          <option name="Prolog" value="Prolog" selected>Prolog</option>
          <option name="Lips" value="Lips">Lips</option>
        {:else}
          <option name="JavaScript" value="JavaScript">JavaScript</option>
          <option name="Prolog" value="Prolog">Prolog</option>
          <option name="Lips" value="Lips" selected>Lips</option>
        {/if}
      </select>
    </div>
  </div>
  <CodeMirror
    height="340px"
    width="980px"
    config={editorConfig}
    {initFinished}
    on:textChange={(event) => {
      textChanged(event.detail.data);
    }}
    on:editorChange={(event) => {
      editorChange(event.detail.data);
    }}
  />
  <div id="buttonRow">
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
      on:click={viewScriptLine}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Script Line
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
    border-radius: 0px 0px 10px 10px;
  }

  #header {
    display: flex;
    flex-direction: column;
    margin: 0px 0px 0px 0px;
    width: 100%;
  }

  #description {
    width: 620px;
    border: solid 1px transparent;
    border-radius: 10px;
    padding: 5px 11px;
    height: 37px;
  }

  #insertChkLab,
  #termScriptLab {
    width: 100px;
    margin: auto 5px auto 10px;
  }

  #insertChk,
  #termScriptChk {
    margin: auto 0px auto 5px;
  }

  :global(.scriptDiv) {
    width: 640px;
  }

  :global(.scriptInput) {
    border-radius: 10px;
  }

  .headerRow {
    display: flex;
    flex-direction: row;
    max-height: 40px;
    margin: 0px 0px 15px 0px;
    width: 100%;
  }

  .headerRow label {
    display: grid;
    justify-content: right;
    margin: auto 20px auto 15px;
    width: 155px;
  }

  #buttonRow {
    display: flex;
    flex-direction: row;
    padding: 20px 10px 10px 10px;
    margin: 5px auto;
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
