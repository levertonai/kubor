package command

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
)

func init() {
	cmd := &Get{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Get struct {
	Command

	Predicate  common.EvaluatingPredicate
	SourceHint bool
}

func (instance *Get) ConfigureCliCommands(context string, hc common.HasCommands, _ string) error {
	if context != "" {
		return nil
	}

	cmd := hc.Command("get", "Get the instances of this project using the provided values.").
		Action(instance.ExecuteFromCli)

	cmd.Flag("sourceHint", "Prints to the output a comment which indicates where the rendered content organically comes from.").
		Envar("KUBOR_SOURCE_HINT").
		Default(fmt.Sprint(instance.SourceHint)).
		BoolVar(&instance.SourceHint)
	cmd.Flag("predicate", "Filters every object that should be listed. Empty allows everything. Pattern: \"[!]<template>=<must match regex>\", Example: \"{{.spec.name}}=Foo.*\"").
		Short('p').
		Envar("KUBOR_PREDICATE").
		SetValue(&instance.Predicate)

	return nil
}

func (instance *Get) RunWithArguments(arguments Arguments) error {
	task := &getTask{
		source:    instance,
		arguments: arguments,
		first:     true,
	}
	oh, err := model.NewObjectHandler(task.onObject, arguments.Project)
	if err != nil {
		return err
	}

	cp, err := arguments.Project.RenderedTemplatesProvider()
	if err != nil {
		return err
	}

	return oh.Handle(cp)
}

type getTask struct {
	source    *Get
	arguments Arguments
	first     bool
}

func (instance *getTask) onObject(source string, _ runtime.Object, unstructured *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(unstructured.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	resource, err := kubernetes.GetObjectResource(unstructured, instance.arguments.DynamicClient, instance.arguments.Project.Scheme)
	if err != nil {
		return err
	}
	ul, err := resource.Get(nil)
	if err != nil {
		return err
	}

	if instance.first {
		instance.first = false
	} else {
		fmt.Print("---\n")
	}
	if instance.source.SourceHint {
		fmt.Printf(sourceHintTemplate, source)
	}

	if err := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: true}).
		Encode(ul, os.Stdout); err != nil {
		return err
	}

	return nil
}
