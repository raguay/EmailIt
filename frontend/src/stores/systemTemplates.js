import { writable } from 'svelte/store';

export const systemTemplates = writable([{
    "name": "Today",
    "description": "Insert today's date.",
    "template": "Today is: {{cDateMDY}}.",
    "system": true
}]);

