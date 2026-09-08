export default {
    async fetch(request, env) {
        const count = parseInt(await env.COUNTER.get("count")) + 1 || 1;
        await env.COUNTER.put("count", count.toString());
        return new Response("Hello, world! You are visitor number " + count + ".");
    },
};
