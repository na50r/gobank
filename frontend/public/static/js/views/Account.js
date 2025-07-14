import AbstractView from "./AbstractView.js";
import { h1Comp, colComp, rowComp, btnComp, containerComp, btnBar } from "../components/Ui.js";
import { getAccount, logout, getImage } from "../util/Calls.js";
import { navigateTo } from "../index.js";
import { accountAccess, cacheAccount, loadAccount } from "../util/Helpers.js";

function renderAccount(account = {}, img = new Image()) {
    const container = containerComp();
    const h1 = h1Comp('Account Details');
    container.append(h1);
    container.append(img);
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
    const btn1 = btnComp("Transfer", () => { navigateTo('/transfer') });
    const btn2 = btnComp('Logout', logout);
    const bar = btnBar([btn1, btn2]);
    container.append(table, bar);
    return container;
}


export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Account");
    }

    async getHtml() {
        if (!accountAccess(this.params.id)) {
            return;
        }
        if (!loadAccount()) {
            const account = await getAccount(this.params.id);
            cacheAccount(account);
        }
        const account = loadAccount();
        const img = await getImage(this.params.id);
        return renderAccount(account, img);
    }
}