const functions = require("@google-cloud/functions-framework");

functions.http("date", (req, res) => {
    res.set("Access-Control-Allow-Origin", "*")
    res.set("Access-Control-Allow-Methods", "GET");

    if (req.method === "OPTIONS") {
        res.status(204).send("");
        return;
    }

    res.set("Content-Type", "application/json");
    res.send(
        JSON.stringify({
            now: Date.now(),
        }),
    );
});
