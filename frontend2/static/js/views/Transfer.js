import AbstractView from "./AbstractView.js";
import { h1Comp, btnComp, formComp, inputComp } from "../components/UI.js";
import { transfer } from "../util/calls.js";

function renderTransfer() {
    const container = document.createElement('div');
    const h1 = h1Comp('Make Transfer');
    container.append(h1);
    const form = formComp(transfer);
    const input1 = inputComp('to', 'Number');
    const input2 = inputComp('amount', 'Amount', 'number');
    const btn = btnComp('Transfer', () => { }, 'submit');
    form.append(input1, input2, btn);
    container.append(form);
    return container;
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Transfer");
    }

    async getHtml() {
        return renderTransfer();
    }
}