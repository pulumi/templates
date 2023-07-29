from datetime import datetime
import json
import azure.functions as func


def main(req: func.HttpRequest) -> func.HttpResponse:

    headers = {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
    }

    if req.method == "OPTIONS":
        return func.HttpResponse("", headers=headers, status_code=204)

    now = datetime.now()
    now_in_ms = int(now.timestamp()) * 1000

    headers["Content-Type"] = "application/json"
    return func.HttpResponse(
        json.dumps({"now": now_in_ms}), headers=headers, status_code=200
    )
