from openai import OpenAI
from openai.types.responses import ResponseTextDeltaEvent

client = OpenAI()



def main():
    response = client.responses.create(
        model="gpt-4o-mini",
        input="What is the capital of the United States?",
    )
    print(response.output_text)

    second_response = client.responses.create(
        model="gpt-4o-mini",
        input="Is this true?",
        previous_response_id=response.id,
    )
    print(second_response.output_text)

    stream = client.responses.create(
        model="gpt-4o-mini",
        input="Tell me another fact",
        previous_response_id=second_response.id,
        stream=True,
    )

    for event in stream:
        if event.type == 'response.output_text.delta' and isinstance(event, ResponseTextDeltaEvent):
            print(event.delta, end="", flush=True)
        elif event.type == 'response.completed':
            print("\nResponse completed")


if __name__ == "__main__":
    main()
