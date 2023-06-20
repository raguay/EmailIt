<script>
  import { onMount } from "svelte";
  import { theme } from "../stores/theme.js";
  import { config } from "../stores/config.js";
  import * as App from '../../wailsjs/go/main/App.js';

  export let show;

  let emails = [];
  let addedit = false;
  let name = "";
  let email = "";

  onMount(async () => {
    await getEmails();
  });

  async function saveEmails() {
    let emailfileloc = await App.AppendPath($config.configDir,"emails.json");
    await App.WriteFile(emailfileloc, JSON.stringify(emails));
  }

  async function getEmails() {
    //
    // Get the emails from the system.
    //
    let emailfileloc = await App.AppendPath($config.configDir,"emails.json");
    if(await App.FileExists(emailfileloc)) {
      let emailfilejson = await App.ReadFile(emailfileloc);
      emails = JSON.parse(emailfilejson);
    }
  }

  function newAddress() {
    addedit = true;
    name = "";
    email = "";
  }

  function closeAddressBook() {
    show = false;
  }

  async function addNew() {
    email = email.trim();
    name = name.trim();
    emails = emails.filter((item) => item.email !== email);
    emails.push({
      name: name,
      email: email,
    });
    emails = emails;
    await saveEmails();
    addedit = false;
  }

  function editEmail(eemail) {
    name = eemail.name;
    email = eemail.email;
  }

  async function deleteEmail(dem) {
    dem.email = dem.email.trim();
    emails = emails.filter((item) => item.email !== dem.email);
    await saveEmails();
  }
</script>

<div
  id="addressBook"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border-color: {$theme.borderColor};"
>
  <div id="tablediv">
    <table>
      <thead>
        <th>Name</th> <th>Address</th> <th /> <th />
      </thead>
      <tbody>
        {#each emails as pemail}
          <tr>
            <td>{pemail.name}</td>
            <td>{pemail.email}</td>
            <td
              class="iconClick"
              on:click={() => {
                editEmail(pemail);
              }}
            >
              🖋
            </td>
            <td
              class="iconClick"
              on:click={() => {
                deleteEmail(pemail);
              }}
            >
              ❌
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
  <div id="buttonrow">
    <button
      on:click={newAddress}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      New
    </button>
    <button
      on:click={closeAddressBook}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Close
    </button>
  </div>
</div>

{#if addedit}
  <div
    id="addeditdiv"
    style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize}; border-color: {$theme.textAreaColor};"
  >
    <div id="addeditrow">
      <label for="aename"> Name: </label>
      <input
        type="text"
        autocomplete="off" spellcheck="false" autocorrect="off"
        id="aename"
        bind:value={name}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
    <div id="addeditrow">
      <label for="aeemail"> Email: </label>
      <input
        type="text"
        autocomplete="off" spellcheck="false" autocorrect="off"
        id="aeemail"
        bind:value={email}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
    <div id="buttonrow">
      <button
        on:click={addNew}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Save
      </button>
      <button
        on:click={() => {
          addedit = false;
        }}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Close
      </button>
    </div>
  </div>
{/if}

<style>
  #addressBook {
    display: flex;
    position: absolute;
    top: 5%;
    left: 5%;
    flex-direction: column;
    border: solid 3px transparent;
    border-radius: 10px;
    max-height: 80%;
    z-index: 100;
  }

  #tablediv {
    overflow: scroll;
  }

  #buttonrow {
    display: flex;
    flex-direction: row;
    margin: 10px auto;
  }

  #buttonrow button {
    border-radius: 10px;
    margin: 10px;
    cursor: pointer;
  }

  #addeditdiv {
    display: flex;
    flex-direction: column;
    position: absolute;
    top: 15%;
    left: 30%;
    border: solid 3px transparent;
    border-radius: 10px;
    z-index: 200;
    padding: 20px 20px 10px 20px;
  }

  #addeditdiv input {
    margin: 10px 0px 10px 0px;
    border-radius: 10px;
  }

  #addeditrow {
    display: flex;
    flex-direction: column;
  }

  .iconClick {
    cursor: pointer;
  }
</style>
