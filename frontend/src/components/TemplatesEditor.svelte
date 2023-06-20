<script>
  import { onMount } from "svelte";
  import CodeMirror from "../components/CodeMirror.svelte";
  import SimpleAutocomplete from "../components/SimpleAutocomplete.svelte";
  import { state } from "../stores/state.js";
  import { theme } from "../stores/theme.js";
  import { showScripts } from "../stores/showScripts.js";
  import { showTemplates } from "../stores/showTemplates.js";
  import { templates } from "../stores/templates.js";
  import { userTemplates } from "../stores/userTemplates.js";
  import { systemTemplates } from "../stores/systemTemplates.js";
  import { templateEditor } from "../stores/templateEditor.js";
  import { config } from "../stores/config.js";
  import * as App from "../../wailsjs/go/main/App.js";

  let editorConfig = {
    language: "javascript",
    lineNumbers: true,
    lineWrapping: true,
    lineHighlight: true,
  };
  let initFinished = false;
  let template = "";
  let templateDescription = "";
  let templateName = "";
  let templateSel;
  let list;

  onMount(() => {
    //
    // Load everything for working with the notes:
    //
    getUserTemplates(() => {
      initFinished = true;
    });
  });

  function getUserTemplates(callback) {
    list = $userTemplates.map((item) => item.name).sort();
    if (typeof callback !== "undefined") callback();
  }

  function getTemplate(name, callback) {
    if (name !== undefined && name !== "") {
      template = $userTemplates.find((item) => item.name === name);
      $templateEditor.setValue(template.template);
      templateDescription = template.description;
      if (typeof callback !== "undefined") callback();
    }
  }

  async function saveUserTemplates() {
    let templateloc = await App.AppendPath($config.configDir, "templates.json");
    await App.WriteFile(templateloc, JSON.stringify($userTemplates));
  }

  function saveTemplate() {
    if (templateName !== undefined && templateName !== "") {
      let newtemplate = $templateEditor.getValue();
      console.log(newtemplate);
      let templatedef = {
        name: templateName,
        description: templateDescription,
        template: newtemplate,
      };
      $userTemplates = $userTemplates.filter(
        (item) => item.name != templateName
      );
      $userTemplates.push(templatedef);
      $templates = $systemTemplates.concat($userTemplates);
      list = $userTemplates.map((item) => item.name).sort();
      template = "";
      templateName = "";
      templateDescription = "";
      templateSel = "";
      $templateEditor.setValue(template);
      saveUserTemplates();
    }
  }

  function deleteTemplate() {
    if (templateName !== undefined && templateName !== "") {
      $userTemplates = $userTemplates.filter(
        (item) => item.name != templateName
      );
      $templates = $systemTemplates.concat($userTemplates);
      list = $userTemplates.map((item) => item.name).sort();
      template = "";
      templateName = "";
      templateDescription = "";
      templateSel = "";
      $templateEditor.setValue(template);
      saveUserTemplates();
    }
  }

  function editorChange(e) {
    $templateEditor = e;
  }

  function textChanged(textCursor) {
    template = textCursor.value;
  }

  function focus() {
    if ($templateEditor !== null) {
      $templateEditor.focus();
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

  function nameChange(newname) {
    if (newname !== undefined) {
      templateName = newname;
      getTemplate(templateName);
    }
  }
</script>

<div
  id="templates"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  <div id="header">
    <div class="headerRow">
      <label for="templateName"> Name: </label>
      <SimpleAutocomplete
        inputId="templateName"
        items={list}
        bind:selectedItem={templateSel}
        inputClassName="templateInput"
        className="templateDiv"
        create="true"
        theme={$theme}
        onChange={nameChange}
        onCreate={(name) => {
          if (name !== undefined) {
            templateName = name;
            templateDescription = "";
            template = "";
            $templateEditor.setValue(template);
          }
        }}
      />
    </div>
    <div class="headerRow">
      <label for="description"> Description: </label>
      <input
        id="description"
        bind:value={templateDescription}
        autocomplete="off" spellcheck="false" autocorrect="off"
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
  </div>
  <CodeMirror
    height="375px"
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
  <div
    id="buttonRow"
    style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
  >
    <button
      on:click={saveTemplate}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Save Template
    </button>
    {#if templateName !== "Defaults"}
      <button
        on:click={deleteTemplate}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Delete Template
      </button>
    {/if}
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
  #templates {
    display: flex;
    flex-direction: column;
    padding: 10px;
    height: 100%;
    border-radius: 0px 0px 10px 10px;
  }

  #header {
    display: flex;
    flex-direction: column;
    margin: 0px 0px 0px 0px;
  }

  #description {
    width: 790px;
    border: solid 1px transparent;
    font-size: inherit;
    border-radius: 10px;
    padding: 5px 11px;
    height: 37px;
  }

  :global(.templateDiv) {
    width: 800px;
  }

  :global(.templateInput) {
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
    margin: 15px auto;
    bottom: 0px;
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
