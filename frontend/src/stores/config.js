import { writable } from 'svelte/store';

let defaultconfig = {
  homeDir: "",
  themeDir: "",
  systemDir: "",
  theme: "",
  configDir: "",
  scriptsDir: "",
  AlfredData: "",
  shell: "",
};

export const config = writable(defaultconfig);

