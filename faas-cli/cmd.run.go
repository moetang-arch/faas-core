package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
)

func run() {
	var gopath string
	flag.StringVar(&gopath, "gopath", "", "gopath for faas-api and faas-service source package location")

	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		fmt.Println("parse command line error:", err)
		os.Exit(4)
	}

	// working directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("get working directory error:", err)
		os.Exit(-1)
	}

	packageName, imports, err := getImports(workDir)
	if err != nil {
		fmt.Println("get imports error:", err)
		os.Exit(-1)
	}
	err = checkImportsValidation(imports, append(golangPkgAllowedList, userPkgAllowedList...))
	if err != nil {
		fmt.Println("check imports error:", err)
		os.Exit(-1)
	}

	// path of faas-core
	// default in gopath
	var faasServiceSourcePath = makeFaasServicePkgPath(getFirstGoPath())

	err = buildBinary(faasServiceSourcePath, workDir, packageName)
	if err != nil {
		fmt.Println("build binary error:", err)
		os.Exit(-1)
	}
}

func buildBinary(serviceSourcePath string, curDir, packageName string) error {
	dir := curDir + string(os.PathSeparator) + ".build"
	defer func() {
		os.RemoveAll(dir)
	}()

	err := prepareBuildSrc(serviceSourcePath, dir)
	if err != nil {
		return err
	}

	err = prepareProjectSrc(curDir, dir)
	if err != nil {
		return err
	}

	err = generateProjectImport(dir, packageName)
	if err != nil {
		return err
	}

	err = buildBinaryFile0(dir, packageName)
	if err != nil {
		return err
	}

	err = copyBinaryToWd(dir, curDir, packageName)
	if err != nil {
		return err
	}

	return nil
}

func copyBinaryToWd(buildSrcDir string, curDir, packageName string) error {
	//TODO
	return nil
}

func buildBinaryFile0(buildSrcDir, packageName string) error {
	gopath := os.Getenv("GOPATH")
	if len(gopath) > 0 {
		gopath = buildSrcDir + string(os.PathListSeparator) + gopath
	} else {
		gopath = buildSrcDir
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := &exec.Cmd{
		Args:   []string{"go", "build", "-o", packageName}, //TODO need to generate os related binary file name
		Env:    append(os.Environ(), "GOPATH="+gopath),
		Dir:    buildSrcDir,
		Stderr: errBuf,
		Stdout: outBuf,
		Stdin:  nil,
	}
	if lp, err := exec.LookPath("go"); err != nil {
		return err
	} else {
		cmd.Path = lp
	}
	err := cmd.Run()
	fmt.Println(outBuf.String())
	fmt.Fprintln(os.Stderr, errBuf.String())
	if err != nil {
		return err
	}
	return nil
}

func generateProjectImport(buildSrcDir, packageName string) error {
	t, err := template.New("generatedSource").Parse(generatedSourceTemplate)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, packageName)
	if err != nil {
		return err
	}

	fileContent := buf.String()
	filePath := buildSrcDir + string(os.PathSeparator) + "genimport.go"
	err = ioutil.WriteFile(filePath, []byte(fileContent), 0664)
	if err != nil {
		return err
	}

	return nil
}

func prepareProjectSrc(curDir string, buildSrcDir string) error {
	//TODO
	return nil
}

func prepareBuildSrc(serviceSourcePath string, dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.Mkdir(dir, 0755)
	if err != nil {
		return err
	}

	// copy dir
	err = copyTo(serviceSourcePath, dir)

	return nil
}
