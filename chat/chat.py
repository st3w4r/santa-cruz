import json
import os
from urllib.parse import urljoin

import requests
from agents import function_tool
from openai import OpenAI
from openai.types.responses import (
    ResponseFunctionCallArgumentsDoneEvent,
    ResponseFunctionToolCall,
    ResponseTextDeltaEvent,
)

DB_API_BASE = "http://localhost:8080"

@function_tool
def fetch_databases():
    endpoint = "databases/"
    url = urljoin(DB_API_BASE, endpoint)
    response = requests.get(url)
    response.raise_for_status()
    return response.json()

@function_tool
def fetch_database_details(database_id: int):
    endpoint = f"databases/{database_id}/"
    url = urljoin(DB_API_BASE, endpoint)
    response = requests.get(url)
    response.raise_for_status()
    return response.json()

@function_tool
def run_query(database_id:int, query: str):
    endpoint = f"databases/{database_id}/run"
    url = urljoin(DB_API_BASE, endpoint)
    params = {'query': query}
    response = requests.get(url, params=params)
    response.raise_for_status()
    return response.json()


def main():
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        print("Please set your OPENAI_API_KEY environment variable.")
        return

    client = OpenAI(api_key=api_key)

    print("Welcome to the streaming chat conversation with conversation state!")
    print("Type 'quit' to exit.")

    previous_response_id = None
    tools = []
    for tool in [
        fetch_databases,
        fetch_database_details,
        run_query,
    ]:
        t = tool.params_json_schema
        t["name"] = tool.name
        t["type"] = "function"
        t["description"] = tool.description
        t["strict_json_schema"] = tool.strict_json_schema
        tools.append(t)

    while True:
        user_input = input("\nYou: ")
        if user_input.lower() in {"quit", "exit"}:
            break

        print("AI: ", end="", flush=True)

        try:
            stream = client.responses.create(
                model="gpt-4o-mini",
                input=user_input,
                stream=True,
                previous_response_id=previous_response_id,
                tools=tools,
            )

            for event in stream:
                print(event.type)
                if event.type == 'response.output_text.delta' and isinstance(event, ResponseTextDeltaEvent):
                    print(event.delta, end="", flush=True)
                elif event.type == 'response.function_call_arguments.done' and isinstance(event, ResponseFunctionCallArgumentsDoneEvent):
                    print(f"\nTool {event.item_id} called with arguments: {event.arguments}")
                elif event.type == "response.output_item.added":
                    print(f"\nOutput item added: {event.item}")
                    if isinstance(event.item, ResponseFunctionToolCall):
                        func_name = event.item.name
                        args = json.loads(event.item.arguments) if event.item.arguments else {}
                        print(f"Calling function {func_name} with arguments: {args}")
                        if func_name == "fetch_databases":
                            result = fetch_databases(**args)
                        elif func_name == "fetch_database_details":
                            result = fetch_database_details(**args)
                        elif func_name == "run_query":
                            result = run_query(**args)
                        else:
                            print(f"Unknown function: {func_name}")
                        print(result)
                elif event.type == 'response.completed':
                    previous_response_id = event.response.id

        except Exception as e:
            print(f"\nAn error occurred: {e}")

    print("\nConversation ended.")

if __name__ == "__main__":
    try:
        main()
    except EOFError:
        print("\nexit")
    except KeyboardInterrupt:
        print("\nexit")
