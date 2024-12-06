<script>
  import { onMount } from "svelte";
  import { createEventDispatcher } from "svelte";
  import { vim } from "@replit/codemirror-vim";
  import {
    lineNumbers,
    highlightActiveLineGutter,
    highlightSpecialChars,
    drawSelection,
    dropCursor,
    rectangularSelection,
    crosshairCursor,
    highlightActiveLine,
    keymap,
    EditorView,
  } from "@codemirror/view";
  import { EditorState } from "@codemirror/state";
  import {
    foldGutter,
    indentOnInput,
    syntaxHighlighting,
    defaultHighlightStyle,
    bracketMatching,
    foldKeymap,
  } from "@codemirror/language";
  import { history, defaultKeymap, historyKeymap } from "@codemirror/commands";
  import { highlightSelectionMatches, searchKeymap } from "@codemirror/search";
  import {
    closeBrackets,
    autocompletion,
    closeBracketsKeymap,
    completionKeymap,
  } from "@codemirror/autocomplete";
  import { lintKeymap } from "@codemirror/lint";
  import { markdown } from "@codemirror/lang-markdown";
  import { javascript } from "@codemirror/lang-javascript";
  import { scheme } from "@codemirror/legacy-modes/mode/scheme";
  import { theme } from "../stores/theme.js";

  const dispatch = createEventDispatcher();

  export let height = 0;
  export let width = 0;
  export let styling = "";
  export let config = {};

  let CodeMirrorEditor = null;
  let edView = null;
  let editorFunctions = null;
  let currentCursor = null;

  //
  // This is the basic theming definitions.
  //
  const editorTheme = /*@__PURE__*/ EditorView.theme(
    {
      "&": {
        color: $theme.textColor,
        backgroundColor: $theme.textAreaColor,
        height: height,
        minHeight: height,
      },
      ".cm-content": {
        caretColor: $theme.Cyan,
        minHeight: height,
      },
      "&.cm-focused .cm-cursor": { borderLeftColor: $theme.Cyan },
      "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, ::selection":
        { backgroundColor: $theme.selectionColor },
      ".cm-panels": {
        backgroundColor: $theme.backgroundColor,
        color: $theme.textColor,
      },
      ".cm-panels.cm-panels-top": { borderBottom: "2px solid black" },
      ".cm-panels.cm-panels-bottom": { borderTop: "2px solid black" },
      ".cm-searchMatch": {
        backgroundColor: "#72a1ff59",
        outline: "1px solid #457dff",
      },
      ".cm-searchMatch.cm-searchMatch-selected": {
        backgroundColor: "#6199ff2f",
      },
      ".cm-activeLine": { backgroundColor: $theme.highlightBackgroundColor },
      ".cm-selectionMatch": { backgroundColor: "#aafe661a" },
      ".cm-matchingBracket, .cm-nonmatchingBracket": {
        backgroundColor: "#bad0f847",
        outline: "1px solid #515a6b",
      },
      ".cm-gutter": {
        minHeight: height,
      },
      ".cm-gutters": {
        backgroundColor: $theme.backgroundColor,
        color: $theme.green,
        border: "none",
      },
      ".cm-activeLineGutter": {
        backgroundColor: $theme.highlightBackgroundColor,
      },
      ".cm-foldPlaceholder": {
        backgroundColor: "transparent",
        border: "none",
        color: "#ddd",
      },
      ".cm-tooltip": {
        border: "1px solid #181a1f",
        backgroundColor: $theme.backgroundColor,
      },
      ".cm-tooltip-autocomplete": {
        "& > ul > li[aria-selected]": {
          backgroundColor: $theme.highlightBackgroundColor,
          color: $theme.textColor,
        },
      },
    },
    { dark: true },
  );

  function fire(name, data) {
    dispatch(name, {
      data: data,
    });
  }

  function setValue(text) {
    //
    // Since we are setting a whole new document, create new editor
    // states and views.
    //
    CreateEditorState(text);
  }

  function CreateEditorState(text) {
    //
    // Clear out the div element in case a previous editor was
    // created.
    //
    if (edView !== null) {
      edView.destroy();
      CodeMirrorEditor.innerHTML = "";
    }

    //
    // Setup the extensions array.
    //
    const exts = [
      vim(),
      highlightSpecialChars(),
      history(),
      drawSelection(),
      dropCursor(),
      EditorState.allowMultipleSelections.of(true),
      indentOnInput(),
      syntaxHighlighting(defaultHighlightStyle, {
        fallback: true,
      }),
      bracketMatching(),
      closeBrackets(),
      autocompletion(),
      rectangularSelection(),
      crosshairCursor(),
      highlightSelectionMatches(),
      keymap.of([
        ...closeBracketsKeymap,
        ...defaultKeymap,
        ...searchKeymap,
        ...historyKeymap,
        ...foldKeymap,
        ...completionKeymap,
        ...lintKeymap,
      ]),
      editorTheme,
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          fire("textChange", {
            value: getValue(),
            cursor: getCursor(),
          });
        }
      }),
    ];

    //
    // Add extensions based on the configuration.
    //
    if (config.lineNumbers) {
      exts.push(foldGutter());
      exts.push(lineNumbers());
    }

    switch (config.language.toLowerCase()) {
      case "markdown":
        exts.push(markdown());
        break;
      case "javascript":
        exts.push(javascript());
        break;

      /*
      case "prolog":
        exts.push(prolog()); // When we get a good prolog highlighter. 
        break;*/

      case "lips":
        exts.push(scheme());
        break;

      default:
        exts.push(markdown());
        break;
    }

    if (config.lineWrapping) {
      exts.push(EditorView.lineWrapping);
    }

    if (config.lineHighlight) {
      exts.push(highlightActiveLine());
    }

    //
    // Create the editor View.
    //
    edView = new EditorView({
      doc: text,
      extensions: exts,
      parent: CodeMirrorEditor,
    });
  }

  onMount(() => {
    //
    // Create the editor.
    //
    CreateEditorState("");

    //
    // Create the editor functions object.
    //
    editorFunctions = {
      getSelection: getSelection,
      getValue: getValue,
      replaceSelection: replaceSelection,
      somethingSelected: somethingSelected,
      setCursor: setCursor,
      getCursor: getCursor,
      setValue: setValue,
      getLine: getLine,
      focus: focus,
      getEdView: getEdView,
      isFocused: isFocused,
      insertAtCursor: insertAtCursor,
    };

    //
    // Give the parent the functions for interacting with the editor.
    //
    fire("editorChange", editorFunctions);

    //
    // Make sure the editor is focused.
    //
    focus();

    //
    // Return a function to run to clean up after mounting.
    //
    return () => {
      // this function runs when the
      // component is destroyed
      edView = null;
      editorFunctions = null;
    };
  });

  function insertAtCursor(text) {
    if (typeof edView !== "undefined") {
      let point = getCursor();
      let transaction = edView.state.update({
        changes: [{ from: point, insert: text }],
      });
      edView.dispatch(transaction);
    }
  }

  function isFocused() {
    if (typeof edView !== "undefined") {
      return edView.hasFocus;
    }
    return false;
  }

  function getLine(pos) {
    if (typeof edView !== "undefined") {
      return edView.docView.domAtPos(pos).node.textContent;
    }
    return "";
  }

  function getSelection() {
    if (typeof edView !== "undefined") {
      return edView.state.sliceDoc(
        edView.state.selection.main.from,
        edView.state.selection.main.to,
      );
    }
  }

  function replaceSelection(newText) {
    if (typeof edView !== "undefined") {
      let point = edView.state.selection.main.from;
      let transaction = edView.state.update({
        changes: [
          {
            from: edView.state.selection.main.from,
            to: edView.state.selection.main.to,
          },
          { from: point, insert: newText },
        ],
      });
      edView.dispatch(transaction);
    }
  }

  function somethingSelected() {
    if (typeof edView !== "undefined") {
      return edView.state.selection.ranges.some((r) => !r.empty);
    }
  }

  function setCursor(pos) {
    if (typeof edView !== "undefined") {
      currentCursor = pos;
      edView.dispatch({ selection: { anchor: currentCursor } });
    }
  }

  function getCursor() {
    if (typeof edView !== "undefined") {
      currentCursor = edView.state.selection.main.head;
      return currentCursor;
    } else {
      return 0;
    }
  }

  function getValue() {
    if (typeof edView !== "undefined") {
      return edView.state.doc.toString();
    }
  }

  function focus() {
    if (typeof edView !== "undefined") {
      edView.focus();
    }
    fire("focus", {});
  }

  function getEdView() {
    return edView;
  }
</script>

<div
  name="editor"
  id="CMeditor"
  bind:this={CodeMirrorEditor}
  on:focus={() => {
    fire("focus", {});
  }}
  style="height: {height}; width: {width}; {styling}"
/>

<style>
  #CMeditor {
    height: 100%;
    width: 100%;
    margin: 0px;
  }

  #CMeditor:focus {
    outline-color: transparent;
  }

  :global(.cm-wrap) {
    height: 100%;
  }

  :global(.cm-scroller) {
    overflow: auto;
  }
</style>
