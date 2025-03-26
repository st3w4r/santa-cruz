import argparse
import asyncio
from dataclasses import dataclass
from datetime import datetime

from agents import Agent, ItemHelpers, Runner, function_tool
from openai.types.responses import ResponseTextDeltaEvent


@dataclass
class UserContext:
    city: str = None
    name: str = None
    age: int = None

    def get_time(self):
        return datetime.now().strftime("%H:%M:%S") 


@function_tool
def get_weather(city: str) -> str:
    return f"The weather in {city} is 42 degrees."


codeAgent = Agent(
    name="Code",
    handoff_description="Code",
    instructions="Help write code!",
)

reviewerAgent = Agent(
    name="Reviewer",
    handoff_description="Code Review",
    instructions="Review code!",
)

pmAgent = Agent(
    name="PM",
    handoff_description="Project Manager",
    instructions="Manage projects!",
)

weatherAgent = Agent[UserContext](
    name="Weather",
    handoff_description="Weather",
    instructions="Get the weather!",
    tools=[get_weather],
)

triageAgent = Agent(
    name="Triager",
    instructions="Call the correct agent!",
    handoffs=[codeAgent, reviewerAgent, pmAgent, weatherAgent],
)

async def main(args):
    result = Runner.run_streamed(triageAgent, input=args.input)
    async for event in result.stream_events():
        if event.type == "raw_response_event" and isinstance(event.data, ResponseTextDeltaEvent):
            # print(event.data.delta, end="", flush=True)
            continue
        elif event.type == "agent_updated_stream_event":
            print(f"Agent updated: {event.new_agent.name}")
            continue
        elif event.type == "run_item_stream_event":
            if event.item.type == "tool_call_item":
                print("-- Tool was called")
            elif event.item.type == "tool_call_output_item":
                print(f"-- Tool output: {event.item.output}")
            elif event.item.type == "message_output_item":
                print(f"-- Message output:\n {ItemHelpers.text_message_output(event.item)}")
            else:
                pass

    # output = result.final_output
    # print(output)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", type=str, default=None)
    args = parser.parse_args()

    asyncio.run(main(args))
