<script>
  import { onMount } from "svelte";
  import ColorPicker from "./ColorPicker.svelte";
  import { config } from "../stores/config.js";
  import { theme } from "../stores/theme.js";
  import * as App from "../../wailsjs/go/main/App.js";

  let colorchange = "";
  let colorID = 0;
  let showPicker = false;
  let buttonColor = false;
  let keepNewInput = false;
  let pickerType = "";
  let style = "Default";
  let explanation;
  let themeList;

  onMount(async () => {
    themeList = await getStyleList();
    style = $config.theme;
  });

  async function styleSelectorChange() {
    if (style === "New") {
      keepNewInput = true;
      $theme.name = style;
    } else {
      $theme = await getStyle(style);
      $theme.name = style;

      //
      // Save the new theme.
      //
      $config.theme = $theme.name;
      let configloc = await App.AppendPath($config.configDir, "config.json");
      await App.WriteFile(configloc, JSON.stringify($config));
    }
  }

  async function getStyle(nm) {
    let theme = {};
    //
    // Load in the themes.
    //
    let thmdir = await App.AppendPath($config.themeDir, nm);
    let thmprojfile = await App.AppendPath(thmdir, "package.json");
    thmprojfile = await App.ReadFile(thmprojfile);
    thmprojfile = JSON.parse(thmprojfile);
    let thmfile = await App.AppendPath(thmdir, thmprojfile.theme.main);
    theme = await App.ReadFile(thmfile);
    theme = JSON.parse(theme);
    return theme;
  }

  async function getStyleList() {
    //
    // get the list of themes.
    //
    let list = await App.ReadDir($config.themeDir);
    list = list.filter((item) => item.IsDir).map((item) => item.Name);
    return list;
  }

  function setColor(id, color) {
    if (buttonColor) {
      $theme.buttons[id].color = color;
    } else {
      switch (id) {
        case "textAreaColor":
          $theme.textAreaColor = color;
          break;

        case "backgroundColor":
          $theme.backgroundColor = color;
          break;

        case "textColor":
          $theme.textColor = color;
          break;

        case "borderColor":
          $theme.borderColor = color;
          break;

        case "Cyan":
          $theme.Cyan = color;
          break;

        case "Green":
          $theme.Green = color;
          break;

        case "Orange":
          $theme.Orange = color;
          break;

        case "Pink":
          $theme.Pink = color;
          break;

        case "Purple":
          $theme.Purple = color;
          break;

        case "Red":
          $theme.Red = color;
          break;

        case "Yellow":
          $theme.Yellow = color;
          break;

        case "functionColor":
          $theme.functionColor = color;
          break;

        case "stringColor":
          $theme.stringColor = color;
          break;

        case "constantColor":
          $theme.constantColor = color;
          break;

        case "keywordColor":
          $theme.keywordColor = color;
          break;

        case "highlightBackgroundColor":
          $theme.highlightBackgroundColor = color;
          break;

        case "selectionColor":
          $theme.selectionColor = color;
          break;

        default:
          console.log("Invalid id.");
          break;
      }
    }
    showPicker = false;
    buttonColor = false;
  }

  function changeColor(button) {
    colorchange = button.color;
    colorID = button.id;
    showPicker = true;
    if (isNaN(colorID)) {
      buttonColor = false;
      pickerType = colorID;
      explanation = `${colorID}`;
    } else {
      buttonColor = true;
      pickerType = "Circle";
      explanation = `#${colorID} Circle`;
    }
  }

  function saveNewTheme() {
    keepNewInput = false;
    updateTheme();
  }

  async function updateTheme() {
    //
    // Save the theme.
    //
    let thmdir = await App.AppendPath($config.themeDir, style);
    if (await App.DirExists(thmdir)) {
      //
      // This is updating an existing theme.
      //
      let thmprojfile = await App.AppendPath(thmdir, "package.json");
      thmprojfile = await App.ReadFile(thmprojfile);
      thmprojfile = JSON.parse(thmprojfile);
      let thmfile = await App.AppendPath(thmdir, thmprojfile.theme.main);
      await App.WriteFile(thmfile, JSON.stringify($theme));
      let basethmfile = await App.AppendPath($config.configDir, "theme.json");
      await App.WriteFile(basethmfile, JSON.stringify($theme));
    } else {
      //
      // Here, we are actually creating a new theme.
      //
      await App.MakeDir(thmdir);
      let thmprojfile = await App.AppendPath(thmdir, "package.json");
      await App.WriteFile(
        thmprojfile,
        `
{
  "name": "${style}",
  "version": "1.0.0",
  "description": "",
  "keywords": [
    "modalfilemanager", "theme"
  ],
  "author": "",
  "license": "MIT",
  "theme": {
    "name": "${style}",
    "description": "",
    "type": 0,
    "github": "",
    "main": "${style}.json"
  }
}
      `
      );
      let thmfile = await App.AppendPath(thmdir, `${style}.json`);
      await App.WriteFile(thmfile, JSON.stringify($theme));
      let basethmfile = await App.AppendPath($config.configDir, "theme.json");
      await App.WriteFile(basethmfile, JSON.stringify($theme));

      //
      // Save the new theme.
      //
      $config.theme = style;
      let configloc = await App.AppendPath($config.configDir, "config.json");
      await App.WriteFile(configloc, JSON.stringify($config));
    }
    themeList = await getStyleList();
  }

  async function deleteStyle() {
    //
    // Delete the theme.
    //
    let thmdir = await App.AppendPath($config.themeDir, style);
    await App.DeleteEntries(thmdir);

    //
    // Get a new list of themes and set the current theme to the first theme.
    //
    themeList = await getStyleList();
    style = themeList[0];

    //
    // Save the new theme in the configuration file.
    //
    $config.theme = style;
    let configloc = await App.AppendPath($config.configDir, "config.json");
    await App.WriteFile(configloc, JSON.stringify($config));

    //
    // Load the new theme.
    //
    $theme = await getStyle(style);
    $theme.name = style;

    //
    // Save to the active theme.
    //
    let basethmfile = await App.AppendPath($config.configDir, "theme.json");
    await App.WriteFile(basethmfile, JSON.stringify($theme));
  }
