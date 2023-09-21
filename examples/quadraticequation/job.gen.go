// This code is automatically generated. EDIT AT YOUR OWN RISK.

package quadraticequation

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/oklog/ulid/v2"
)

type QuadraticEquationJob struct {
	JobArgs        []goodjob.JobArg
	Tasks          []goodjob.Task
	TaskArgs       map[goodjob.TaskID][]*goodjob.TaskArg
	Visible        bool
	LastTask       goodjob.Task
	LastTaskPos    int
	LastTaskResult *goodjob.TaskResult
	ID             goodjob.JobID
}

func NewQuadraticEquationJob(a float64, b float64, c float64) goodjob.Job {
	jobID := goodjob.JobID(fmt.Sprintf("quadratic equation (a * x^2) + (b * x) + c = 0 (%s)", ulid.Make()))
	jobArgs := []goodjob.JobArg{
		{
			Name:  "a",
			Value: a,
		},
		{
			Name:  "b",
			Value: b,
		},
		{
			Name:  "c",
			Value: c,
		},
	}

	taskIDs := map[string]goodjob.TaskID{
		"calculate b^2":                                      goodjob.TaskID(fmt.Sprintf("calculate b^2 (%s)", ulid.Make())),
		"calculate 4 * a * c":                                goodjob.TaskID(fmt.Sprintf("calculate 4 * a * c (%s)", ulid.Make())),
		"calculate b^2 - (4 * a * c)":                        goodjob.TaskID(fmt.Sprintf("calculate b^2 - (4 * a * c) (%s)", ulid.Make())),
		"calculate sqrt(b^2 - (4 * a * c))":                  goodjob.TaskID(fmt.Sprintf("calculate sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate -b":                                       goodjob.TaskID(fmt.Sprintf("calculate -b (%s)", ulid.Make())),
		"calculate -b - sqrt(b^2 - (4 * a * c))":             goodjob.TaskID(fmt.Sprintf("calculate -b - sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate -b + sqrt(b^2 - (4 * a * c))":             goodjob.TaskID(fmt.Sprintf("calculate -b + sqrt(b^2 - (4 * a * c)) (%s)", ulid.Make())),
		"calculate 2 * a":                                    goodjob.TaskID(fmt.Sprintf("calculate 2 * a (%s)", ulid.Make())),
		"calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)": goodjob.TaskID(fmt.Sprintf("calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a) (%s)", ulid.Make())),
		"calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)": goodjob.TaskID(fmt.Sprintf("calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a) (%s)", ulid.Make())),
		"solution x1, x2":                                    goodjob.TaskID(fmt.Sprintf("solution x1, x2 (%s)", ulid.Make())),
	}
	tasks := []goodjob.Task{
		MultiplyNumbersTask{
			ID:    taskIDs["calculate b^2"],
			JobID: jobID,
		},
		MultiplyNumbersTask{
			ID:    taskIDs["calculate 4 * a * c"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID:    taskIDs["calculate b^2 - (4 * a * c)"],
			JobID: jobID,
		},
		SquareRootOfNumberTask{
			ID:    taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID:    taskIDs["calculate -b"],
			JobID: jobID,
		},
		SubtractNumbersTask{
			ID:    taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		AddNumbersTask{
			ID:    taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"],
			JobID: jobID,
		},
		MultiplyNumbersTask{
			ID:    taskIDs["calculate 2 * a"],
			JobID: jobID,
		},
		DivideNumbersTask{
			ID:    taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			JobID: jobID,
		},
		DivideNumbersTask{
			ID:    taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			JobID: jobID,
		},
		GetEquationSolutionTask{
			ID:    taskIDs["solution x1, x2"],
			JobID: jobID,
		},
	}
	taskArgs := map[goodjob.TaskID][]*goodjob.TaskArg{
		taskIDs["calculate b^2"]: {
			{
				Name:  "b",
				Value: b,
			},
			{
				Name:  "b",
				Value: b,
			},
		},
		taskIDs["calculate 4 * a * c"]: {
			{
				Name:  "4",
				Value: 4.0,
			},
			{
				Name:  "a",
				Value: a,
			},
			{
				Name:  "c",
				Value: c,
			},
		},
		taskIDs["calculate b^2 - (4 * a * c)"]: {
			{
				Name:      "b^2",
				ValueFrom: taskIDs["calculate b^2"],
			},
			{
				Name:      "4 * a * c",
				ValueFrom: taskIDs["calculate 4 * a * c"],
			},
		},
		taskIDs["calculate sqrt(b^2 - (4 * a * c))"]: {
			{
				Name:      "sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate b^2 - (4 * a * c)"],
			},
		},
		taskIDs["calculate -b"]: {
			{
				Name:  "0",
				Value: 0.0,
			},
			{
				Name:  "b",
				Value: b,
			},
		},
		taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"]: {
			{
				Name:      "-b",
				ValueFrom: taskIDs["calculate -b"],
			},
			{
				Name:      "sqrt(b^2 - (4 * a * c)",
				ValueFrom: taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			},
		},
		taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"]: {
			{
				Name:      "-b",
				ValueFrom: taskIDs["calculate -b"],
			},
			{
				Name:      "sqrt(b^2 - (4 * a * c)",
				ValueFrom: taskIDs["calculate sqrt(b^2 - (4 * a * c))"],
			},
		},
		taskIDs["calculate 2 * a"]: {
			{
				Name:  "2",
				Value: 2.0,
			},
			{
				Name:  "a",
				Value: a,
			},
		},
		taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"]: {
			{
				Name:      "-b - sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate -b - sqrt(b^2 - (4 * a * c))"],
			},
			{
				Name:      "2 * a",
				ValueFrom: taskIDs["calculate 2 * a"],
			},
		},
		taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"]: {
			{
				Name:      "-b + sqrt(b^2 - (4 * a * c))",
				ValueFrom: taskIDs["calculate -b + sqrt(b^2 - (4 * a * c))"],
			},
			{
				Name:      "2 * a",
				ValueFrom: taskIDs["calculate 2 * a"],
			},
		},
		taskIDs["solution x1, x2"]: {
			{
				Name:      "x1",
				ValueFrom: taskIDs["calculate (-b + sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			},
			{
				Name:      "x2",
				ValueFrom: taskIDs["calculate (-b - sqrt(b^2 - (4 * a * c))) / (2 * a)"],
			},
		},
	}
	return &QuadraticEquationJob{
		JobArgs:        jobArgs,
		LastTaskPos:    0,
		Tasks:          tasks,
		TaskArgs:       taskArgs,
		Visible:        true,
		LastTask:       nil,
		LastTaskResult: nil,
		ID:             jobID,
	}
}

func (j *QuadraticEquationJob) GetID() goodjob.JobID {
	return j.ID
}

func (j *QuadraticEquationJob) GetTasks() []goodjob.Task {
	return j.Tasks
}

func (j *QuadraticEquationJob) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *QuadraticEquationJob) GetVisible() bool {
	return j.Visible
}

func (j *QuadraticEquationJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *QuadraticEquationJob) GetLastTask() goodjob.Task {
	return j.LastTask
}

func (j *QuadraticEquationJob) SetLastTask(task goodjob.Task) {
	j.LastTask = task
}

func (j *QuadraticEquationJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *QuadraticEquationJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *QuadraticEquationJob) GetLastTaskResult() *goodjob.TaskResult {
	return j.LastTaskResult
}

func (j *QuadraticEquationJob) SetLastTaskResult(result *goodjob.TaskResult) {
	j.LastTaskResult = result
}
