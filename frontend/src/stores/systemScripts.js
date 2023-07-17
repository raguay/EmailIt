import { writable } from "svelte/store";

export const systemScripts = writable([
  {
    name: "Upper Case",
    insert: false,
    description: "Makes the text all Upper Cased.",
    help: "",
    script: "SP.text = SP.text.toUpperCase()",
    termscript: false,
    system: true,
  },
  {
    name: "Lower Case",
    insert: false,
    description: "Makes the text all lower cased.",
    help: "",
    script: "SP.text = SP.text.toLowerCase()",
    termscript: false,
    system: true,
  },
  {
    name: "Title Case",
    insert: false,
    description: "Makes the text Title cased.",
    help: "",
    script:
      "function TitleCase(t){var a=t.split(' ');var e=['to','an','and','at','as','but','by','for','if','in','on','or','is','with','a','the','of','vs','vs.','via','via','en','I','II','III','IV','V','VI','VII','VIII','IX','X','HTML','CSS','AT&T','PHP','Python','JavaScript','IBM'];for(var r=0;r<a.length;r++){var o=a[r];var s=false;for(var I=0;I<e.length;I++){if(o.toLowerCase()===e[I].toLowerCase()){o=e[I];s=true}}if(!s){o=o.charAt(0).toUpperCase()+o.substr(1).toLowerCase()}a[r]=o}a[0]=a[0].charAt(0).toUpperCase()+a[0].substr(1).toLowerCase();return a.join(' ')}SP.text=TitleCase(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Snake Case",
    insert: false,
    description: "Makes the text Snake cased, or underlines instead of spaces.",
    help: "",
    script:
      "function SnakeCase(e){if(e.indexOf(' ')!==-1)return e.split(' ').join('_');else return e.split('_').join(' ')}SP.text=SnakeCase(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Lower Camel Case",
    insert: false,
    description:
      "Makes all the words together with the first word all lower cased and each other word capitalized.",
    help: "",
    script:
      "function LowerCamelCase(e){var a=e.split(' ');for(var r=0;r<a.length;r++){var t=a[r];if(r!==0){t=t.charAt(0).toUpperCase()+t.substr(1).toLowerCase()}else{t=t.toLowerCase()}a[r]=t}return a.join('')}SP.text=LowerCamelCase(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Upper Camel Case",
    insert: false,
    description: "Makes all the words capitalized and without spaces.",
    help: "",
    script:
      "function UpperCamelCase(e){var a=e.split(' ');for(var r=0;r<a.length;r++){var t=a[r];t=t.charAt(0).toUpperCase()+t.substr(1).toLowerCase();a[r]=t}return a.join('')}; SP.text=UpperCamelCase(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Hyphen Case",
    insert: false,
    description:
      "Makes all the spaces into hyphens. Running it again changes all the hyphens to spaces.",
    help: "",
    script:
      "function HyphenCase(e){if(e.indexOf(' ')!==-1)return e.split(' ').join('-');else return e.split('-').join(' ')}SP.text=HyphenCase(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "URI Encode",
    insert: false,
    description: "Encodes the text for a URI.",
    help: "",
    script: "SP.text = encodeURIComponent( SP.text )",
    termscript: false,
    system: true,
  },
  {
    name: "URI Decode",
    insert: false,
    description: "Decodes URI encoded text.",
    help: "",
    script: "SP.text = decodeURIComponent( SP.text )",
    termscript: false,
    system: true,
  },
  {
    name: "HTML Escape",
    insert: false,
    description: "Escapes special characters using HTML escaping.",
    help: "",
    script:
      "var e = document.createElement( 'DIV' ); var t = document.createTextNode( SP.text );e.appendChild( t );SP.text = e.innerHTML;",
    termscript: false,
    system: true,
  },
  {
    name: "Undo HTML Escaping",
    insert: false,
    description: "Undoes HTML escaping.",
    help: "",
    script:
      "var e = document.createElement( 'DIV' );e.innerHTML = SP.text; SP.text = e.childNodes.length === 0 ? ' : e.childNodes[ 0 ].nodeValue;",
    termscript: false,
    system: true,
  },
  {
    name: "Surround with ()",
    insert: false,
    description: "Surrounds the text in ().",
    help: "",
    script: "SP.text = '(' + SP.text + ')'",
    termscript: false,
    system: true,
  },
  {
    name: "Surround with []",
    insert: false,
    description: "Surrounds the text in [].",
    help: "",
    script: "SP.text = '[' + SP.text + ']'",
    termscript: false,
    system: true,
  },
  {
    name: "Surround with {}",
    insert: false,
    description: "Surrounds the text in {}.",
    help: "",
    script: "SP.text = '{' + SP.text + '}'",
    termscript: false,
    system: true,
  },
  {
    name: "Surround with **",
    insert: false,
    description: "Surrounds the text in '**'.",
    help: "",
    script: "SP.text = '**' + SP.text + '**'",
    termscript: false,
    system: true,
  },
  {
    name: "Bullet Lines with *",
    insert: false,
    description: "Puts a '* ' at the front of each line of text.",
    help: "",
    script:
      "function BulletItems(t){var l=t.match(/^.*((\\r\\n|\\n|\\r)|$)/gm);for(var n=0;n<l.length;n++){l[n]='* '+l[n]}return l.join('')}SP.text=BulletItems(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Bullet Lines with Numbers",
    insert: false,
    description:
      "Puts a number, period, and a space in the front of each line of text. The numbers are sequential from 1.",
    help: "",
    script:
      "function NumberItems(t){var r=t.match(/^.*((\\r\\n|\\n|\\r)|$)/gm);var n=1;for(var e=0;e<r.length;e++){r[e]=n.toString()+'. '+r[e];n++}return r.join('')}SP.text=NumberItems(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Undo Line Numbering",
    insert: false,
    description: "Removes the line numbering from the beginning of each line.",
    help: "",
    script:
      "for(var lines=SP.text.split('\\n'),i=0;i<lines.length;i++){var match=lines[i].match(/^\\d+\\.\\ (.*)$/);null!=match&&(lines[i]=match[1])}SP.text=lines.join('\\n');",
    termscript: false,
    system: true,
  },
  {
    name: "Insert Date & Time",
    insert: true,
    description: "Inserts the current date and time at the cursor point.",
    help: "",
    script: "SP.text = SP.DateTime.now().toFormat('LLLL dd, yyyy, h:mm:ss a');",
    termscript: false,
    system: true,
  },
  {
    name: "Insert Date",
    insert: true,
    description: "Inserts the current date at the cursor point.",
    help: "",
    script: "SP.text = SP.DateTime.now().toFormat('LLLL dd, yyyy');",
    termscript: false,
    system: true,
  },
  {
    name: "Insert Time",
    insert: true,
    description: "Inserts the current time at the cursor point.",
    help: "",
    script: "SP.text = SP.DateTime.now().toFormat('h:mm:ss a');",
    termscript: false,
    system: true,
  },
  {
    name: "Insert Day of Week",
    insert: true,
    description:
      "Inserts the current name of the day of the week at the cursor point.",
    help: "",
    script: "SP.text = SP.DateTime.now().toFormat('cccc');",
    termscript: false,
    system: true,
  },
  {
    name: "Insert Date/Time by Format",
    insert: false,
    description:
      "Takes the selected text as formatting for the date/time function.",
    help: "",
    script: "SP.text = SP.DateTime.now().toFormat(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Number of Days from Now",
    insert: false,
    description:
      "Takes the selected number and gives the date that many days in the future.",
    help: "",
    script:
      "SP.text = SP.DateTime.now().plus({ days: SP.text}).toFormat('LLLL dd, yyyy');",
    termscript: false,
    system: true,
  },
  {
    name: "Number of Days Ago",
    insert: false,
    description:
      "Takes the selected number and gives the date that many days in the past.",
    help: "",
    script:
      "SP.text = SP.DateTime.now().minus({ days: SP.text}).toFormat('LLLL dd, yyyy');",
    termscript: false,
    system: true,
  },
  {
    name: "Basic Math",
    insert: false,
    description: "Takes the selection and processes as math.",
    help: "",
    script: "SP.text = SP.ProcessMathSelection(SP.text)",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Page for Math",
    insert: false,
    description:
      "Takes the entire note and processes it for math if the line starts with '>'.",
    help: "",
    script: "SP.text = SP.ProcessMathNotes(SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 1 as Script",
    insert: false,
    description:
      "Takes the first note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(1), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 2 as Script",
    insert: false,
    description:
      "Takes the second note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(2), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 3 as Script",
    insert: false,
    description:
      "Takes the third note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(3), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 4 as Script",
    insert: false,
    description:
      "Takes the fourth note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(4), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 5 as Script",
    insert: false,
    description:
      "Takes the fifth note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(5), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 6 as Script",
    insert: false,
    description:
      "Takes the sixth note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(6), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 7 as Script",
    insert: false,
    description:
      "Takes the seventh note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(7), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 8 as Script",
    insert: false,
    description:
      "Takes the eighth note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(8), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "Evaluate Note 9 as Script",
    insert: false,
    description:
      "Takes the nineth note and runs it as the script for the text in the current note.",
    help: "",
    script: "SP.text = SP.runScript(SP.returnNote(9), SP.text);",
    termscript: false,
    system: true,
  },
  {
    name: "help",
    insert: false,
    description: "Show the help command in the terminal.",
    help:
      "help <command>?\n\r    If no command is given, it gives a list of built in terminal commands.\n\r    If given the name of a built in termianl command, it displays that command's help.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `help ${SP.text}`,\n  lines: [],\n});\n",
    termscript: true,
  },
  {
    name: "ls",
    insert: false,
    description: "  List the items in the specified directory.",
    help:
      "ls <command>?\n\r    List the contents of the current directory or the directory given.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `ls ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "cd",
    insert: false,
    description: "  A terminal script for changing directories.",
    help: "cd <directory>\n\r    directory = The directory to change to.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `cd '${SP.text}'`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "open",
    insert: false,
    description: "Open the file given.",
    help:
      "open <file>\n\r    This command opens the file given in the default program for the file.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `open ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "runscript",
    insert: false,
    description: "Run a script on a file or text string.",
    help:
      "runscript <scriptname>,<file>|<text>\n\r    This command will run the script on the file given or on the text given. The two argument are separated by a comma. The text has to be surrounded by quotes.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `runscript ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "edit",
    insert: false,
    description: "Edit the given file.",
    help: "edit <file>\n\r    This command will edit the given file.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `edit ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "alias",
    insert: false,
    description:
      "The alais command allows you to substitute a single command for a list of commands.",
    help:
      "alias <name>=<commands>\n\r     This command will run the comma separated list of commands when the name is ran.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `alias ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "hist",
    insert: false,
    description:
      "The hist command lists the previous command lines that worked. You can activate one to rerun the command.",
    help:
      "hist <depth>\n\r     This command displays the previously ran commands to select one for running again. If a depth is given, it will try to display that many commands. The default depth is 5.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `hist ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "rm",
    insert: false,
    description:
      "  The rm command will delete the given file or directory. If nothing is given, it will list the current directory and an item can be delete using the command mode and r.",
    help:
      "rm <filedir>\n\r     If filedir is blank, then the current directory is listed and an item can be deleted using the command mode and r. Otherwise, the given file/directory is deleted.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `rm ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "mkdir",
    insert: false,
    description:
      "The mkdir command makes the given directory if it doesn't already exist.",
    help:
      "mkdir <dir>\n\r     Create the given dir if it doesn't already exist.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `mkdir ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "mkfile",
    insert: false,
    description:
      "The mkfile command makes the given file if it doesn't already exist.",
    help:
      "mkfile <file>\n\r     Create the given file if it doesn't already exist.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `mkfile ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "rmalias",
    insert: false,
    description:
      "The rmalias command removes the alias from the list of aliases.",
    help: "rmalias <name>\n\r     Remove the alias given.",
    script:
      "SP.text = JSON.stringify({\n  tcommand: `rmalias ${SP.text}`,\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "notes",
    insert: false,
    description: "Show the notes.",
    help: "notes",
    script:
      "SP.text = JSON.stringify({\n  tcommand: 'notes',\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "emailit",
    insert: false,
    description: "Show the EmailIt interface.",
    help: "EmailIt Interface",
    script:
      "SP.text = JSON.stringify({\n  tcommand: 'emailit',\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "editscript",
    insert: false,
    description: "Show the Script Editor.",
    help: "Go to the Script Editor",
    script:
      "SP.text = JSON.stringify({\n  tcommand: 'editscript',\n  lines: []\n});",
    termscript: true,
  },
  {
    name: "quit",
    insert: false,
    description: "Quit the EmailIt program.",
    help: "Quit the EmailIt program",
    script: "SP.text = JSON.stringify({\n  tcommand: 'quit',\n  lines: []\n});",
    termscript: true,
  },
]);
