import { writable } from 'svelte/store';

let defaultTheme = {
  name: "default",
  font: "Fira Code, Menlo",
  fontSize: "16pt",
  textAreaColor: '#454158',
  backgroundColor: '#22212C',
  textColor: '#80ffea',
  borderColor: '#1B1A23',
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
      color: '#80FFEA',
      id: 0
    },
    {
      color: '#8AFF80',
      id: 1
    },
    {
      color: '#FFCA80',
      id: 2
    },
    {
      color: '#FF80BF',
      id: 3
    },
    {
      color: '#9580FF',
      id: 4
    },
    {
      color: '#FF9580',
      id: 5
    },
    {
      color: 'blue',
      id: 6
    },
    {
      color: 'green',
      id: 7
    },
    {
      color: 'red',
      id: 8
    },
    {
      color: 'purple',
      id: 9
    },
  ]
};

export const theme = writable( defaultTheme );