</script>

<div id="themeName">
  <label>Name of theme: </label>
  {#if keepNewInput}
    <input
      class="prefInput"
      autocomplete="off" spellcheck="false" autocorrect="off"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor};"
      bind:value={style}
      on:change={() => {
        keepNewInput = true;
      }}
    />
    <button
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
      on:click={saveNewTheme}
    >
      Save New
    </button>
  {:else}
    <select
      class="prefSelector"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor};"
      bind:value={style}
      on:change={() => {
        styleSelectorChange();
      }}
    >
      {#if themeList !== undefined}
        {#each themeList as item}
          <option value={item}>{item}</option>
        {/each}
      {/if}
      <option value="New">New</option>
    </select>
    <button
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
      on:click={updateTheme}
    >
      Update Theme
    </button>
    <button
      style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
      on:click={deleteStyle}
    >
      Delete Theme
    </button>
  {/if}
</div>
<h3>Circle Colors</h3>
<div class="circlePickersWrap">
  {#each $theme.buttons as button, key}
    <div class="circlePickerWrap">
      <label class="circlePickerLabel">#{key}</label>
      <div
        class="circlePicker"
        on:click={(event) => {
          changeColor(button);
        }}
        style="background-color: {button.color};"
      />
      <label class="circlePickerLabel">{button.color}</label>
    </div>
  {/each}
</div>
<h3>Various Other Colors</h3>
<div id="variousOtherColorsDiv">
  <label class="variousPickerLabel1"> textAreaColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "textAreaColor",
        color: $theme.textAreaColor,
      });
    }}
    style="background-color: {$theme.textAreaColor};"
  />
  <label class="variousPickerLabel2">{$theme.textAreaColor}</label>
  <label class="variousPickerLabel1"> backgroundColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "backgroundColor",
        color: $theme.backgroundColor,
      });
    }}
    style="background-color: {$theme.backgroundColor};"
  />
  <label class="variousPickerLabel2">{$theme.backgroundColor}</label>
  <label class="variousPickerLabel1"> textColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "textColor",
        color: $theme.textColor,
      });
    }}
    style="background-color: {$theme.textColor};"
  />
  <label class="variousPickerLabel2">{$theme.textColor}</label>
  <label class="variousPickerLabel1"> borderColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "borderColor",
        color: $theme.borderColor,
      });
    }}
    style="background-color: {$theme.borderColor};"
  />
  <label class="variousPickerLabel2">{$theme.borderColor}</label>
  <label class="variousPickerLabel1"> functionColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "functionColor",
        color: $theme.functionColor,
      });
    }}
    style="background-color: {$theme.functionColor};"
  />
  <label class="variousPickerLabel2">{$theme.functionColor}</label>
  <label class="variousPickerLabel1"> stringColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "stringColor",
        color: $theme.stringColor,
      });
    }}
    style="background-color: {$theme.stringColor};"
  />
  <label class="variousPickerLabel2">{$theme.stringColor}</label>
  <label class="variousPickerLabel1"> constantColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "constantColor",
        color: $theme.constantColor,
      });
    }}
    style="background-color: {$theme.constantColor};"
  />
  <label class="variousPickerLabel2">{$theme.constantColor}</label>
  <label class="variousPickerLabel1"> keywordColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "keywordColor",
        color: $theme.keywordColor,
      });
    }}
    style="background-color: {$theme.keywordColor};"
  />
  <label class="variousPickerLabel2">{$theme.keywordColor}</label>
  <label class="variousPickerLabel1"> highlightBackgroundColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "highlightBackgroundColor",
        color: $theme.highlightBackgroundColor,
      });
    }}
    style="background-color: {$theme.highlightBackgroundColor};"
  />
  <label class="variousPickerLabel2">{$theme.highlightBackgroundColor}</label>
  <label class="variousPickerLabel1"> selectionColor </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "selectionColor",
        color: $theme.selectionColor,
      });
    }}
    style="background-color: {$theme.selectionColor};"
  />
  <label class="variousPickerLabel2">{$theme.selectionColor}</label>
  <label class="variousPickerLabel1"> Cyan </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Cyan",
        color: $theme.Cyan,
      });
    }}
    style="background-color: {$theme.Cyan};"
  />
  <label class="variousPickerLabel2">{$theme.Cyan}</label>
  <label class="variousPickerLabel1"> Green </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Green",
        color: $theme.Green,
      });
    }}
    style="background-color: {$theme.Green};"
  />
  <label class="variousPickerLabel2">{$theme.Green}</label>
  <label class="variousPickerLabel1"> Orange </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Orange",
        color: $theme.Orange,
      });
    }}
    style="background-color: {$theme.Orange};"
  />
  <label class="variousPickerLabel2">{$theme.Orange}</label>
  <label class="variousPickerLabel1"> Pink </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Pink",
        color: $theme.Pink,
      });
    }}
    style="background-color: {$theme.Pink};"
  />
  <label class="variousPickerLabel2">{$theme.Pink}</label>
  <label class="variousPickerLabel1"> Purple </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Purple",
        color: $theme.Purple,
      });
    }}
    style="background-color: {$theme.Purple};"
  />
  <label class="variousPickerLabel2">{$theme.Purple}</label>
  <label class="variousPickerLabel1"> Red </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Red",
        color: $theme.Red,
      });
    }}
    style="background-color: {$theme.Red};"
  />
  <label class="variousPickerLabel2">{$theme.Red}</label>
  <label class="variousPickerLabel1"> Yellow </label>
  <div
    class="circlePicker"
    on:click={(event) => {
      changeColor({
        id: "Yellow",
        color: $theme.Yellow,
      });
    }}
    style="background-color: {$theme.Yellow};"
  />
  <label class="variousPickerLabel2">{$theme.Yellow}</label>
