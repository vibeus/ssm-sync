package main

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Document struct {
	Resources []Resource `hcl:"resource,block"`
}

type Resource struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}

type SsmResource struct {
	Name  string         `hcl:"name"`
	Type  string         `hcl:"type"`
	Value hcl.Expression `hcl:"value"`
}

type App struct {
	AwsRegion string
	ResFile   string
	VarFile   string
	TfVarFile string
	DataDir   string
	DataFile  string
}

func (app *App) Run() error {
	src, err := os.ReadFile(app.ResFile)
	if err != nil {
		return err
	}

	file, d := hclsyntax.ParseConfig(src, app.ResFile, hcl.InitialPos)
	if d.HasErrors() {
		return d
	}

	var doc Document
	d = gohcl.DecodeBody(file.Body, nil, &doc)
	if d.HasErrors() {
		return d
	}

	var writers []OutputWriter

	varWriter, err := NewVarFileWriter(app.VarFile)
	if err != nil {
		return err
	}
	defer varWriter.Close()
	writers = append(writers, varWriter)

	tfVarWriter, err := NewTfVarFileWriter(app.TfVarFile, app.AwsRegion)
	if err != nil {
		return err
	}
	defer tfVarWriter.Close()
	writers = append(writers, tfVarWriter)

	dataWriter, err := NewDataFileWriter(app.DataDir, app.DataFile)
	if err != nil {
		return err
	}
	defer dataWriter.Close()
	writers = append(writers, dataWriter)

	for _, resource := range doc.Resources {
		if resource.Type != "aws_ssm_parameter" {
			continue
		}

		var ssmRes SsmResource
		d = gohcl.DecodeBody(resource.Config, nil, &ssmRes)
		if d.HasErrors() {
			return d
		}

		resName := resource.Name
		ssmName := ssmRes.Name
		varName := ssmRes.Value.Variables()[0][1].(hcl.TraverseAttr).Name

		for _, writer := range writers {
			if err := writer.Write(resName, ssmName, varName); err != nil {
				return err
			}
		}
	}

	return nil
}
