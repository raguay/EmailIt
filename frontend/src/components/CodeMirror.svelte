<div 
  name="editor" 
  id="CMeditor" 
  bind:this='{CodeMirrorEditor}' 
  style="height: {height}; width: {width}; {styling}"
>
</div>

<style>
  #CMeditor {
    height: 100%;
    width: 100%;
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

<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import { markdown } from "@codemirror/lang-markdown";
  import { HighlightStyle, tags } from '@codemirror/highlight';
  import { highlightSpecialChars, drawSelection, highlightActiveLine, keymap, EditorView } from '@codemirror/view';
  import { EditorState, Prec } from '@codemirror/state';
  import { history, historyKeymap } from '@codemirror/history';
  import { indentOnInput } from '@codemirror/language';
  import { lineNumbers } from '@codemirror/gutter';
  import { defaultKeymap } from '@codemirror/commands';
  import { bracketMatching } from '@codemirror/matchbrackets';
  import { closeBrackets, closeBracketsKeymap } from '@codemirror/closebrackets';
  import { highlightSelectionMatches, searchKeymap } from '@codemirror/search';
  import { autocompletion, completionKeymap } from '@codemirror/autocomplete';
  import { commentKeymap } from '@codemirror/comment';
  import { rectangularSelection } from '@codemirror/rectangular-selection';
  import { defaultHighlightStyle } from '@codemirror/highlight';
  import { lintKeymap } from '@codemirror/lint';
  import { javascript } from '@codemirror/lang-javascript';
  import { theme } from '../stores/theme.js';

  const dispatch = createEventDispatcher();
  
  export let height;
  export let width;
  export let styling = '';
  export let config;
  export let initFinished = false;

  let CodeMirrorEditor;
  let edState;
  let edView;
  let editorFunctions;
  let currentCursor;

  //
  // This is the basic theming definitions.
  //
  const editorTheme = /*@__PURE__*/EditorView.theme({
      "&": {
          color: $theme.textColor,
          backgroundColor: $theme.textAreaColor
      },
      ".cm-content": {
          caretColor: $theme.Cyan
      },
      "&.cm-focused .cm-cursor": { borderLeftColor: $theme.Cyan },
      "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, ::selection": { backgroundColor: $theme.selectionColor },
      ".cm-panels": { backgroundColor: $theme.backgroundColor, color: $theme.textColor },
      ".cm-panels.cm-panels-top": { borderBottom: "2px solid black" },
      ".cm-panels.cm-panels-bottom": { borderTop: "2px solid black" },
      ".cm-searchMatch": {
          backgroundColor: "#72a1ff59",
          outline: "1px solid #457dff"
      },
      ".cm-searchMatch.cm-searchMatch-selected": {
          backgroundColor: "#6199ff2f"
      },
      ".cm-activeLine": { backgroundColor: $theme.highlightBackgroundColor },
      ".cm-selectionMatch": { backgroundColor: "#aafe661a" },
      ".cm-matchingBracket, .cm-nonmatchingBracket": {
          backgroundColor: "#bad0f847",
          outline: "1px solid #515a6b"
      },
      ".cm-gutters": {
          backgroundColor: $theme.backgroundColor,
          color: $theme.green,
          border: "none"
      },
      ".cm-activeLineGutter": {
          backgroundColor: $theme.highlightBackgroundColor
      },
      ".cm-foldPlaceholder": {
          backgroundColor: "transparent",
          border: "none",
          color: "#ddd"
      },
      ".cm-tooltip": {
          border: "1px solid #181a1f",
          backgroundColor: $theme.backgroundColor
      },
      ".cm-tooltip-autocomplete": {
          "& > ul > li[aria-selected]": {
              backgroundColor: $theme.highlightBackgroundColor,
              color: $theme.textColor
          }
      }
  }, { dark: true });
  /**
  The highlighting style for code.
  */
  const editorHighlightStyle = /*@__PURE__*/HighlightStyle.define([
      { tag: tags.keyword,
          color: $theme.keywordColor },
      { tag: [tags.name, tags.deleted, tags.character, tags.propertyName, tags.macroName],
          color: $theme.Pink },
      { tag: [/*@__PURE__*/tags.function(tags.variableName), tags.labelName],
          color: $theme.functionColor },
      { tag: [tags.color, /*@__PURE__*/tags.constant(tags.name), /*@__PURE__*/tags.standard(tags.name)],
          color: $theme.constantColor },
      { tag: [/*@__PURE__*/tags.definition(tags.name), tags.separator],
          color: $theme.textColor },
      { tag: [tags.typeName, tags.className, tags.number, tags.changed, tags.annotation, tags.modifier, tags.self, tags.namespace],
          color: $theme.Yellow },
      { tag: [tags.operator, tags.operatorKeyword, tags.url, tags.escape, tags.regexp, tags.link, /*@__PURE__*/tags.special(tags.string)],
          color: $theme.Cyan },
      { tag: [tags.meta, tags.comment],
          color: $theme.green },
      { tag: tags.strong,
          fontWeight: "bold" },
      { tag: tags.emphasis,
          fontStyle: "italic" },
      { tag: tags.link,
          color: $theme.green,
          textDecoration: "underline" },
      { tag: tags.heading,
          fontWeight: "bold",
          color: $theme.Pink },
      { tag: [tags.atom, tags.bool, /*@__PURE__*/tags.special(tags.variableName)],
          color: $theme.constantColor },
      { tag: [tags.processingInstruction, tags.string, tags.inserted],
          color: $theme.stringColor },
      { tag: tags.invalid,
          color: $theme.Red },
  ]);
  /**
  Extension to enable the theme.
  */
  const editor = [editorTheme, editorHighlightStyle];

  function fire(name, data) {
    dispatch(name, {
      data: data
    });
  }

  function setValue(text) {
    // 
    // Since we are setting a whole new document, create new editor 
    // states and views.
    // 
    if(initFinished) {
      CreateEditorState(text);
    }
  }

  function CreateEditorState(text) {
    // 
    // Clear out the div element in case a previous editor was
    // created.
    //
    CodeMirrorEditor.innerHTML = '';

    //
    // Setup the extensions array.
    //
    const exts = [
      highlightSpecialChars(),
      history(),
      drawSelection(),
      EditorState.allowMultipleSelections.of(true),
      indentOnInput(),
      Prec.fallback(defaultHighlightStyle),
      bracketMatching(),
      closeBrackets(),
      autocompletion(),
      rectangularSelection(),
      highlightSelectionMatches(),
      keymap.of([
          ...closeBracketsKeymap,
          ...defaultKeymap,
          ...searchKeymap,
          ...historyKeymap,
          ...commentKeymap,
          ...completionKeymap,
          ...lintKeymap
      ]),
      editor,
      EditorView.updateListener.of(update => {
        if(update.docChanged) {
          fire('textChange', {
            value: getValue(),
            cursor: getCursor()
          })
        }
      })
    ];

    // 
    // Add extensions based on the configuration.
    // 
    if(config.lineNumbers) {
      exts.push(lineNumbers());
    }

    switch(config.language) {
      case 'markdown':
        exts.push(markdown());
        break;
      case 'javascript':
        exts.push(javascript());
        break;
      default: 
        exts.push(markdown());
        break;
    }

    if(config.lineWrapping) {
      exts.push(EditorView.lineWrapping);
    }

    if(config.lineHighlight) {
      exts.push(highlightActiveLine());
    }
    
    // 
    // Create the editor state.
    //
    edState = EditorState.create({
      doc: text,
      extensions: exts
    });

    // 
    // Create the editor View.
    // 
    edView = new EditorView({
      state: edState,
      parent: CodeMirrorEditor
    });
  }

  onMount(() => {
    // 
    // Create the editor.
    // 
    CreateEditorState('');

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
      getEdState: getEdState,
      isFocused: isFocused,
      insertAtCursor: insertAtCursor
    };
    
    //
    // Give the parent the functions for interacting with the editor.
    //
    fire('editorChange', editorFunctions);

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
    };
  });

  function insertAtCursor(text) {
    console.log(text);
    if(typeof edView !== 'undefined') {
      let point = getCursor();
      let transaction = edView.state.update({changes: [{from: point, insert: text}]});
      edView.dispatch(transaction);
    }
  }

  function isFocused() {
    if(typeof edView !== 'undefined') {
      return(edView.hasFocus);
    }
    return(false);
  }

  function getLine(pos) {
    if(typeof edView !== 'undefined') {
      return(edView.docView.domAtPos(pos).node.textContent);
    }
    return('');
  }

  function getSelection() {
    if(typeof edView !== 'undefined') {
      return edView.state.sliceDoc(
        edView.state.selection.main.from,
        edView.state.selection.main.to);
    }
  }

  function replaceSelection(newText) {
    console.log(newText);
    if(typeof edView !== 'undefined') {
      let point = edView.state.selection.main.from;
      let transaction = edView.state.update({changes: [{from: edView.state.selection.main.from, to: edView.state.selection.main.to}, {from: point, insert: newText}]});
      edView.dispatch(transaction);
    }
  }

  function somethingSelected() {
    if(typeof edView !== 'undefined') {
      return edView.state.selection.ranges.some(r => !r.empty);
    }
  }

  function setCursor(pos) {
    if(typeof edView !== 'undefined') {
      currentCursor = pos;
      edView.dispatch({selection: {anchor: currentCursor}})
    }
  }

  function getCursor() {
    if(typeof edView !== 'undefined') {
      currentCursor = edView.state.selection.main.head;
      return(currentCursor);
    } else {
      return(0);
    }
  }

  function getValue() {
    if(typeof edView !== 'undefined') {
      return edView.state.doc.toString();
    }
  }

  function focus() {
    if(typeof edView !== 'undefined') {
      edView.focus();
    }
  }

  function getEdView() {
    return(edView);
  }

  function getEdState() {
    return(edState);
  }
</script>
