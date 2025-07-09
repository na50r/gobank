function h1Comp(text) {
    const h1 = document.createElement('h1');
    h1.innerText = text;
    return h1;
}

function colComp(text) {
    const col = document.createElement('td');
    col.innerText = text;
    return col;
}

// Spreading, equivalent to Python's **cols
function rowComp(cols) {
    const row = document.createElement('tr');
    row.append(...cols);
    return row;
}

function formComp(onsubmit) {
    const form = document.createElement('form');
    form.addEventListener('submit', onsubmit);
    return form;
}

function btnComp(text, onclick = function () { }, type = 'button', route=undefined) {
    const btn = document.createElement('button');
    btn.innerText = text;
    btn.onclick = onclick;
    btn.classList.add('btn');
    btn.type = type;
    return btn;
}

function inputComp(name, placeholder, type = 'text') {
    const input = document.createElement('input');
    input.name = name;
    input.placeholder = placeholder;
    input.required = true;
    input.type = type;
    return input;
}


export { h1Comp, colComp, rowComp, formComp, btnComp, inputComp };