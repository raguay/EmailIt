<script>
  import { onMount } from "svelte";
  import showdown from "showdown";
  import EmailIt from "./components/EmailIt.svelte";
  import Notes from "./components/Notes.svelte";
  import ScriptMenu from "./components/ScriptMenu.svelte";
  import TemplateMenu from "./components/TemplateMenu.svelte";
  import ScriptsEditor from "./components/ScriptsEditor.svelte";
  import TemplatesEditor from "./components/TemplatesEditor.svelte";
  import Preferences from "./components/Preferences.svelte";
  import ScriptLine from "./components/ScriptLine.svelte";
  import { commands } from "./stores/commands.js";
  import { state } from "./stores/state.js";
  import { sp } from "./stores/sp.js";
  import { email } from "./stores/email.js";
  import { emails } from "./stores/emails.js";
  import { wd } from "./stores/wd.js";
  import { scripts } from "./stores/scripts.js";
  import { termscripts } from "./stores/termscripts.js";
  import { showScripts } from "./stores/showScripts.js";
  import { templates } from "./stores/templates.js";
  import { showTemplates } from "./stores/showTemplates.js";
  import { theme } from "./stores/theme.js";
  import { config } from "./stores/config.js";
  import { extScripts } from "./stores/extScripts.js";
  import { systemScripts } from "./stores/systemScripts.js";
  import { userScripts } from "./stores/userScripts.js";
  import { userTemplates } from "./stores/userTemplates.js";
  import { systemTemplates } from "./stores/systemTemplates.js";
  import { runscript } from "./stores/runscript.js";
  import { notes } from "./stores/notes.js";
  import { runhandlebars } from "./stores/runhandlebars.js";
  import { runtemplate } from "./stores/runtemplate.js";
  import { noteEditor } from "./stores/noteEditor.js";
  import { currentNote } from "./stores/currentNote.js";
  import { account } from "./stores/account.js";
  import { emailaccounts } from "./stores/emailaccounts.js";
  import { DateTime } from "luxon";
  import { create, all } from "mathjs";
  import Handlebars from "handlebars";
  import * as App from "../wailsjs/go/main/App.js";
  import * as rt from "../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  let notestruct = {
    notes: ["", "", "", "", "", "", "", "", "", ""],
    loadNotes: async function () {
      //
      // Load the notes.json file.
      //
      let notejsonloc = await App.AppendPath($config.configDir, "notes.json");
      if (await App.FileExists(notejsonloc)) {
        let tmp = await App.ReadFile(notejsonloc);
        notestruct.notes = JSON.parse(tmp);
      } else {
        notestruct.saveNotes();
      }
    },
    saveNotes: async function () {
      //
      // Save the notes.json file.
      //
      let notejsonloc = await App.AppendPath($config.configDir, "notes.json");
      await App.WriteFile(notejsonloc, JSON.stringify(notestruct.notes));
    },
    getNote: function (noteid) {
      return notestruct.notes[noteid];
    },
    putNote: async function (noteid, note) {
      notestruct.notes[noteid] = note;
      await notestruct.saveNotes();
    },
  };
  const mathconfig = {
    epsilon: 1e-12,
    matrix: "Matrix",
    number: "number",
    precision: 64,
    predictable: false,
    randomSeed: null,
  };
  const mathjs = create(all, mathconfig);
  const fs = {
    SystemOpenFile: async function (prog) {
      await App.SystemOpenFile(prog);
    },
    Getwd: async function () {
      let wd = await App.Getwd();
      return wd;
    },
    ReadFile: async function (path) {
      let contents = await App.ReadFile(path);
      return contents;
    },
    WriteFile: async function (path, contents) {
      await App.WriteFile(path, contents);
    },
    GetHomeDir: async function () {
      let hd = await App.GetHomeDir();
      return hd;
    },
    FileExists: async function (path) {
      let exists = await App.FileExists(path);
      return exists;
    },
    DirExists: async function (path) {
      let exists = await App.DirExists(path);
      return exists;
    },
    SplitFile: async function (path) {
      let parts = await App.SplitFile(path);
      return parts;
    },
    ReadDir: async function (path) {
      let items = await App.ReadDir(path);
      return items;
    },
    MakeDir: async function (path) {
      await App.MakeDir(path);
    },
    MakeFile: async function (path) {
      await App.MakeFile(path);
    },
    MoveEntries: async function (from, to) {
      await App.MoveEntries(from, to);
    },
    RenameEntry: async function (from, to) {
      await App.RenameEntry(from, to);
    },
    CopyEntries: async function (from, to) {
      await App.CopyEntries(from, to);
    },
    DeleteEntries: async function (path) {
      await App.DeleteEntries(path);
    },
    CreateTempFile: async function (contents) {
      let fname = await App.CreateTempFile(contents);
      return fname;
    },
    Chmod: async function (path, mode) {
      App.Chmod(path, mode);
    },
    RunCommandLine: async function (cmd, args, env, dir) {
      let output = await App.RunCommandLine(cmd, args, env, dir);
      return output;
    },
    GetClip: async function () {
      let output = await App.GetClip();
      return output;
    },
    SetClip: async function (msg) {
      await App.SetClip(msg);
    },
    GetEnvironment: async function () {
      let output = await App.GetEnvironment();
      return output;
    },
    GetEnv: async function (name) {
      let envname = await App.GetEnv(name);
      return envname;
    },
    AppendPath: async function (dir, name) {
      let result = await App.AppendPath(dir, name);
      return result;
    },
  };
  $sp = {
    text: "",
    data: {
      result: "", // Contains the result from the last command.
      registers: [], // Temporary storage registers.
      line: "", // Running command line.
      previousLines: [], // former command lines.
    },
    DateTime: DateTime,
    mathjs: mathjs,
    mathParser: mathjs.parser(),
    that: $sp,
    Handlebars: $runhandlebars,
    fs: fs,
    ProcessMath: function (txt) {
      var result = this.mathParser.evaluate(txt);
      return result;
    },
    ProcessMathNotes: function (Note) {
      var result = "";
      this.mathParser.clear();
      let lines = Note.match(/^.*((\\r\\n|\\n|\\r)|$)/gm);
      let numLines = lines.length;
      for (let i = 0; i < lines.length; i++) {
        var line = lines[i];
        var lineResult = "";
        line = line.trim();
        if (line.indexOf(">") === 0) {
          line = line.substr(2);
          var inx = line.indexOf("|");
          if (inx !== -1) {
            line = line.substr(0, inx - 1);
            line = line.trim();
          }
          lineResult = this.mathParser.evaluate(line);
          if (typeof lineResult === "function") {
            lineResult = "Definition";
          }
          var whiteSP = 32 - (line.length + 3);
          if (whiteSP < 1) {
            whiteSP = 1;
          }
          result +=
            "> " +
            line +
            this.insertCharacters(whiteSP, " ") +
            " | " +
            lineResult;
        } else {
          result += line;
        }
        if (--numLines !== 0) result += "\n";
      }
      return result;
    },
    runScript: function (script, text) {
      let result =
        "Script not found! Only JavaScript scripts can be ran this way. No external scripts!";
      let scriptIndex = $userScripts.find((ele) => {
        return ele.name === script;
      });
      if (typeof scriptIndex !== "undefined") {
        script = scriptIndex.script;
        result = runJavaScript(script, text);
      } else {
        scriptIndex = $systemScripts.find((ele) => {
          return ele.name === script;
        });
        if (typeof scriptIndex !== "undefined") {
          script = scriptIndex.script;
          result = runJavaScript(script, text);
        }
      }
      return result;
    },
    returnNote: function (id) {
      var result = "";
      if (id >= 0 && id <= 9) result = $notes.getNote(id);
      return result;
    },
    insertCharacters: function (num, ichar) {
      var result = "";
      if (num < 0) return result;
      for (var i = 0; i < num; i++) {
        result += ichar;
      }
      return result;
    },
  };

  onMount(async () => {
    //
    // Set the state to emailit.
    //
    $state = "emailit";
    $commands = [];

    //
    // Get stuff from the server.
    //
    await getConfig();
    await getScriptsList();
    getTermScriptsList();
    await getTemplatesList();
    await getAccounts();
    await getTheme();
    await getNotes();
    setupHandlebarHelpers();

    $runscript = runScript;
    $runhandlebars = Handlebars;
    $sp.Handlebars = Handlebars;
    $runtemplate = runTemplate;

    //
    // Run the StartUp script if one exists.
    //
    let startup = $scripts.find((item) => item.name === "StartUp");
    if (typeof startup !== "undefined") {
      //
      // A StartUp script is defined. Run it!
      //
      await $runscript(startup, "");
    }

    await getEmails();

    //
    // Start the listeners from the Server.
    //
    rt.EventsOn("changewd", async (msg) => {
      $wd = msg.wd;
    });

    rt.EventsOn("notechange", async (msg) => {
      let nid = parseInt(msg.noteid);
      let result = await $notes.getNote(nid);
      if (msg.aw === "w") {
        result = msg.note;
      } else {
        result = `${result}\n${msg.note}`;
      }
      await $notes.putNote(nid, result);
      if ($noteEditor !== null && $currentNote === nid) {
        $noteEditor.setValue(result);
      }
    });
    rt.EventsOn("envList", async (msg) => {
      let envlistloc = await App.AppendPath(
        $config.configDir,
        "environments.json"
      );
      let envlist = await App.ReadFile(envlistloc);
      envlist = JSON.parse(envlist);
      await rt.EventsEmit(
        msg.returnMsg,
        envlist.map((item) => item.name)
      );
    });

    rt.EventsOn("scriptList", async (msg) => {
      await rt.EventsEmit(
        msg.returnMsg,
        $scripts.map((item) => item.name)
      );
    });

    rt.EventsOn("scriptRun", async (msg) => {
      let result = "";
      if (msg.commandLine == null || msg.commandLine.trim().length === 0) {
        result = await $runscript(msg.script, msg.text, msg.envVar);
      } else {
        result = await $runscript(msg.commandLine, msg.text, msg.envVar);
      }
      await rt.EventsEmit(msg.returnMsg, result);
    });

    rt.EventsOn("templateList", async (msg) => {
      await rt.EventsEmit(
        msg.returnMsg,
        $templates.map((item) => item.name)
      );
    });

    rt.EventsOn("templateRun", async (msg) => {
      let result = await $runtemplate(msg.template, msg.text);
      await rt.EventsEmit(msg.returnMsg, result);
    });

    rt.EventsOn("emailSend", async (msg) => {
      let result = "Sent successfully!";
      let processText = await $runtemplate("given", msg.body);
      if (msg.account === "") {
        //
        // Use the current account to send it.
        //
        let processHTML = makeHtml($account, processText);
        result = await App.SendEmail(
          $account.username,
          $account.from,
          $account.password,
          $account.smtpserver,
          $account.port,
          msg.to,
          processHTML,
          processText,
          msg.subject,
          msg.attachment
        );
      } else if (msg.account === "Default") {
        //
        // Get the default account.
        //
        let acc = $emailaccounts.filter((item) => item.default)[0];
        let processHTML = makeHtml(acc, processText);
        result = await App.SendEmail(
          acc.username,
          acc.from,
          acc.password,
          acc.smtpserver,
          acc.port,
          msg.to,
          processHTML,
          processText,
          msg.subject,
          msg.attachment
        );
      } else {
        //
        // It has a different account to send by. Look it up and use it.
        //
        let acc = $emailaccounts.filter((item) => item.name === msg.account);
        if (acc.length > 0) {
          acc = acc[0];
          let processHTML = makeHtml(acc, processText);
          result = await App.SendEmail(
            acc.username,
            acc.from,
            acc.password,
            acc.smtpserver,
            acc.port,
            msg.to,
            processHTML,
            processText,
            msg.subject,
            msg.attachment
          );
        } else {
          result = "Account not found.";
        }
      }
      await rt.EventsEmit(msg.returnMsg, result);
    });
    rt.EventsOn("EditEmail", async (msg) => {
      $email = {
        to: msg.to,
        subject: msg.subject,
        body: msg.body,
        new: true,
      };
      $state = "emailit";
      await rt.WindowSetSize($config.width, $config.height);
      await rt.Show();
      await rt.EventsEmit(msg.returnMsg, "Okay");
    });

    //
    // The next three are for switching program states. Since we aren't checking current state,
    // just set the right window size just to be sure.
    //
    rt.EventsOn("EmailIt", async (msg) => {
      $state = "emailit";
      await rt.WindowSetSize($config.width, $config.height);
      await rt.Show();
      await rt.EventsEmit(msg.returnMsg, "Okay");
    });
    rt.EventsOn("ScriptLine", async (msg) => {
      $state = "scriptline";
      await rt.WindowSetSize($config.width, $config.height);
      await rt.Show();
      await rt.EventsEmit(msg.returnMsg, "Okay");
    });
    rt.EventsOn("Notes", async (msg) => {
      $state = "notes";
      await rt.WindowSetSize($config.width, $config.height);
      await rt.Show();
      await rt.EventsEmit(msg.returnMsg, "Okay");
    });
  });

  async function getEmails() {
    //
    // Get the emails from the system.
    //
    let emailfileloc = await App.AppendPath($config.configDir, "emails.json");
    if (await App.FileExists(emailfileloc)) {
      let emailfilejson = await App.ReadFile(emailfileloc);
      $emails = JSON.parse(emailfilejson);
    }
  }

  function makeHtml(acc, text) {
    let result = "";
    var converter = new showdown.Converter({
      extensions: [],
    });
    converter.setOption("tables", true);
    result = converter.makeHtml(text + acc.signiture);
    if (typeof acc.headerHTML !== "undefined") {
      result = acc.headerHTML + result + acc.footerHTML;
    }
    return result;
  }

  async function getAccounts() {
    //
    // Get the accounts from the system.
    //
    let accountfileloc = await App.AppendPath(
      $config.configDir,
      "emailaccounts.json"
    );
    if (await App.FileExists(accountfileloc)) {
      $emailaccounts = await App.ReadFile(accountfileloc);
      $emailaccounts = JSON.parse($emailaccounts);
      if ($emailaccounts.length > 0) {
        $account = $emailaccounts.find((item) => item.default);
      }
    }
  }

  function setupHandlebarHelpers() {
    //
    // Create the helpers functions for Handlebars.
    //
    Handlebars.registerHelper("date", function (dFormat) {
      return DateTime.now().toFormat(dFormat);
    });

    Handlebars.registerHelper("cdate", function (cTime, dFormat) {
      return DateTime.fromISO(cTime).toFormat(dFormat);
    });

    Handlebars.registerHelper("last", function (weeks, fmat) {
      return DateTime.now().minus({ weeks: weeks }).toFormat(fmat);
    });

    Handlebars.registerHelper("next", function (weeks, fmat) {
      return DateTime.now().plus({ weeks: weeks }).toFormat(fmat);
    });
  }

  async function getCopyClip(ccadd) {
    let result = `CopyClip ${ccadd} not set.`;
    let ccdir = `${$config.homeDir}${$config.AlfredData}/com.customct.CopyClips/copyclips`;
    let ccfile = `${ccdir}/copy-${ccadd}.txt`;
    if (await App.FileExists(ccfile)) {
      result = await App.ReadFile(ccfile);
    }
    return result;
  }

  async function runTemplate(template, text) {
    //
    // process the template.
    //
    let result = "";
    try {
      //
      // Create various default targets for the templater. These have
      // to be created since they are various types of time/date stamps.
      //
      var data = [];
      data["cDateMDY"] = DateTime.now().toFormat("LLLL dd, yyyy");
      data["cDateDMY"] = DateTime.now().toFormat("dd LLLL yyyy");
      data["cDateDOWDMY"] = DateTime.now().toFormat("cccc, dd LLLL yyyy");
      data["cDateDOWMDY"] = DateTime.now().toFormat("cccc LLLL dd, yyyy");
      data["cDay"] = DateTime.now().toFormat("dd");
      data["cMonth"] = DateTime.now().toFormat("LLLL");
      data["cYear"] = DateTime.now().toFormat("yyyy");
      data["cMonthShort"] = DateTime.now().toFormat("LLL");
      data["cYearShort"] = DateTime.now().toFormat("yy");
      data["cDOW"] = DateTime.now().toFormat("cccc");
      data["cMDthYShort"] = DateTime.now().toFormat("LLL d yy");
      data["cMDthY"] = DateTime.now().toFormat("LLLL d yyyy");
      data["cHMSampm"] = DateTime.now().toFormat("h:mm:ss a");
      data["cHMampm"] = DateTime.now().toFormat("h:mm a");
      data["cHMS24"] = DateTime.now().toFormat("H:mm:ss");
      data["cHM24"] = DateTime.now().toFormat("H:mm");
      data["cc0"] = await getCopyClip("0");
      data["cc1"] = await getCopyClip("1");
      data["cc2"] = await getCopyClip("2");
      data["cc3"] = await getCopyClip("3");
      data["cc4"] = await getCopyClip("4");
      data["cc5"] = await getCopyClip("5");
      data["cc6"] = await getCopyClip("6");
      data["cc7"] = await getCopyClip("7");
      data["cc8"] = await getCopyClip("8");
      data["cc9"] = await getCopyClip("9");

      //
      // Get the User's default data.
      //
      var defaultData = $templates.find((item) => item.name === "Defaults");
      if (template === "given") {
        template = text;
        data["text"] = "";
      } else {
        template = $templates.find((item) => item.name === template);
        data["text"] = text;
        if (typeof template !== "undefined") {
          template = template.template;
        }
      }
      if (defaultData === undefined) {
        defaultData = {};
        defaultData.name = "Defaults";
        defaultData.template = `{
   "projectName": "This is a default test project",
   "email": "Your Email",
   "firstName": "Your Frist Name",
   "lastName": "You Last Name",
   "home": {
      "address1": "",
      "address2": "",
      "city": "",
      "state": "",
      "zip": ""
   },
   "office": {
      "address1": "",
      "address2": "",
      "city": "",
      "state": "",
      "zip": ""
   }}`;
        $userTemplates.push(defaultData);
        saveUserTemplates();
      }
      data = MergeRecursive(data, JSON.parse(defaultData.template));

      //
      // Parse the Template
      //
      var tpTemplate = $runhandlebars.compile(template);

      //
      // Run the template with the data.
      //
      result = tpTemplate(data);

      //
      // Make sure we have a string and not an array or object.
      //
      if (typeof result != "string") {
        result = result.toString();
      }
    } catch (error) {
      console.error("Handlebars Error: " + error);
      result = "There was an error.";
    }
    return result;
  }

  //
  //  Function:        MergeRecursive
  //
  //  Description:     Recursively merge properties of two objects
  //
  //  Inputs:
  //                   obj1    The first object to merge
  //                   obj2    The second object to merge
  //
  function MergeRecursive(obj1, obj2) {
    for (var p in obj2) {
      try {
        // Property in destination object set; update its value.
        if (obj2[p].constructor == Object) {
          obj1[p] = MergeRecursive(obj1[p], obj2[p]);
        } else {
          obj1[p] = obj2[p];
        }
      } catch (e) {
        // Property in destination object not set; create it and set its value.
        obj1[p] = obj2[p];
      }
    }
    return obj1;
  }

  function openNote(id) {
    $currentNote = id;
    $noteEditor.setValue($notes.notes[$currentNote]);
  }

  function wait(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(ms);
      }, ms);
    });
  }

  async function getNotes() {
    $notes = notestruct;
    await $notes.loadNotes();
  }

  async function getConfig() {
    //
    // Get the config and store in in the config store.
    //
    let hmdir = await App.GetHomeDir();
    let configdir = await App.AppendPath(hmdir, ".config");
    let configloc = "";
    configdir = await App.AppendPath(configdir, "EmailIt");
    if (await App.DirExists(configdir)) {
      configloc = await App.AppendPath(configdir, "config.json");
      if (await App.FileExists(configloc)) {
        let configJSON = await App.ReadFile(configloc);
        $config = JSON.parse(configJSON);
        let thmDir = new String($config.themeDir);
        if (thmDir.includes("themes")) {
          //
          // The config file has the error I introduced on earlier versions. Fix it.
          //
          $config.themeDir = thmDir.replaceAll("themes", "styles");
          await App.WriteFile(configloc, JSON.stringify($config));
        }
      } else {
        //
        // Create the default config file.
        //
        await createDefaultConfig(hmdir, configdir);
      }
    } else {
      //
      // The config directory needs to be created and populated.
      //
      await App.MakeDir(configdir);
      await createDefaultConfig(hmdir, configdir);
    }
    //
    // TODO: This fixes old configs and needs to be made more generic. This needs to be
    // made configuratable in the preferences.
    //
    if (typeof $config.AlfredData === "undefined") {
      $config.AlfredData = "/Library/Application Support/Alfred/Workflow Data";
      await App.WriteFile(configloc, JSON.stringify($config));
    }
    if (typeof $config.shell === "undefined") {
      $config.shell = "/bin/zsh";
      await App.WriteFile(configloc, JSON.stringify($config));
    }
    if (typeof $config.height === "undefined") {
      $config.height = 608;
      $config.width = 1022;
      await App.WriteFile(configloc, JSON.stringify($config));
    }
  }

  async function createDefaultConfig(homedir, configdir) {
    //
    // Create the default Configuration file.
    //
    let themedir = await App.AppendPath(configdir, "styles");
    let scriptsdir = await App.AppendPath(configdir, "scripts");
    let configloc = await App.AppendPath(configdir, "config.json");
    let emailaccountloc = await App.AppendPath(configdir, "emailacounts.json");
    let emailsloc = await App.AppendPath(configdir, "emails.json");
    let envloc = await App.AppendPath(configdir, "environments.json");
    let extscriptloc = await App.AppendPath(configdir, "extscripts.json");
    let scriptsloc = await App.AppendPath(configdir, "scripts.json");
    let currentthemeloc = await App.AppendPath(configdir, "theme.json");
    let systemdir = await App.GetExecutable();
    $config = {
      homeDir: homedir,
      systemDir: systemdir,
      themeDir: themedir,
      theme: "Default",
      configDir: configdir,
      scriptsDir: scriptsdir,
      AlfredData: "/Library/Application Support/Alfred/Workflow Data/",
      shell: "/bin/zsh",
      height: 608,
      width: 1022,
    };
    //
    // TODO: The AlfredData, shell needs to be made configuratable in the preferences.
    //
    await App.WriteFile(configloc, JSON.stringify($config));

    //
    // Create the different default files that are empty.
    //
    await App.MakeDir(scriptsdir);
    await App.WriteFile(emailaccountloc, JSON.stringify([]));
    await App.WriteFile(emailsloc, JSON.stringify([]));
    await App.WriteFile(envloc, JSON.stringify([]));
    await App.WriteFile(extscriptloc, JSON.stringify([]));
    await App.WriteFile(scriptsloc, JSON.stringify([]));

    //
    // Create the default Notes.
    //
    var notesloc = await App.AppendPath(configdir, "notes.json");
    await App.WriteFile(
      notesloc,
      JSON.stringify(["", "", "", "", "", "", "", "", "", ""])
    );

    //
    // Create the default Theme.
    //
    let defaultThemeDir = await App.AppendPath($config.themeDir, "Default");
    if (!(await App.DirExists(defaultThemeDir))) {
      await App.MakeDir(defaultThemeDir);
    }
    let defaultThemefile = await App.AppendPath(
      defaultThemeDir,
      "Default.json"
    );
    $theme = {
      name: "Default",
      font: "Fira Code, Menlo",
      fontSize: "16pt",
      textAreaColor: "#454158",
      backgroundColor: "#22212C",
      textColor: "#80ffea",
      borderColor: "#1B1A23",
      Cyan: "#80FFEA",
      Green: "#8AFF80",
      Orange: "#FFCA80",
      Pink: "#FF80BF",
      Purple: "#9580FF",
      Red: "#FF9580",
      Yellow: "#FFFF80",
      functionColor: "#9580FF",
      stringColor: "#8AFF80",
      constantColor: "#FFCA80",
      keywordColor: "#FFFF80",
      highlightBackgroundColor: "#4f4f5f",
      selectionColor: "#22212C",
      buttons: [
        {
          color: "#80FFEA",
          id: 0,
        },
        {
          color: "#8AFF80",
          id: 1,
        },
        {
          color: "#FFCA80",
          id: 2,
        },
        {
          color: "#FF80BF",
          id: 3,
        },
        {
          color: "#9580FF",
          id: 4,
        },
        {
          color: "#FF9580",
          id: 5,
        },
        {
          color: "blue",
          id: 5,
        },
        {
          color: "green",
          id: 7,
        },
        {
          color: "red",
          id: 8,
        },
        {
          color: "purple",
          id: 9,
        },
      ],
    };
    await App.WriteFile(defaultThemefile, JSON.stringify($theme));
    await App.WriteFile(currentthemeloc, JSON.stringify($theme));
    let themecnfg = await App.AppendPath(defaultThemeDir, "package.json");
    await App.WriteFile(
      themecnfg,
      JSON.stringify({
        name: "Default",
        version: "1.0.0",
        description: "The EmailIt program default theme.",
        keywords: ["emailit", "theme"],
        author: "Richard Guay",
        license: "MIT",
        theme: {
          name: "Default",
          description: "The EmailIt program default theme.",
          type: 0,
          github: "",
          main: "Default.json",
        },
      })
    );
  }

  async function getTheme() {
    let thm = await App.AppendPath($config.configDir, "theme.json");
    let thmjson = await App.ReadFile(thm);
    $theme = JSON.parse(thmjson);
  }

  async function getScriptsList() {
    let userScriptsLoc = await App.AppendPath(
      $config.configDir,
      "scripts.json"
    );
    let extScriptsLoc = await App.AppendPath(
      $config.configDir,
      "extscripts.json"
    );
    let userScriptFile = {};
    if (await App.FileExists(userScriptsLoc)) {
      //
      // Load the user scripts file.
      //
      userScriptFile = await App.ReadFile(userScriptsLoc);
      $userScripts = JSON.parse(userScriptFile);
    } else {
      //
      // Save a default user scripts file.
      //
      $userScripts = [];
    }
    if (await App.FileExists(extScriptsLoc)) {
      //
      // Load the external scripts file.
      //
      let extScriptFile = await App.ReadFile(extScriptsLoc);
      $extScripts = JSON.parse(extScriptFile);
    } else {
      //
      // Save a default external scripts file.
      //
      $extScripts = [];
    }
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
          })
      )
      .concat(
        $extScripts
          .filter((value) => value.termscript === false)
          .map((value) => {
            return { name: value.name, insert: false };
          })
      );
  }

  function getTermScriptsList() {
    $termscripts = $userScripts
      .filter((value) => value.termscript === true)
      .concat($systemScripts.filter((value) => value.termscript === true))
      .concat($extScripts.filter((value) => value.termscript === true));
  }

  async function getTemplatesList() {
    let templateloc = await App.AppendPath($config.configDir, "templates.json");
    if (await App.FileExists(templateloc)) {
      let templatefile = await App.ReadFile(templateloc);
      $userTemplates = JSON.parse(templatefile);
      $templates = $systemTemplates.concat($userTemplates);
    } else {
      var defaultData = {};
      defaultData.name = "Defaults";
      defaultData.template = `{
   "projectName": "This is a default test project",
   "email": "Your Email",
   "firstName": "Your Frist Name",
   "lastName": "You Last Name",
   "home": {
      "address1": "",
      "address2": "",
      "city": "",
      "state": "",
      "zip": ""
   },
   "office": {
      "address1": "",
      "address2": "",
      "city": "",
      "state": "",
      "zip": ""
   }}`;
      $userTemplates = [];
      $userTemplates.push(defaultData);
      $templates = $systemTemplates.concat($userTemplates);
      await saveUserTemplates();
    }
  }

  async function saveUserTemplates() {
    let templateloc = await App.AppendPath($config.configDir, "templates.json");
    await App.WriteFile(templateloc, JSON.stringify($userTemplates));
  }

  //
  // Function:         runJavaScriptScriptFile
  //
  // Description:      This will run the JavaScript function on the contents of a file.
  //
  // Inputs:
  //                   script          The script.
  //                   text            The text to process.
  //
  async function runJavaScriptFile(script, file) {
    let result = "Invalid File";
    if (await App.FileExists(file)) {
      result = await App.ReadFile(file);
      result = runJavaScript(script, result);
      let parts = await App.SplitFile(file);
      let nfile = await App.AppendPath(
        parts.Dir,
        `${parts.Name}-processed${parts.Extension}`
      );
      await App.WriteFile(nfile, result);
      result = `${parts.Name}-processed${parts.Extension} was created!`;
    }
    return result;
  }

  //
  // Function:         runJavaScript
  //
  // Description:      This will run some given text with a script.
  //
  // Inputs:
  //                   script          The script.
  //                   text            The text to process.
  //
  function runJavaScript(script, text) {
    $sp.data.originalText = text;
    $sp.that = $sp;
    $sp.text = text;
    $sp.data.text = text;

    //
    // Try to evaluate the expression.
    //
    try {
      var scriptFunction = new Function("SP", `${script} ; return SP;`);
      $sp.text = scriptFunction($sp).text;
    } catch (error) {
      console.error(error);
      $sp.text = $sp.data.originalText;
    }
    //
    // Make sure we have a string and not an array or object.
    //
    if (typeof $sp.text != "string") {
      $sp.text = $sp.text.toString();
    }
    return $sp.text;
  }

  async function runExtScript(extScrpt, text, scriptEnv) {
    //
    // extScript.name     - File name of the script
    // extScript.script   - User name for the script
    // extScript.path     - directory of the script
    // extScript.env      - name of the environment
    // extScript.termscript - Boolean true if a script terminal script.
    // extScript.description - Description of what the script does.
    // extScript.help     - A help message for the script
    //
    var result = "";
    var env = {};

    //
    // Get the environment.
    //
    if (extScrpt.env !== "") {
      let envlistloc = await App.AppendPath(
        $config.configDir,
        "environments.json"
      );
      let envlist = await App.ReadFile(envlistloc);
      envlist = JSON.parse(envlist);
      env = envlist.find((item) => item.name === extScrpt.env);
      if (env !== "undefined") {
        env = env.envVar;
      } else {
        env = {};
      }
    }
    if (scriptEnv !== null) {
      for (let i = 0; i < scriptEnv.length; i++) {
        let parts = scriptEnv[i].split("=");
        env[parts[0]] = parts[1];
      }
    }

    //
    // Add the $sp.data structure to the environment.
    //
    env["RUNNINGDATA"] = JSON.stringify($sp.data);

    try {
      let args = [];
      args.push(text);
      var strEnv = [];
      for (const key in env) {
        strEnv.push(`${key}=${env[key]}`);
      }
      result = await App.RunCommandLine(
        extScrpt.script,
        args,
        strEnv,
        extScrpt.path
      );
      if (result !== null && typeof result === "object")
        result = result.toString();
    } catch (e) {
      result = "Error: " + e.message;
    }
    return result;
  }

  async function runScript(script, text, env) {
    var result = "Script not found!";

    //
    // Check to see if it is text or a file or directory name.
    //
    let isLongText = text.includes("\n");
    let isfile = false;
    let isDir = false;
    if (!isLongText) {
      isfile = await App.FileExists(text);
      isDir = await App.DirExists(text);
    }
    if (typeof script === "object") {
      if (typeof script.name !== "undefined") {
        script = script.name;
      }
    }

    //
    // Find the script and run it.
    //
    var scriptIndex = $userScripts.find((ele) => {
      return ele.name === script;
    });
    if (typeof scriptIndex !== "undefined") {
      script = scriptIndex.script;
      if (isfile && !isDir && !isLongText) {
        result = await runJavaScriptFile(script, text);
      } else {
        result = runJavaScript(script, text);
      }
    } else {
      scriptIndex = $systemScripts.find((ele) => {
        return ele.name === script;
      });
      if (typeof scriptIndex !== "undefined") {
        script = scriptIndex.script;
        if (isfile && !isDir && !isLongText) {
          result = await runJavaScriptFile(script, text);
        } else {
          result = runJavaScript(script, text);
        }
      } else {
        scriptIndex = $extScripts.find((ele) => {
          return ele.name === script;
        });
        if (typeof scriptIndex !== "undefined") {
          //
          // It's an external script.
          //
          result = await runExtScript(scriptIndex, text, env);
        } else {
          //
          // It is a command line. Run it directly.
          //
          let tmpscript = await App.CreateTempFile(`#!${$config.shell}

${text}
`);

          //
          // Change the mode.
          //
          await App.Chmod(tmpscript, 0o777);

          //
          // Run the script.
          //
          result = await runExtScript(
            {
              name: "",
              script: tmpscript,
              env: "Default",
              path: "",
            },
            "",
            env
          );
          result = `\n${result}`;

          //
          // Remove the temperary script.
          //
          await App.DeleteEntries(tmpscript);
        }
      }
    }
    return result;
  }

  async function keyDownProcessor(e) {
    if (e.metaKey && e.key === ",") {
      e.preventDefault();
      $state = "preferences";
    } else if (e.ctrlKey) {
      switch (e.key) {
        case "e":
          $state = "emailit";
          await rt.WindowSetSize($config.width, $config.height);
          e.preventDefault();
          break;

        case "n":
          $state = "notes";
          await rt.WindowSetSize($config.width, $config.height);
          e.preventDefault();
          break;

        case "m":
          $showScripts = !$showScripts;
          await rt.WindowSetSize($config.width, $config.height);
          e.preventDefault();
          break;

        case "t":
          $showTemplates = !$showTemplates;
          await rt.WindowSetSize($config.width, $config.height);
          e.preventDefault();
          break;

        case "p":
          $state = "preferences";
          await rt.WindowSetSize($config.width, $config.height);
          e.preventDefault();
          break;

        case "s":
          $state = "scriptline";
          e.preventDefault();
          break;

        case "1":
          openNote(0);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "2":
          openNote(1);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "3":
          openNote(2);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "4":
          openNote(3);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "5":
          openNote(4);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "6":
          openNote(5);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "7":
          openNote(6);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "8":
          openNote(7);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "9":
          openNote(8);
          if ($state === "notes") $noteEditor.focus();
          break;

        case "0":
          openNote(9);
          if ($state === "notes") $noteEditor.focus();
          break;
      }
    }
  }
