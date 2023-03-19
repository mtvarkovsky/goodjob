// This code is automatically generated. EDIT AT YOUR OWN RISK.

package quadraticequation

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
)

type QuadraticEquationJob struct {
	Visible bool
	LastTask interfaces.Task
	LastTaskPos int
	LastTaskResult *interfaces.TaskResult
	ID interfaces.JobID
	JobArgs []interfaces.JobArg
	Tasks []interfaces.Task
	TaskArgs map[interfaces.TaskID][]*interfaces.TaskArg
}

func NewQuadraticEquationJob(a float64, b float64, c float64) interfaces.Job {
	jobID := interfaces.JobID(fmt.Sprintf("quadratic equation (a * x^2) + (b * x) + c = 0 (%s)", ulid.Make()))
	jobArgs := []interfaces.JobArg{
		{
			Name: "a",
			Value: a,
		},
		{
			Name: "b",
			Value: b,
		},
		{
			Name: "c",
			Value: c,
		},
	}

	taskIDs := map[string]interfaces.TaskID{
		"calculate b^2": interfaces.TaskID(fmt.Sprintf("calculate b^2 (%s)", ulid.Make())),
		"calculate 4 * a * c": interfaces.TaskID(fmt.Sprintf("calculate 4 * a * c (%s)", ulid.Make())),
		"calculate b^2 - (4 * a * c)": interfaces.TaskID(fmt.Sprintf("calculate b^2 - (4 * a * c) (%s)", ulid.Make())),
		"calculate sqrt(b^2 - (4 * a * c))": interfaces.TaskID(fmt.Sprintf("calculate sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate -b": interfaces.TaskID(fmt.Sprintf("calculate -b (%s)", ulid.Make())),
		"calculate -b - sqrt(b^2 - (4 * a * c))": interfaces.TaskID(fmt.Sprintf("calculate -b - sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate -b + sqrt(b^2 - (4 * a * c))": interfaces.TaskID(fmt.Sprintf("calculate -b + sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate 2 * a": interfaces.TaskID(fmt.Sprintf("calculate 2 * a (%s)", ulid.Make())),
		"calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)": interfaces.TaskID(fmt.Sprintf("calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a) (%s)", ulid.Make())),
		"calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)": interfaces.TaskID(fmt.Sprintf("calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a) (%s)", ulid.Make())),
		"solution x1, x2": interfaces.TaskID(fmt.Sprintf("solution x1, x2 (%s)", ulid.Make())),
	}
	tasks := []interfaces.Task{
		MultiplyNumbersTask{
			ID: taskIDs["calculate b^2"],
			JobID: jobID,
		},
		MultiplyNumbersTask{
			ID: taskIDs["calculate 4 * a * c"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID: taskIDs["calculate b^2 - (4 * a * c)"],
			JobID: jobID,
		},
		SquareRootOfNumberTask{
			ID: taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID: taskIDs["calculate -b"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID: taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		AddNumbersTask{
			ID: taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		MultiplyNumbersTask{
			ID: taskIDs["calculate 2 * a"],
			JobID: jobID,
		},
		DivideNumbersTask{
			ID: taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			JobID: jobID,
		},
		DivideNumbersTask{
			ID: taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			JobID: jobID,
		},
		GetEquationSolutionTask{
			ID: taskIDs["solution x1, x2"],
			JobID: jobID,
		},
	}
	taskArgs := map[interfaces.TaskID][]*interfaces.TaskArg{
		taskIDs["calculate b^2"]: {
			{
				Name: "b",
				Value: b,
			},
			{
				Name: "b",
				Value: b,
			},
		},
		taskIDs["calculate 4 * a * c"]: {
			{
				Name: "4",
				Value: 4.0,
			},
			{
				Name: "a",
				Value: a,
			},
			{
				Name: "c",
				Value: c,
			},
		},
		taskIDs["calculate b^2 - (4 * a * c)"]: {
			{
				Name: "b^2",
				ValueFrom: taskIDs["calculate b^2"],
			},
			{
				Name: "4 * a * c",
				ValueFrom: taskIDs["calculate 4 * a * c"],
			},
		},
		taskIDs["calculate sqrt(b^2 - (4 * a * c))"]: {
			{
				Name: "sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate b^2 - (4 * a * c)"],
			},
		},
		taskIDs["calculate -b"]: {
			{
				Name: "0",
				Value: 0.0,
			},
			{
				Name: "b",
				Value: b,
			},
		},
		taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"]: {
			{
				Name: "-b",
				ValueFrom: taskIDs["calculate -b"],
			},
			{
				Name: "sqrt(b^2 - (4 * a * c)",
				ValueFrom: taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			},
		},
		taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"]: {
			{
				Name: "-b",
				ValueFrom: taskIDs["calculate -b"],
			},
			{
				Name: "sqrt(b^2 - (4 * a * c)",
				ValueFrom: taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			},
		},
		taskIDs["calculate 2 * a"]: {
			{
				Name: "2",
				Value: 2.0,
			},
			{
				Name: "a",
				Value: a,
			},
		},
		taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"]: {
			{
				Name: "-b - sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"],
			},
			{
				Name: "2 * a",
				ValueFrom: taskIDs["calculate 2 * a"],
			},
		},
		taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"]: {
			{
				Name: "-b + sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"],
			},
			{
				Name: "2 * a",
				ValueFrom: taskIDs["calculate 2 * a"],
			},
		},
		taskIDs["solution x1, x2"]: {
			{
				Name: "x1",
				ValueFrom: taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			},
			{
				Name: "x2",
				ValueFrom: taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			},
		},
	}
	return &QuadraticEquationJob{
		JobArgs: jobArgs,
		LastTaskPos: 0,
		Tasks: tasks,
		TaskArgs: taskArgs,
		Visible: true,
		LastTask: nil,
		LastTaskResult: nil,
		ID: jobID,
	}
}

func (j *QuadraticEquationJob) GetID() interfaces.JobID {
	return j.ID
}

func (j *QuadraticEquationJob) GetTasks() []interfaces.Task {
	return j.Tasks
}

func (j *QuadraticEquationJob) GetTaskArgs(taskID interfaces.TaskID) []*interfaces.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *QuadraticEquationJob) GetVisible() bool {
	return j.Visible
}

func (j *QuadraticEquationJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *QuadraticEquationJob) GetLastTask() interfaces.Task {
	return j.LastTask
}

func (j *QuadraticEquationJob) SetLastTask(task interfaces.Task) {
	j.LastTask = task
}

func (j *QuadraticEquationJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *QuadraticEquationJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *QuadraticEquationJob) GetLastTaskResult() *interfaces.TaskResult {
	return j.LastTaskResult
}

func (j *QuadraticEquationJob) SetLastTaskResult(result *interfaces.TaskResult) {
	j.LastTaskResult = result
}

