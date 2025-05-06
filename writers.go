package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const GeneratedContentNotice = "# DO NOT ADD NEW VARIABLES TO THIS FILE. THEY ARE AUTO GENERATED. READ: https://github.com/vibeus/ssm-sync"

type OutputWriter interface {
	io.Closer
	Write(resName, ssmName, varName string) error
}

type VarFileWriter struct {
	out io.WriteCloser
}

func NewVarFileWriter(file string) (*VarFileWriter, error) {
	out, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(out, GeneratedContentNotice)

	return &VarFileWriter{out: out}, nil
}

func (w *VarFileWriter) Close() error {
	return w.out.Close()
}

func (w *VarFileWriter) Write(resName, ssmName, varName string) error {
	fmt.Fprintln(w.out)
	fmt.Fprintf(w.out, `variable "%s" {
	type = string
	sensitive = true
	description = "A secret store in SSM. Use ssm-sync to sync before apply."
}`, varName)
	fmt.Fprintln(w.out)
	return nil
}

type TfVarFileWriter struct {
	out    io.WriteCloser
	ssmSvc *ssm.SSM
}

func NewTfVarFileWriter(file, region string) (*TfVarFileWriter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	ssmSvc := ssm.New(sess)

	out, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	return &TfVarFileWriter{out: out, ssmSvc: ssmSvc}, nil
}

func (w *TfVarFileWriter) Close() error {
	return w.out.Close()
}

func (w *TfVarFileWriter) Write(resName, ssmName, varName string) error {
	param, err := w.ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(ssmName),
		WithDecryption: aws.Bool(true),
	})

	var varValue string

	if err != nil {
		if _, ok := err.(*ssm.ParameterNotFound); ok {
			fmt.Printf("%s is a new variable. Please update its value and run 'terraform apply'.\n", varName)
			varValue = "-"
		} else {
			return err
		}
	} else {
		varValue = *param.Parameter.Value
	}

	fmt.Fprintf(w.out, `%s = %s`, varName, strconv.Quote(varValue))
	fmt.Fprintln(w.out)
	return nil
}

type DataFileWriter struct {
	out io.WriteCloser
}

func NewDataFileWriter(dir, file string) (*DataFileWriter, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	out, err := os.Create(path.Join(dir, file))
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(out, GeneratedContentNotice)
	fmt.Fprintln(out)

	return &DataFileWriter{out: out}, nil
}

func (w *DataFileWriter) Close() error {
	return w.out.Close()
}

func (w *DataFileWriter) Write(resName, ssmName, varName string) error {
	fmt.Fprintf(w.out, `data "aws_ssm_parameter" "%s" { name = "%s" }`, resName, ssmName)
	fmt.Fprintln(w.out)
	fmt.Fprintf(w.out, `output "%s" { value = data.aws_ssm_parameter.%s }`, resName, resName)
	fmt.Fprintln(w.out)
	return nil
}
