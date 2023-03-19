package main

import (
	"fmt"
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

var (
	defaultImports = []string{
		"fmt",
		"github.com/oklog/ulid/v2",
		"github.com/mtvarkovsky/goodjob/pkg/interfaces",
	}

	jobStructFieldToType = map[string]string{
		"JobID":          "interfaces.JobID",
		"TaskID":         "interfaces.TaskID",
		"Tasks":          "[]interfaces.Task",
		"TaskArgs":       "map[interfaces.TaskID]*interfaces.TaskArg",
		"JobArgs":        "[]interfaces.JobArg",
		"LastTask":       "interfaces.Task",
		"LastTaskPos":    "int",
		"LastTaskResult": "*interfaces.TaskResult",
		"Visible":        "bool",
		"Job":            "interfaces.Job",
	}

	jobStructFields = map[string]string{
		"ID":             "interfaces.JobID",
		"JobArgs":        "[]interfaces.JobArg",
		"Tasks":          "[]interfaces.Task",
		"TaskArgs":       "map[interfaces.TaskID][]*interfaces.TaskArg",
		"Visible":        "bool",
		"LastTask":       "interfaces.Task",
		"LastTaskPos":    "int",
		"LastTaskResult": "*interfaces.TaskResult",
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
		os.Exit(65)
		return
	}
	defer f.Close()

	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(65)
		return
	}

	spec := jobSpec{}
	err = yaml.Unmarshal(specBytes, &spec)
	if err != nil {
		os.Exit(65)
		return
	}

	gen := newGenerator(pkg, destination, &spec)
	gen.generateJob()

	_, _ = fmt.Fprint(f, gen.builder.String())
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
	g.builder.WriteString("}\n\n")

	// define job constructor
	g.builder.WriteString(fmt.Sprintf("func New%s(", g.jobSpec.Name))
	for i, arg := range g.jobSpec.Arguments {
		g.builder.WriteString(fmt.Sprintf("%s %s", arg.Name, arg.Type))
		if i != len(g.jobSpec.Arguments)-1 {
			g.builder.WriteString(", ")
		}
	}
	g.builder.WriteString(fmt.Sprintf(") %s {\n", jobStructFieldToType[g.jobSpec.Type]))
	// prepare job ID
	g.builder.WriteString(fmt.Sprintf("\tjobID := interfaces.JobID(fmt.Sprintf(\"%s", g.jobSpec.ID))
	g.builder.WriteString(" (%s)\", ulid.Make()))\n")
	// prepare job args
	g.builder.WriteString("\tjobArgs := []interfaces.JobArg{\n")
	for _, jobArg := range g.jobSpec.Arguments {
		g.builder.WriteString("\t\t{\n")
		g.builder.WriteString(fmt.Sprintf("\t\t\tName: \"%s\",\n", jobArg.Name))
		g.builder.WriteString(fmt.Sprintf("\t\t\tValue: %s,\n", jobArg.Name))
		g.builder.WriteString("\t\t},\n")
	}
	g.builder.WriteString("\t}\n\n")
	// prepare task IDs
	g.builder.WriteString("\ttaskIDs := map[string]interfaces.TaskID{\n")
	for _, task := range g.jobSpec.Tasks {
		g.builder.WriteString(fmt.Sprintf("\t\t\"%s\": interfaces.TaskID(fmt.Sprintf(\"%s", task.ID, task.ID))
		g.builder.WriteString(" (%s)\", ulid.Make())),\n")
	}
	g.builder.WriteString("\t}\n")

	// prepare tasks
	g.builder.WriteString("\ttasks := []interfaces.Task{\n")
	for _, task := range g.jobSpec.Tasks {
		g.builder.WriteString(fmt.Sprintf("\t\t%s{\n", task.Task))
		g.builder.WriteString(fmt.Sprintf("\t\t\tID: taskIDs[\"%s\"],\n", task.ID))
		g.builder.WriteString(fmt.Sprintf("\t\t\tJobID: jobID,\n"))
		g.builder.WriteString("\t\t},\n")
	}
	g.builder.WriteString("\t}\n")

	// prepare task args
	g.builder.WriteString("\ttaskArgs := map[interfaces.TaskID][]*interfaces.TaskArg{\n")
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
	g.builder.WriteString("\t}\n")

	g.builder.WriteString("}\n\n")

	// implement job methods

	// implement GetID()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetID() interfaces.JobID {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.ID\n")
	g.builder.WriteString("}\n\n")

	// implement GetTasks()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTasks() []interfaces.Task {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.Tasks\n")
	g.builder.WriteString("}\n\n")

	// implement GetTaskArgs()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetTaskArgs(taskID interfaces.TaskID) []*interfaces.TaskArg {\n", g.jobSpec.Name))
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
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetLastTask() interfaces.Task {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.LastTask\n")
	g.builder.WriteString("}\n\n")

	// implement SetLastTask()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetLastTask(task interfaces.Task) {\n", g.jobSpec.Name))
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
	g.builder.WriteString(fmt.Sprintf("func (j *%s) GetLastTaskResult() *interfaces.TaskResult {\n", g.jobSpec.Name))
	g.builder.WriteString("\treturn j.LastTaskResult\n")
	g.builder.WriteString("}\n\n")

	// implement SetLastTaskResult()
	g.builder.WriteString(fmt.Sprintf("func (j *%s) SetLastTaskResult(result *interfaces.TaskResult) {\n", g.jobSpec.Name))
	g.builder.WriteString("\tj.LastTaskResult = result\n")
	g.builder.WriteString("}\n\n")

	return nil
}
