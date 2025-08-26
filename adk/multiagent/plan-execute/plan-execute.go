package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt"

	"github.com/cloudwego/eino-examples/adk/common/prints"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute/agent"
)

func main() {
	ctx := context.Background()
	planAgent, err := agent.NewPlanner(ctx)
	if err != nil {
		log.Fatalf("agent.NewPlanner failed, err: %v", err)
	}

	executeAgent, err := agent.NewExecutor(ctx)
	if err != nil {
		log.Fatalf("agent.NewExecutor failed, err: %v", err)
	}

	replanAgent, err := agent.NewCritiqueAgent(ctx)
	if err != nil {
		log.Fatalf("agent.NewCritiqueAgent failed, err: %v", err)
	}

	entryAgent, err := prebuilt.NewPlanExecuteAgent(ctx, &prebuilt.PlanExecuteConfig{
		Planner:   planAgent,
		Executor:  executeAgent,
		Replanner: replanAgent,
	})
	if err != nil {
		log.Fatalf("NewPlanExecuteAgent failed, err: %v", err)
	}

	r := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: entryAgent,
	})

	iter := r.Query(ctx, "Plan a 3-day trip to Tokyo in March. I need flights from New York, hotel recommendations, and must-see attractions.")

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		prints.Event(event)
	}
}
