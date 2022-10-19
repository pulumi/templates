"use strict";

const express = require("express");
const app = express();

app.get("/", (req, res) => {
    res.json({ message: "Hello, world! 👋" });
});

const port = process.env.PORT;
app.listen(port, () => {
    console.log(`Listening on port ${port}`);
});
