import { writable } from 'svelte/store';

let defaultconfig = {
  homeDir: "",
  themeDir: "",
  systemDir: "",
  theme: "",
  configDir: "",
  scriptsDir: ""
};

export const config = writable(defaultconfig);

