import argparse
import asyncio
from dataclasses import dataclass
from datetime import datetime

from agents import Agent, Runner, function_tool


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
    result = await Runner.run(triageAgent, input=args.input)
    output = result.final_output
    print(output)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", type=str, default=None)
    args = parser.parse_args()

    asyncio.run(main(args))
