import os

from openai import OpenAI
from openai.types.responses import ResponseTextDeltaEvent


def main():
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        print("Please set your OPENAI_API_KEY environment variable.")
        return

    client = OpenAI(api_key=api_key)

    print("Welcome to the streaming chat conversation with conversation state!")
    print("Type 'quit' to exit.")

    previous_response_id = None

    while True:
        user_input = input("\nYou: ")
        if user_input.lower() in {"quit", "exit"}:
            break

        print("AI: ", end="", flush=True)

        try:
            params = {
                "model": "gpt-4o-mini",
                "input": user_input,
                "stream": True,
                "store": True,
                "previous_response_id": previous_response_id,
            }

            stream = client.responses.create(**params)

            for event in stream:
                if event.type == 'response.output_text.delta' and isinstance(event, ResponseTextDeltaEvent):
                    print(event.delta, end="", flush=True)
                elif event.type == 'response.completed':
                    previous_response_id = event.response.id

        except Exception as e:
            print(f"\nAn error occurred: {e}")

    print("\nConversation ended.")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nexit")
