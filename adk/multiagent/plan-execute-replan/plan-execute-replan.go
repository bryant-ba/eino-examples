package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt"
	"github.com/cloudwego/eino/components/tool"

	"github.com/cloudwego/eino-examples/adk/common/prints"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute-replan/agent"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute-replan/store"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute-replan/tools"
	"github.com/cloudwego/eino-examples/adk/multiagent/plan-execute-replan/trace"
)

func main() {
	ctx := context.Background()

	traceCloseFn, client := trace.AppendCozeLoopCallbackIfConfigured(ctx)
	defer traceCloseFn(ctx)

	planAgent, err := agent.NewPlanner(ctx)
	if err != nil {
		log.Fatalf("agent.NewPlanner failed, err: %v", err)
	}

	executeAgent, err := agent.NewExecutor(ctx)
	if err != nil {
		log.Fatalf("agent.NewExecutor failed, err: %v", err)
	}

	replanAgent, err := agent.NewReplanAgent(ctx)
	if err != nil {
		log.Fatalf("agent.NewReplanAgent failed, err: %v", err)
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
		Agent:           entryAgent,
		CheckPointStore: store.NewInMemoryStore(),
	})

	query := `Plan a 3-day trip to Beijing in Next Month. I need flights from New York, hotel recommendations, and must-see attractions.
Today is 2025-09-09.`
	ctx, finishFn := trace.StartRootSpan(client, ctx, query)
	checkPointID := "per-id-123456"
	isToResume := false
	var iter *adk.AsyncIterator[*adk.AgentEvent]
	for {
		if !isToResume {
			iter = r.Query(ctx, query, adk.WithCheckPointID(checkPointID))
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("\nyour input here: ")
			scanner.Scan()
			fmt.Println()
			nInput := scanner.Text()
			var err_ error
			iter, err_ = r.Resume(ctx, checkPointID, adk.WithToolOptions([]tool.Option{tools.WithNewInput(nInput)}))
			if err_ != nil {
				log.Fatalf("agent Resume failed, err: %v", err_)
			}
		}
		isFinished := false
		for {
			event, ok := iter.Next()
			if !ok {
				isFinished = true
				break
			}

			prints.Event(event)

			if event.Action != nil && event.Action.Interrupted != nil {
				isToResume = true
				break
			}
		}

		if isFinished {
			break
		}

	}

	finishFn(ctx, "end")
}
