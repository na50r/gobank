import AbstractView from "./AbstractView.js";
import { h1Comp, btnComp, formComp, inputComp, containerComp, btnBar } from "../components/Ui.js";
import { transfer } from "../util/Calls.js";
import { navigateTo } from "../index.js";

function renderTransfer() {
    const number = localStorage.getItem('number');
    const container = containerComp();
    const h1 = h1Comp('Make Transfer');
    container.append(h1);
    const form = formComp(transfer);
    const input1 = inputComp('to', 'Number');
    const input2 = inputComp('amount', 'Amount', 'number');
    const btn = btnComp('Transfer', () => { }, 'submit');
    const btn2 = btnComp('Account', () => {navigateTo(`/account/${number}`)});
    const bar = btnBar([btn, btn2]);
    form.append(input1, input2, bar);
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