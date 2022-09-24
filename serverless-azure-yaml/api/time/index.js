module.exports = async function (context, req) {

    context.res = {
        body: JSON.stringify({
            now: Date.now(),
        }),
    };
};
