<script>
  import { onMount, afterUpdate } from "svelte";
  import { theme } from "../stores/theme.js";
  import { termscripts } from "../stores/termscripts.js";
  import { scripts } from "../stores/scripts.js";
  import { config } from "../stores/config.js";
  import * as ap from '../../wailsjs/go/main/App.js';

  let repos = null;
  let themes = null;
  let msgs = [];
  let pickerDOM;
  let hiddenInput;
  let timeOut;
  let loading = true;
  let cfg = {
    configDir: "",
  };

  onMount(async () => {
    let hdir = await ap.GetHomeDir();
    cfg.configDir = await ap.AppendPath(
      hdir,
      ".config/scriptserver"
    );
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
      repos = await ap.GetGitHubScripts();
      themes = await ap.GetGitHubThemes();
      for (var i = 0; i < repos.length; i++) {
        repos[i].loaded = await extExists(repos[i]);
      }
      for (var i = 0; i < themes.length; i++) {
        themes[i].loaded = await themeExists(themes[i]);
      }
      loading = false;
    }
    repos = repos;
    themes = themes;
  }

  async function installTheme(thm) {
    var thmDir = await ap.AppendPath(
      cfg.configDir,
      "styles"
    );
    if (!(await ap.DirExists(thmDir))) {
      await ap.MakeDir(thmDir);
    }
    thmDir = await ap.AppendPath(thmDir, thm.name);
    if (!(await ap.DirExists(thmDir))) {
      await ap.MakeDir(thmDir);
    }
    await runCommandLine(
      `git clone '${thm.git_url}' '${thmDir}'`,
      [],
      () => {
        //
        // The clone should be there. Let's load the new theme.
        //
        loadTheme(thm);
        themes = themes.map((item) => {
          if (item.name === thm.name) {
            item.loaded = true;
          }
          return item;
        });
        loadRepoInfo();
      },
      "."
    );
  }

  async function loadTheme(thm) {
    var thmDir = await ap.AppendPath(
      cfg.configDir,
      "styles"
    );
    thmDir = await ap.AppendPath(thmDir, thm.name);
    const pfile = await ap.AppendPath(thmDir, "package.json");
    if (await ap.FileExists(pfile)) {
      var manifest = await ap.ReadFile(pfile);
      manifest = JSON.parse(manifest);
      const mfile = await ap.AppendPath(
        thmDir,
        manifest.theme.main
      );
      var newTheme = await ap.ReadFile(mfile);
      newTheme = JSON.parse(newTheme);
      $theme = newTheme;
      addMsg(thm, "This theme is now being used.");
    } else {
      addMsg(thm, "The theme doesn't have a package.json file.");
    }
  }

  async function themeExists(thm) {
    var thmDir = await ap.AppendPath(
      cfg.configDir,
      "styles"
    );
    thmDir = await ap.AppendPath(thmDir, thm.name);
    var result = await ap.DirExists(thmDir);
    return result;
  }

  async function deleteTheme(thm) {
    var thmDir = await ap.AppendPath(
      cfg.configDir,
      "styles"
    );
    var tpath = await ap.AppendPath(thmDir, thm.name);
    await ap.DeleteEntries(tpath);
    themes = themes.map((item) => {
      if (item.name === thm.name) {
        item.loaded = false;
      }
      return item;
    });
    loadRepoInfo();
  }

  async function installExtension(ext) {
    var extDir = await ap.AppendPath(
      cfg.configDir,
      "scripts"
    );
    if (!(await ap.DirExists(extDir))) {
      await ap.MakeDir(extDir);
    }
    extDir = await ap.AppendPath(extDir, ext.name);
    if (!(await ap.DirExists(extDir))) {
      await ap.MakeDir(extDir);
    }
    await runCommandLine(
      `git clone '${ext.git_url}' '${extDir}'`,
      [],
      async () => {
        let cfgloc = await ap.AppendPath(
          extDir,
          "package.json"
        );
        let cfg = await ap.ReadFile(cfgloc);
        cfg = JSON.parse(cfg);
        let script = {
          name: cfg.script.name,
          script: cfg.script.script,
          path: extDir,
          env: "Default",
          termscript: cfg.script.termscript,
          description: cfg.script.description,
          help: cfg.script.help,
        };
        //
        // Save the script #TODO
        //
        if (script.termscript) {
          $termscripts.push(script);
        } else {
          script.insert = false;
          $scripts.push(script);
        }
        repos = repos.map((item) => {
          if (item.name === ext.name) {
            item.loaded = true;
          }
          return item;
        });
        addMsg(ext, `${ext.name} external script has been downloaded.`);
        loadRepoInfo();
      },
      "."
    );
  }

  async function extExists(ext) {
    var extDir = await ap.AppendPath(
      cfg.configDir,
      "scripts"
    );
    extDir = await ap.AppendPath(extDir, ext.name);
    var flag = await ap.DirExists(extDir);
    return flag;
  }

  async function deleteExtension(ext) {
    var extDir = await ap.AppendPath(
      cfg.configDir,
      "scripts"
    );
    let epath = await ap.AppendPath(extDir, ext.name);
    let cfgloc = await ap.AppendPath(epath, "package.json");
    let cfg = await ap.ReadFile(cfgloc);
    cfg = JSON.parse(cfg);
    await ap.DeleteEntries(epath);
    repos = repos.map((item) => {
      if (item.name === ext.name) {
        item.loaded = false;
      }
      return item;
    });
    loadRepoInfo();

    //
    // Remove from the external scripts list.
    //
    if (cfg.termscript) {
      $termscripts = $termscripts.filter((item) => item.name !== ext.name);
    } else {
      $scripts = $scripts.filter((item) => item.name !== ext.name);
    }
    addMsg(ext, "This external script has been removed.");
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

	async function runCommandLine(line, rEnv, callback, dir) {
	  //
	  // Get the environment to use.
    //
    let envlistloc = await ap.AppendPath($config.configDir, "environments.json");
    let envlist = await ap.ReadFile(envlistloc);
    envlist = JSON.parse(envlist);
    let nEnv = envlist.map(item => 'Default');

    //
    // Get the default environment.
    //
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
	  var result = await ap.RunCommandLine(cmd, args, penv, dir);
	  var err = await ap.GetError();
	  //
	  // If callback is given, call it with the results.
	  //
	  if (typeof callback !== "undefined" || callback !== null) {
	    callback(err, result);
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
              {repo.stars} ⭐️s
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
                on:click={() => {
                  deleteExtension(repo);
                }}
                style="background-color: {$theme.Red};"
              >
                Delete
              </button>
            {:else}
              <button
                on:click={() => {
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
              {thm.stars} ⭐️s
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
                on:click={() => {
                  loadTheme(thm);
                }}
                style="background-color: {$theme.Green};"
              >
                Load
              </button>
              <button
                on:click={() => {
                  deleteTheme(thm);
                }}
                style="background-color: {$theme.Red};"
              >
                Delete
              </button>
            {:else}
              <button
                on:click={() => {
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
