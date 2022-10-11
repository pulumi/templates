"use strict";

const express = require("express");
const app = express();

app.get("/", (req, res) => {
    res.send("Hello, world! 👋");
});

const host = "0.0.0.0";
const port = process.env.PORT;

app.listen(port, host, () => {
    console.log(`Running on http://${host}:${port}`);
});