</script>

<svelte:window on:keydown={keyDownProcessor} />
<div
  id="dragbar"
  style="background-color: {$theme.backgroundColor}; 
         font-family: {$theme.font}; 
         color: {$theme.textColor}; 
         font-size: {$theme.fontSize};"
>
  <div
    id="close"
    on:click={() => {
      rt.Quit();
    }}
  >
    <span style="color: {$theme.Red};">X</span>
  </div>
</div>

{#if $state === "emailit"}
  <EmailIt />
{:else if $state === "notes"}
  <Notes />
{:else if $state === "scripts"}
  <ScriptsEditor />
{:else if $state === "templates"}
  <TemplatesEditor />
{:else if $state === "preferences"}
  <Preferences />
{:else if $state === "scriptline"}
  <ScriptLine />
{:else}
  <div id="error">
    <h1>Something went wrong!</h1>
  </div>
{/if}

<ScriptMenu />
<TemplateMenu />

<style>
  :global(body) {
    background-color: rgba(0, 0, 0, 0);
  }

  #dragbar {
    height: 20px;
    width: 1022px;
    border-radius: 10px 10px 0px 0px;
    border: 0px transparent;
    --wails-draggable: drag;
  }
  #close {
    font-size: 13px;
    height: 15px;
    width: 15px;
    margin: 0px;
    padding: 5px 0px 0px 7px;
    cursor: pointer;
  }
</style>
