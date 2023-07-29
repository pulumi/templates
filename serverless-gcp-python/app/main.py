from datetime import datetime
from flask import jsonify
import functions_framework


@functions_framework.http
def data(request):

    headers = {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
    }

    if request.method == "OPTIONS":
        return "", 204, headers

    now = datetime.now()
    now_in_ms = int(now.timestamp()) * 1000

    headers["Content-Type"] = "application/json"
    return jsonify({"now": now_in_ms}), 200, headers
