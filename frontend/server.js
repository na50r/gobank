const express = require('express');
const path = require('path');
const app = express();
const port = process.env.PORT || 3030;

app.use(express.static(path.resolve(__dirname, "public")));
app.get('/{*any}', (req, res) => {
    res.sendFile(path.resolve(__dirname, "public", "index.html"));
});


app.listen(port, () => console.log("Server running..."));
