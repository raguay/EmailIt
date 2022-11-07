<script>
  import { afterUpdate, onMount } from "svelte";
  import CodeMirror from "../components/CodeMirror.svelte";
  import AddressBook from "../components/AddressBook.svelte";
  import showdown from "showdown";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";
  import { account } from "../stores/account.js";
  import { emailaccounts } from "../stores/emailaccounts.js";
  import { email } from "../stores/email.js";
  import { emailEditor } from "../stores/emailEditor.js";
  import { commandLineEmail } from "../stores/commandLineEmail.js";
  import { showScripts } from "../stores/showScripts.js";
  import { showTemplates } from "../stores/showTemplates.js";
  import { config } from "../stores/config.js";
  import { runtemplate } from "../stores/runtemplate.js";
  import * as App from '../../wailsjs/go/main/App.js';
  import * as rt from "../../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  let receiver = "";
  let subject = "";
  let emailState = "edit";
  let showChangeAccount = false;
  let showNewAccount = false;
  let showPreview = false;
  let origAccount;
  let accountFrom = "";
  let accountName = "";
  let accountUsername = "";
  let accountSmptServer = "";
  let accountPassword = "";
  let accountPort = "";
  let accountSig = "";
  let accountHeaderHTML = "";
  let accountFooterHTML = "";
  let accountDefault = false;
  let previewHTML = "";
  let bodyValue = "";
  let oldState = "";
  let starting = false;
  let emails = [];
  let editorConfig = {
    language: "markdown",
    lineNumbers: false,
    lineWrapping: true,
    lineHighlight: true,
  };
  let initFinished = true;
  let badEmails = [];
  let elist = [];
  let alertTitle = "";
  let alertMsg = "";
  let showAlert = false;
  let showEmailList = false;
  let showAddressB = false;
  let receiverDOM;

  onMount(async () => {
    await getEmails();
    emailState = "edit";
    oldState = "edit";
    initFinished = true;
  });

  afterUpdate(() => {
    if ($commandLineEmail !== undefined && $commandLineEmail !== "") {
      receiver = $commandLineEmail;
      $email.to = $commandLineEmail;
      $commandLineEmail = undefined;
    }
    if (starting || $email.new) {
      receiver = $email.to;
      var rec = document.getElementById("receiverInput");
      rec.value = $email.to;
      subject.value = $email.subject;
      starting = false;
      $emailEditor.setValue($email.body);
      rt.WindowShow();
      $email = {
        account: 'Default',
        to: '',
        from: '',
        subject: '',
        new: false
      };
    }
    if ($emailEditor !== null) {
      if (emailState === "edit" && oldState === "preview") {
        $emailEditor.setValue(bodyValue);
        oldState = "edit";
      }
    }
  });

  function generateEmailList(e) {
    if (e !== undefined && (e.key === "Escape" || e.key === "Tab")) {
      showEmailList = false;
      if (e.key !== "Tab" && e.key !== "Escape") e.preventDefault();
    } else {
      var fullLine = receiver.toString().toLowerCase();
      var currentPart = "";
      var parts = fullLine.split(",");
      if (parts.length > 1) {
        currentPart = parts[parts.length - 1];
        elist = emails.filter((item) =>
          item.toString().toLowerCase().includes(currentPart.trim())
        );
      } else {
        elist = emails.filter((item) =>
          item.toString().toLowerCase().includes(fullLine)
        );
      }
    }
  }

  function addToInput(newEmail) {
    var parts = receiver.split(",");
    if (parts.length > 1) {
      receiver =
        parts
          .slice(0, parts.length - 1)
          .map((item) => item.trim())
          .join(",") +
        ", " +
        newEmail;
      receiverDOM.selectionStart = receiver.length;
    } else {
      receiver = newEmail;
    }
  }

  function editorChange(e) {
    $emailEditor = e;
  }

  function textChanged(textCursor) {
    $email.body = textCursor.value;
    bodyValue = textCursor.value;
  }

  async function getEmails() {
    //
    // Get the emails from the system.
    //
    let emailfileloc = await App.AppendPath($config.configDir,"emails.json");
    if(await App.FileExists(emailfileloc)) {
      let emailfilejson = await App.ReadFile(emailfileloc);
      emails = JSON.parse(emailfilejson);
      emails = emails.map((item) => {
       if (item.name === "") {
         return item.email;
       } else {
         return `${item.name.trim()} <${item.email.trim()}>`;
       }
      });
    }
  }

  function changeAccount() {
    origAccount = $account;
    showChangeAccount = !showChangeAccount;
  }

  function saveAccount() {
    if (emailState === "preview") {
      makeHtml();
    }
    showChangeAccount = false;
  }

  function clearFormData() {
    accountName = "";
    accountDefault = false;
    accountFrom = "";
    accountUsername = "";
    accountSmptServer = "";
    accountPassword = "";
    accountPort = "";
    accountSig = "";
    accountHeaderHTML = "";
    accountFooterHTML = "";
  }

  function newAccount() {
    clearFormData();
    showNewAccount = true;
  }

  function cancelAccountChange() {
    $account = origAccount;
    if (emailState === "edit") {
      bodyValue = $emailEditor.getValue();
    } else {
      makeHtml();
    }
    showChangeAccount = false;
  }

  function changeActiveAccount(acc) {
    //
    // Set new account and create previews.
    //
    $account = acc;
    if (emailState === "edit") {
      bodyValue = $emailEditor.getValue();
    } else {
      makeHtml();
    }
  }

  function createPreview() {
    var toAddress;
    if (typeof $email.to !== "undefined") {
      toAddress = $email.to;
    } else {
      var em = document.getElementById("receiverInput").value;
      toAddress = em;
    }
    if (validate(toAddress)) {
      //
      // Keep a copy of the body value for when we exit preview mode.
      //
      bodyValue = $emailEditor.getValue();

      //
      // Set to preview and keep a copy of the new state.
      //
      emailState = "preview";
      oldState = emailState;

      //
      // Creat a preview of the email.
      //
      makeHtml();

      //
      // Show the preview.
      //
      showPreview = true;
    } else {
      showInvalidEmails();
    }
  }

  function makeHtml() {
    var converter = new showdown.Converter({
      extensions: [],
    });
    converter.setOption("tables", true);
    let newBody = $runtemplate("given", bodyValue);
    previewHTML = converter.makeHtml(newBody + $account.signiture);
    if (typeof $account.headerHTML !== "undefined") {
      previewHTML = $account.headerHTML + previewHTML + $account.footerHTML;
    }
  }

  async function editEmail() {
    emailState = "edit";
    showPreview = false;
  }

  function cleanTags(msg) {
    return msg.replace(/<hr>/gimu, "\n").replace(/<[/]*[^>]+>/gimu, "");
  }

  async function sendEmail() {
    var toAddress;
    if (typeof $email.to !== "undefined") {
      toAddress = $email.to;
    } else {
      var em = document.getElementById("receiverInput").value;
      toAddress = em;
    }
    if (validate(toAddress)) {
      addToEmails(toAddress);
      var bodyText = bodyValue + cleanTags($account.signiture);
      bodyText = $runtemplate('given',bodyText);
      showPreview = false;
      emailState = "edit";

      // 
      // Send the email.
      //
      let result = await App.SendEmail($account.username, $account.from, $account.password, $account.smtpserver, $account.port, toAddress, previewHTML, bodyText, subject.value);
      if(result !== "Success") {
        alertTitle = "Sending Email Failed";
        alertMsg = result;
        showAlert = true;
      }
    } else {
      showInvalidEmails();
    }
  }

  function showInvalidEmails() {
    alertTitle = "Invalid Emails";
    alertMsg = "Please fix these emails!";
    showAlert = true;
  }

  function validate(emailLine) {
    var valid = true;
    badEmails = [];
    emailLine
      .toString()
      .split(",")
      .forEach((item) => {
        const svalid = validateSingle(item.trim());
        if (!svalid) {
          valid = false;
          badEmails.push(item);
        }
      });
    return valid;
  }

  function validateSingle(email) {
    if (email.includes("<")) {
      //
      // It's an email with a name infront. Get the email to check.
      //
      const nameRegexp = new RegExp(/^[^<]*\<([^\>]*)\>/);
      const matches = email.match(nameRegexp);
      email = matches[1];
    }

    //
    // Check the email itself.
    //
    const emailRegexp = new RegExp(
      /^[a-zA-Z0-9][\-_\.\+\!\#\$\%\&\'\*\/\=\?\^\`\{\|]{0,1}([a-zA-Z0-9][\-_\.\+\!\#\$\%\&\'\*\/\=\?\^\`\{\|]{0,1})*[a-zA-Z0-9]@[a-zA-Z0-9][-\.]{0,1}([a-zA-Z][-\.]{0,1})*[a-zA-Z0-9]\.[a-zA-Z0-9]{1,}([\.\-]{0,1}[a-zA-Z]){0,}[a-zA-Z0-9]{0,}$/i
    );

    return emailRegexp.test(email);
  }

  function addToEmails(emailLine) {
    emailLine
      .toString()
      .split(",")
      .forEach((item) => {
        var email = item;
        var name = "";
        if (item.includes("<")) {
          //
          // It's an email with a name infront. Get the email to check.
          //
          const nameRegexp = new RegExp(/^([^<]*)\<([^\>]*)\>/);
          const matches = item.match(nameRegexp);
          name = matches[1];
          email = matches[2];
        }
        addToEmailsSingle(name, email);
      });
  }

  function addToEmailsSingle(name, email) {
    email = email.trim();
    name = name.trim();
    emails = emails.filter((item) => item.email !== email);
    emails.push({
      name: name,
      email: email,
    });
    emails = emails;
  }

  function clearEmail() {
    var rec = document.getElementById("receiverInput");
    rec.value = "";
    receiver = "";
    subject.value = "";
    $emailEditor.setValue("");
    bodyValue = "";
    $email.to = "";
    $email.subject = "";
    $email.body = "";
    showPreview = false;
    emailState = "edit";
    oldState = "edit";
  }

  function saveNewAccount() {
    var acc = {
      name: accountName,
      default: accountDefault,
      from: accountFrom,
      username: accountUsername,
      smtpserver: accountSmptServer,
      port: accountPort,
      password: accountPassword,
      signiture: accountSig,
      headerHTML: accountHeaderHTML,
      footerHTML: accountFooterHTML,
    };

    //
    // If this is to be the default account, make sure all the others
    // are false.
    //
    if (accountDefault) {
      $emailaccounts = $emailaccounts.map((item) => {
        item.default = false;
        return item;
      });
    }
    var orig = $emailaccounts.filter((item) => item.name === acc.name);
    if (orig.length > 0) {
      $emailaccounts = $emailaccounts.filter((item) => item.name !== acc.name);
      $emailaccounts.push(acc);
    }
    $account = acc;
    $emailaccounts = $emailaccounts;
    saveAccounts();
    if (emailState === "preview") makeHtml();
    clearFormData();
    showNewAccount = false;
  }

  function deleteAccount() {
    var acc = $account;
    $emailaccounts = $emailaccounts.filter((item) => item.name !== acc.name);
    if ($emailaccounts.length > 0) {
      $account = $emailaccounts[0];
      origAccount = $emailaccounts[0];
      if (emailState === "preview") makeHtml();
    } else {
      showNewAccount = false;
      $account = undefined;
      origAccount = undefined;
    }
    saveAccounts();
  }

  async function saveAccounts() {
    //
    // This will save the accounts to the harddrive.
    //
    let accountfileloc = await App.AppendPath($config.configDir,"emailaccounts.json");
    await App.WriteFile(accountfileloc, JSON.stringify($emailaccounts));
  }

  function cancelNewAccount() {
    clearFormData();
    showNewAccount = false;
  }

  function editAccountChange() {
    accountName = $account.name;
    accountDefault = $account.default;
    accountFrom = $account.from;
    accountUsername = $account.username;
    accountSmptServer = $account.smtpserver;
    accountPassword = $account.password;
    accountPort = $account.port;
    accountSig = $account.signiture;
    accountHeaderHTML = $account.headerHTML;
    accountFooterHTML = $account.footerHTML;
    showNewAccount = true;
  }

  function viewNotes() {
    saveEmailState();
    state.set("notes");
  }

  function saveEmailState() {
    if ($emailEditor !== null) {
      var rec = document.getElementById("receiverInput");
      $email.to = rec.value;
      $email.body = $emailEditor.getValue();
      $email.subject = subject.value;
    }
  }

  function viewScriptsMenu() {
    $showScripts = !$showScripts;
  }

  function viewTemplateMenu() {
    $showTemplates = !$showTemplates;
  }

  function cumulativeOffset(element) {
    var top = 0,
      left = 0;
    do {
      top += element.offsetTop || 0;
      left += element.offsetLeft || 0;
      element = element.offsetParent;
    } while (element);
    return {
      top: top,
      left: left,
    };
  }

  function showAddressBook() {
    showAddressB = !showAddressB;
  }
</script>

<div
  id="main"
  style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
>
  {#if showNewAccount}
    <div
      id="newAccountDialog"
      style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; border-color: {$theme.borderColor}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
    >
      <label for="accountDefaultInput" class="newAccountLabel">
        Default:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountDefaultInput"
        type="checkbox"
        bind:checked={accountDefault}
      />
      <label for="accountNameInput" class="newAccountLabel">
        Name of Account:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountNameInput"
        bind:value={accountName}
      />
      <label for="accountFromInput" class="newAccountLabel">
        From Email:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; border-color: {$theme.borderColor}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
        id="accountFromInput"
        bind:value={accountFrom}
      />
      <label for="accountSmptServerInput" class="newAccountLabel">
        Address of SMPT Server:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountSmptServerInput"
        bind:value={accountSmptServer}
      />
      <label for="accountPortInput" class="newAccountLabel">
        SMPT Server Port:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; border-color: {$theme.borderColor}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
        id="accountPortInput"
        bind:value={accountPort}
      />
      <label for="accountUsernameInput" class="newAccountLabel">
        User Name:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountUsernameInput"
        bind:value={accountUsername}
      />
      <label for="accountPasswordInput" class="newAccountLabel">
        Password:
      </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountPasswordInput"
        bind:value={accountPassword}
      />
      <label for="accountSigInput" class="newAccountLabel"> Signiture: </label>
      <textarea
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountSigInput"
        bind:value={accountSig}
      />
      <label for="accountHeaderHTMLInput" class="newAccountLabel">
        Header HTML:
      </label>
      <textarea
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountHeaderHTMLInput"
        bind:value={accountHeaderHTML}
      />
      <label for="accountFooterHTMLInput" class="newAccountLabel">
        Footer HTML:
      </label>
      <textarea
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        id="accountFooterHTMLInput"
        bind:value={accountFooterHTML}
      />
      <div id="buttonrow" style="grid-column: 1 / span 2;">
        <button
          id="save"
          on:click={saveNewAccount}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Save
        </button>
        <button
          id="cancel"
          on:click={cancelNewAccount}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Cancel
        </button>
      </div>
    </div>
  {/if}
  {#if showChangeAccount}
    <div
      id="accountsDiv"
      style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      <h1>Email It - Change Account</h1>
      {#if $account !== undefined}
        <h2>Current Account: {$account.name}</h2>
      {:else}
        <h2>Current Account: Please Create an Account</h2>
      {/if}
      <div id="AccountsList">
        {#each $emailaccounts as acc}
          <button
            on:click={() => {
              changeActiveAccount(acc);
            }}
            style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
          >
            {acc.name}
          </button>
        {/each}
      </div>
      <div id="buttonrow">
        <button
          id="save"
          on:click={saveAccount}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Save
        </button>
        <button
          id="new"
          on:click={newAccount}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          New
        </button>
        <button
          id="edit"
          on:click={editAccountChange}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Edit
        </button>
        <button
          id="cancel"
          on:click={cancelAccountChange}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Cancel
        </button>
        <button
          id="delete"
          on:click={deleteAccount}
          style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          Delete
        </button>
      </div>
    </div>
  {/if}
  <h1>Email It</h1>
  <div id="header">
    <div class="headerRow">
      <label for="receiverInput"> To: </label>
      <input
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; border-color: {$theme.borderColor}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
        id="receiverInput"
        class="receiverInput"
        bind:value={receiver}
        bind:this={receiverDOM}
        on:blur={() => {
          $email.to = receiver;
          //showEmailList = false;
        }}
        on:focus={() => {
          showEmailList = true;
          generateEmailList();
          receiverDOM.selectionStart = receiver.length;
        }}
        on:keydown={(e) => {
          generateEmailList(e);
        }}
        on:keyup={(e) => {
          generateEmailList(e);
        }}
        on:change={() => {
          generateEmailList();
        }}
      />
      {#if showEmailList && elist.length > 0}
        <div
          id="elist"
          style="width: {receiverDOM.offsetWidth}px; left: {cumulativeOffset(
            receiverDOM
          ).left}px; top: {cumulativeOffset(receiverDOM).top +
            40}px; background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
        >
          <ul>
            {#each elist as item}
              <li
                on:click={() => {
                  addToInput(item);
                  receiverDOM.focus();
                }}
              >
                {item}
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="headerRow">
      <label for="subject"> Subject: </label>
      <input
        id="subject"
        bind:this={subject}
        on:blur={saveEmailState}
        style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; border-color: {$theme.borderColor}; font-size: {$theme.fontSize};"
        on:focus={()=>{
          showEmailList = false;
        }}
      />
    </div>
  </div>
  {#if showPreview}
    <div id="preview" style="border-color: {$theme.textColor};">
      {@html previewHTML}
    </div>
  {:else}
    <CodeMirror
      height="330px"
      width="983px"
      config={editorConfig}
      {initFinished}
      styling="position: relative; margin-bottom: 20px; border: solid 1px transparent; border-radius: 20px; overflow: hidden;"
      on:textChange={(event) => {
        showEmailList = false;
        textChanged(event.detail.data);
      }}
      on:editorChange={(event) => {
        showEmailList = false;
        editorChange(event.detail.data);
      }}
      on:focus={() => {
        showEmailList = false;
      }}
    />
  {/if}
  <div id="buttonrow">
    <button
      on:click={viewNotes}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Notes
    </button>
    {#if emailState === "edit"}
      <button
        id="showPreview"
        on:click={createPreview}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Preview
      </button>
    {:else}
      <button
        id="edit"
        on:click={editEmail}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Edit
      </button>
      <button
        id="send"
        on:click={sendEmail}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Send
      </button>
    {/if}
    {#if $account === undefined}
      <button
        id="account"
        on:click={newAccount}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        Create a New Account
      </button>
    {:else}
      <button
        id="account"
        on:click={changeAccount}
        style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
      >
        {$account.name}
      </button>
    {/if}
    <button
      id="clear"
      on:click={clearEmail}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Clear
    </button>
    <button
      on:click={viewTemplateMenu}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Templates
    </button>
    <button
      on:click={viewScriptsMenu}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Scripts
    </button>
    <button
      id="addressBook"
      on:click={showAddressBook}
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
    >
      Address Book
    </button>
  </div>
</div>

{#if showAlert}
  <div
    id="alert"
    style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor};"
  >
    <h1>{alertTitle}</h1>
    <p>{alertMsg}</p>
    {#if badEmails.length > 0}
      <ul>
        {#each badEmails as item}
          <li>{item}</li>
        {/each}
      </ul>
    {/if}
    <button
      class="alertbutton"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor}; border-color: {$theme.borderColor}; border-color: {$theme.backgroundColor};"
      on:click={() => {
        showAlert = false;
      }}
    >
      Close
    </button>
  </div>
{/if}

{#if showAddressB}
  <AddressBook bind:show={showAddressB} />
{/if}

<style>
  :global(body) {
    margin: 0px;
    padding: 0px;
    overflow: hidden;
  }

  :global(table) {
    border-collapse: separate;
    border-spacing: 20px 5px;
  }

  #main {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px 20px 0px 20px;
  }

  #alert {
    display: flex;
    flex-direction: column;
    position: absolute;
    top: 10%;
    left: 30%;
    padding: 20px;
    max-height: 70%;
    border: solid 3px;
    border-radius: 10px;
    z-index: 400;
  }

  #alert ul {
    height: 100%;
    overflow-y: scroll;
  }

  #alert button {
    margin: 20px 10px 0px 10px;
    cursor: pointer;
  }

  .headerRow {
    display: flex;
    flex-direction: row;
    max-height: 40px;
    margin: 0px 0px 20px 0px;
    width: 100%;
  }

  #buttonrow {
    display: flex;
    flex-direction: row;
    margin: auto;
  }

  #buttonrow button {
    border-radius: 10px;
    padding: 5px 10px 5px 10px;
    margin: 0px 5px;
    width: 100%;
    max-height: 40px;
    height: 40px;
    width: auto;
    cursor: pointer;
  }

  #header {
    display: flex;
    flex-direction: column;
    margin: 0px 0px 0px 0px;
    width: 100%;
  }

  #accountsDiv {
    display: flex;
    flex-direction: column;
    z-index: 100;
    position: absolute;
    width: 80%;
    padding: 20px;
    border: solid 3px;
    border-radius: 10px;
    top: 10%;
    left: 10%;
    background-color: lightblue;
  }

  #accountsDiv button {
    border-radius: 20px;
    cursor: pointer;
  }

  #AccountsList {
    display: flex;
    flex-direction: column;
    height: 250px;
    width: 100%;
    margin: 10px;
    overflow-y: auto;
  }

  #newAccountDialog {
    display: grid;
    grid-template-columns: 1fr 2fr;
    grid-gap: 20px;
    width: 93%;
    z-index: 200;
    position: absolute;
    top: 30px;
    left: 20px;
    border-radius: 10px;
    padding: 20px;
    height: 90%;
    overflow-y: auto;
    overflow-x: hidden;
  }

  #newAccountDialog textarea {
    height: 100px;
  }

  #preview {
    height: 270px;
    max-height: 270px;
    border-radius: 20px;
    overflow-y: auto;
    padding: 10px;
    border: solid 3px;
    border-radius: 20px;
    margin: 0px 0px 15px 0px;
  }

  #elist {
    position: absolute;
    z-index: 100;
    max-height: 400px;
    overflow: scroll;
    border: solid 1px transparent;
    border-radius: 10px;
  }

  ul {
    padding: 10px;
    margin: 0px;
  }

  li {
    list-style: none;
    padding: 0px;
    margin: 5px;
    cursor: pointer;
  }

  :global(.receiverInput) {
    margin: 0px;
    padding: 0px;
    width: 100%;
    border-radius: 10px;
    border-color: transparent;
  }

  :global(.input-container) {
    background-color: transparent !important;
    width: 100% !important;
    padding: 0px !important;
    margin: 0px !important;
    border-color: transparent !important;
  }

  :global(.receiverDiv) {
    background-color: transparent;
    width: 100%;
  }

  .newAccountLabel {
    width: 310px;
  }

  .alertbutton {
    border: solid 1px transparent;
    border-radius: 10px;
    padding: 5px;
    cursor: pointer;
  }

  label {
    display: grid;
    justify-content: right;
    margin: auto 20px auto 0px;
    width: 150px;
  }

  input {
    margin: 0px;
    padding: 8px;
    width: 100%;
  }

  :global(*:focus) {
    outline: none;
  }

  h1 {
    margin: 0px auto 10px auto;
  }

  h2 {
    margin: 10px;
  }

  textarea {
    height: 300px;
    max-height: 300px;
    border-radius: 20px;
    padding: 10px;
  }

  input {
    border-radius: 10px;
  }
</style>
