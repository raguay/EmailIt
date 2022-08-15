<script>
  import { onMount, afterUpdate } from "svelte";
  import { theme } from "../stores/theme.js";
  import { Octokit } from "@octokit/rest";

  let octok;
  let repos = null;
  let themes = null;
  let msgs = [];
  let pickerDOM;
  let hiddenInput;
  let timeOut;
  let loading = true;
  let config = {
    configDir: "",
  };

  onMount(async () => {
    let hdir = await window.go.main.App.GetHomeDir();
    config.configDir = await window.go.main.App.AppendPath(
      hdir,
      ".config/scriptserver"
    );
    octok = new Octokit();
    await loadRepoInfo();
    timeOut = setTimeout(focusInput, 1000);
    return () => {
      hiddenInput = null;
      clearTimeout(timeOut);
    };
  });

  afterUpdate(() => {
    if (typeof hiddenInput !== "undefined") hiddenInput.focus();
  });

  function focusInput() {
    clearTimeout(timeOut);
    if (typeof hiddenInput !== "undefined" && hiddenInput !== null)
      hiddenInput.focus();
    timeOut = setTimeout(focusInput, 1000);
  }

  async function loadRepoInfo() {
    if (repos === null && themes === null) {
      loading = true;
      if (typeof repos !== "undefined") {
        repos = {};
      }
      if (typeof themes !== "undefined") {
        themes = {};
      }
      var repost = await octok.search.repos({
        q: "topic:emailit+topic:script",
      });
      if (check(repost)) {
        repost = repost.data.items;
        for (var i = 0; i < repost.length; i++) {
          repost[i].loaded = await extExists(repost[i]);
        }
        repos = repost;
      }
      var themest = await octok.search.repos({
        q: "topic:emailit+topic:theme",
      });
      if (check(themest)) {
        themest = themest.data.items;
        for (var i = 0; i < themest.length; i++) {
          themest[i].loaded = await themeExists(themest[i]);
        }
        themes = themest;
      }
      loading = false;
    }
    repos = repos;
    themes = themes;
  }

  function check(val) {
    return (
      typeof val !== "undefined" &&
      typeof val.data !== "undefined" &&
      typeof val.data.items !== "undefined"
    );
  }

  async function runCommandLine(line, rEnv, callback, dir) {
    //
    // Get the environment to use.
    //
    console.log(line, rEnv, dir);

    var resp = await fetch(`http://localhost:9978/api/scripts/env/Default`);
    var nEnv = await resp.json();
    console.table(nEnv);
    if (typeof rEnv !== "undefined") {
      nEnv = { ...nEnv, ...rEnv };
    }

    //
    // Fix the environment from a map to an array of strings.
    //
    var penv = [];
    for (const key in nEnv) {
      penv.push(`${key}=${nEnv[key]}`);
    }

    //
    // Make sure dir has a value.
    //
    if (typeof dir === "undefined") dir = ".";

    //
    // Run the command line in a shell. #TODO: make the shell configurable.
    //
    var args = ["/bin/zsh", "-c", line];
    var cmd = "/bin/zsh";

    //
    // Run the command line.
    //
    var result = await window.go.main.App.RunCommandLine(cmd, args, penv, dir);
    var err = await window.go.main.App.GetError();
    //
    // If callback is given, call it with the results.
    //
    if (typeof callback !== "undefined" || callback !== null) {
      callback(err, result);
    }
  }

  async function installTheme(thm) {
    var thmDir = await window.go.main.App.AppendPath(
      config.configDir,
      "styles"
    );
    thmDir = await window.go.main.App.AppendPath(thmDir, thm.name);
    if (!(await window.go.main.App.DirExists(thmDir))) {
      await window.go.main.App.MakeDir(thmDir);
    }
    await runCommandLine(
      `git clone '${thm.git_url}' '${thmDir}'`,
      [],
      () => {
        //
        // The clone should be there. Let's load the new theme.
        //
        loadTheme(thm);
        loadRepoInfo();
      },
      "."
    );
  }

  async function loadTheme(thm) {
    var thmDir = await window.go.main.App.AppendPath(
      config.configDir,
      "styles"
    );
    thmDir = await window.go.main.App.AppendPath(thmDir, thm.name);
    const pfile = await window.go.main.App.AppendPath(thmDir, "package.json");
    if (await window.go.main.App.FileExists(pfile)) {
      var manifest = await window.go.main.App.ReadFile(pfile);
      manifest = JSON.parse(manifest);
      const mfile = await window.go.main.App.AppendPath(
        thmDir,
        manifest.theme.main
      );
      var newTheme = await window.go.main.App.ReadFile(mfile);
      newTheme = JSON.parse(newTheme);
      $theme = newTheme;
      addMsg(thm, "This theme is now being used.");
    } else {
      addMsg(thm, "The theme doesn't have a package.json file.");
    }
  }

  async function themeExists(thm) {
    var thmDir = await window.go.main.App.AppendPath(
      config.configDir,
      "styles"
    );
    thmDir = await window.go.main.App.AppendPath(thmDir, thm.name);
    var result = await window.go.main.App.DirExists(thmDir);
    return result;
  }

  async function deleteTheme(thm) {
    var thmDir = await window.go.main.App.AppendPath(
      config.configDir,
      "styles"
    );
    var tpath = await window.go.main.App.AppendPath(thmDir, thm.name);
    await window.go.main.App.DeleteEntries(tpath);
    themes = themes.map((item) => {
      if (item.name === thm.name) {
        item.loaded = false;
      }
      return item;
    });
  }

  async function installExtension(ext) {
    var extDir = await window.go.main.App.AppendPath(
      config.configDir,
      "scripts"
    );
    extDir = await window.go.main.App.AppendPath(extDir, ext.name);
    if (!(await window.go.main.App.DirExists(extDir))) {
      await window.go.main.App.MakeDir(extDir);
    }
    await runCommandLine(
      `git clone '${ext.git_url}' '${extDir}'`,
      [],
      () => {
        addMsg(ext, "Restart the program to use this extension.");
        loadRepoInfo();
      },
      "."
    );
  }

  async function extExists(ext) {
    var extDir = await window.go.main.App.AppendPath(
      config.configDir,
      "scripts"
    );
    extDir = await window.go.main.App.AppendPath(extDir, ext.name);
    var flag = await window.go.main.App.DirExists(extDir);
    return flag;
  }

  async function deleteExtension(ext) {
    var extDir = await window.go.main.App.AppendPath(
      config.configDir,
      "scripts"
    );
    let epath = await window.go.main.App.AppendPath(extDir, ext.name);
    await window.go.main.App.DeleteEntries(epath);
    repos = repos.map((item) => {
      if (item.name === ext.name) {
        item.loaded = false;
      }
      return item;
    });
    addMsg(ext, "Rerun the application to remove the script from memory.");
  }

  function hasMsg(rp) {
    if (msgs.length > 0) {
      return msgs.find((item) => item.name === rp.name) !== "undefined";
    } else {
      return false;
    }
  }

  function getMsg(rp) {
    if (hasMsg(rp)) {
      var item = msgs.find((item) => item.name === rp.name);
      if (typeof item !== "undefined") {
        return item.msg;
      }
    }
    return "";
  }

  function addMsg(rp, msg) {
    msgs.push({
      name: rp.name,
      msg: msg,
    });
    msgs = msgs;
    themes = themes;
    repos = repos;
  }

  function inputChange(e) {
    if (e.key === "ArrowUp" || e.key === "k") {
      //
      // Go up the list. Zero is at the top.
      //
      scrollDOM(-1);
    } else if (e.key === "ArrowDown" || e.key === "j") {
      //
      // Go down the list. The largest index is at the bottom.
      //
      scrollDOM(1);
    }
  }

  function scrollDOM(amount) {
    var adj = amount * 20;

    if (pickerDOM !== null) {
      pickerDOM.scrollTop += adj;
      if (pickerDOM.scrollTop < 0) pickerDOM.scrollTop = 0;
    }
  }
