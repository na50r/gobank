import AbstractView from "./AbstractView.js";
import { getElement } from "../util/Calls.js";
//Based on: https://www.youtube.com/watch?v=Ot5FQobG33A&lc=UgyRK4esogM-Vqi3Ofp4AaABAg
// Prepare Game HTML
function createGame() {
    const container = document.createElement('div');
    container.id = 'gamebox';
    container.classList.add('gamebox');
    const section = document.createElement('section');
    container.appendChild(section);
    const aside = document.createElement('aside');
    container.appendChild(aside);
    return container;
}

// Creates an Element in HTML
function CreateElem(emoji, name) {
    const elem = document.createElement("div");
    elem.setAttribute("data-emoji", emoji);
    elem.setAttribute("data-elem", name);
    elem.classList.add("element");
    const span1 = document.createElement("span");
    const span2 = document.createElement("span");
    span1.innerText = emoji;
    span2.innerText = name;
    twemoji.parse(span1);
    elem.append(span1, span2);
    return elem;
}

async function Merge(elemA, elemB) {
    console.log(elemA.getAttribute("data-elem"), elemB.getAttribute("data-elem"));
    const res = await getElement(elemA.getAttribute("data-elem"), elemB.getAttribute("data-elem"));
    const capitalized = res.charAt(0).toUpperCase() + res.slice(1);
    console.log(elemA.getAttribute("data-emoji"), elemB.getAttribute("data-emoji"));
    return { emoji: '⚫', name: capitalized };
}

// Creates bar, must be rendered constantly
function renderBar(elem1, elem2, out) {
    const plus = document.createElement("span");
    plus.textContent = "+";
    plus.className = "symbol";

    const equals = document.createElement("span");
    equals.textContent = "=";
    equals.className = "symbol";

    const bar = document.createElement("div");
    bar.id = "bar";
    const div = document.createElement("div");
    bar.appendChild(elem1 ? elem1 : div);
    bar.appendChild(plus);
    bar.appendChild(elem2 ? elem2 : div);
    bar.appendChild(equals);
    bar.appendChild(out ? out : div);
    return bar;
}

// Renders Element list
function renderElems(state, elems) {
    state.aside.innerHTML = "";
    elems.forEach(element => {
        element = CreateElem(element.emoji, element.name);
        state.aside.appendChild(element);
    });
}

// Renders Game
async function renderGame(game) {
    const section = game.querySelector('section');
    const aside = game.querySelector('aside');
    const startList = [
        { emoji: '🔥', name: 'Fire' },
        { emoji: '💧', name: 'Water' },
        { emoji: '🌪️', name: 'Wind' },
        { emoji: '🌎', name: 'Earth' },
    ]
    const state = {
        section: section,
        aside: aside,
        selected: [],
        elems: startList,

    };
    renderElems(state, startList);
    const bar = renderBar(...state.selected);
    state.section.append(bar);
    document.addEventListener("click", async (e) => {
        const copy = e.target.closest('.element');
        if (copy) {
            state.selected.push(copy.cloneNode(true));
        }
        if (state.selected.length === 2) {
            const out = await Merge(state.selected[0], state.selected[1]);
            const merged = CreateElem(out.emoji, out.name);
            if (!state.elems.find(elem => elem.name === out.name)) {
                state.elems.push(out);
                renderElems(state, state.elems);
            }
            state.selected.push(merged);

        }
        console.log(state.selected);
        if (state.selected.length > 3) {
            state.selected = [];
            state.selected.push(copy.cloneNode(true));
        }
        state.section.replaceChild(renderBar(...state.selected), state.section.querySelector('#bar'));
    });
    return game;
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Game");
    }

    async getHtml() {
        const game = createGame();
        return renderGame(game);
    }
}