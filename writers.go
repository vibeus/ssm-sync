package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

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

	return &VarFileWriter{out: out}, nil
}

func (w *VarFileWriter) Close() error {
	return w.out.Close()
}

func (w *VarFileWriter) Write(resName, ssmName, varName string) error {
	fmt.Fprintf(w.out, `variable "%s" { type = string }`, varName)
	fmt.Fprintln(w.out)
	return nil
}

type TfVarFileWriter struct {
	out    io.WriteCloser
	ssmSvc *ssm.SSM
}

func NewTfVarFileWriter(file, region string) (*TfVarFileWriter, error) {
	out, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	ssmSvc := ssm.New(sess)

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

	if err != nil {
		return err
	}

	fmt.Fprintf(w.out, `%s = "%s"`, varName, *param.Parameter.Value)
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
