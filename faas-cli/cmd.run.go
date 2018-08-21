package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
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
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(dir)
	}()
	srcDir := dir + string(os.PathSeparator) + "src"

	err = prepareBuildSrc(serviceSourcePath, srcDir)
	if err != nil {
		return err
	}

	err = prepareProjectSrc(curDir, srcDir, packageName)
	if err != nil {
		return err
	}

	err = generateProjectImport(srcDir, packageName)
	if err != nil {
		return err
	}

	binaryName := packageName + "_bin"
	err = buildBinaryFile0(srcDir, packageName, dir, binaryName)
	if err != nil {
		return err
	}

	runBinary(srcDir, binaryName)

	return nil
}

func runBinary(dir string, binaryName string) error {

	//TODO passing parameter to run
	cmd := exec.Command(dir + string(os.PathSeparator) + binaryName)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errOut, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for s := range signalChan {
			switch s {
			case os.Interrupt:
				cmd.Process.Signal(os.Interrupt)
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := out.Read(buf)
			if n > 0 {
				fmt.Fprint(os.Stdout, string(buf[:n]))
			}
			if err != nil {
				if err == io.EOF {
					return
				} else {
					panic(err)
				}
			}
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := errOut.Read(buf)
			if n > 0 {
				fmt.Fprint(os.Stderr, string(buf[:n]))
			}
			if err != nil {
				if err == io.EOF {
					return
				} else {
					panic(err)
				}
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}
	fmt.Println("starting service... pid:" + strconv.Itoa(cmd.Process.Pid))
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func buildBinaryFile0(buildSrcDir, packageName string, altGoPath string, binaryName string) error {
	gopath := os.Getenv("GOPATH")
	if len(gopath) > 0 {
		gopath = altGoPath + string(os.PathListSeparator) + gopath
	} else {
		gopath = altGoPath
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := &exec.Cmd{
		Args:   []string{"go", "build", "-o", binaryName}, //TODO need to generate os related binary file name
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

func prepareProjectSrc(curDir string, buildSrcDir, packageName string) error {
	packageFolder := buildSrcDir + string(os.PathSeparator) + packageName
	err := os.MkdirAll(packageFolder, 0755)
	if err != nil {
		return err
	}
	// copy all file in working directory into build directory
	err = filepath.Walk(curDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if curDir == path {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}

		err = copyFile(path, packageFolder+string(os.PathSeparator)+info.Name())
		return err
	})

	return err
}

func prepareBuildSrc(serviceSourcePath string, dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// copy dir files
	err = copyToWithOnlyFiles(serviceSourcePath, dir)

	return nil
}
