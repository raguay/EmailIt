const util = {
  runCommandLine: runCommandLine,
}

async function runCommandLine(line, rEnv, callback, dir) {
  //
  // Get the environment to use.
  //
  var resp = await fetch(`http://localhost:9978/api/scripts/env/Default`);
  var nEnv = await resp.json();
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

export default util;
