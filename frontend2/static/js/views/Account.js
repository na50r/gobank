import AbstractView from "./AbstractView.js";
import { h1Comp, colComp, rowComp, btnComp } from "../components/UI.js";
import { getAccount, logout } from "../util/calls.js";


function renderAccount(account = {}) {
    const container = document.createElement('div');
    const h1 = h1Comp('Account Details');
    container.append(h1);
    const table = document.createElement('table');
    const info = [
        { name: 'Number', value: account.number },
        { name: 'First Name', value: account.first_name },
        { name: 'Last Name', value: account.last_name },
        { name: 'Balance', value: account.balance }
    ];
    info.forEach((item) => {
        const col1 = colComp(item.name);
        const col2 = colComp(item.value ?? 'Loading...');
        const row = rowComp([col1, col2]);
        table.append(row);
    });
    const btn2 = btnComp('Logout', logout);
    container.append(btn2);
    container.append(table);
    return container;
}


export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Account");
    }

    async getHtml() {
        const account = await getAccount();
        return renderAccount(account);
    }
}