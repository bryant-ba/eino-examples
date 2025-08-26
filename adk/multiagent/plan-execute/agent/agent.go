package agent

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt"
	"github.com/cloudwego/eino/compose"

	"github.com/cloudwego/eino-examples/adk/common/model"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute/tools"
)

func NewPlanner(ctx context.Context) (adk.Agent, error) {

	return prebuilt.NewPlanner(ctx, &prebuilt.PlannerConfig{
		ToolCallingChatModel: model.NewChatModel(),
		Instruction: `You are a travel planning expert. Create a detailed step-by-step plan for travel requests.
Break down complex travel planning into specific, actionable tasks like:
1. Check weather conditions
2. Search for flights
3. Find accommodations
4. Research attractions
5. Create itinerary
Each step should be clear and executable. Do not skip essential planning steps.`,
	})
}

func NewExecutor(ctx context.Context) (adk.Agent, error) {
	// Get travel tools for the executor
	travelTools, err := tools.GetAllTravelTools(ctx)
	if err != nil {
		return nil, err
	}

	return prebuilt.NewExecutor(ctx, &prebuilt.ExecutorConfig{
		Model: model.NewChatModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: travelTools,
			},
		},

		Instruction: `You are a travel research executor. Execute each planning step by using available tools.
For weather queries, use get_weather tool.
For flight searches, use search_flights tool.
For hotel searches, use search_hotels tool.
For attraction research, use search_attractions tool.
Provide detailed results for each task.`,
	})
}

func NewCritiqueAgent(ctx context.Context) (adk.Agent, error) {
	return prebuilt.NewReplanner(ctx, &prebuilt.ReplannerConfig{
		ChatModel: model.NewChatModel(),
		Instruction: `You are a travel planning critic and replanner. Review the travel planning progress and adjust the plan.

Original travel request: {UserInput}
Original plan: {Plan}
Completed steps: {ExecuteResults}

Analyze what has been accomplished and what still needs to be done.
If the travel plan is complete with all necessary information (weather, flights, hotels, attractions), provide a final summary.
Otherwise, identify missing steps and create an updated plan with only the remaining tasks.
Focus on practical travel planning needs.`,
	})
}
