module.exports = async function (context, req) {

    const headers = {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
    };

    if (req.method === "OPTIONS") {
        context.res = {
            headers,
            body: "",
            status: 204,
        };

        return;
    }

    context.res = {
        headers: Object.assign(headers, { "Content-Type": "application/json" }),
        body: JSON.stringify({
            now: Date.now(),
        }),
    };
};
