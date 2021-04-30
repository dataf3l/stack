package echo

import (
	"fmt"
	"log"
	"path"
	"strings"

	//	"runtime/debug"
	//"reflect"
	"os"
	"runtime"
)

// Empty is used to determine the current package name
type Empty struct{}

/*

The input to this function kinda looks like this:

So,

goroutine 1 [running]:
main.echo(0x10cd19d, 0x4)			        <- given this line (which says main.echo)
	/Users/b/work/ntutree/stack/stack.go:20 +0x6f
main.CFUNCTIONNAME(...)					<- return this line (which is the user's function name)
	/Users/b/work/ntutree/stack/stack.go:34
main.b(...)
	/Users/b/work/ntutree/stack/stack.go:16
main.a(...)
	/Users/b/work/ntutree/stack/stack.go:13
main.main()
	/Users/b/work/ntutree/stack/stack.go:37 +0x3c


However, our function name could change, so we use reflection for that.

*/

func getCallerName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next() //surprisingly, not an err
	//return fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
	return fmt.Sprintf("%s", frame.Function)
}

// Errcho prints an error
func Errcho(a error) {
	log.Println(a)
}

// No is just an alias for Echo
func No(str string) {
	Echo(str)
}

// Echo will print a nice error message with the function name that called it
func Echo(str string) {
	programName := path.Base(os.Args[0])
	thisFunctionName := getCallerName()

	//currentPackage := reflect.TypeOf(Empty{}).PkgPath() // just in case somebody renames this

	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	//log.Printf("%s", buf)
	sbuf := string(buf)
	l := strings.Split(sbuf, "\n")
	// parse line
	//-main.echo(0x10cd19d, 0x4)
	//2021/04/28 20:18:52 -	/Users/b/work/ntutree/stack/stack.go:20 +0x6f

	pos := 9999
	foundPos := -1
	for i := range l {
		if foundPos+2 == pos {
			//l[i]
			// get function name from line
			parts := strings.Split(l[i], "(")
			fullFunctionName := strings.Split(parts[0], ".")
			functionName := "N/A"
			if len(fullFunctionName) >= 2 {
				_ = fullFunctionName[0] // module name is here
				functionName = fullFunctionName[1]
			}
			log.Println(programName + "." + functionName + "() : " + str) // todo add function name? moduleName + "/" +
			break
		}
		searchFor := thisFunctionName + "(" // i.e. stack/echo.Echo
		if strings.HasPrefix(l[i], searchFor) {
			foundPos = pos
		}
		//log.Println("-" + l[i])
		pos++
	}

	//debug.PrintStack()
	//log.Println(">>" + str)
}

// https://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
