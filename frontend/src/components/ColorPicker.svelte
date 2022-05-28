{#if show}
    <div 
        id="colorPicker"
        style="background-color: {$theme.backgroundColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
    >
        <p>The color for {explainText}:</p>
        <canvas id="picker" bind:this="{canvas}" ></canvas><br>
        <input
            id="color"  
            bind:this="{input}"
            value='{color}'
            on:keyup="{(event) => { processKey(event);}}"
        />
        <button
            on:click="{(event) => {saveColor()}}" 
            style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
        >
            Select
        </button>
        <button 
            on:click="{(event) => {quitColorPicker();}}" 
            style="background-color: {$theme.textAreaColor}; font-family: {$theme.font}; color: {$theme.textColor}; font-size: {$theme.fontSize};"
        >
            Quit
        </button>
    </div>
{/if}

<style>
  #colorPicker button {
    border-radius: 5px;
    border-color: black;
    font-size: 15px;
    height: 30px;
    outline: none;
    margin: 0px 10px;
    padding: 3px 6px 6px 6px;
    user-select: none;
    outline-style:none;
    cursor: pointer;
  }


   #colorPicker {
        position: absolute;
        top: 20%;
        left: 20%;
        z-index: 100;
        padding: 20px;
        border: 3px solid;
        border-radius: 10px;
    }
</style>

<script>
    import { createEventDispatcher, afterUpdate } from 'svelte';
    import { KellyColorPicker } from './html5kellycolorpicker.js';
    import { theme } from '../stores/theme.js';
    export let color = "";
    export let id = 0;
    export let show = false;
    export let item = 'circle';
    export let explainText = '';
    let canvas;
    let input;
    let picker = {};
	const dispatch = createEventDispatcher();

    afterUpdate(() => {
        if(show) {
            picker = new KellyColorPicker({
                place: canvas, 
                input: input, 
                colorSaver: true,
                methodSwitch: true
            });
        }
    });

    function fire(name, data) {
        dispatch(name, {
            data: data
        });
    }

    function quitColorPicker() {    
        fire('quitColorPicker',{})
    }

    function saveColor() {
        document.body.style.cursor = 'default';
        fire('colorChanged',{color: input.value, id: id});
    }

    function processKey(event) {
        if(event.key === 'Enter') {
            event.preventDefault();
            event.stopPropagation();
            saveColor();
        }
    }
</script>
