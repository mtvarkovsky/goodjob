package main

import (
	"fmt"
	"go/format"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type (
	generator struct {
		builder     strings.Builder
		pkg         string
		destination string
		jobSpec     *jobSpec
	}

	jobSpec struct {
		Imports   []string `yaml:"imports"`
		Name      string   `yaml:"name"`
		ID        string   `yaml:"jobID"`
		Type      string   `yaml:"type"`
		Arguments []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		}
		Tasks []struct {
			Task      string `yaml:"task"`
			ID        string `yaml:"taskID"`
			Arguments []struct {
				Name             string  `yaml:"name"`
				ValueFromLiteral *string `yaml:"valueFromLiteral"`
				ValueFromJob     *string `yaml:"valueFromJob"`
				ValueFromTask    *string `yaml:"valueFromTask"`
			}
		}
	}
)

const (
	RetryableJob           = "RetryableJob"
	RevertibleJob          = "RevertibleJob"
	RetryableRevertibleJob = "RetryableRevertibleJob"
)

var (
	defaultImports = []string{
		"fmt",
		"github.com/oklog/ulid/v2",
		"github.com/mtvarkovsky/goodjob/pkg/goodjob",
	}

	jobStructFieldToType = map[string]string{
		"JobID":                  "goodjob.JobID",
		"TaskID":                 "goodjob.TaskID",
		"Tasks":                  "[]goodjob.Task",
		"TaskArgs":               "map[goodjob.TaskID]*goodjob.TaskArg",
		"JobArgs":                "[]goodjob.JobArg",
		"LastTask":               "goodjob.Task",
		"LastTaskPos":            "int",
		"LastTaskResult":         "*goodjob.TaskResult",
		"Visible":                "bool",
		"Job":                    "goodjob.Job",
		"RetryableJob":           "goodjob.RetryableJob",
		"RevertibleJob":          "goodjob.RevertibleJob",
		"RetryableRevertibleJob": "goodjob.RetryableRevertibleJob",
	}

	jobStructFields = map[string]string{
		"ID":             "goodjob.JobID",
		"JobArgs":        "[]goodjob.JobArg",
		"Tasks":          "[]goodjob.Task",
		"TaskArgs":       "map[goodjob.TaskID][]*goodjob.TaskArg",
		"Visible":        "bool",
		"LastTask":       "goodjob.Task",
		"LastTaskPos":    "int",
		"LastTaskResult": "*goodjob.TaskResult",
	}

	retryableJobStructFields = map[string]string{
		"RetryThreshold":      "int",
		"RetryThresholdCount": "map[goodjob.TaskID]int",
	}

	revertibleJobStructFields = map[string]string{
		"RevertState": "bool",
	}

	retryableRevertibleJobStructFields = map[string]string{
		"ForwardRetryThreshold":       "int",
		"ForwardRetryThresholdCount":  "map[goodjob.TaskID]int",
		"BackwardRetryThreshold":      "int",
		"BackwardRetryThresholdCount": "map[goodjob.TaskID]int",
		"RevertState":                 "bool",
	}

	allowedJobTypes = map[string]bool{
		"Job":                    true,
		"RetryableJob":           true,
		"RevertibleJob":          true,
		"RetryableRevertibleJob": true,
	}
)

func main() {
	if len(os.Args) != 4 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: generateast {specification} {package name} {destination folder} {file name}")
		os.Exit(65)
		return
	}
	specPath := os.Args[1]
	pkg := os.Args[2]
	destination := os.Args[3]
	f, err := os.Create(destination)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}
	defer f.Close()

	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}

	spec := jobSpec{}
	err = yaml.Unmarshal(specBytes, &spec)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}

	gen := newGenerator(pkg, destination, &spec)
	err = gen.generateJob()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}

	formattedCode, err := format.Source([]byte(gen.builder.String()))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}
	_, err = fmt.Fprint(f, string(formattedCode))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
		return
	}
}

