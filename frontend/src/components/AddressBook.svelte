<script>
  import { onMount } from "svelte";
  import { theme } from "../stores/theme.js";

  export let show;

  let emails = [];
  let addedit = false;
  let name = "";
  let email = "";

  onMount(() => {
    getEmails();
  });

  function getEmails(callback) {
    //
    // Get the emails from the server.
    //
    fetch("http://localhost:9978/api/emailit/emails", {
      method: "GET",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        emails = data.emails;
        if (typeof callback !== "undefined") callback();
      });
  }

  function newAddress() {
    console.log("add new...");
    addedit = true;
    name = "";
    email = "";
  }

  function closeAddressBook() {
    show = false;
  }

  function addNew() {
    email = email.trim();
    name = name.trim();
    emails = emails.filter((item) => item.email !== email);
    emails.push({
      name: name,
      email: email,
    });
    emails = emails;
    fetch("http://localhost:9978/api/emailit/addEmail", {
      method: "PUT",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        email: email,
      }),
    });
    addedit = false;
  }

  function editEmail(eemail) {
    console.log(eemail);
    addedit = true;
    name = eemail.name;
    email = eemail.email;
  }

  function deleteEmail(dem) {
    dem.email = dem.email.trim();
    emails = emails.filter((item) => item.email !== dem.email);
    fetch("http://localhost:9978/api/emailit/addEmail", {
      method: "DELETE",
      headers: {
        "Content-type": "application/json",
      },
      body: {
        name: "",
        email: dem.email,
      },
    });
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
        id="aename"
        bind:value={name}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      />
    </div>
    <div id="addeditrow">
      <label for="aeemail"> Email: </label>
      <input
        type="text"
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
