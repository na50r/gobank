import AbstractView from "./AbstractView.js";
//Based on: https://www.youtube.com/watch?v=Ot5FQobG33A&lc=UgyRK4esogM-Vqi3Ofp4AaABAg

function createGame() {
    const game = document.createElement('div');
    game.id = 'scoped-app';
    const container = document.createElement('div');
    container.classList.add('container');
    game.appendChild(container);
    const section = document.createElement('section');
    container.appendChild(section);
    const aside = document.createElement('aside');
    container.appendChild(aside);
    return game;
}

function AddElem(state, emoji, name) {
    let elem = CreateElem(emoji, name);
    elem.addEventListener("mousedown", (e) => {
        AddDraggable(state, emoji, name);

    });
    state.aside.appendChild(elem);
}

function AddDraggable(state, emoji, name) {
    let newElem = CreateElem(emoji, name);
    state.section.appendChild(newElem);
    newElem.classList.add("movable");
    state.isDragged = newElem;
    newElem.classList.add("dragged");

    newElem.addEventListener("mousedown", (e) => {
        if (newElem.classList.contains("loading")) {
            return;
        }

        state.isDragged = newElem;
        newElem.classList.add("dragged");
    });

    newElem.addEventListener("contextmenu", (e) => {
        if (e.ctrlKey) return;
        e.preventDefault();
        newElem.remove();
    });

    newElem.addEventListener("mouseenter", (e) => {
        if (newElem.classList.contains("loading")) {
            return;
        }
        state.isHovered = newElem;
    });

    newElem.addEventListener("mouseleave", (e) => {
        if (newElem.classList.contains("loading")) {
            return;
        }
        state.isHovered = null;
    });
    return newElem;
}

function CreateElem(emoji, name) {
    const elem = document.createElement("div");
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

function Merge(elemA, elemB) {
    return { emoji: '⭐', name: 'Star' };
}

function renderGame(game) {
    const section = game.querySelector('section');
    const aside = game.querySelector('aside');
    const state = {
        section: section,
        aside: aside,
        isDragged: null,
        isHovered: null,
    };
    AddElem(state, '😊', 'Smile')
    AddElem(state, '🏀', 'Basketball')
    document.addEventListener("mouseup", async (e) => {
        console.log(state.isDragged, state.isHovered);

        if (state.isDragged && state.isHovered) {
            let elemA = state.isDragged.getAttribute("data-elem");
            let elemB = state.isHovered.getAttribute("data-elem");

            state.isDragged.classList.add("loading")
            state.isHovered.classList.add("loading")
            let copyA = state.isDragged
            let copyB = state.isHovered
            const res = Merge(elemA, elemB);
            if (res) {
                const newElem = AddDraggable(state, res.emoji, res.name);
                section.appendChild(newElem);
                newElem.style.top = copyA.style.top;
                newElem.style.left = copyA.style.left;
                copyA.remove();
                copyB.remove();
            } else {
                state.isDragged.classList.remove("loading")
                state.isHovered.classList.remove("loading")
            }

            if (state.isDragged) {
                state.isDragged.classList.remove("dragged");
            }
            state.isDragged = null;
            state.isHovered = null;
        }
        if (state.isDragged) {
            state.isDragged.classList.remove("dragged");
        }
        state.isDragged = null;
    });

    document.addEventListener("mousemove", (e) => {
        if (!state.isDragged) return;
        state.isDragged.style.top = e.clientY - state.isDragged.clientHeight / 2 + "px";
        state.isDragged.style.left = e.clientX - state.isDragged.clientWidth / 2 + "px";
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