func newGenerator(pkg string, destination string, jobSpec *jobSpec) *generator {
	jobSpec.Imports = append(jobSpec.Imports, defaultImports...)
	return &generator{
		pkg:         pkg,
		destination: destination,
		builder:     strings.Builder{},
		jobSpec:     jobSpec,
	}
}

func (g *generator) generateJob() error {
	// check job type
	if _, found := allowedJobTypes[g.jobSpec.Type]; !found {
		return fmt.Errorf("unknown job type %s", g.jobSpec.Type)
	}

	// add header
	g.builder.WriteString("// This code is automatically generated. EDIT AT YOUR OWN RISK.\n\n")

	// define package
	g.builder.WriteString("package ")
	g.builder.WriteString(g.pkg)
	g.builder.WriteString("\n\n")

	// add imports
	g.builder.WriteString("import (\n")
	for _, imprt := range g.jobSpec.Imports {
		g.builder.WriteString(fmt.Sprintf("\t\"%s\"\n", imprt))
	}
	g.builder.WriteString(")\n\n")

	// define job struct
	g.builder.WriteString(fmt.Sprintf("type %s struct {\n", g.jobSpec.Name))
	for fName, fType := range jobStructFields {
		g.builder.WriteString(fmt.Sprintf("\t%s %s\n", fName, fType))
	}
	if g.jobSpec.Type == RetryableJob {
		for fName, fType := range retryableJobStructFields {
			g.builder.WriteString(fmt.Sprintf("\t%s %s\n", fName, fType))
		}
	}
	if g.jobSpec.Type == RevertibleJob {
		for fName, fType := range revertibleJobStructFields {
			g.builder.WriteString(fmt.Sprintf("\t%s %s\n", fName, fType))
		}
	}
	if g.jobSpec.Type == RetryableRevertibleJob {
		for fName, fType := range retryableRevertibleJobStructFields {
			g.builder.WriteString(fmt.Sprintf("\t%s %s\n", fName, fType))
		}
	}
	g.builder.WriteString("}\n\n")

	// define job constructor
	g.builder.WriteString(fmt.Sprintf("func New%s(", g.jobSpec.Name))
	if g.jobSpec.Type == RetryableJob {
		g.jobSpec.Arguments = append(
			g.jobSpec.Arguments,
			struct {
				Name string `yaml:"name"`
				Type string `yaml:"type"`
			}{
				Name: "retryThreshold",
				Type: "int",
			},
		)
	}
	if g.jobSpec.Type == RetryableRevertibleJob {
		g.jobSpec.Arguments = append(
			g.jobSpec.Arguments,
			[]struct {
				Name string `yaml:"name"`
				Type string `yaml:"type"`
			}{
				{
					Name: "forwardRetryThreshold",
					Type: "int",
				},
				{
					Name: "backwardRetryThreshold",
					Type: "int",
				},
			}...,
		)
	}
	for i, arg := range g.jobSpec.Arguments {
		g.builder.WriteString(fmt.Sprintf("%s %s", arg.Name, arg.Type))
		if i != len(g.jobSpec.Arguments)-1 {
			g.builder.WriteString(", ")
		}
	}
	g.builder.WriteString(fmt.Sprintf(") %s {\n", jobStructFieldToType[g.jobSpec.Type]))
	// prepare job ID
	g.builder.WriteString(fmt.Sprintf("\tjobID := goodjob.JobID(fmt.Sprintf(\"%s", g.jobSpec.ID))
	g.builder.WriteString(" (%s)\", ulid.Make()))\n")
	// prepare job args
	g.builder.WriteString("\tjobArgs := []goodjob.JobArg{\n")
	for _, jobArg := range g.jobSpec.Arguments {
		g.builder.WriteString("\t\t{\n")
		g.builder.WriteString(fmt.Sprintf("\t\t\tName: \"%s\",\n", jobArg.Name))
		g.builder.WriteString(fmt.Sprintf("\t\t\tValue: %s,\n", jobArg.Name))
		g.builder.WriteString("\t\t},\n")
	}
	g.builder.WriteString("\t}\n\n")
	// prepare task IDs
	g.builder.WriteString("\ttaskIDs := map[string]goodjob.TaskID{\n")
	for _, task := range g.jobSpec.Tasks {
		g.builder.WriteString(fmt.Sprintf("\t\t\"%s\": goodjob.TaskID(fmt.Sprintf(\"%s", task.ID, task.ID))
		g.builder.WriteString(" (%s)\", ulid.Make())),\n")
	}
	g.builder.WriteString("\t}\n")

	// prepare tasks
	g.builder.WriteString("\ttasks := []goodjob.Task{\n")
	for _, task := range g.jobSpec.Tasks {
		g.builder.WriteString(fmt.Sprintf("\t\t%s{\n", task.Task))
		g.builder.WriteString(fmt.Sprintf("\t\t\tID: taskIDs[\"%s\"],\n", task.ID))
		g.builder.WriteString(fmt.Sprintf("\t\t\tJobID: jobID,\n"))
		g.builder.WriteString("\t\t},\n")
	}
	g.builder.WriteString("\t}\n")

	// prepare task args
	g.builder.WriteString("\ttaskArgs := map[goodjob.TaskID][]*goodjob.TaskArg{\n")
	for _, task := range g.jobSpec.Tasks {
		g.builder.WriteString(fmt.Sprintf("\t\ttaskIDs[\"%s\"]: {\n", task.ID))
		for _, arg := range task.Arguments {
			g.builder.WriteString("\t\t\t{\n")
			g.builder.WriteString(fmt.Sprintf("\t\t\t\tName: \"%s\",\n", arg.Name))
			if arg.ValueFromLiteral != nil {
				g.builder.WriteString(fmt.Sprintf("\t\t\t\tValue: %s,\n", *arg.ValueFromLiteral))
			}
			if arg.ValueFromJob != nil {
				g.builder.WriteString(fmt.Sprintf("\t\t\t\tValue: %s,\n", *arg.ValueFromJob))
			}
			if arg.ValueFromTask != nil {
				g.builder.WriteString(fmt.Sprintf("\t\t\t\tValueFrom: taskIDs[\"%s\"],\n", *arg.ValueFromTask))
			}
			g.builder.WriteString("\t\t\t},\n")
		}
		g.builder.WriteString("\t\t},\n")
	}
	g.builder.WriteString("\t}\n")

	// return job struct
	g.builder.WriteString(fmt.Sprintf("\treturn &%s{\n", g.jobSpec.Name))
	g.builder.WriteString("\t\tJobArgs: jobArgs,\n")
	g.builder.WriteString("\t\tLastTaskPos: 0,\n")
	g.builder.WriteString("\t\tTasks: tasks,\n")
	g.builder.WriteString("\t\tTaskArgs: taskArgs,\n")
	g.builder.WriteString("\t\tVisible: true,\n")
	g.builder.WriteString("\t\tLastTask: nil,\n")
	g.builder.WriteString("\t\tLastTaskResult: nil,\n")
	g.builder.WriteString("\t\tID: jobID,\n")

	if g.jobSpec.Type == RetryableJob {
		g.builder.WriteString("\t\tRetryThreshold: retryThreshold,\n")
		g.builder.WriteString("\t\tRetryThresholdCount: make(map[goodjob.TaskID]int),\n")
	}

	if g.jobSpec.Type == RevertibleJob {
		g.builder.WriteString("\t\tRevertState: false,\n")
	}

	if g.jobSpec.Type == RetryableRevertibleJob {
		g.builder.WriteString("\t\tForwardRetryThreshold: forwardRetryThreshold,\n")
		g.builder.WriteString("\t\tForwardRetryThresholdCount: make(map[goodjob.TaskID]int),\n")
		g.builder.WriteString("\t\tBackwardRetryThreshold: backwardRetryThreshold,\n")
		g.builder.WriteString("\t\tBackwardRetryThresholdCount: make(map[goodjob.TaskID]int),\n")
		g.builder.WriteString("\t\tRevertState: false,\n")
	}

	g.builder.WriteString("\t}\n")

	g.builder.WriteString("}\n\n")

	// implement job methods

	// implement GetID()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetID() goodjob.JobID {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.ID\n")
	g.builder.WriteString("}\n\n")

	// implement GetTasks()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasks() []goodjob.Task {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.Tasks\n")
	g.builder.WriteString("}\n\n")

	// implement GetTaskArgs()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.TaskArgs[taskID]\n")
	g.builder.WriteString("}\n\n")

	// implement GetVisible()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetVisible() bool {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.Visible\n")
	g.builder.WriteString("}\n\n")

	// implement SetVisible()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetVisible(visible bool) {\n", g.jobSpec.Name))
	g.builder.WriteString("\tj.Visible = visible\n")
	g.builder.WriteString("}\n\n")

	// implement GetLastTask()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetLastTask() goodjob.Task {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.LastTask\n")
	g.builder.WriteString("}\n\n")

	// implement SetLastTask()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetLastTask(task goodjob.Task) {\n", g.jobSpec.Name))
	g.builder.WriteString("\tj.LastTask = task\n")
	g.builder.WriteString("}\n\n")

	// implement GetLastTaskPos()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetLastTaskPos() int {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.LastTaskPos\n")
	g.builder.WriteString("}\n\n")

	// implement SetLastTaskPos()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetLastTaskPos(pos int) {\n", g.jobSpec.Name))
	g.builder.WriteString("\tj.LastTaskPos = pos\n")
	g.builder.WriteString("}\n\n")

	// implement GetLastTaskResult()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetLastTaskResult() *goodjob.TaskResult {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.LastTaskResult\n")
	g.builder.WriteString("}\n\n")

	// implement SetLastTaskResult()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetLastTaskResult(result *goodjob.TaskResult) {\n", g.jobSpec.Name))
	g.builder.WriteString("\tj.LastTaskResult = result\n")
	g.builder.WriteString("}\n\n")

	// implement RetryableJob methods
	if g.jobSpec.Type == RetryableJob {
		// implement GetTasksToRetry()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasksToRetry() []goodjob.Task {\n", g.jobSpec.Name))
		g.builder.WriteString("\tif j.LastTaskResult.Err != nil {\n")
		g.builder.WriteString("\t\treturn j.Tasks[j.LastTaskPos:]\n")
		g.builder.WriteString("\t}\n")
		g.builder.WriteString("\treturn nil\n")
		g.builder.WriteString("}\n\n")

		// implement IncreaseRetryCount()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) IncreaseRetryCount(taskID goodjob.TaskID) {\n", g.jobSpec.Name))
		g.builder.WriteString("\tj.RetryThresholdCount[taskID]++\n")
		g.builder.WriteString("}\n\n")

		// implement RetryThresholdReached()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) RetryThresholdReached(taskID goodjob.TaskID) bool {\n", g.jobSpec.Name))
		g.builder.WriteString("\treturn j.RetryThresholdCount[taskID] >= j.RetryThreshold")
		g.builder.WriteString("}\n\n")
	}

	// implement RevertibleJob methods
	if g.jobSpec.Type == RevertibleJob {
		// implement GetTasksToRevert()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasksToRevert() []goodjob.RevertibleTask {\n", g.jobSpec.Name))
		g.builder.WriteString("\tvar tasksToRevert []goodjob.RevertibleTask\n")
		g.builder.WriteString("\tfor i := j.LastTaskPos - 1; i >= 0; i-- {\n")
		g.builder.WriteString("\t\ttask := j.Tasks[i]\n")
		g.builder.WriteString("\t\tswitch t := task.(type) {\n")
		g.builder.WriteString("\t\t\tcase goodjob.RevertibleTask:\n")
		g.builder.WriteString("\t\t\t\ttasksToRevert = append(tasksToRevert, t)\n")
		g.builder.WriteString("\t\t\t}\n")
		g.builder.WriteString("\t\t}\n")
		g.builder.WriteString("\treturn tasksToRevert\n")
		g.builder.WriteString("}\n\n")

		// implement GetRevertState()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetRevertState() bool {\n", g.jobSpec.Name))
		g.builder.WriteString("\treturn j.RevertState\n")
		g.builder.WriteString("}\n\n")

		// implement SetRevertState()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) SetRevertState(revert bool) {\n", g.jobSpec.Name))
		g.builder.WriteString("\tj.RevertState = revert\n")
		g.builder.WriteString("}\n\n")
	}

	// implement RevertibleRevertibleJob methods
	if g.jobSpec.Type == RetryableRevertibleJob {
		// implement GetTasksToRetry()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasksToRetry() []goodjob.Task {\n", g.jobSpec.Name))
		g.builder.WriteString("\tif j.LastTaskResult.Err != nil {\n")
		g.builder.WriteString("\t\treturn j.Tasks[j.LastTaskPos:]\n")
		g.builder.WriteString("\t}\n")
		g.builder.WriteString("\treturn nil\n")
		g.builder.WriteString("}\n\n")

		// implement IncreaseRetryCount()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) IncreaseRetryCount(taskID goodjob.TaskID) {\n", g.jobSpec.Name))
		g.builder.WriteString("\tif j.RevertState {\n")
		g.builder.WriteString("\t\tj.BackwardRetryThresholdCount[taskID]++\n")
		g.builder.WriteString("\t} else {\n")
		g.builder.WriteString("\t\tj.ForwardRetryThresholdCount[taskID]++\n")
		g.builder.WriteString("\t}\n")
		g.builder.WriteString("}\n\n")

		// implement RetryThresholdReached()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) RetryThresholdReached(taskID goodjob.TaskID) bool {\n", g.jobSpec.Name))
		g.builder.WriteString("\tif j.RevertState {\n")
		g.builder.WriteString("\t\treturn j.BackwardRetryThresholdCount[taskID] >= j.BackwardRetryThreshold\n")
		g.builder.WriteString("\t} else {\n")
		g.builder.WriteString("\t\treturn j.ForwardRetryThresholdCount[taskID] >= j.ForwardRetryThreshold\n")
		g.builder.WriteString("\t}\n")
		g.builder.WriteString("}\n\n")

		// implement GetTasksToRevert()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasksToRevert() []goodjob.RevertibleTask {\n", g.jobSpec.Name))
		g.builder.WriteString("\tvar tasksToRevert []goodjob.RevertibleTask\n")
		g.builder.WriteString("\tfor i := j.LastTaskPos - 1; i >= 0; i-- {\n")
		g.builder.WriteString("\t\ttask := j.Tasks[i]\n")
		g.builder.WriteString("\t\tswitch t := task.(type) {\n")
		g.builder.WriteString("\t\t\tcase goodjob.RevertibleTask:\n")
		g.builder.WriteString("\t\t\t\ttasksToRevert = append(tasksToRevert, t)\n")
		g.builder.WriteString("\t\t\t}\n")
		g.builder.WriteString("\t\t}\n")
		g.builder.WriteString("\treturn tasksToRevert\n")
		g.builder.WriteString("}\n\n")

		// implement GetRevertState()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) GetRevertState() bool {\n", g.jobSpec.Name))
		g.builder.WriteString("\treturn j.RevertState\n")
		g.builder.WriteString("}\n\n")

		// implement SetRevertState()
		g.builder.WriteString(fmt.Sprintf("func (j *%s) SetRevertState(revert bool) {\n", g.jobSpec.Name))
		g.builder.WriteString("\tj.RevertState = revert\n")
		g.builder.WriteString("}\n\n")
	}

	return nil
}