</div>
{#if $theme !== undefined}
  <ColorPicker
    item={pickerType}
    explainText={explanation}
    color={colorchange}
    id={colorID}
    show={showPicker}
    on:colorChanged={(event) => {
      setColor(event.detail.data.id, event.detail.data.color);
    }}
    on:quitColorPicker={(event) => {
      showPicker = false;
    }}
  />
{/if}

<style>
  #themeName {
    display: flex;
    flex-direction: row;
    align-items: center;
  }

  #variousOtherColorsDiv {
    display: grid;
    grid-auto-flow: row dense;
    grid-template-columns: 400px 40px 80px;
    grid-column-gap: 10px;
    grid-row-gap: 20px;
    margin-left: 10px;
  }

  button {
    border-radius: 10px;
    border-color: transparent;
    outline: none;
    margin: 0px 10px;
    padding: 5px 10px 5px 10px;
    user-select: none;
    outline-style: none;
    cursor: pointer;
  }

  .prefInput {
    font-size: 15px;
    border-radius: 5px;
    box-shadow: inset 0px 0px 5px 2px #0f0a16;
    border: 2px #0f0a16;
    min-height: 20px;
    padding: 10px 10px 10px 10px;
    margin: 5px 5px 0px 5px;
    outline: none;
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    -o-user-select: text;
    user-select: text;
  }

  .prefSelector {
    font-size: 15px;
    border-radius: 5px;
    box-shadow: inset 0px 0px 5px 2px #0f0a16;
    border: 2px #0f0a16;
    min-height: 20px;
    padding: 10px 10px 10px 10px;
    margin: 5px 5px 0px 5px;
    outline: none;
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    -o-user-select: text;
    user-select: text;
  }

  .prefColorDiv {
    display: flex;
    flex-direction: column;
    margin: 20px 0px 0px 0px;
  }

  .prefColorHolder {
    display: flex;
    flex-flow: row;
  }

  .circlePickersWrap {
    display: flex;
    flex-flow: column wrap;
    padding: 0px 0px 0px 10px;
    height: 165px;
    align-items: left;
    overflow-y: auto;
  }

  .circlePickerWrap {
    display: flex;
    flex-direction: row;
    margin: 10px auto;
  }

  .variousPickerLabel1 {
    font-size: 25px;
    margin: 0px 10px 0px 10px;
    grid-column: 1;
    text-align: right;
  }

  .variousPickerLabel2 {
    font-size: 25px;
    margin: 0px 10px 0px 10px;
    user-select: text;
    grid-column: 3;
  }

  .circlePickerLabel {
    font-size: 25px;
    margin: 0px 10px 0px 10px;
    user-select: text;
  }

  .styleLabel {
    font-size: 18px;
    margin: 0px 0px 0px 0px;
    user-select: none;
  }

  .circlePicker {
    height: 30px;
    min-height: 30px;
    width: 30px;
    min-width: 30px;
    border-radius: 30px;
    cursor: pointer;
    border: solid 1px white;
    grid-column: 2;
  }
</style>