</script>

<div
  id="GitHub"
  style="background-color: {$theme.backgroundColor};
         color: {$theme.textColor};"
>
  <div id="GitHubHeader">
    <h3>GitHub Themes and Scripts Importer</h3>
    <input id="inputHidden" bind:this={hiddenInput} on:keydown={inputChange} />
  </div>
  {#if loading}
    <h1>Loading....</h1>
  {:else}
    <div id="GitHubList" bind:this={pickerDOM}>
      {#each repos as repo}
        <div class="repoblock">
          <div class="reporow">
            <p class="reponame">
              {repo.name}
            </p>
            <p class="repostars" style="color: {$theme.Yellow};">
              {repo.stargazers_count} ⭐️s
            </p>
          </div>
          <div class="reporow">
            <p class="repodisc">
              {repo.description}
            </p>
          </div>
          {#if hasMsg(repo)}
            <div class="reporow" style="color: {$theme.Red};">
              {@html getMsg(repo)}
            </div>
          {/if}
          <div class="repobuttons">
            {#if repo.loaded}
              <button
                on:click={(e) => {
                  deleteExtension(repo);
                }}
                style="background-color: {$theme.Red};"
              >
                Delete
              </button>
            {:else}
              <button
                on:click={(e) => {
                  installExtension(repo);
                }}
                style="background-color: {$theme.Green};"
              >
                Install
              </button>
            {/if}
          </div>
        </div>
      {/each}
      {#each themes as thm}
        <div class="repoblock">
          <div class="reporow">
            <p class="reponame">
              {thm.name}
            </p>
            <p class="repostars" style="color: {$theme.Yellow};">
              {thm.stargazers_count} ⭐️ s
            </p>
          </div>
          <div class="reporow">
            <p class="repodisc">
              {thm.description}
            </p>
          </div>
          {#if hasMsg(thm)}
            <div class="reporow" style="color: {$theme.Red};">
              {@html getMsg(thm)}
            </div>
          {/if}
          <div class="repobuttons">
            {#if thm.loaded}
              <button
                on:click={(e) => {
                  loadTheme(thm);
                }}
                style="background-color: {$theme.Green};"
              >
                Load
              </button>
              <button
                on:click={(e) => {
                  deleteTheme(thm);
                }}
                style="background-color: {$theme.Red};"
              >
                Delete
              </button>
            {:else}
              <button
                on:click={(e) => {
                  installTheme(thm);
                }}
                style="background-color: {$theme.Green};"
              >
                Install
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  #GitHub {
    display: flex;
    flex-direction: column;
  }

  #GitHubHeader {
    display: flex;
    flex-direction: row;
    margin: 10px;
  }

  #GitHubHeader h3 {
    margin: 0px auto 0px 0px;
  }

  #GitHubList {
    margin: 5px 10px;
    overflow-y: auto;
    overflow-x: hidden;
  }

  #inputHidden {
    width: 0px;
    height: 0px;
    margin: 0px;
    padding: 0px;
    border: 0px solid transparent;
  }

  .reporow {
    display: flex;
    flex-direction: row;
    margin: 0px;
  }

  .reponame {
    margin: 0px auto 0px 0px;
  }

  .repostars {
    margin: 0px 0px 0px auto;
  }

  .repodisc {
    margin: 0px 0px 0px 15px;
  }

  .repoblock {
    display: flex;
    flex-direction: column;
    margin: 5px 0px;
  }

  .repobuttons {
    display: flex;
    flex-direction: row;
    margin: 5px auto;
  }

  .repobuttons button {
    margin: 0px 10px;
    border-radius: 5px;
  }
</style>
